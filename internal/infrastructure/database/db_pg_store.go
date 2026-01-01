package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"obsidianGoNaive/internal/domain"

	"github.com/google/uuid"
)

type PgDB struct {
	DB *sql.DB
}

func (p *PgDB) Insert(ctx context.Context, note domain.Note) (uuid.UUID, error) {

	newId := uuid.New()

	newNote := dbNote{}
	newNote = nm.DomainToDatabase(note)
	newNote.Id = newId

	query := `INSERT INTO notes (id, title, path, class, tags, links, content, create_time, update_time) 
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := p.DB.ExecContext(
		ctx,
		query,
		newNote.Id,
		newNote.Title,
		newNote.Path,
		newNote.Class,
		newNote.Tags,
		newNote.Links,
		newNote.Content,
		newNote.CreateTime,
		newNote.UpdateTime,
	)

	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to insert note into database: %w", err)
	}

	return newId, nil
}

func (p *PgDB) GetByID(ctx context.Context, id uuid.UUID) (domain.Note, error) {
	query := `SELECT id, title, path, class, tags, links, content, create_time, update_time
              FROM notes WHERE id = $1`
	row := p.DB.QueryRowContext(ctx, query, id)

	note := dbNote{}
	err := row.Scan(
		&note.Id,
		&note.Title,
		&note.Path,
		&note.Class,
		&note.Tags,
		&note.Links,
		&note.Content,
		&note.CreateTime,
		&note.UpdateTime,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return domain.Note{}, fmt.Errorf("note with id=%s not found in database: %w", id, err)
	}
	if err != nil {
		return domain.Note{}, fmt.Errorf("database error: %w", err)
	}

	return nm.DatabaseToDomain(note), nil
}

func (p *PgDB) GetAll(ctx context.Context) ([]domain.Note, error) {
	notes := make([]domain.Note, 0)

	query := `SELECT id, title, path, class, tags, links, content, create_time, update_time 
			  FROM notes`
	rows, err := p.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed get all notes from database: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		n := dbNote{}
		err := rows.Scan(
			&n.Id,
			&n.Title,
			&n.Path,
			&n.Class,
			&n.Tags,
			&n.Links,
			&n.Content,
			&n.CreateTime,
			&n.UpdateTime,
		)

		if err != nil {
			return nil, fmt.Errorf("failed push another note to notes list from database: %w", err)
		}

		notes = append(notes, nm.DatabaseToDomain(n))
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed: %w", err)
	}

	return notes, nil
}

func (p *PgDB) UpdateById(ctx context.Context, n domain.Note) error {
	newNote := nm.DomainToDatabase(n)
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

	_, err := p.DB.ExecContext(
		ctx,
		query,
		newNote.Title,
		newNote.Path,
		newNote.Class,
		newNote.Tags,
		newNote.Links,
		newNote.Content,
		newNote.CreateTime,
		newNote.UpdateTime,
		newNote.Id,
	)

	if err != nil {
		return fmt.Errorf("failed to update note in database with id %s: %w", newNote.Id, err)
	}

	return nil
}

func (p *PgDB) DeleteById(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM notes WHERE id = $1`
	_, err := p.DB.ExecContext(ctx, query, id)

	if err != nil {
		return fmt.Errorf("failed delete note in database with id=%s: %w", id, err)
	}

	return nil
}

func (p *PgDB) FindByName(ctx context.Context, name string) (domain.Note, error) {
	newNote := dbNote{}

	query := `SELECT id, title, path, class, tags, links, content, create_time, update_time 
              FROM notes 
              WHERE right(path, strpos(reverse(path), '/') - 1) = $1`

	err := p.DB.QueryRowContext(ctx, query, name+".md").Scan(
		&newNote.Id,
		&newNote.Title,
		&newNote.Path,
		&newNote.Class,
		&newNote.Tags,
		&newNote.Links,
		&newNote.Content,
		&newNote.CreateTime,
		&newNote.UpdateTime,
	)

	if err != nil {
		return domain.Note{}, fmt.Errorf("failed to get note '%s' from database: %w", name, err)
	}

	fmt.Print(newNote)
	return nm.DatabaseToDomain(newNote), nil
}

// func (p *PgDB) FindAncestors(ctx context.Context, name string) (domain.Note, error) {}

func (p *PgDB) FindByAncestor(ctx context.Context, ancestor string) ([]domain.Note, error) {

	notes := make([]domain.Note, 0)

	query := `SELECT id, title, path, class, tags, links, content, create_time, update_time  
              FROM notes 
              WHERE right(path, strpos(reverse(path), '/') - 1) ILIKE $1`

	rows, err := p.DB.QueryContext(ctx, query, ancestor+"%")
	if err != nil {
		return nil, fmt.Errorf("failed to find descendant of %s: %w", ancestor, err)
	}
	defer rows.Close()

	for rows.Next() {
		note := dbNote{}
		err := rows.Scan(
			&note.Id,
			&note.Title,
			&note.Path,
			&note.Class,
			&note.Tags,
			&note.Links,
			&note.Content,
			&note.CreateTime,
			&note.UpdateTime,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to push note to notes from database in finding descendents of %s: %w", ancestor, err)
		}
		notes = append(notes, nm.DatabaseToDomain(note))
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed: %w", err)
	}

	return notes, nil
}
