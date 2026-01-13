package use_case

import (
	"context"
	"fmt"
	cmn "obsidianGoNaive/protos/gen/common"
	pbn "obsidianGoNaive/protos/gen/notes"
	pbu "obsidianGoNaive/protos/gen/updater"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UpdaterService struct {
	pbu.UnimplementedUpdaterServer
	client pbn.NotesClient
}

func NewUpdaterService(client pbn.NotesClient) *UpdaterService {

	return &UpdaterService{client: client}
}

func (us *UpdaterService) Update(ctx context.Context, r *pbu.UpdateRequest) (*pbu.UpdateResponse, error) {

	oldNote, _ := ProtoToNote(r.Note)

	note := ReturnTypesNote(oldNote)
	obsiLog.Debug("Update ReturnTypesNote", fmt.Sprintf("%T", note))

	links, err := us.LinksFormat(ctx, note)
	if err != nil {

		obsiLog.Error("Update links Note ERROR", "note", oldNote, "error", err)
		return &pbu.UpdateResponse{}, fmt.Errorf("update links note ERROR: %w", err)
	}

	tags, err := us.TagsFormat(ctx, note)
	if err != nil {

		obsiLog.Error("Update tags Note ERROR", "note", oldNote, "error", err)
		return &pbu.UpdateResponse{}, fmt.Errorf("update tags note ERROR: %w", err)
	}

	tempNote := Note{oldNote.Id, oldNote.Title, oldNote.Path, oldNote.Class, tags, links, oldNote.Content, oldNote.CreateTime, oldNote.UpdateTime}
	newNote := NoteToProto(tempNote)
	_, err = us.client.UpdateById(ctx, &pbn.UpdateByIdRequest{Note: &newNote})
	if err != nil {

		obsiLog.Error("Repo Update Note ERROR", "note", newNote, "error", err)
		return &pbu.UpdateResponse{}, fmt.Errorf("repo update note ERROR: %w", err)
	}

	obsiLog.Debug("Successful Note Update", "note", newNote)
	return &pbu.UpdateResponse{}, nil
}

func ProtoToNote(protoNote *cmn.Note) (Note, error) {
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

func NoteToProto(note Note) cmn.Note {
	return cmn.Note{
		Id:         note.Id.String(),
		Title:      note.Title,
		Path:       note.Path,
		Class:      note.Class,
		Tags:       note.Tags,
		Links:      note.Links,
		Content:    note.Content,
		CreateTime: timestamppb.New(note.CreateTime),
		UpdateTime: timestamppb.New(note.UpdateTime),
	}
}
