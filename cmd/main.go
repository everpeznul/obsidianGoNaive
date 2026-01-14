// main.go
package main

import (
	"log/slog"
	"net/http"
	pbn "obsidianGoNaive/pkg/protos/gen/notes"
	pbu "obsidianGoNaive/pkg/protos/gen/updater"
	"os"
	"time"

	obsiHttp "obsidianGoNaive/gatewate/http"

	"google.golang.org/grpc"
)

var base = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
var obsiLog = base.With("package", "main gateway")

type HTTPServer struct {
	notesClient   pbn.NotesClient
	updaterClient pbu.UpdaterClient
}

func main() {

	var httpLog = base.With("package", "http")
	obsiHttp.SetLog(httpLog)

	// Получаем адреса из переменных окружения
	notesAddr := os.Getenv("NOTES_ADDR")
	if notesAddr == "" {
		notesAddr = "localhost:9001" // для локального запуска
	}

	updaterAddr := os.Getenv("UPDATER_ADDR")
	if updaterAddr == "" {
		updaterAddr = "localhost:9002"
	}

	notesConn, _ := grpc.NewClient(notesAddr, grpc.WithInsecure())
	notesClient := pbn.NewNotesClient(notesConn)
	updaterConn, _ := grpc.NewClient(updaterAddr, grpc.WithInsecure())
	updaterClient := pbu.NewUpdaterClient(updaterConn)

	// server := &HTTPServer{notesClient, updaterClient}

	gateway := obsiHttp.NewGateway(notesClient, updaterClient)

	// создание сервера
	mux := http.NewServeMux()
	//http.HandleFunc("/", obsiHttp.HomeHandler)
	fs := http.FileServer(http.Dir("./web"))
	mux.Handle("/", fs)
	mux.HandleFunc("/notes/{id}", gateway.NotesUUIDHandler)
	mux.HandleFunc("/notes", gateway.NotesHandler)
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
