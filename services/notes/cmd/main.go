package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net"
	pb "obsidianGoNaive/protos/gen/notes"
	"obsidianGoNaive/services/config"
	"obsidianGoNaive/services/notes/database"
	"obsidianGoNaive/services/notes/domain"
	"os"

	"google.golang.org/grpc"
)

var base = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
var obsiLog = base.With("package", "main notes")

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

	domainLog := obsiLog.With("package", "domain")
	domain.SetLog(domainLog)

	notesService := domain.NewNoteService()

	server := grpc.NewServer()
	pb.RegisterNotesServer(server, notesService)

	// Прокидывание логгеров
	domain.SetLog(base.With("package", "domain"))

	// получение конфигов
	/*wd, err := os.Getwd()
	if err != nil {

		obsiLog.Error("Getwd ERROR", "error", err)
		panic(nil)
	}
		cfgPathRel := "/configs/config.yaml"
		cfgPath := filepath.Clean(filepath.Join(wd, cfgPathRel))
		cfg, err := config.Load(cfgPath)
		if err != nil {

			obsiLog.Error("config.Load ERROR", "error", err)
		}*/
	cfg := config.LoadDBFromEnv()

	// подключение к базе данных
	db, err := createDatabaseConnection(cfg.DB)
	if err != nil {

		obsiLog.Error("cannot connect to database", "error", err)
		panic(nil)
	}
	defer db.Close()

	// создание и инициализация репозитория и помошников
	repo := &database.PgDB{DB: db}
	domain.SetRepo(repo)

	lis, _ := net.Listen("tcp", ":9001")
	server.Serve(lis)
}
