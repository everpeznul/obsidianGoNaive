// main.go
package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"obsidianGoNaive/internal/config"
	"obsidianGoNaive/internal/domain"
	"obsidianGoNaive/internal/infrastructure/database"
	obsiHttp "obsidianGoNaive/internal/infrastructure/http"
	"obsidianGoNaive/internal/use_case"
	"os"
	"path/filepath"
	"time"

	_ "github.com/lib/pq"
)

var base = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
var obsiLog = base.With("package", "main")

func createDatabaseConnection(DB config.DBConfig) (*sql.DB, error) {

	// подключение к базе данных
	connStr := DB.DSN()
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия соединения: %w", err)
	}

	// проверка связи с базой данных
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("ошибка ping базы данных: %w", err)
	}

	obsiLog.Info("Successful connect to database")
	return db, nil
}

func main() {

	// инициализация логеров в пакетах
	use_case.SetLog(base.With("package", "use_case"))
	obsiHttp.SetLog(base.With("package", "http"))
	domain.SetLog(base.With("package", "domain"))

	// получение конфигов
	wd, err := os.Getwd()
	if err != nil {

		obsiLog.Error("Getwd ERROR", "error", err)
		panic(nil)
	}
	cfgPathRel := "/configs/config.yaml"
	cfgPath := filepath.Clean(filepath.Join(wd, cfgPathRel))
	cfg, err := config.Load(cfgPath)
	if err != nil {

		obsiLog.Error("config.Load ERROR", "error", err)
	}

	// подключение к базе данных
	db, err := createDatabaseConnection(cfg.DB)
	if err != nil {

		obsiLog.Error("cannot connect to database", "error", err)
		panic(nil)
	}
	defer db.Close()

	// создание и инициализация репозитория и помошников
	repo := &database.PgDB{DB: db}
	domain.InitRepo(repo)
	use_case.InitUpdater(repo)

	// создание сервера
	mux := http.NewServeMux()
	//http.HandleFunc("/", obsiHttp.HomeHandler)
	fs := http.FileServer(http.Dir("./web"))
	mux.Handle("/", fs)
	mux.HandleFunc("/notes/{id}", obsiHttp.NotesUUIDHandler)
	mux.HandleFunc("/notes", obsiHttp.NotesHandler)
	handler := obsiHttp.TimeoutMiddleware(
		obsiHttp.JsonMiddleware(
			mux))
	srv := &http.Server{
		Addr:              ":8080",
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// запуск сервера
	obsiLog.Info("Server start")
	srv.ListenAndServe()
}
