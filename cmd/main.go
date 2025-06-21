// main.go
package main

import (
	"database/sql"
	"fmt"
	"log"
	"obsidianGoNaive/internal/infrastructure/database"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"

	"obsidianGoNaive/internal/domain"
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

	// Запуск тестов
	runRepositoryTests(noteRepo)
}

func createDatabaseConnection() (*sql.DB, error) {
	connStr := "user=postgres password=mypass dbname=postgres sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия соединения: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("ошибка ping базы данных: %w", err)
	}

	fmt.Println("✓ Успешное подключение к базе данных")
	return db, nil
}

func runRepositoryTests(repo domain.NoteRepository) {
	fmt.Println("\n=== Тестирование репозитория заметок ===")

	// Создание тестовых данных
	testNotes := createTestNotes()
	var createdIDs []uuid.UUID

	// Тест создания заметок
	fmt.Println("\n1. Создание тестовых заметок...")
	for i, note := range testNotes {
		id, err := repo.Insert(note)
		if err != nil {
			log.Printf("Ошибка создания заметки %d: %v", i+1, err)
			continue
		}
		createdIDs = append(createdIDs, id)
		fmt.Printf("✓ Заметка '%s' создана с ID: %s\n", note.Title, id)
	}

	if len(createdIDs) == 0 {
		log.Fatal("Не удалось создать ни одной заметки")
	}

	// Тест получения по ID
	testGetByID(repo, createdIDs[0])

	// Тест получения всех заметок
	testGetAll(repo)

	// Тест поиска по заголовку
	testGetByTitle(repo, testNotes[0].Title)

	// Тест поиска по предку
	testGetByAncestor(repo, "Тест")

	// Тест обновления
	testUpdate(repo, createdIDs[0])

	// Тест удаления
	testDelete(repo, createdIDs)

	fmt.Println("\n=== Все тесты завершены ===")
}

func createTestNotes() []domain.Note {
	now := time.Now()

	return []domain.Note{
		{
			Title:      "Тестовая заметка 1",
			Path:       "/test/note1.md",
			Class:      "test",
			Tags:       []string{"golang", "тест", "репозиторий"},
			Links:      []string{"link1", "link2"},
			Content:    "Содержимое первой заметки",
			CreateTime: now,
			UpdateTime: now,
		},
		{
			Title:      "Тестовая заметка 2",
			Path:       "/test/note2.md",
			Class:      "example",
			Tags:       []string{"пример", "база данных"},
			Links:      []string{"external-link"},
			Content:    "Пример содержимого",
			CreateTime: now,
			UpdateTime: now,
		},
		{
			Title:      "Другая заметка",
			Path:       "/other/note.md",
			Class:      "other",
			Tags:       []string{"другое"},
			Links:      []string{},
			Content:    "Другое содержимое",
			CreateTime: now,
			UpdateTime: now,
		},
	}
}

func testGetByID(repo domain.NoteRepository, id uuid.UUID) {
	fmt.Println("\n2. Тестирование получения по ID...")
	note, err := repo.GetByID(id)
	if err != nil {
		log.Printf("Ошибка получения заметки: %v", err)
		return
	}
	fmt.Printf("✓ Получена заметка: '%s'\n", note.Title)
	printNoteDetails(note)
}

func testGetAll(repo domain.NoteRepository) {
	fmt.Println("\n3. Тестирование получения всех заметок...")
	notes, err := repo.GetAll()
	if err != nil {
		log.Printf("Ошибка получения всех заметок: %v", err)
		return
	}
	fmt.Printf("✓ Получено заметок: %d\n", len(notes))
	for i, note := range notes {
		fmt.Printf("  %d. %s\n", i+1, note.Title)
	}
}

func testGetByTitle(repo domain.NoteRepository, title string) {
	fmt.Println("\n4. Тестирование поиска по заголовку...")
	note, err := repo.FindByTitle(title)
	if err != nil {
		log.Printf("Ошибка поиска по заголовку: %v", err)
		return
	}
	fmt.Printf("✓ Найдена заметка: '%s'\n", note.Title)
}

func testGetByAncestor(repo domain.NoteRepository, ancestor string) {
	fmt.Println("\n5. Тестирование поиска по предку...")
	notes, err := repo.FindByAncestor(ancestor)
	if err != nil {
		log.Printf("Ошибка поиска по предку: %v", err)
		return
	}
	fmt.Printf("✓ Найдено заметок с предком '%s': %d\n", ancestor, len(notes))
	for _, note := range notes {
		fmt.Printf("  - %s\n", note.Title)
	}
}

func testUpdate(repo domain.NoteRepository, id uuid.UUID) {
	fmt.Println("\n6. Тестирование обновления...")

	// Получаем заметку
	note, err := repo.GetByID(id)
	if err != nil {
		log.Printf("Ошибка получения заметки для обновления: %v", err)
		return
	}

	// Обновляем данные
	note.Title = "Обновленная заметка"
	note.Content = note.Content + " Добавленная строка"
	note.UpdateTime = time.Now()

	// Сохраняем изменения
	err = repo.UpdateById(note)
	if err != nil {
		log.Printf("Ошибка обновления заметки: %v", err)
		return
	}

	fmt.Println("✓ Заметка успешно обновлена")

	// Проверяем обновление
	updatedNote, err := repo.GetByID(id)
	if err != nil {
		log.Printf("Ошибка проверки обновления: %v", err)
		return
	}
	fmt.Printf("✓ Новый заголовок: '%s'\n", updatedNote.Title)
}

func testDelete(repo domain.NoteRepository, ids []uuid.UUID) {
	fmt.Println("\n7. Тестирование удаления...")

	for i, id := range ids {
		err := repo.DeleteByID(id)
		if err != nil {
			log.Printf("Ошибка удаления заметки %d: %v", i+1, err)
			continue
		}
		fmt.Printf("✓ Заметка %d удалена\n", i+1)

		// Проверяем удаление
		_, err = repo.GetByID(id)
		if err != nil {
			fmt.Printf("✓ Подтверждено: заметка %d не найдена\n", i+1)
		} else {
			fmt.Printf("⚠ Предупреждение: заметка %d все еще существует\n", i+1)
		}
	}
}

func printNoteDetails(note domain.Note) {
	fmt.Printf("  ID: %s\n", note.Id)
	fmt.Printf("  Путь: %s\n", note.Path)
	fmt.Printf("  Класс: %s\n", note.Class)
	fmt.Printf("  Теги: %v\n", note.Tags)
	fmt.Printf("  Ссылки: %v\n", note.Links)
	fmt.Printf("  Содержимое: %s\n", note.Content)
}
