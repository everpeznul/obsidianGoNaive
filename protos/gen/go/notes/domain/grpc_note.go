package domain

import (
	"context"
	pb "obsidianGoNaive/protos/gen/go/notes"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type NoteServer struct {
	pb.UnimplementedNotesServer
}

func NewNoteServer() *NoteServer {

	return &NoteServer{}
}

func (ns *NoteServer) Create(ctx context.Context, r *pb.CreateRequest) (*pb.CreateResponse, error) {

	domainNote, _ := ProtoToNote(r.Note)
	id, err := Repo.Insert(ctx, domainNote)

	return &pb.CreateResponse{Note: &pb.Note{Id: id.String()}}, err
}
func (ns *NoteServer) GetByID(ctx context.Context, r *pb.GetByIdRequest) (*pb.GetByIdResponse, error) {

	id, err := uuid.Parse(r.Id)
	domainNote, err := Repo.GetByID(ctx, id)
	protoNote := NoteToProto(domainNote)

	return &pb.GetByIdResponse{Note: protoNote}, err
}
func (ns *NoteServer) Find(ctx context.Context, r *pb.FindRequest) (*pb.FindResponse, error) {

	if r.Limit == 0 {

		domainNotes, err := Repo.GetAll(ctx)
		protoNotes := NotesToProto(domainNotes)

		return &pb.FindResponse{Note: protoNotes}, err
	}
	if r.Name != "" {

		domainNote, err := Repo.FindByName(ctx, r.Name)
		protoNote := NoteToProto(domainNote)

		return &pb.FindResponse{Note: []*pb.Note{protoNote}}, err
	}
	if r.Ancestor != "" {

		domainNote, err := Repo.FindByAncestor(ctx, r.Ancestor)
		protoNote := NotesToProto(domainNote)

		return &pb.FindResponse{Note: protoNote}, err
	}
	return nil, status.Error(codes.NotFound, "note not found")
}
func (ns *NoteServer) UpdateById(ctx context.Context, r *pb.UpdateByIdRequest) (*pb.UpdateByIdResponse, error) {

	protoNote := r.Note
	domainNote, err := ProtoToNote(protoNote)
	err = Repo.UpdateById(ctx, domainNote)

	return &pb.UpdateByIdResponse{}, err

}
func (ns *NoteServer) DeleteById(ctx context.Context, r *pb.DeleteByIdRequest) (*pb.DeleteByIdResponse, error) {

	id, err := uuid.Parse(r.Id)
	err = Repo.DeleteById(ctx, id)

	return &pb.DeleteByIdResponse{}, err
}
