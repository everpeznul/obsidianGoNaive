package transport

import (
	"context"
	"fmt"
	"log/slog"
	"obsidianGoNaive/pkg/protos/gen/notes"
	"obsidianGoNaive/pkg/protos/gen/updater"
	"obsidianGoNaive/services/updater/internal/domain"
)

type UpdaterService struct {
	updater.UnimplementedUpdaterServer
	Log    *slog.Logger
	Client notes.NotesClient
}

func NewUpdaterService(client notes.NotesClient, log *slog.Logger) *UpdaterService {

	return &UpdaterService{Client: client, Log: log}
}

func (us *UpdaterService) Update(ctx context.Context, r *updater.UpdateRequest) (*updater.UpdateResponse, error) {

	oldNote, _ := ProtoToNote(r.Note)

	note := domain.ReturnTypesNote(oldNote)
	us.Log.Debug("Update ReturnTypesNote", fmt.Sprintf("%T", note))

	links, err := domain.LinksFormat(ctx, note, us, us.Log)
	if err != nil {

		us.Log.Error("Update links Note ERROR", "note", oldNote, "error", err)
		return &updater.UpdateResponse{}, fmt.Errorf("update links note ERROR: %w", err)
	}

	tags, err := domain.TagsFormat(ctx, note, us)
	if err != nil {

		us.Log.Error("Update tags Note ERROR", "note", oldNote, "error", err)
		return &updater.UpdateResponse{}, fmt.Errorf("update tags note ERROR: %w", err)
	}

	tempNote := &domain.Note{oldNote.Id, oldNote.Title, oldNote.Path, oldNote.Class, tags, links, oldNote.Content, oldNote.CreateTime, oldNote.UpdateTime}
	newNote := NoteToProto(tempNote)
	_, err = us.Client.UpdateById(ctx, &notes.UpdateByIdRequest{Note: newNote})
	if err != nil {

		us.Log.Error("Repo Update Note ERROR", "note", newNote, "error", err)
		return &updater.UpdateResponse{}, fmt.Errorf("repo update note ERROR: %w", err)
	}

	us.Log.Debug("Successful Note Update", "note", newNote)
	return &updater.UpdateResponse{}, nil
}
