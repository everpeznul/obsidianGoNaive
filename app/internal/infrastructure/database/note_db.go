package database

import (
	"obsidianGoNaive/internal/domain"
	"strings"
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

var nm noteMapperDb

type noteMapperDb struct{}

func (nm *noteMapperDb) DatabaseToDomain(note dbNote) domain.Note {

	return domain.Note{
		Id:         note.Id,
		Title:      note.Title,
		Path:       note.Path,
		Class:      note.Class,
		Tags:       []string(note.Tags),
		Links:      []string(note.Links),
		Content:    strings.Join(note.Content, ""),
		CreateTime: note.CreateTime,
		UpdateTime: note.UpdateTime,
	}
}

func (nm *noteMapperDb) DatabaseToDomainSlice(notes []dbNote) []domain.Note {
	domainNotes := make([]domain.Note, 0, len(notes))

	for _, note := range notes {
		domainNotes = append(domainNotes, domain.Note{
			Id:         note.Id,
			Title:      note.Title,
			Path:       note.Path,
			Class:      note.Class,
			Tags:       []string(note.Tags),
			Links:      []string(note.Links),
			Content:    strings.Join(note.Content, ""),
			CreateTime: note.CreateTime,
			UpdateTime: note.UpdateTime,
		})
	}

	return domainNotes
}

func (nm *noteMapperDb) DomainToDatabase(note domain.Note) dbNote {

	return dbNote{
		note.Id,
		note.Title,
		note.Path,
		note.Class,
		pq.StringArray(note.Tags),
		pq.StringArray(note.Links),
		pq.StringArray{note.Content},
		note.CreateTime,
		note.UpdateTime,
	}
}

func (nm *noteMapperDb) DomainToDatabaseSlice(notes []domain.Note) []dbNote {
	dbNotes := make([]dbNote, 0, len(notes))

	for _, note := range notes {
		dbNotes = append(dbNotes, dbNote{
			Id:         note.Id,
			Title:      note.Title,
			Path:       note.Path,
			Class:      note.Class,
			Tags:       pq.StringArray(note.Tags),
			Links:      pq.StringArray(note.Links),
			Content:    pq.StringArray{note.Content},
			CreateTime: note.CreateTime,
			UpdateTime: note.UpdateTime,
		})
	}

	return dbNotes
}
