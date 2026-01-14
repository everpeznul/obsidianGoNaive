package app

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net"
	pb "obsidianGoNaive/protos/gen/notes"
	"obsidianGoNaive/services/notes/internal/config"
	"obsidianGoNaive/services/notes/internal/repository"
	"obsidianGoNaive/services/notes/internal/repository/postgres"
	"obsidianGoNaive/services/notes/internal/transport"
	"os"
	"path/filepath"

	"google.golang.org/grpc"
)

type App struct {
	Config     *config.Config
	Log        *slog.Logger
	Db         *sql.DB
	Repo       *repository.Repository
	GRPCServer *grpc.Server
	Listener   net.Listener
}

func NewApp() *App {
	config := InitEnvConfig()
	log := InitLog(&config.Log)
	db, err := InitPgDatabase(log, &config.DB)
	if err != nil {

		log.Error("Failed to init database", "error", err)
		os.Exit(1)
	}
	repo := InitPgRepo(db)
	lis, grpcServer, err := InitGRPCServer(log, &config.Net, repo)
	if err != nil {

		log.Error("Failed to init grpc server", "error", err)
		os.Exit(1)
	}
	return &App{
		Config:     config,
		Log:        log,
		Db:         db,
		Repo:       repo,
		GRPCServer: grpcServer,
		Listener:   lis,
	}
}

func InitFileConfig() *config.Config {

	wd, err := os.Getwd()
	if err != nil {

		os.Exit(1)
	}

	cfgPathRel := "/configs/config.yaml"
	path := filepath.Clean(filepath.Join(wd, cfgPathRel))

	cfg, err := config.LoadFileConfig(path)
	if err != nil {

		os.Exit(1)
	}

	return cfg
}

func InitEnvConfig() *config.Config {
	return config.LoadEnvConfig()
}

func InitLog(cfgLog *config.LogConfig) *slog.Logger {

	base := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: cfgLog.NotesLevel}))

	obsiLog := base.With("package", "notes main")

	return obsiLog
}

func InitPgDatabase(obsiLog *slog.Logger, cfgDB *config.DBConfig) (*sql.DB, error) {

	connStr := cfgDB.DSN()
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		obsiLog.Error("Failed to connect to database", "error", err)
		return nil, err
	}

	if err = db.Ping(); err != nil {
		obsiLog.Error("Bad ping database", "error", err)
		return nil, err
	}

	obsiLog.Info("Successful connect to database")

	return db, nil
}

func InitPgRepo(db *sql.DB) *repository.Repository {
	pg := postgres.NewPostgres(db, &postgres.PostgresMapper{})
	return repository.NewRepository(pg)
}

func InitGRPCServer(obsiLog *slog.Logger, cfg *config.NetConfig, repo *repository.Repository) (net.Listener, *grpc.Server, error) {

	notesService := transport.NewNoteService(repo)

	server := grpc.NewServer()
	pb.RegisterNotesServer(server, notesService)

	address := fmt.Sprintf(":%d", cfg.ServerPort)

	lis, err := net.Listen("tcp", address)
	if err != nil {

		obsiLog.Error("Failed to listen", "error", err)
		return nil, nil, err
	}

	return lis, server, nil
}
