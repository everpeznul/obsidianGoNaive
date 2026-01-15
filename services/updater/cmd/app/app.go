package app

import (
	"fmt"
	"log/slog"
	"net"
	pbn "obsidianGoNaive/pkg/protos/gen/notes"
	pbu "obsidianGoNaive/pkg/protos/gen/updater"
	"obsidianGoNaive/services/updater/internal/config"
	"obsidianGoNaive/services/updater/internal/domain"
	"obsidianGoNaive/services/updater/internal/transport"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type App struct {
	Config      *config.Config
	Log         *slog.Logger
	Notes       *domain.Notes
	NotesClient pbn.NotesClient
	GRPCServer  *grpc.Server
	Listener    net.Listener
}

func NewApp() *App {
	cfg := InitConfig()
	appLog, domainLog, transportLog := InitLog(cfg.Log)
	notes := InitNotes(domainLog)
	notesClient := InitNotesClient(cfg.Net)
	lis, server, err := InitGRPCServer(notesClient, transportLog, cfg.Net)
	if err != nil {

		appLog.Error("Error initializing gRPC server", "error", err)
		os.Exit(1)
	}

	return &App{
		Config:      cfg,
		Log:         appLog,
		Notes:       notes,
		NotesClient: notesClient,
		GRPCServer:  server,
		Listener:    lis,
	}
}

func InitConfig() *config.Config { return config.LoadEnvConfig() }

func InitLog(cfgLog *config.LogConfig) (*slog.Logger, *slog.Logger, *slog.Logger) {
	base := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: cfgLog.UpdaterLevel}))

	appLog := base.With("package", "updater main")
	domainLog := base.With("package", "updater domain")
	transportLog := base.With("package", "updater transport")

	return appLog, domainLog, transportLog
}

func InitNotes(domainLog *slog.Logger) *domain.Notes {
	return domain.NewNotes(domainLog)
}

func InitNotesClient(cfg *config.NetConfig) pbn.NotesClient {
	host := os.Getenv("NOTES_ADDR")
	if host == "" {
		host = fmt.Sprintf("localhost:%d", cfg.ClientPort)
	}
	conn, _ := grpc.NewClient(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	client := pbn.NewNotesClient(conn)

	return client
}

func InitGRPCServer(client pbn.NotesClient, transportLog *slog.Logger, cfg *config.NetConfig) (net.Listener, *grpc.Server, error) {
	updaterService := transport.NewUpdaterService(client, transportLog)

	server := grpc.NewServer()
	pbu.RegisterUpdaterServer(server, updaterService)

	address := fmt.Sprintf(":%d", cfg.ServerPort)

	lis, err := net.Listen("tcp", address)
	if err != nil {

		transportLog.Error("Failed to listen", "error", err)
		return nil, nil, err
	}

	return lis, server, nil
}
