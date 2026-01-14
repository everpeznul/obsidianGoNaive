package repository

import (
	"context"
	"obsidianGoNaive/services/notes/internal/domain"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Repository struct {
	DB Database
}

func NewRepository(db Database) *Repository {
	return &Repository{DB: db}
}

type Database interface {
	Insert(context.Context, *domain.Note) (uuid.UUID, error)
	GetByID(context.Context, uuid.UUID) (*domain.Note, error)
	GetAll(context.Context) ([]*domain.Note, error)
	UpdateById(context.Context, *domain.Note) error
	DeleteById(context.Context, uuid.UUID) error
	FindByName(context.Context, string) (*domain.Note, error)
	FindByAncestor(context.Context, string) ([]*domain.Note, error)
}

type Note struct {
	Id         uuid.UUID
	Title      string
	Path       string
	Class      string
	Tags       pq.StringArray
	Links      pq.StringArray
	Content    pq.StringArray
	CreateTime time.Time
	UpdateTime time.Time
}

type Mapper interface {
	DomainToRepo(*domain.Note) *Note
	DomainToRepoSlice([]*domain.Note) []*Note
	RepoToDomain(*Note) *domain.Note
	RepoToDomainSlice([]*Note) []*domain.Note
}
