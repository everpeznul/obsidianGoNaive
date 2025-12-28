package http

import (
	"obsidianGoNaive/internal/domain"
	"time"

	"github.com/google/uuid"
)

type httpNote struct {
	Id         uuid.UUID `json:"id,omitempty"`
	Title      string    `json:"title"`
	Path       string    `json:"path"`
	Class      string    `json:"class"`
	Tags       []string  `json:"tags"`
	Links      []string  `json:"links"`
	Content    string    `json:"content"`
	CreateTime time.Time `json:"create_time,omitempty"`
	UpdateTime time.Time `json:"update_time,omitempty"`
}

var nm = noteMapperHttp{}

type noteMapperHttp struct{}

func (nm *noteMapperHttp) HTTPToDomain(note httpNote) domain.Note {

	return domain.Note{
		Id:         note.Id,
		Title:      note.Title,
		Path:       note.Path,
		Class:      note.Class,
		Tags:       note.Tags,
		Links:      note.Links,
		Content:    note.Content,
		CreateTime: note.CreateTime,
		UpdateTime: note.UpdateTime,
	}
}

func (nm *noteMapperHttp) HTTPToDomainSlice(notes []httpNote) []domain.Note {

	domainNotes := make([]domain.Note, 0, len(notes))
	for _, note := range notes {

		domainNotes = append(domainNotes, domain.Note{
			Id:         note.Id,
			Title:      note.Title,
			Path:       note.Path,
			Class:      note.Class,
			Tags:       note.Tags,
			Links:      note.Links,
			Content:    note.Content,
			CreateTime: note.CreateTime,
			UpdateTime: note.UpdateTime,
		})
	}
	return domainNotes
}

func (nm *noteMapperHttp) DomainToHTTP(note domain.Note) httpNote {

	return httpNote{
		Id:         note.Id,
		Title:      note.Title,
		Path:       note.Path,
		Class:      note.Class,
		Tags:       note.Tags,
		Links:      note.Links,
		Content:    note.Content,
		CreateTime: note.CreateTime,
		UpdateTime: note.UpdateTime,
	}
}

func (nm *noteMapperHttp) DomainToHTTPSlice(notes []domain.Note) []httpNote {

	httpNotes := make([]httpNote, 0, len(notes))
	for _, note := range notes {

		httpNotes = append(httpNotes, httpNote{
			Id:         note.Id,
			Title:      note.Title,
			Path:       note.Path,
			Class:      note.Class,
			Tags:       note.Tags,
			Links:      note.Links,
			Content:    note.Content,
			CreateTime: note.CreateTime,
			UpdateTime: note.UpdateTime,
		})
	}
	return httpNotes
}
