package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"obsidianGoNaive/services/notes/internal/domain"
	"obsidianGoNaive/services/notes/internal/repository"

	"github.com/google/uuid"
)

type Postgres struct {
	DB   *sql.DB
	Mapp *PostgresMapper
}

func NewPostgres(db *sql.DB, mapp *PostgresMapper) *Postgres {
	return &Postgres{DB: db, Mapp: mapp}
}

func (p *Postgres) Insert(ctx context.Context, note *domain.Note) (uuid.UUID, error) {

	newId := uuid.New()

	newNote := p.Mapp.DomainToRepo(note)
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

func (p *Postgres) GetByID(ctx context.Context, id uuid.UUID) (*domain.Note, error) {
	query := `SELECT id, title, path, class, tags, links, content, create_time, update_time
              FROM notes WHERE id = $1`
	row := p.DB.QueryRowContext(ctx, query, id)

	note := &repository.Note{}
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
		return nil, fmt.Errorf("note with id=%s not found in database: %w", id, err)
	}
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	return p.Mapp.RepoToDomain(note), nil
}

func (p *Postgres) GetAll(ctx context.Context) ([]*domain.Note, error) {
	notes := make([]*domain.Note, 0)

	query := `SELECT id, title, path, class, tags, links, content, create_time, update_time 
			  FROM notes`
	rows, err := p.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed get all notes from database: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		n := &repository.Note{}
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

		notes = append(notes, p.Mapp.RepoToDomain(n))
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed: %w", err)
	}

	return notes, nil
}

func (p *Postgres) UpdateById(ctx context.Context, n *domain.Note) error {
	newNote := p.Mapp.DomainToRepo(n)
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

func (p *Postgres) DeleteById(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM notes WHERE id = $1`
	_, err := p.DB.ExecContext(ctx, query, id)

	if err != nil {
		return fmt.Errorf("failed delete note in database with id=%s: %w", id, err)
	}

	return nil
}

func (p *Postgres) FindByName(ctx context.Context, name string) (*domain.Note, error) {
	newNote := &repository.Note{}

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
		return nil, fmt.Errorf("failed to get note '%s' from database: %w", name, err)
	}

	return p.Mapp.RepoToDomain(newNote), nil
}

// func (p *PgDB) FindAncestors(ctx context.Context, name string) (domain.Note, error) {}

func (p *Postgres) FindByAncestor(ctx context.Context, ancestor string) ([]*domain.Note, error) {

	notes := make([]*domain.Note, 0)

	query := `SELECT id, title, path, class, tags, links, content, create_time, update_time  
              FROM notes 
              WHERE right(path, strpos(reverse(path), '/') - 1) ILIKE $1`

	rows, err := p.DB.QueryContext(ctx, query, ancestor+"%")
	if err != nil {
		return nil, fmt.Errorf("failed to find descendant of %s: %w", ancestor, err)
	}
	defer rows.Close()

	for rows.Next() {
		note := &repository.Note{}
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
		notes = append(notes, p.Mapp.RepoToDomain(note))
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed: %w", err)
	}

	return notes, nil
}
