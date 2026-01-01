package database

import (
	"obsidianGoNaive/internal/domain"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

func TestNoteMapperDb_DatabaseToDomain(t *testing.T) {
	tests := []struct {
		name     string
		input    dbNote
		expected domain.Note
	}{
		{
			name: "converts dbNote to domain.Note with all fields",
			input: dbNote{
				Id:         uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				Title:      "Test Note",
				Path:       "/test/path",
				Class:      "article",
				Tags:       pq.StringArray{"tag1", "tag2"},
				Links:      pq.StringArray{"link1", "link2"},
				Content:    pq.StringArray{"Line 1", "Line 2", "Line 3"},
				CreateTime: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
				UpdateTime: time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC),
			},
			expected: domain.Note{
				Id:         uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				Title:      "Test Note",
				Path:       "/test/path",
				Class:      "article",
				Tags:       []string{"tag1", "tag2"},
				Links:      []string{"link1", "link2"},
				Content:    "Line 1Line 2Line 3",
				CreateTime: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
				UpdateTime: time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "converts dbNote with empty arrays",
			input: dbNote{
				Id:         uuid.MustParse("223e4567-e89b-12d3-a456-426614174001"),
				Title:      "Empty Note",
				Path:       "/empty",
				Class:      "note",
				Tags:       pq.StringArray{},
				Links:      pq.StringArray{},
				Content:    pq.StringArray{},
				CreateTime: time.Date(2023, 2, 1, 10, 0, 0, 0, time.UTC),
				UpdateTime: time.Date(2023, 2, 1, 10, 0, 0, 0, time.UTC),
			},
			expected: domain.Note{
				Id:         uuid.MustParse("223e4567-e89b-12d3-a456-426614174001"),
				Title:      "Empty Note",
				Path:       "/empty",
				Class:      "note",
				Tags:       []string{},
				Links:      []string{},
				Content:    "",
				CreateTime: time.Date(2023, 2, 1, 10, 0, 0, 0, time.UTC),
				UpdateTime: time.Date(2023, 2, 1, 10, 0, 0, 0, time.UTC),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := nm.DatabaseToDomain(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("DatabaseToDomain() = %+v, want %+v", result, tt.expected)
			}
		})
	}
}

func TestNoteMapperDb_DatabaseToDomainSlice(t *testing.T) {
	tests := []struct {
		name     string
		input    []dbNote
		expected []domain.Note
	}{
		{
			name: "converts multiple dbNotes to domain.Notes",
			input: []dbNote{
				{
					Id:         uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
					Title:      "Note 1",
					Path:       "/path1",
					Class:      "article",
					Tags:       pq.StringArray{"tag1"},
					Links:      pq.StringArray{"link1"},
					Content:    pq.StringArray{"Content 1"},
					CreateTime: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
					UpdateTime: time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC),
				},
			},
			expected: []domain.Note{
				{
					Id:         uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
					Title:      "Note 1",
					Path:       "/path1",
					Class:      "article",
					Tags:       []string{"tag1"},
					Links:      []string{"link1"},
					Content:    "Content 1",
					CreateTime: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
					UpdateTime: time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			name:     "converts empty slice",
			input:    []dbNote{},
			expected: []domain.Note{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := nm.DatabaseToDomainSlice(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("DatabaseToDomainSlice() = %+v, want %+v", result, tt.expected)
			}
		})
	}
}

func TestNoteMapperDb_DomainToDatabase(t *testing.T) {
	tests := []struct {
		name     string
		input    domain.Note
		expected dbNote
	}{
		{
			name: "converts domain.Note to dbNote with all fields",
			input: domain.Note{
				Id:         uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				Title:      "Test Note",
				Path:       "/test/path",
				Class:      "article",
				Tags:       []string{"tag1", "tag2"},
				Links:      []string{"link1", "link2"},
				Content:    "Complete content string",
				CreateTime: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
				UpdateTime: time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC),
			},
			expected: dbNote{
				Id:         uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				Title:      "Test Note",
				Path:       "/test/path",
				Class:      "article",
				Tags:       pq.StringArray{"tag1", "tag2"},
				Links:      pq.StringArray{"link1", "link2"},
				Content:    pq.StringArray{"Complete content string"},
				CreateTime: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
				UpdateTime: time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := nm.DomainToDatabase(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("DomainToDatabase() = %+v, want %+v", result, tt.expected)
			}
		})
	}
}

func TestNoteMapperDb_DomainToDatabaseSlice(t *testing.T) {
	tests := []struct {
		name     string
		input    []domain.Note
		expected []dbNote
	}{
		{
			name:     "converts empty slice",
			input:    []domain.Note{},
			expected: []dbNote{},
		},
		{
			name: "converts domain.Note to dbNote with all fields",
			input: []domain.Note{
				domain.Note{
					Id:         uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
					Title:      "Test Note",
					Path:       "/test/path",
					Class:      "article",
					Tags:       []string{"tag1", "tag2"},
					Links:      []string{"link1", "link2"},
					Content:    "Complete content string",
					CreateTime: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
					UpdateTime: time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC),
				},
			},
			expected: []dbNote{
				dbNote{
					Id:         uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
					Title:      "Test Note",
					Path:       "/test/path",
					Class:      "article",
					Tags:       pq.StringArray{"tag1", "tag2"},
					Links:      pq.StringArray{"link1", "link2"},
					Content:    pq.StringArray{"Complete content string"},
					CreateTime: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
					UpdateTime: time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := nm.DomainToDatabaseSlice(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("DomainToDatabaseSlice() = %+v, want %+v", result, tt.expected)
			}
		})
	}
}

// Benchmark-тесты остаются без изменений
func BenchmarkDatabaseToDomain(b *testing.B) {
	note := dbNote{
		Id:         uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
		Title:      "Benchmark Note",
		Path:       "/benchmark",
		Class:      "test",
		Tags:       pq.StringArray{"tag1", "tag2"},
		Links:      pq.StringArray{"link1"},
		Content:    pq.StringArray{"Content"},
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = nm.DatabaseToDomain(note)
	}
}
