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

func main() {
	wd, err := os.Getwd()
	cfgRel := "/configs/config.yaml"
	cfgPath := filepath.Clean(filepath.Join(wd, cfgRel))
	cnfg, _ := config.Load(cfgPath)

	base := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	obsiLog := base.With("package", "main")

	use_case.UseCaseSetLog(base.With("package", "use_case"))
	obsiHttp.HttpSetLog(base.With("package", "http"))
	domain.DomainSetLog(base.With("package", "domain"))

	// Создание подключения к базе данных
	db, err := createDatabaseConnection(cnfg.DB)
	if err != nil {

		obsiLog.Error("cannot connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// Создание репозитория через композицию
	noteRepo := &database.PgDB{DB: db}
	domain.InitRepo(noteRepo)
	use_case.InitUpdater(noteRepo)

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

	obsiLog.Info("Server start")
	srv.ListenAndServe()
}

func createDatabaseConnection(DB config.DBConfig) (*sql.DB, error) {

	connStr := DB.DSN()

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия соединения: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("ошибка ping базы данных: %w", err)
	}

	fmt.Println("Успешное подключение к базе данных")

	return db, nil
}
