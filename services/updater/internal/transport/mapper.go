package transport

import (
	"fmt"
	"obsidianGoNaive/pkg/protos/gen/common"
	"obsidianGoNaive/services/updater/internal/domain"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Из protobuf в доменную модель
func ProtoToNote(protoNote *common.Note) (*domain.Note, error) {
	id, err := uuid.Parse(protoNote.Id)
	if err != nil {
		return nil, fmt.Errorf("invalid UUID: %w", err)
	}

	return &domain.Note{
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
func NoteToProto(note *domain.Note) *common.Note {
	return &common.Note{
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
func ProtoToNotes(protoNotes []*common.Note) ([]*domain.Note, error) {
	notes := make([]*domain.Note, len(protoNotes))

	for i := range protoNotes {
		note, err := ProtoToNote(protoNotes[i])
		if err != nil {
			return nil, fmt.Errorf("failed to convert note: %w", err)
		}
		notes[i] = note
	}

	return notes, nil
}

// Из доменного слайса в protobuf слайс
func NotesToProto(notes []*domain.Note) []*common.Note {
	protoNotes := make([]*common.Note, len(notes))

	for i := range notes {
		protoNotes[i] = NoteToProto(notes[i])
	}

	return protoNotes
}
