package domain

import "github.com/google/uuid"

type NoteRepository interface {
	Insert(note *Note) error
	GetByID(id uuid.UUID) (*Note, error)
	GetAll() ([]*Note, error)
	Update(note *Note) error
	Delete(id string) error
	FindByTitle(title string) ([]*Note, error)
	FindByAncestor(title string) ([]*Note, error)
}
