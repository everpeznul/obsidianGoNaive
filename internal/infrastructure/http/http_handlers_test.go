package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"obsidianGoNaive/internal/domain"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

// Mock репозитория для тестов
type mockRepo struct {
	getNoteFunc        func(id uuid.UUID) (domain.Note, error)
	getAllFunc         func() ([]domain.Note, error)
	insertFunc         func(note domain.Note) (uuid.UUID, error)
	updateFunc         func(note domain.Note) error
	deleteFunc         func(id uuid.UUID) error
	findByNameFunc     func(name string) (domain.Note, error)
	findByAncestorFunc func(ancestor string) ([]domain.Note, error)
}

func (m *mockRepo) GetByID(id uuid.UUID) (domain.Note, error) {
	if m.getNoteFunc != nil {
		return m.getNoteFunc(id)
	}
	return domain.Note{}, nil
}

func (m *mockRepo) GetAll() ([]domain.Note, error) {
	if m.getAllFunc != nil {
		return m.getAllFunc()
	}
	return []domain.Note{}, nil
}

func (m *mockRepo) Insert(note domain.Note) (uuid.UUID, error) {
	if m.insertFunc != nil {
		return m.insertFunc(note)
	}
	return uuid.New(), nil
}

func (m *mockRepo) UpdateById(note domain.Note) error {
	if m.updateFunc != nil {
		return m.updateFunc(note)
	}
	return nil
}

func (m *mockRepo) DeleteById(id uuid.UUID) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(id)
	}
	return nil
}

func (m *mockRepo) FindByName(name string) (domain.Note, error) {
	if m.findByNameFunc != nil {
		return m.findByNameFunc(name)
	}
	return domain.Note{}, nil
}

func (m *mockRepo) FindByAncestor(ancestor string) ([]domain.Note, error) {
	if m.findByAncestorFunc != nil {
		return m.findByAncestorFunc(ancestor)
	}
	return []domain.Note{}, nil
}

func TestHomeHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	HomeHandler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	body := w.Body.String()
	if body != "Home page" {
		t.Errorf("expected 'Home page', got %s", body)
	}
}

func TestNotesUUIDHandler_GET(t *testing.T) {
	testID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")

	t.Run("successful get by id", func(t *testing.T) {
		oldRepo := domain.Repo
		defer func() { domain.Repo = oldRepo }()

		domain.Repo = &mockRepo{
			getNoteFunc: func(id uuid.UUID) (domain.Note, error) {
				return domain.Note{
					Id:         testID,
					Title:      "Test Note",
					Path:       "/test.md",
					Content:    "Test content",
					CreateTime: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
					UpdateTime: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
				}, nil
			},
		}

		req := httptest.NewRequest(http.MethodGet, "/notes/"+testID.String(), nil)
		req.SetPathValue("id", testID.String())
		w := httptest.NewRecorder()

		NotesUUIDHandler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})

	t.Run("invalid uuid format", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/notes/invalid-uuid", nil)
		req.SetPathValue("id", "invalid-uuid")
		w := httptest.NewRecorder()

		NotesUUIDHandler(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})

	t.Run("note not found", func(t *testing.T) {
		oldRepo := domain.Repo
		defer func() { domain.Repo = oldRepo }()

		domain.Repo = &mockRepo{
			getNoteFunc: func(id uuid.UUID) (domain.Note, error) {
				return domain.Note{}, errors.New("not found: sql: no rows in result set")
			},
		}

		req := httptest.NewRequest(http.MethodGet, "/notes/"+testID.String(), nil)
		req.SetPathValue("id", testID.String())
		w := httptest.NewRecorder()

		NotesUUIDHandler(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status 500, got %d", w.Code)
		}
	})
}

func TestNotesUUIDHandler_DELETE(t *testing.T) {
	testID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")

	t.Run("successful delete", func(t *testing.T) {
		oldRepo := domain.Repo
		defer func() { domain.Repo = oldRepo }()

		domain.Repo = &mockRepo{
			deleteFunc: func(id uuid.UUID) error {
				return nil
			},
		}

		req := httptest.NewRequest(http.MethodDelete, "/notes/"+testID.String(), nil)
		req.SetPathValue("id", testID.String())
		w := httptest.NewRecorder()

		NotesUUIDHandler(w, req)

		if w.Code != http.StatusNoContent {
			t.Errorf("expected status 204, got %d", w.Code)
		}
	})

	t.Run("delete with error", func(t *testing.T) {
		oldRepo := domain.Repo
		defer func() { domain.Repo = oldRepo }()

		domain.Repo = &mockRepo{
			deleteFunc: func(id uuid.UUID) error {
				return errors.New("database error")
			},
		}

		req := httptest.NewRequest(http.MethodDelete, "/notes/"+testID.String(), nil)
		req.SetPathValue("id", testID.String())
		w := httptest.NewRecorder()

		NotesUUIDHandler(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status 500, got %d", w.Code)
		}
	})
}

func TestNotesUUIDHandler_MethodNotAllowed(t *testing.T) {
	testID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")

	req := httptest.NewRequest(http.MethodPatch, "/notes/"+testID.String(), nil)
	req.SetPathValue("id", testID.String())
	w := httptest.NewRecorder()

	NotesUUIDHandler(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %d", w.Code)
	}
}

func TestNotesHandler_GET_All(t *testing.T) {
	t.Run("get all notes", func(t *testing.T) {
		oldRepo := domain.Repo
		defer func() { domain.Repo = oldRepo }()

		domain.Repo = &mockRepo{
			getAllFunc: func() ([]domain.Note, error) {
				return []domain.Note{
					{
						Id:      uuid.New(),
						Title:   "Note 1",
						Content: "Content 1",
					},
					{
						Id:      uuid.New(),
						Title:   "Note 2",
						Content: "Content 2",
					},
				}, nil
			},
		}

		req := httptest.NewRequest(http.MethodGet, "/notes", nil)
		w := httptest.NewRecorder()

		NotesHandler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})

	t.Run("get all with error", func(t *testing.T) {
		oldRepo := domain.Repo
		defer func() { domain.Repo = oldRepo }()

		domain.Repo = &mockRepo{
			getAllFunc: func() ([]domain.Note, error) {
				return nil, errors.New("database error")
			},
		}

		req := httptest.NewRequest(http.MethodGet, "/notes", nil)
		w := httptest.NewRecorder()

		NotesHandler(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status 500, got %d", w.Code)
		}
	})
}

func TestNotesHandler_GET_ByName(t *testing.T) {
	t.Run("find by name", func(t *testing.T) {
		oldRepo := domain.Repo
		defer func() { domain.Repo = oldRepo }()

		domain.Repo = &mockRepo{
			findByNameFunc: func(name string) (domain.Note, error) {
				return domain.Note{
					Id:      uuid.New(),
					Title:   "Found Note",
					Content: "Content",
				}, nil
			},
		}

		req := httptest.NewRequest(http.MethodGet, "/notes?name=test", nil)
		w := httptest.NewRecorder()

		NotesHandler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})

	t.Run("find by name with error", func(t *testing.T) {
		oldRepo := domain.Repo
		defer func() { domain.Repo = oldRepo }()

		domain.Repo = &mockRepo{
			findByNameFunc: func(name string) (domain.Note, error) {
				return domain.Note{}, errors.New("not found")
			},
		}

		req := httptest.NewRequest(http.MethodGet, "/notes?name=nonexistent", nil)
		w := httptest.NewRecorder()

		NotesHandler(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status 500, got %d", w.Code)
		}
	})
}

func TestNotesHandler_GET_ByAncestor(t *testing.T) {
	t.Run("find by ancestor", func(t *testing.T) {
		oldRepo := domain.Repo
		defer func() { domain.Repo = oldRepo }()

		domain.Repo = &mockRepo{
			findByAncestorFunc: func(ancestor string) ([]domain.Note, error) {
				return []domain.Note{
					{Id: uuid.New(), Title: "Child 1"},
					{Id: uuid.New(), Title: "Child 2"},
				}, nil
			},
		}

		req := httptest.NewRequest(http.MethodGet, "/notes?ancestor=parent", nil)
		w := httptest.NewRecorder()

		NotesHandler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})
}

func TestNotesHandler_GET_InvalidParams(t *testing.T) {
	t.Run("unknown parameter", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/notes?invalid=value", nil)
		w := httptest.NewRecorder()

		NotesHandler(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}

		body := w.Body.String()
		if !strings.Contains(body, "Unknown parameter") {
			t.Errorf("expected 'Unknown parameter' error, got %s", body)
		}
	})

	t.Run("two parameters", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/notes?name=test&ancestor=parent", nil)
		w := httptest.NewRecorder()

		NotesHandler(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})

	t.Run("too many parameters", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/notes?a=1&b=2&c=3", nil)
		w := httptest.NewRecorder()

		NotesHandler(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})
}

func TestNotesHandler_POST(t *testing.T) {
	t.Run("successful create", func(t *testing.T) {
		oldRepo := domain.Repo
		defer func() { domain.Repo = oldRepo }()

		domain.Repo = &mockRepo{
			insertFunc: func(note domain.Note) (uuid.UUID, error) {
				return uuid.New(), nil
			},
		}

		note := map[string]interface{}{
			"title":   "New Note",
			"content": "New content",
		}
		body, _ := json.Marshal(note)

		req := httptest.NewRequest(http.MethodPost, "/notes", bytes.NewReader(body))
		w := httptest.NewRecorder()

		NotesHandler(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("expected status 201, got %d", w.Code)
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/notes", strings.NewReader("invalid json"))
		w := httptest.NewRecorder()

		NotesHandler(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})

	t.Run("insert error", func(t *testing.T) {
		oldRepo := domain.Repo
		defer func() { domain.Repo = oldRepo }()

		domain.Repo = &mockRepo{
			insertFunc: func(note domain.Note) (uuid.UUID, error) {
				return uuid.Nil, errors.New("database error")
			},
		}

		note := map[string]interface{}{"title": "Test"}
		body, _ := json.Marshal(note)

		req := httptest.NewRequest(http.MethodPost, "/notes", bytes.NewReader(body))
		w := httptest.NewRecorder()

		NotesHandler(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status 500, got %d", w.Code)
		}
	})
}

func TestNotesHandler_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodPatch, "/notes", nil)
	w := httptest.NewRecorder()

	NotesHandler(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %d", w.Code)
	}
}
