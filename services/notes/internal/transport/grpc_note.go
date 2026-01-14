package transport

import (
	"context"
	"obsidianGoNaive/protos/gen/common"
	pb "obsidianGoNaive/protos/gen/notes"
	"obsidianGoNaive/services/notes/internal/repository"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type NoteService struct {
	pb.UnimplementedNotesServer
	Repo *repository.Repository
}

func NewNoteService(repo *repository.Repository) *NoteService {

	return &NoteService{Repo: repo}
}

func (ns *NoteService) Create(ctx context.Context, r *pb.CreateRequest) (*pb.CreateResponse, error) {

	domainNote, _ := ProtoToNote(r.Note)
	id, err := ns.Repo.DB.Insert(ctx, domainNote)

	return &pb.CreateResponse{Note: &common.Note{Id: id.String()}}, err
}
func (ns *NoteService) GetByID(ctx context.Context, r *pb.GetByIdRequest) (*pb.GetByIdResponse, error) {

	id, err := uuid.Parse(r.Id)
	domainNote, err := ns.Repo.DB.GetByID(ctx, id)
	protoNote := NoteToProto(domainNote)

	return &pb.GetByIdResponse{Note: protoNote}, err
}
func (ns *NoteService) Find(ctx context.Context, r *pb.FindRequest) (*pb.FindResponse, error) {

	if r.Name != "" {

		domainNote, err := ns.Repo.DB.FindByName(ctx, r.Name)
		protoNote := NoteToProto(domainNote)

		return &pb.FindResponse{Note: []*common.Note{protoNote}}, err
	}
	if r.Ancestor != "" {

		domainNote, err := ns.Repo.DB.FindByAncestor(ctx, r.Ancestor)
		protoNote := NotesToProto(domainNote)

		return &pb.FindResponse{Note: protoNote}, err
	}
	if r.Limit == 0 {

		domainNotes, err := ns.Repo.DB.GetAll(ctx)
		protoNotes := NotesToProto(domainNotes)

		return &pb.FindResponse{Note: protoNotes}, err
	}
	return nil, status.Error(codes.NotFound, "note not found")
}
func (ns *NoteService) UpdateById(ctx context.Context, r *pb.UpdateByIdRequest) (*pb.UpdateByIdResponse, error) {

	protoNote := r.Note
	domainNote, err := ProtoToNote(protoNote)
	err = ns.Repo.DB.UpdateById(ctx, domainNote)

	return &pb.UpdateByIdResponse{}, err

}
func (ns *NoteService) DeleteById(ctx context.Context, r *pb.DeleteByIdRequest) (*pb.DeleteByIdResponse, error) {

	id, err := uuid.Parse(r.Id)
	err = ns.Repo.DB.DeleteById(ctx, id)

	return &pb.DeleteByIdResponse{}, err
}
