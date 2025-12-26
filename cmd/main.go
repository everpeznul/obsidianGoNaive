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

	//http.HandleFunc("/", obsiHttp.HomeHandler)
	http.HandleFunc("/notes/{id}", obsiHttp.NotesUUIDHandler)
	http.HandleFunc("/notes", obsiHttp.NotesHandler)

	fs := http.FileServer(http.Dir("./web"))
	http.Handle("/", fs)

	http.ListenAndServe(":8080", nil)
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
