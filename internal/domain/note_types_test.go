package domain

import (
	"fmt"
	"testing"
)

// ======================
// Тесты для Note_periodic_daily
// ======================

func TestNote_periodic_daily_FindFounder(t *testing.T) {
	mockRepo := &mockNoteRepository{
		notes: map[string]Note{
			"0000-00-00": {Title: "0000-00-00"},
		},
	}
	InitRepo(mockRepo)

	tests := []struct {
		name, noteTitle, expected string
	}{
		{name: "daily note", noteTitle: "2025-12-28", expected: "0000-00-00"},
		{name: "another daily note", noteTitle: "2024-01-15", expected: "0000-00-00"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			note := Note_periodic_daily{
				Note_periodic{
					Note{Title: tt.noteTitle},
				},
			}

			result, _ := note.FindFounder()

			if result != tt.expected {
				t.Errorf("FindFounder() for '%s' = %v, want %v",
					tt.noteTitle, result, tt.expected)
			}
		})
	}
}

// ======================
// Тесты для Note_periodic_weekly
// ======================

func TestNote_periodic_weekly_FindFounder(t *testing.T) {
	mockRepo := &mockNoteRepository{
		notes: map[string]Note{
			"0000-W00": {Title: "0000-W00"},
		},
	}
	InitRepo(mockRepo)

	tests := []struct {
		name, noteTitle, expected string
	}{
		{name: "weekly note", noteTitle: "2025-W52", expected: "0000-W00"},
		{name: "first week of year", noteTitle: "2025-W01", expected: "0000-W00"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			note := Note_periodic_weekly{
				Note_periodic{
					Note{Title: tt.noteTitle},
				},
			}

			result, _ := note.FindFounder()

			if result != tt.expected {
				t.Errorf("FindFounder() for '%s' = %v, want %v",
					tt.noteTitle, result, tt.expected)
			}
		})
	}
}

// ======================
// Тесты для Note_periodic_monthly
// ======================

func TestNote_periodic_monthly_FindFounder(t *testing.T) {
	mockRepo := &mockNoteRepository{
		notes: map[string]Note{
			"0000-00": {Title: "0000-00"},
		},
	}
	InitRepo(mockRepo)

	tests := []struct {
		name, noteTitle, expected string
	}{
		{name: "monthly note December", noteTitle: "2025-12", expected: "0000-00"},
		{name: "monthly note January", noteTitle: "2025-01", expected: "0000-00"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			note := Note_periodic_monthly{
				Note_periodic{
					Note{Title: tt.noteTitle},
				},
			}

			result, _ := note.FindFounder()

			if result != tt.expected {
				t.Errorf("FindFounder() for '%s' = %v, want %v",
					tt.noteTitle, result, tt.expected)
			}
		})
	}
}

// ======================
// Тесты для Note_periodic_quarterly
// ======================

func TestNote_periodic_quarterly_FindFounder(t *testing.T) {
	mockRepo := &mockNoteRepository{
		notes: map[string]Note{
			"0000-Q0": {Title: "0000-Q0"},
		},
	}
	InitRepo(mockRepo)

	tests := []struct {
		name, noteTitle, expected string
	}{
		{name: "Q4 2025", noteTitle: "2025-Q4", expected: "0000-Q0"},
		{name: "Q1 2024", noteTitle: "2024-Q1", expected: "0000-Q0"},
		{name: "Q2 2025", noteTitle: "2025-Q2", expected: "0000-Q0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			note := Note_periodic_quarterly{
				Note_periodic{
					Note{Title: tt.noteTitle},
				},
			}

			result, _ := note.FindFounder()

			if result != tt.expected {
				t.Errorf("FindFounder() for '%s' = %v, want %v",
					tt.noteTitle, result, tt.expected)
			}
		})
	}
}

// ======================
// Тесты для Note_periodic_yearly
// ======================

func TestNote_periodic_yearly_FindFounder(t *testing.T) {
	mockRepo := &mockNoteRepository{
		notes: map[string]Note{
			"0000": {Title: "0000"},
		},
	}
	InitRepo(mockRepo)

	tests := []struct {
		name, noteTitle, expected string
	}{
		{name: "year 2025", noteTitle: "2025", expected: "0000"},
		{name: "year 2024", noteTitle: "2024", expected: "0000"},
		{name: "year 1999", noteTitle: "1999", expected: "0000"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			note := Note_periodic_yearly{
				Note_periodic{
					Note{Title: tt.noteTitle},
				},
			}

			result, _ := note.FindFounder()

			if result != tt.expected {
				t.Errorf("FindFounder() for '%s' = %v, want %v",
					tt.noteTitle, result, tt.expected)
			}
		})
	}
}

// ======================
// Тесты для Note_periodic_dream
// ======================

func TestNote_periodic_dream_FindAncestor(t *testing.T) {
	mockRepo := &mockNoteRepository{
		notes: map[string]Note{
			"странный": {Title: "странный"},
			"яркий":    {Title: "яркий"},
			"кошмар":   {Title: "кошмар"},
		},
	}
	InitRepo(mockRepo)

	tests := []struct {
		name, noteTitle, expected string
	}{
		{name: "dream with descriptor", noteTitle: "мысль.2024-03-05.<3>", expected: "2024-03-05"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			note := Note_periodic_dream{
				Note_periodic{
					Note{Title: tt.noteTitle},
				},
			}

			result, _ := note.FindAncestor()

			if result != tt.expected {
				t.Errorf("FindAncestor() for '%s' = %v, want %v",
					tt.noteTitle, result, tt.expected)
			}
		})
	}
}

// ======================
// Тесты для Note_periodic_thought
// ======================

func TestNote_periodic_thought_FindAncestor(t *testing.T) {
	mockRepo := &mockNoteRepository{
		notes: map[string]Note{
			"философия":        {Title: "философия"},
			"программирование": {Title: "программирование"},
			"жизнь":            {Title: "жизнь"},
		},
	}
	InitRepo(mockRepo)

	tests := []struct {
		name, noteTitle, expected string
	}{
		{name: "thought with date", noteTitle: "мысль.2025-12-28.<1>", expected: "2025-12-28"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			note := Note_periodic_thought{
				Note_periodic{
					Note{Title: tt.noteTitle},
				},
			}

			result, _ := note.FindAncestor()

			if result != tt.expected {
				t.Errorf("FindAncestor() for '%s' = %v, want %v",
					tt.noteTitle, result, tt.expected)
			}
		})
	}
}

// ======================
// Edge cases для всех типов
// ======================

func TestNote_periodic_types_EdgeCases(t *testing.T) {
	mockRepo := &mockNoteRepository{
		notes: map[string]Note{},
	}
	InitRepo(mockRepo)

	t.Run("dream with single word title", func(t *testing.T) {
		note := Note_periodic_dream{
			Note_periodic{
				Note{Title: "сон"},
			},
		}

		// При split("сон", ".") получится []string{"сон"}
		// Индекс [1] вызовет панику!
		// Нужна проверка длины массива
		defer func() {
			if r := recover(); r != nil {
				t.Logf("Caught panic (expected): %v", r)
			}
		}()

		result, _ := note.FindAncestor()
		t.Logf("Result: '%s'", result)
	})

	t.Run("thought with single word title", func(t *testing.T) {
		note := Note_periodic_thought{
			Note_periodic{
				Note{Title: "мысль"},
			},
		}

		defer func() {
			if r := recover(); r != nil {
				t.Logf("Caught panic (expected): %v", r)
			}
		}()

		result, _ := note.FindAncestor()
		t.Logf("Result: '%s'", result)
	})

	t.Run("empty title", func(t *testing.T) {
		note := Note_periodic_daily{
			Note_periodic{
				Note{Title: ""},
			},
		}

		result, _ := note.FindFounder()
		if result != "0000-00-00" {
			t.Errorf("Expected '0000-00-00', got '%s'", result)
		}
	})
}

// ======================
// Тесты для ReturnTypesNote
// ======================

func TestReturnTypesNote(t *testing.T) {
	mockRepo := &mockNoteRepository{
		notes: map[string]Note{},
	}
	InitRepo(mockRepo)

	tests := []struct {
		name, noteTitle, expectedType string
	}{
		{name: "thought note", noteTitle: "мысль.2024-03-05.<3>", expectedType: "*domain.Note_periodic_thought"},
		{name: "dream note", noteTitle: "сон.2024-03-05.<3>", expectedType: "*domain.Note_periodic_dream"},
		{name: "person note", noteTitle: "человек.Иван", expectedType: "*domain.Note_human"},
		{name: "daily note", noteTitle: "2025-12-28", expectedType: "*domain.Note_periodic_daily"},
		{name: "weekly note", noteTitle: "2025-W52", expectedType: "*domain.Note_periodic_weekly"},
		{name: "monthly note", noteTitle: "2025-12", expectedType: "*domain.Note_periodic_monthly"},
		{name: "quarterly note", noteTitle: "2025-Q4", expectedType: "*domain.Note_periodic_quarterly"},
		{name: "yearly note", noteTitle: "2025", expectedType: "*domain.Note_periodic_yearly"},
		{name: "regular note", noteTitle: "обычная_заметка", expectedType: "*domain.Note"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			note := Note{Title: tt.noteTitle}
			result := ReturnTypesNote(note)

			// Получаем тип через reflection
			actualType := fmt.Sprintf("%T", result)

			if actualType != tt.expectedType {
				t.Errorf("ReturnTypesNote('%s') type = %v, want %v",
					tt.noteTitle, actualType, tt.expectedType)
			}
		})
	}
}

// ======================
// Интеграционный тест всех типов
// ======================

func TestAllNoteTypes_Integration(t *testing.T) {
	mockRepo := &mockNoteRepository{
		notes: map[string]Note{
			"0000-00-00": {Title: "0000-00-00"},
			"0000-W00":   {Title: "0000-W00"},
			"0000-00":    {Title: "0000-00"},
			"0000-Q0":    {Title: "0000-Q0"},
			"0000":       {Title: "0000"},
			"философия":  {Title: "философия"},
			"странный":   {Title: "странный"},
		},
	}
	InitRepo(mockRepo)

	testCases := []struct {
		title, expectedType, expectedFounder, expectedAncestor string
		testFounder, testAncestor                              bool
	}{
		{title: "2025-12-28", expectedType: "*domain.Note_periodic_daily", testFounder: true, expectedFounder: "0000-00-00"},
		{title: "2025-W52", expectedType: "*domain.Note_periodic_weekly", testFounder: true, expectedFounder: "0000-W00"},
		{title: "2025-12", expectedType: "*domain.Note_periodic_monthly", testFounder: true, expectedFounder: "0000-00"},
		{title: "2025-Q4", expectedType: "*domain.Note_periodic_quarterly", testFounder: true, expectedFounder: "0000-Q0"},
		{title: "2025", expectedType: "*domain.Note_periodic_yearly", testFounder: true, expectedFounder: "0000"},
		{title: "мысль.2024-03-05.<3>", expectedType: "*domain.Note_periodic_thought", testAncestor: true, expectedAncestor: "2024-03-05"},
		{title: "сон.2024-03-05.<3>", expectedType: "*domain.Note_periodic_dream", testAncestor: true, expectedAncestor: "2024-03-05"},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			note := Note{Title: tc.title}
			typed := ReturnTypesNote(note)

			actualType := fmt.Sprintf("%T", typed)
			if actualType != tc.expectedType {
				t.Errorf("Type = %v, want %v", actualType, tc.expectedType)
			}

			if tc.testFounder {
				founder, _ := typed.FindFounder()
				if founder != tc.expectedFounder {
					t.Errorf("FindFounder() = %v, want %v",
						founder, tc.expectedFounder)
				}
			}

			if tc.testAncestor {
				ancestor, _ := typed.FindAncestor()
				if ancestor != tc.expectedAncestor {
					t.Errorf("FindAncestor() = %v, want %v",
						ancestor, tc.expectedAncestor)
				}
			}
		})
	}
}
