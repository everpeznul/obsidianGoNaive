package domain

import "testing"

func TestIsDay(t *testing.T) {
	// Таблица тестовых случаев
	tests := []struct {
		name     string // описание теста
		input    string // входные данные
		expected bool   // ожидаемый результат
	}{
		{name: "valid day format", input: "2025-12-28", expected: true},
		{name: "invalid month", input: "2025-13-01", expected: false},
		{name: "invalid day", input: "2025-12-32", expected: false},
		{name: "invalid format", input: "2025/12/28", expected: false},
		{name: "empty string", input: "", expected: false},
	}

	// Запускаем каждый тест
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsDay(tt.input)
			if result != tt.expected {
				t.Errorf("IsDay(%q) = %v, want %v",
					tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsMonth(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"valid month", "2025-12", true},
		{"valid month with zero", "2025-01", true},
		{"invalid month", "2025-13", false},
		{"invalid format", "2025-1", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsMonth(tt.input); got != tt.expected {
				t.Errorf("IsMonth(%q) = %v, want %v",
					tt.input, got, tt.expected)
			}
		})
	}
}
