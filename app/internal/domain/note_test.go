package domain

import (
	"testing"

	"github.com/google/uuid"
)

// Mock репозиторий для тестов
type mockNoteRepository struct {
	notes map[string]Note
}

func (m *mockNoteRepository) FindByName(name string) (Note, error) {
	note, exists := m.notes[name]
	if !exists {
		return Note{}, nil // возвращаем пустую заметку если не найдена
	}
	return note, nil
}

// Заглушки для остальных методов интерфейса
func (m *mockNoteRepository) Insert(note Note) (uuid.UUID, error) {
	return uuid.New(), nil
}

func (m *mockNoteRepository) GetByID(id uuid.UUID) (Note, error) {
	return Note{}, nil
}

func (m *mockNoteRepository) GetAll() ([]Note, error) {
	return []Note{}, nil
}

func (m *mockNoteRepository) UpdateById(note Note) error {
	return nil
}

func (m *mockNoteRepository) DeleteById(id uuid.UUID) error {
	return nil
}

func (m *mockNoteRepository) FindByAncestor(ancestor string) ([]Note, error) {
	return []Note{}, nil
}

// Тесты для FindFounder
func TestNote_FindFounder(t *testing.T) {
	// Инициализируем mock репозиторий с тестовыми данными
	mockRepo := &mockNoteRepository{
		notes: map[string]Note{
			"мысль":   {Title: "мысль"},
			"сон":     {Title: "сон"},
			"человек": {Title: "человек"},
			"2025":    {Title: "2025"},
		},
	}
	InitRepo(mockRepo)

	tests := []struct {
		name, noteTitle, expected string
	}{
		{name: "simple thought note", noteTitle: "мысль", expected: "мысль"},
		{name: "nested thought note", noteTitle: "мысль.философия", expected: "мысль"},
		{name: "deeply nested note", noteTitle: "мысль.философия.экзистенциализм", expected: "мысль"},
		{name: "dream note", noteTitle: "сон.странный", expected: "сон"},
		{name: "person note", noteTitle: "человек.Иван.Петров", expected: "человек"},
		{name: "year note", noteTitle: "2025", expected: "2025"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			note := Note{Title: tt.noteTitle}
			result, _ := note.FindFounder()

			if result != tt.expected {
				t.Errorf("FindFounder() for '%s' = %v, want %v",
					tt.noteTitle, result, tt.expected)
			}
		})
	}
}

// Тесты для FindFather
func TestNote_FindFather(t *testing.T) {
	mockRepo := &mockNoteRepository{
		notes: map[string]Note{
			"мысль":           {Title: "мысль"},
			"мысль.философия": {Title: "мысль.философия"},
		},
	}
	InitRepo(mockRepo)

	tests := []struct {
		name, noteTitle, expected string
	}{
		{name: "two level note", noteTitle: "мысль.философия", expected: "мысль"},
		{name: "three level note", noteTitle: "мысль.философия.экзистенциализм", expected: "мысль.философия"},
		{name: "single level returns empty", noteTitle: "мысль", expected: "мысль"},
		{name: "four level note", noteTitle: "сон.странный.место.дом", expected: "сон.странный.место"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			note := Note{Title: tt.noteTitle}
			result, _ := note.FindFather()

			if result != tt.expected {
				t.Errorf("FindFather() for '%s' = %v, want %v",
					tt.noteTitle, result, tt.expected)
			}
		})
	}
}

// Тесты для FindAncestor
func TestNote_FindAncestor(t *testing.T) {
	mockRepo := &mockNoteRepository{
		notes: map[string]Note{
			"мысль":           {Title: "мысль"},
			"мысль.философия": {Title: "мысль.философия"},
		},
	}
	InitRepo(mockRepo)

	tests := []struct {
		name, noteTitle, expected string
	}{
		{name: "three level note", noteTitle: "мысль.философия.экзистенциализм", expected: "мысль.философия"},
		{name: "another threee level note", noteTitle: "мысль.мяу.мур", expected: "мысль.мяу"},
		{name: "two level note", noteTitle: "мысль.философия", expected: "мысль"},
		{name: "single level returns empty", noteTitle: "мысль", expected: "мысль"},
		{name: "with percentage sign", noteTitle: "сон.%2025-12-28.странный", expected: "сон"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			note := Note{Title: tt.noteTitle}
			result, _ := note.FindAncestor()

			if result != tt.expected {
				t.Errorf("FindAncestor() for '%s' = %v, want %v",
					tt.noteTitle, result, tt.expected)
			}
		})
	}
}
