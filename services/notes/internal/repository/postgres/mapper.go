package postgres

import (
	"obsidianGoNaive/services/notes/internal/domain"
	"obsidianGoNaive/services/notes/internal/repository"
	"strings"

	"github.com/lib/pq"
)

type PostgresMapper struct{}

func (pm *PostgresMapper) RepoToDomain(note *repository.Note) *domain.Note {

	return &domain.Note{
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

func (pm *PostgresMapper) RepoToDomainSlice(notes []*repository.Note) []*domain.Note {
	domainNotes := make([]*domain.Note, len(notes))

	for i := range notes {
		domainNotes[i] = pm.RepoToDomain(notes[i])
	}

	return domainNotes
}

func (pm *PostgresMapper) DomainToRepo(note *domain.Note) *repository.Note {

	return &repository.Note{
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

func (pm *PostgresMapper) DomainToRepoSlice(notes []*domain.Note) []*repository.Note {
	dbNotes := make([]*repository.Note, len(notes))

	for i := range notes {
		dbNotes[i] = pm.DomainToRepo(notes[i])
	}

	return dbNotes
}
