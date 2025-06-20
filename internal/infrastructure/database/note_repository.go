package database

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"log"
	"obsidianGoNaive/internal/domain"
	"time"
)

var nt = noteTranslater{}

type PgDB struct {
	sql.DB
}

type note struct {
	Id         uuid.UUID
	Title      string
	Path       string
	Class      string
	Tags       []string
	Links      []string
	Content    []string
	CreateTime time.Time
	UpdateTime time.Time
}

func someFunc() {
	connStr := "user=postgres password=mypass dbname=productdb sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
}

func (db *PgDB) Insert(note domain.Note) (uuid.UUID, error) {
	query := `INSERT INTO notes (title, path, class, tags, links, content, create_time, update_time) 
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`

	tempNote := nt.DomainToDatabase(note)
	var pk uuid.UUID
	err := db.QueryRow(query,
		tempNote.Title,
		tempNote.Path,
		tempNote.Class,
		pq.Array(tempNote.Tags),
		pq.Array(tempNote.Links),
		pq.Array(tempNote.Content),
		tempNote.CreateTime,
		tempNote.UpdateTime,
	).Scan(&pk)

	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to insert note: %w", err)
	}

	return pk, nil
}

func (db *PgDB) GetByID(id uuid.UUID) (domain.Note, error) {
	query := `SELECT * FROM notes WHERE id = $1`
	row := db.QueryRow(query, id)

	note := note{}
	err := row.Scan(
		&note.Id,
		&note.Title,
		&note.Path,
		&note.Class,
		pq.Array(&note.Tags),
		pq.Array(&note.Links),
		pq.Array(&note.Content),
		&note.CreateTime,
		&note.UpdateTime,
	)

	if err != nil {
		return domain.Note{}, fmt.Errorf("failed: %w", err)
	}

	return nt.DatabaseToDomain(note), nil
}

func (db *PgDB) GetAll() ([]domain.Note, error) {
	var notes []domain.Note
	query := `SELECT * FROM notes`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var n note
		err := rows.Scan(
			&n.Id,
			&n.Title,
			&n.Path,
			&n.Class,
			pq.Array(&n.Tags),
			pq.Array(&n.Links),
			pq.Array(&n.Content),
			&n.CreateTime,
			&n.UpdateTime,
		)
		if err != nil {
			return nil, fmt.Errorf("failed: %w", err)
		}
		notes = append(notes, nt.DatabaseToDomain(n))
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed: %w", err)
	}

	return notes, nil
}

func (db *PgDB) UpdateById(n domain.Note) error {
	nTemp := nt.DomainToDatabase(n)
	query := `UPDATE notes 
              SET title = $1,
                  path = $2,
                  class = $3,
                  tags = $4,
                  links = $5,
                  content = $6,
                  create_time = $7,
                  update_time = $8
              WHERE id = $9`

	_, err := db.Exec(query,
		nTemp.Title,
		nTemp.Path,
		nTemp.Class,
		pq.Array(nTemp.Tags),
		pq.Array(nTemp.Links),
		pq.Array(nTemp.Content),
		nTemp.CreateTime,
		nTemp.UpdateTime,
		nTemp.Id,
	)

	return err
}

func (db *PgDB) DeleteByID(id uuid.UUID) error {
	query := `DELETE FROM notes WHERE id = $1`
	_, err := db.Exec(query, id)
	return err
}

func (db *PgDB) GetByTitle(title string) (domain.Note, error) {
	var note note
	query := `SELECT id, title, path, class, tags, links, content, create_time, update_time 
              FROM notes 
              WHERE title = $1`

	err := db.QueryRow(query, title).Scan(
		&note.Id,
		&note.Title,
		&note.Path,
		&note.Class,
		pq.Array(&note.Tags),
		pq.Array(&note.Links),
		pq.Array(&note.Content),
		&note.CreateTime,
		&note.UpdateTime,
	)

	if err != nil {
		return domain.Note{}, fmt.Errorf("failed: %w", err)
	}

	return nt.DatabaseToDomain(note), nil
}

func (db *PgDB) GetByAncestor(ancestor string) ([]domain.Note, error) {
	var notes []domain.Note
	query := `SELECT id, title, path, class, tags, links, content, create_time, update_time  
              FROM notes 
              WHERE title ILIKE $1`

	rows, err := db.Query(query, ancestor+"%")
	if err != nil {
		return nil, fmt.Errorf("failed: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var note note
		err := rows.Scan(
			&note.Id,
			&note.Title,
			&note.Path,
			&note.Class,
			pq.Array(&note.Tags),
			pq.Array(&note.Links),
			pq.Array(&note.Content),
			&note.CreateTime,
			&note.UpdateTime,
		)
		if err != nil {
			return nil, fmt.Errorf("failed: %w", err)
		}
		notes = append(notes, nt.DatabaseToDomain(note))
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed: %w", err)
	}

	return notes, nil
}
