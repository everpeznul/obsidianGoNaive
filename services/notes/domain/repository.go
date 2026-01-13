package domain

import (
	"context"

	"github.com/google/uuid"
)

type NoteRepository interface {
	Insert(ctx context.Context, note Note) (uuid.UUID, error)
	GetByID(ctx context.Context, id uuid.UUID) (Note, error)
	GetAll(ctx context.Context) ([]Note, error)
	UpdateById(ctx context.Context, note Note) error
	DeleteById(ctx context.Context, id uuid.UUID) error
	FindByName(ctx context.Context, name string) (Note, error)
	FindByAncestor(ctx context.Context, ancestor string) ([]Note, error)
}

var Repo NoteRepository

func SetRepo(r NoteRepository) {

	Repo = r
}
