package use_case

import (
	"context"
	"fmt"
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
	obsiLog.Debug("Update ReturnTypesNote", fmt.Sprintf("%T", note))

	links, err := u.Linker.Format(ctx, note)
	if err != nil {

		obsiLog.Error("Update links Note ERROR", "note", oldNote, "error", err)
		return fmt.Errorf("update links note ERROR: %w", err)
	}

	tags, err := u.Tager.Format(ctx, note)
	if err != nil {

		obsiLog.Error("Update tags Note ERROR", "note", oldNote, "error", err)
		return fmt.Errorf("update tags note ERROR: %w", err)
	}

	newNote := &domain.Note{oldNote.Id, oldNote.Title, oldNote.Path, oldNote.Class, tags, links, oldNote.Content, oldNote.CreateTime, oldNote.UpdateTime}

	err = u.Repo.UpdateById(ctx, *newNote)
	if err != nil {

		obsiLog.Error("Repo Update Note ERROR", "note", newNote, "error", err)
		return fmt.Errorf("repo update note ERROR: %w", err)
	}

	obsiLog.Debug("Successful Note Update", "note", newNote)
	return nil
}

func (u *Updater) Update_all(note domain.Note) {
}
