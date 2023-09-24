package document

import (
	"context"
	"demo-system-document-module/db"
	"demo-system-document-module/prisma"
	"demo-system-document-module/proto"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	proto.UnimplementedDocumentServiceServer
}

func (s *Server) GetDocumentList(ctx context.Context, req *proto.GetDocumentListRequest) (*proto.GetDocumentListReply, error) {

	client := prisma.GetPrismaClient()

	var documents = []*proto.Document{}

	records, err := client.Document.FindMany().Exec(ctx)
	if err != nil {
		panic(err)
	}

	for _, record := range records {

		documents = append(documents, &proto.Document{
			Id:        record.ID,
			Title:     record.Title,
			Content:   record.Content,
			CreatedAt: timestamppb.New(record.CreatedAt),
			UpdatedAt: timestamppb.New(record.UpdatedAt),
		})
	}
	return &proto.GetDocumentListReply{Documents: documents}, nil
}

func (s *Server) GetDocument(ctx context.Context, req *proto.GetDocumentRequest) (*proto.GetDocumentReply, error) {

	client := prisma.GetPrismaClient()

	record, err := client.Document.FindUnique(
		db.Document.ID.Equals(req.Id),
	).Exec(ctx)

	if err != nil {
		panic(err)
	}

	return &proto.GetDocumentReply{Document: &proto.Document{
		Id:        record.ID,
		Title:     record.Title,
		Content:   record.Content,
		CreatedAt: timestamppb.New(record.CreatedAt),
		UpdatedAt: timestamppb.New(record.UpdatedAt),
	}}, nil
}

func (s *Server) SaveDocument(ctx context.Context, req *proto.SaveDocumentRequest) (*proto.SaveDocumentReply, error) {
	id := uuid.New().String()

	client := prisma.GetPrismaClient()

	_, err := client.Document.CreateOne(
		db.Document.Title.Set(req.Title),
		db.Document.Content.Set(req.Content),
		db.Document.CreatedAt.Set(time.Now()),
		db.Document.UpdatedAt.Set(time.Now()),
	).Exec(ctx)

	if err != nil {
		panic(err)
	}

	return &proto.SaveDocumentReply{Id: id}, nil
}

func (s *Server) UpdateDocument(ctx context.Context, req *proto.UpdateDocumentRequest) (*proto.UpdateDocumentReply, error) {
	// Implement the UpdateDocument method here
	return &proto.UpdateDocumentReply{Id: req.Id}, nil
}

func (s *Server) DeleteDocument(ctx context.Context, req *proto.DeleteDocumentRequest) (*proto.DeleteDocumentReply, error) {
	// Implement the DeleteDocument method here
	return &proto.DeleteDocumentReply{Id: req.Id}, nil
}
