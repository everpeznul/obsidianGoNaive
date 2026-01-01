package http

import (
	"obsidianGoNaive/internal/domain"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNoteMapperHttp_HTTPToDomain(t *testing.T) {
	tests := []struct {
		name     string
		input    httpNote
		expected domain.Note
	}{
		{
			name: "converts httpNote to domain.Note with all fields",
			input: httpNote{
				Id:         uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				Title:      "Test Note",
				Path:       "/test/path.md",
				Class:      "article",
				Tags:       []string{"tag1", "tag2"},
				Links:      []string{"link1", "link2"},
				Content:    "Test content",
				CreateTime: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
				UpdateTime: time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC),
			},
			expected: domain.Note{
				Id:         uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				Title:      "Test Note",
				Path:       "/test/path.md",
				Class:      "article",
				Tags:       []string{"tag1", "tag2"},
				Links:      []string{"link1", "link2"},
				Content:    "Test content",
				CreateTime: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
				UpdateTime: time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "converts httpNote with empty slices",
			input: httpNote{
				Id:         uuid.MustParse("223e4567-e89b-12d3-a456-426614174001"),
				Title:      "Empty Note",
				Path:       "/empty.md",
				Class:      "note",
				Tags:       []string{},
				Links:      []string{},
				Content:    "",
				CreateTime: time.Date(2023, 2, 1, 10, 0, 0, 0, time.UTC),
				UpdateTime: time.Date(2023, 2, 1, 10, 0, 0, 0, time.UTC),
			},
			expected: domain.Note{
				Id:         uuid.MustParse("223e4567-e89b-12d3-a456-426614174001"),
				Title:      "Empty Note",
				Path:       "/empty.md",
				Class:      "note",
				Tags:       []string{},
				Links:      []string{},
				Content:    "",
				CreateTime: time.Date(2023, 2, 1, 10, 0, 0, 0, time.UTC),
				UpdateTime: time.Date(2023, 2, 1, 10, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "converts httpNote with nil uuid and zero times",
			input: httpNote{
				Id:         uuid.Nil,
				Title:      "New Note",
				Path:       "/new.md",
				Class:      "draft",
				Tags:       []string{"draft"},
				Links:      []string{},
				Content:    "Draft content",
				CreateTime: time.Time{},
				UpdateTime: time.Time{},
			},
			expected: domain.Note{
				Id:         uuid.Nil,
				Title:      "New Note",
				Path:       "/new.md",
				Class:      "draft",
				Tags:       []string{"draft"},
				Links:      []string{},
				Content:    "Draft content",
				CreateTime: time.Time{},
				UpdateTime: time.Time{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := nm.HTTPToDomain(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("HTTPToDomain() = %+v, want %+v", result, tt.expected)
			}
		})
	}
}

func TestNoteMapperHttp_HTTPToDomainSlice(t *testing.T) {
	tests := []struct {
		name     string
		input    []httpNote
		expected []domain.Note
	}{
		{
			name: "converts multiple httpNotes to domain.Notes",
			input: []httpNote{
				{
					Id:         uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
					Title:      "Note 1",
					Path:       "/path1.md",
					Class:      "article",
					Tags:       []string{"tag1"},
					Links:      []string{"link1"},
					Content:    "Content 1",
					CreateTime: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
					UpdateTime: time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC),
				},
				{
					Id:         uuid.MustParse("223e4567-e89b-12d3-a456-426614174001"),
					Title:      "Note 2",
					Path:       "/path2.md",
					Class:      "note",
					Tags:       []string{"tag2", "tag3"},
					Links:      []string{},
					Content:    "Content 2",
					CreateTime: time.Date(2023, 2, 1, 10, 0, 0, 0, time.UTC),
					UpdateTime: time.Date(2023, 2, 1, 10, 0, 0, 0, time.UTC),
				},
			},
			expected: []domain.Note{
				{
					Id:         uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
					Title:      "Note 1",
					Path:       "/path1.md",
					Class:      "article",
					Tags:       []string{"tag1"},
					Links:      []string{"link1"},
					Content:    "Content 1",
					CreateTime: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
					UpdateTime: time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC),
				},
				{
					Id:         uuid.MustParse("223e4567-e89b-12d3-a456-426614174001"),
					Title:      "Note 2",
					Path:       "/path2.md",
					Class:      "note",
					Tags:       []string{"tag2", "tag3"},
					Links:      []string{},
					Content:    "Content 2",
					CreateTime: time.Date(2023, 2, 1, 10, 0, 0, 0, time.UTC),
					UpdateTime: time.Date(2023, 2, 1, 10, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			name:     "converts empty slice",
			input:    []httpNote{},
			expected: []domain.Note{},
		},
		{
			name: "converts single element slice",
			input: []httpNote{
				{
					Id:         uuid.MustParse("323e4567-e89b-12d3-a456-426614174002"),
					Title:      "Solo Note",
					Path:       "/solo.md",
					Class:      "memo",
					Tags:       []string{},
					Links:      []string{},
					Content:    "Solo content",
					CreateTime: time.Date(2023, 3, 1, 9, 0, 0, 0, time.UTC),
					UpdateTime: time.Date(2023, 3, 1, 9, 0, 0, 0, time.UTC),
				},
			},
			expected: []domain.Note{
				{
					Id:         uuid.MustParse("323e4567-e89b-12d3-a456-426614174002"),
					Title:      "Solo Note",
					Path:       "/solo.md",
					Class:      "memo",
					Tags:       []string{},
					Links:      []string{},
					Content:    "Solo content",
					CreateTime: time.Date(2023, 3, 1, 9, 0, 0, 0, time.UTC),
					UpdateTime: time.Date(2023, 3, 1, 9, 0, 0, 0, time.UTC),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := nm.HTTPToDomainSlice(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("HTTPToDomainSlice() = %+v, want %+v", result, tt.expected)
			}
		})
	}
}

func TestNoteMapperHttp_DomainToHTTP(t *testing.T) {
	tests := []struct {
		name     string
		input    domain.Note
		expected httpNote
	}{
		{
			name: "converts domain.Note to httpNote with all fields",
			input: domain.Note{
				Id:         uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				Title:      "Test Note",
				Path:       "/test/path.md",
				Class:      "article",
				Tags:       []string{"tag1", "tag2"},
				Links:      []string{"link1", "link2"},
				Content:    "Test content",
				CreateTime: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
				UpdateTime: time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC),
			},
			expected: httpNote{
				Id:         uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				Title:      "Test Note",
				Path:       "/test/path.md",
				Class:      "article",
				Tags:       []string{"tag1", "tag2"},
				Links:      []string{"link1", "link2"},
				Content:    "Test content",
				CreateTime: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
				UpdateTime: time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "converts domain.Note with empty slices",
			input: domain.Note{
				Id:         uuid.MustParse("223e4567-e89b-12d3-a456-426614174001"),
				Title:      "Empty Note",
				Path:       "/empty.md",
				Class:      "note",
				Tags:       []string{},
				Links:      []string{},
				Content:    "",
				CreateTime: time.Date(2023, 2, 1, 10, 0, 0, 0, time.UTC),
				UpdateTime: time.Date(2023, 2, 1, 10, 0, 0, 0, time.UTC),
			},
			expected: httpNote{
				Id:         uuid.MustParse("223e4567-e89b-12d3-a456-426614174001"),
				Title:      "Empty Note",
				Path:       "/empty.md",
				Class:      "note",
				Tags:       []string{},
				Links:      []string{},
				Content:    "",
				CreateTime: time.Date(2023, 2, 1, 10, 0, 0, 0, time.UTC),
				UpdateTime: time.Date(2023, 2, 1, 10, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "converts domain.Note with multiline content",
			input: domain.Note{
				Id:         uuid.MustParse("323e4567-e89b-12d3-a456-426614174002"),
				Title:      "Multiline Note",
				Path:       "/multiline.md",
				Class:      "article",
				Tags:       []string{"important"},
				Links:      []string{"ref1", "ref2", "ref3"},
				Content:    "Line 1\nLine 2\nLine 3",
				CreateTime: time.Date(2023, 3, 1, 9, 0, 0, 0, time.UTC),
				UpdateTime: time.Date(2023, 3, 1, 9, 30, 0, 0, time.UTC),
			},
			expected: httpNote{
				Id:         uuid.MustParse("323e4567-e89b-12d3-a456-426614174002"),
				Title:      "Multiline Note",
				Path:       "/multiline.md",
				Class:      "article",
				Tags:       []string{"important"},
				Links:      []string{"ref1", "ref2", "ref3"},
				Content:    "Line 1\nLine 2\nLine 3",
				CreateTime: time.Date(2023, 3, 1, 9, 0, 0, 0, time.UTC),
				UpdateTime: time.Date(2023, 3, 1, 9, 30, 0, 0, time.UTC),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := nm.DomainToHTTP(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("DomainToHTTP() = %+v, want %+v", result, tt.expected)
			}
		})
	}
}

func TestNoteMapperHttp_DomainToHTTPSlice(t *testing.T) {
	tests := []struct {
		name     string
		input    []domain.Note
		expected []httpNote
	}{
		{
			name: "converts multiple domain.Notes to httpNotes",
			input: []domain.Note{
				{
					Id:         uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
					Title:      "Note 1",
					Path:       "/path1.md",
					Class:      "article",
					Tags:       []string{"tag1"},
					Links:      []string{"link1"},
					Content:    "Content 1",
					CreateTime: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
					UpdateTime: time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC),
				},
				{
					Id:         uuid.MustParse("223e4567-e89b-12d3-a456-426614174001"),
					Title:      "Note 2",
					Path:       "/path2.md",
					Class:      "note",
					Tags:       []string{"tag2", "tag3"},
					Links:      []string{},
					Content:    "Content 2",
					CreateTime: time.Date(2023, 2, 1, 10, 0, 0, 0, time.UTC),
					UpdateTime: time.Date(2023, 2, 1, 10, 0, 0, 0, time.UTC),
				},
			},
			expected: []httpNote{
				{
					Id:         uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
					Title:      "Note 1",
					Path:       "/path1.md",
					Class:      "article",
					Tags:       []string{"tag1"},
					Links:      []string{"link1"},
					Content:    "Content 1",
					CreateTime: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
					UpdateTime: time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC),
				},
				{
					Id:         uuid.MustParse("223e4567-e89b-12d3-a456-426614174001"),
					Title:      "Note 2",
					Path:       "/path2.md",
					Class:      "note",
					Tags:       []string{"tag2", "tag3"},
					Links:      []string{},
					Content:    "Content 2",
					CreateTime: time.Date(2023, 2, 1, 10, 0, 0, 0, time.UTC),
					UpdateTime: time.Date(2023, 2, 1, 10, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			name:     "converts empty slice",
			input:    []domain.Note{},
			expected: []httpNote{},
		},
		{
			name: "converts single element slice",
			input: []domain.Note{
				{
					Id:         uuid.MustParse("323e4567-e89b-12d3-a456-426614174002"),
					Title:      "Solo Note",
					Path:       "/solo.md",
					Class:      "memo",
					Tags:       []string{},
					Links:      []string{},
					Content:    "Solo content",
					CreateTime: time.Date(2023, 3, 1, 9, 0, 0, 0, time.UTC),
					UpdateTime: time.Date(2023, 3, 1, 9, 0, 0, 0, time.UTC),
				},
			},
			expected: []httpNote{
				{
					Id:         uuid.MustParse("323e4567-e89b-12d3-a456-426614174002"),
					Title:      "Solo Note",
					Path:       "/solo.md",
					Class:      "memo",
					Tags:       []string{},
					Links:      []string{},
					Content:    "Solo content",
					CreateTime: time.Date(2023, 3, 1, 9, 0, 0, 0, time.UTC),
					UpdateTime: time.Date(2023, 3, 1, 9, 0, 0, 0, time.UTC),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := nm.DomainToHTTPSlice(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("DomainToHTTPSlice() = %+v, want %+v", result, tt.expected)
			}
		})
	}
}

// Benchmark тесты
func BenchmarkHTTPToDomain(b *testing.B) {
	note := httpNote{
		Id:         uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
		Title:      "Benchmark Note",
		Path:       "/benchmark.md",
		Class:      "test",
		Tags:       []string{"tag1", "tag2", "tag3"},
		Links:      []string{"link1", "link2"},
		Content:    "Benchmark content string",
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = nm.HTTPToDomain(note)
	}
}

func BenchmarkHTTPToDomainSlice(b *testing.B) {
	notes := make([]httpNote, 100)
	for i := 0; i < 100; i++ {
		notes[i] = httpNote{
			Id:         uuid.New(),
			Title:      "Note",
			Path:       "/path.md",
			Class:      "test",
			Tags:       []string{"tag1", "tag2"},
			Links:      []string{"link1"},
			Content:    "Content",
			CreateTime: time.Now(),
			UpdateTime: time.Now(),
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = nm.HTTPToDomainSlice(notes)
	}
}

func BenchmarkDomainToHTTP(b *testing.B) {
	note := domain.Note{
		Id:         uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
		Title:      "Benchmark Note",
		Path:       "/benchmark.md",
		Class:      "test",
		Tags:       []string{"tag1", "tag2", "tag3"},
		Links:      []string{"link1", "link2"},
		Content:    "Benchmark content string",
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = nm.DomainToHTTP(note)
	}
}

func BenchmarkDomainToHTTPSlice(b *testing.B) {
	notes := make([]domain.Note, 100)
	for i := 0; i < 100; i++ {
		notes[i] = domain.Note{
			Id:         uuid.New(),
			Title:      "Note",
			Path:       "/path.md",
			Class:      "test",
			Tags:       []string{"tag1", "tag2"},
			Links:      []string{"link1"},
			Content:    "Content",
			CreateTime: time.Now(),
			UpdateTime: time.Now(),
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = nm.DomainToHTTPSlice(notes)
	}
}
