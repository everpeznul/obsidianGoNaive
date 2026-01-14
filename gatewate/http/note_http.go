package http

import (
	cmn "obsidianGoNaive/pkg/protos/gen/common"
	pbn "obsidianGoNaive/pkg/protos/gen/notes"
	pbu "obsidianGoNaive/pkg/protos/gen/updater"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Gateway struct {
	notesClient   pbn.NotesClient
	updaterClient pbu.UpdaterClient
}

func NewGateway(notesClient pbn.NotesClient, updaterClient pbu.UpdaterClient) *Gateway {

	return &Gateway{notesClient: notesClient, updaterClient: updaterClient}
}

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

func (nm *noteMapperHttp) HTTPToProto(note httpNote) cmn.Note {
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

func (nm *noteMapperHttp) HTTPToProtoSlice(notes []httpNote) []*cmn.Note {
	commonNotes := make([]*cmn.Note, 0, len(notes))
	for _, note := range notes {
		commonNotes = append(commonNotes, &cmn.Note{
			Id:         note.Id.String(),
			Title:      note.Title,
			Path:       note.Path,
			Class:      note.Class,
			Tags:       note.Tags,
			Links:      note.Links,
			Content:    note.Content,
			CreateTime: timestamppb.New(note.CreateTime),
			UpdateTime: timestamppb.New(note.UpdateTime),
		})
	}
	return commonNotes
}

func (nm *noteMapperHttp) ProtoToHTTP(note *cmn.Note) (httpNote, error) {
	id, err := uuid.Parse(note.Id)
	if err != nil {
		return httpNote{}, err
	}

	return httpNote{
		Id:         id,
		Title:      note.Title,
		Path:       note.Path,
		Class:      note.Class,
		Tags:       note.Tags,
		Links:      note.Links,
		Content:    note.Content,
		CreateTime: note.CreateTime.AsTime(),
		UpdateTime: note.UpdateTime.AsTime(),
	}, nil
}

func (nm *noteMapperHttp) ProtoToHTTPSlice(notes []*cmn.Note) ([]httpNote, error) {
	httpNotes := make([]httpNote, 0, len(notes))
	for _, note := range notes {
		hn, err := nm.ProtoToHTTP(note)
		if err != nil {
			return nil, err
		}
		httpNotes = append(httpNotes, hn)
	}
	return httpNotes, nil
}
