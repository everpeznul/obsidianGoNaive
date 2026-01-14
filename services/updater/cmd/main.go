package main

import (
	"log/slog"
	"net"
	pbn "obsidianGoNaive/pkg/protos/gen/notes"
	pbu "obsidianGoNaive/pkg/protos/gen/updater"
	"obsidianGoNaive/services/updater/use_case"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var base = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
var obsiLog = base.With("package", "main updater")

func main() {

	updaterLog := obsiLog.With("package", "updater")
	use_case.SetLog(updaterLog)

	notesAddr := os.Getenv("NOTES_ADDR")
	if notesAddr == "" {
		notesAddr = "localhost:9001" // для локального запуска
	}

	conn, _ := grpc.NewClient(notesAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	client := pbn.NewNotesClient(conn)

	updater := use_case.NewUpdaterService(client)

	server := grpc.NewServer()
	pbu.RegisterUpdaterServer(server, updater)

	lis, _ := net.Listen("tcp", ":9002")
	server.Serve(lis)

}
