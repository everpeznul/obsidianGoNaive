package domain

import (
	"fmt"
	cmn "obsidianGoNaive/protos/gen/common"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Из protobuf в доменную модель
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

// Из доменной модели в protobuf
func NoteToProto(note Note) *cmn.Note {
	return &cmn.Note{
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

// Из protobuf слайса в доменный слайс
func ProtoToNotes(protoNotes []*cmn.Note) ([]Note, error) {
	notes := make([]Note, 0, len(protoNotes))

	for _, protoNote := range protoNotes {
		note, err := ProtoToNote(protoNote)
		if err != nil {
			return nil, fmt.Errorf("failed to convert note: %w", err)
		}
		notes = append(notes, note)
	}

	return notes, nil
}

// Из доменного слайса в protobuf слайс
func NotesToProto(notes []Note) []*cmn.Note {
	protoNotes := make([]*cmn.Note, 0, len(notes))

	for _, note := range notes {
		protoNotes = append(protoNotes, NoteToProto(note))
	}

	return protoNotes
}
