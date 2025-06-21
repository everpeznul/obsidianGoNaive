package domain

import "github.com/google/uuid"

type NoteRepository interface {
	Insert(note Note) (uuid.UUID, error)
	GetByID(id uuid.UUID) (Note, error)
	GetAll() ([]Note, error)
	UpdateById(note Note) error
	DeleteByID(id uuid.UUID) error
	FindByTitle(title string) (Note, error)
	FindByAncestor(title string) ([]Note, error)
}
