package database

import (
	"obsidianGoNaive/internal/domain"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type dbNote struct {
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

var nt noteTranslate

type noteTranslate struct{}

func (nt *noteTranslate) DatabaseToDomain(n dbNote) domain.Note {

	return domain.Note{
		Id:         n.Id,
		Title:      n.Title,
		Path:       n.Path,
		Class:      n.Class,
		Tags:       []string(n.Tags),
		Links:      []string(n.Links),
		Content:    n.Content[0],
		CreateTime: n.CreateTime,
		UpdateTime: n.UpdateTime,
	}
}

func (nt *noteTranslate) DomainToDatabase(n domain.Note) dbNote {

	return dbNote{
		n.Id,
		n.Title,
		n.Path,
		n.Class,
		pq.StringArray(n.Tags),
		pq.StringArray(n.Links),
		pq.StringArray{n.Content},
		n.CreateTime,
		n.UpdateTime,
	}
}
