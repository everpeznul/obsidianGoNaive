package transport

import (
	"context"
	"obsidianGoNaive/pkg/protos/gen/common"
	"obsidianGoNaive/pkg/protos/gen/notes"
	"obsidianGoNaive/services/notes/internal/repository"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type NoteService struct {
	notes.UnimplementedNotesServer
	Repo *repository.Repository
}

func NewNotesService(repo *repository.Repository) *NoteService {

	return &NoteService{Repo: repo}
}

func (ns *NoteService) Create(ctx context.Context, r *notes.CreateRequest) (*notes.CreateResponse, error) {

	domainNote, _ := ProtoToNote(r.Note)
	id, err := ns.Repo.DB.Insert(ctx, domainNote)

	return &notes.CreateResponse{Note: &common.Note{Id: id.String()}}, err
}
func (ns *NoteService) GetByID(ctx context.Context, r *notes.GetByIdRequest) (*notes.GetByIdResponse, error) {

	id, err := uuid.Parse(r.Id)
	domainNote, err := ns.Repo.DB.GetByID(ctx, id)
	protoNote := NoteToProto(domainNote)

	return &notes.GetByIdResponse{Note: protoNote}, err
}
func (ns *NoteService) Find(ctx context.Context, r *notes.FindRequest) (*notes.FindResponse, error) {

	if r.Name != "" {

		domainNote, err := ns.Repo.DB.FindByName(ctx, r.Name)
		protoNote := NoteToProto(domainNote)

		return &notes.FindResponse{Note: []*common.Note{protoNote}}, err
	}
	if r.Ancestor != "" {

		domainNote, err := ns.Repo.DB.FindByAncestor(ctx, r.Ancestor)
		protoNote := NotesToProto(domainNote)

		return &notes.FindResponse{Note: protoNote}, err
	}
	if r.Limit == 0 {

		domainNotes, err := ns.Repo.DB.GetAll(ctx)
		protoNotes := NotesToProto(domainNotes)

		return &notes.FindResponse{Note: protoNotes}, err
	}
	return nil, status.Error(codes.NotFound, "note not found")
}
func (ns *NoteService) UpdateById(ctx context.Context, r *notes.UpdateByIdRequest) (*notes.UpdateByIdResponse, error) {

	protoNote := r.Note
	domainNote, err := ProtoToNote(protoNote)
	err = ns.Repo.DB.UpdateById(ctx, domainNote)

	return &notes.UpdateByIdResponse{}, err

}
func (ns *NoteService) DeleteById(ctx context.Context, r *notes.DeleteByIdRequest) (*notes.DeleteByIdResponse, error) {

	id, err := uuid.Parse(r.Id)
	err = ns.Repo.DB.DeleteById(ctx, id)

	return &notes.DeleteByIdResponse{}, err
}
