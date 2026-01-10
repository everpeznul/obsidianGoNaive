package use_case

import (
	"context"
	"fmt"
	pbn "obsidianGoNaive/protos/gen/go/notes"
	domain2 "obsidianGoNaive/protos/gen/go/notes/domain"
	pbu "obsidianGoNaive/protos/gen/go/updater"

	"github.com/google/uuid"
)

type UpdaterService struct {
	pbu.UnimplementedUpdaterServer
	pbn
	Linker
	Tager
}

func (us *UpdaterService) Update(ctx context.Context, r *pb.UpdateRequest) (*pb.UpdateResponse, error) {

	note, _ := ProtoToNote(r.Note)

	note := domain2.ReturnTypesNote(oldNote)
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

	newNote := &domain2.Note{oldNote.Id, oldNote.Title, oldNote.Path, oldNote.Class, tags, links, oldNote.Content, oldNote.CreateTime, oldNote.UpdateTime}

	// client.notes.Update() что-нибудь такое сделать
	err = u.Repo.UpdateById(ctx, *newNote)
	if err != nil {

		obsiLog.Error("Repo Update Note ERROR", "note", newNote, "error", err)
		return fmt.Errorf("repo update note ERROR: %w", err)
	}

	obsiLog.Debug("Successful Note Update", "note", newNote)
	return nil
}

func ProtoToNote(protoNote *pb.Note) (Note, error) {
	id, err := uuid.Parse(protoNote.Id)
	if err != nil {
		return Note{}, fmt.Errorf("invalid UUID: %w", err)
	}

	return Note{
		Id:         id,
		Title:      protoNote.Title,
		Path:       protoNote.Path,
		Class:      protoNote.Class,
		Tags:       protoNote.Tags,
		Links:      protoNote.Links,
		Content:    protoNote.Content,
		CreateTime: protoNote.CreateTime.AsTime(),
		UpdateTime: protoNote.UpdateTime.AsTime(),
	}, nil
}
