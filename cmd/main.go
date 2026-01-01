// main.go
package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"obsidianGoNaive/internal/domain"
	"obsidianGoNaive/internal/infrastructure/database"
	obsiHttp "obsidianGoNaive/internal/infrastructure/http"
	"obsidianGoNaive/internal/use_case"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	// Создание подключения к базе данных
	db, err := createDatabaseConnection()
	if err != nil {
		log.Fatal("Ошибка создания подключения:", err)
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
	srv.ListenAndServe()
}

func createDatabaseConnection() (*sql.DB, error) {

	connStr := "host=localhost port=5432 user=postgres password=mypass dbname=postgres sslmode=disable"

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
