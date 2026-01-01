package database

import (
	"database/sql"
	"errors"
	"obsidianGoNaive/internal/domain"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

func TestPgDB_Insert(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	pgDB := &PgDB{DB: db}

	note := domain.Note{
		Title:      "Test Note",
		Path:       "/test/path.md",
		Class:      "article",
		Tags:       []string{"tag1", "tag2"},
		Links:      []string{"link1"},
		Content:    "Test content",
		CreateTime: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
		UpdateTime: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
	}

	t.Run("successful insert", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO notes").
			WithArgs(
				sqlmock.AnyArg(),
				note.Title,
				note.Path,
				note.Class,
				pq.StringArray(note.Tags),
				pq.StringArray(note.Links),
				pq.StringArray{note.Content},
				note.CreateTime,
				note.UpdateTime,
			).
			WillReturnResult(sqlmock.NewResult(1, 1))

		id, err := pgDB.Insert(note)
		if err != nil {
			t.Errorf("Insert() error = %v, want nil", err)
		}
		if id == uuid.Nil {
			t.Error("Insert() returned nil UUID")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %v", err)
		}
	})

	t.Run("insert with database error", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO notes").
			WillReturnError(errors.New("database connection error"))

		_, err := pgDB.Insert(note)
		if err == nil {
			t.Error("Insert() expected error, got nil")
		}
	})
}

func TestPgDB_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	pgDB := &PgDB{DB: db}

	testID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
	testTime := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)

	t.Run("successful get by id", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "title", "path", "class", "tags", "links", "content", "create_time", "update_time"}).
			AddRow(
				testID,
				"Test Note",
				"/test/path.md",
				"article",
				pq.StringArray{"tag1", "tag2"},
				pq.StringArray{"link1"},
				pq.StringArray{"Content line 1", "Content line 2"},
				testTime,
				testTime,
			)

		mock.ExpectQuery("SELECT (.+) FROM notes WHERE id").
			WithArgs(testID).
			WillReturnRows(rows)

		note, err := pgDB.GetByID(testID)
		if err != nil {
			t.Errorf("GetByID() error = %v, want nil", err)
		}
		if note.Id != testID {
			t.Errorf("GetByID() id = %v, want %v", note.Id, testID)
		}
		if note.Title != "Test Note" {
			t.Errorf("GetByID() title = %v, want 'Test Note'", note.Title)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %v", err)
		}
	})

	t.Run("note not found", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM notes WHERE id").
			WithArgs(testID).
			WillReturnError(sql.ErrNoRows)

		_, err := pgDB.GetByID(testID)
		if err == nil {
			t.Error("GetByID() expected error, got nil")
		}
	})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM notes WHERE id").
			WithArgs(testID).
			WillReturnError(errors.New("connection error"))

		_, err := pgDB.GetByID(testID)
		if err == nil {
			t.Error("GetByID() expected error, got nil")
		}
	})
}

func TestPgDB_GetAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	pgDB := &PgDB{DB: db}
	testTime := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)

	t.Run("successful get all notes", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "title", "path", "class", "tags", "links", "content", "create_time", "update_time"}).
			AddRow(
				uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				"Note 1",
				"/path1.md",
				"article",
				pq.StringArray{"tag1"},
				pq.StringArray{"link1"},
				pq.StringArray{"Content 1"},
				testTime,
				testTime,
			).
			AddRow(
				uuid.MustParse("223e4567-e89b-12d3-a456-426614174001"),
				"Note 2",
				"/path2.md",
				"note",
				pq.StringArray{"tag2"},
				pq.StringArray{},
				pq.StringArray{"Content 2"},
				testTime,
				testTime,
			)

		mock.ExpectQuery("SELECT (.+) FROM notes").
			WillReturnRows(rows)

		notes, err := pgDB.GetAll()
		if err != nil {
			t.Errorf("GetAll() error = %v, want nil", err)
		}
		if len(notes) != 2 {
			t.Errorf("GetAll() returned %d notes, want 2", len(notes))
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %v", err)
		}
	})

	t.Run("get all with empty result", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "title", "path", "class", "tags", "links", "content", "create_time", "update_time"})

		mock.ExpectQuery("SELECT (.+) FROM notes").
			WillReturnRows(rows)

		notes, err := pgDB.GetAll()
		if err != nil {
			t.Errorf("GetAll() error = %v, want nil", err)
		}
		if len(notes) != 0 {
			t.Errorf("GetAll() returned %d notes, want 0", len(notes))
		}
	})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM notes").
			WillReturnError(errors.New("database error"))

		_, err := pgDB.GetAll()
		if err == nil {
			t.Error("GetAll() expected error, got nil")
		}
	})
}

func TestPgDB_UpdateById(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	pgDB := &PgDB{DB: db}

	note := domain.Note{
		Id:         uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
		Title:      "Updated Note",
		Path:       "/updated/path.md",
		Class:      "article",
		Tags:       []string{"newtag"},
		Links:      []string{"newlink"},
		Content:    "Updated content",
		CreateTime: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
		UpdateTime: time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC),
	}

	t.Run("successful update", func(t *testing.T) {
		mock.ExpectExec("UPDATE notes").
			WithArgs(
				note.Title,
				note.Path,
				note.Class,
				pq.StringArray(note.Tags),
				pq.StringArray(note.Links),
				pq.StringArray{note.Content},
				note.CreateTime,
				note.UpdateTime,
				note.Id,
			).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := pgDB.UpdateById(note)
		if err != nil {
			t.Errorf("UpdateById() error = %v, want nil", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %v", err)
		}
	})

	t.Run("update with database error", func(t *testing.T) {
		mock.ExpectExec("UPDATE notes").
			WillReturnError(errors.New("update failed"))

		err := pgDB.UpdateById(note)
		if err == nil {
			t.Error("UpdateById() expected error, got nil")
		}
	})
}

func TestPgDB_DeleteById(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	pgDB := &PgDB{DB: db}
	testID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")

	t.Run("successful delete", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM notes WHERE id").
			WithArgs(testID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := pgDB.DeleteById(testID)
		if err != nil {
			t.Errorf("DeleteById() error = %v, want nil", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %v", err)
		}
	})

	t.Run("delete with database error", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM notes WHERE id").
			WithArgs(testID).
			WillReturnError(errors.New("delete failed"))

		err := pgDB.DeleteById(testID)
		if err == nil {
			t.Error("DeleteById() expected error, got nil")
		}
	})
}

func TestPgDB_FindByName(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	pgDB := &PgDB{DB: db}
	testTime := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	testID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")

	t.Run("successful find by name", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "title", "path", "class", "tags", "links", "content", "create_time", "update_time"}).
			AddRow(
				testID,
				"Test Note",
				"/test/path/testfile.md",
				"article",
				pq.StringArray{"tag1"},
				pq.StringArray{"link1"},
				pq.StringArray{"Content"},
				testTime,
				testTime,
			)

		mock.ExpectQuery("SELECT (.+) FROM notes WHERE right").
			WithArgs("testfile.md").
			WillReturnRows(rows)

		note, err := pgDB.FindByName("testfile")
		if err != nil {
			t.Errorf("FindByName() error = %v, want nil", err)
		}
		if note.Id != testID {
			t.Errorf("FindByName() id = %v, want %v", note.Id, testID)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %v", err)
		}
	})

	t.Run("note not found by name", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM notes WHERE right").
			WithArgs("nonexistent.md").
			WillReturnError(sql.ErrNoRows)

		_, err := pgDB.FindByName("nonexistent")
		if err == nil {
			t.Error("FindByName() expected error, got nil")
		}
	})
}

func TestPgDB_FindByAncestor(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	pgDB := &PgDB{DB: db}
	testTime := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)

	t.Run("successful find by ancestor", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "title", "path", "class", "tags", "links", "content", "create_time", "update_time"}).
			AddRow(
				uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				"Child Note 1",
				"/parent/child1.md",
				"article",
				pq.StringArray{"tag1"},
				pq.StringArray{},
				pq.StringArray{"Content 1"},
				testTime,
				testTime,
			).
			AddRow(
				uuid.MustParse("223e4567-e89b-12d3-a456-426614174001"),
				"Child Note 2",
				"/parent/child2.md",
				"note",
				pq.StringArray{},
				pq.StringArray{},
				pq.StringArray{"Content 2"},
				testTime,
				testTime,
			)

		mock.ExpectQuery("SELECT (.+) FROM notes WHERE right").
			WithArgs("parent%").
			WillReturnRows(rows)

		notes, err := pgDB.FindByAncestor("parent")
		if err != nil {
			t.Errorf("FindByAncestor() error = %v, want nil", err)
		}
		if len(notes) != 2 {
			t.Errorf("FindByAncestor() returned %d notes, want 2", len(notes))
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %v", err)
		}
	})

	t.Run("find by ancestor with no results", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "title", "path", "class", "tags", "links", "content", "create_time", "update_time"})

		mock.ExpectQuery("SELECT (.+) FROM notes WHERE right").
			WithArgs("nonexistent%").
			WillReturnRows(rows)

		notes, err := pgDB.FindByAncestor("nonexistent")
		if err != nil {
			t.Errorf("FindByAncestor() error = %v, want nil", err)
		}
		if len(notes) != 0 {
			t.Errorf("FindByAncestor() returned %d notes, want 0", len(notes))
		}
	})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM notes WHERE right").
			WithArgs("test%").
			WillReturnError(errors.New("query failed"))

		_, err := pgDB.FindByAncestor("test")
		if err == nil {
			t.Error("FindByAncestor() expected error, got nil")
		}
	})
}
