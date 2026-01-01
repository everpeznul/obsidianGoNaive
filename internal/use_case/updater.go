package use_case

import (
	"context"
	"obsidianGoNaive/internal/domain"
	"obsidianGoNaive/internal/infrastructure/database"
)

var Updtr Updater

func InitUpdater(repo *database.PgDB) {

	Updtr = Updater{repo, Linker{}, Tager{}}
}

type Updater struct {
	Repo domain.NoteRepository
	Linker
	Tager
}

func (u *Updater) Update(ctx context.Context, oldNote domain.Note) error {

	note := domain.ReturnTypesNote(oldNote)
	links := u.Linker.Format(note)
	tags := u.Tager.Format(note)

	newNote := &domain.Note{oldNote.Id, oldNote.Title, oldNote.Path, oldNote.Class, tags, links, oldNote.Content, oldNote.CreateTime, oldNote.UpdateTime}

	u.Repo.UpdateById(ctx, *newNote)

	return nil
}

func (u *Updater) Update_all(note domain.Note) {
}
