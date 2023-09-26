package document

import (
	"context"
	"demo-system-document-module/db"
	"demo-system-document-module/prisma"
	"demo-system-document-module/proto"
	"time"

	"github.com/golang/protobuf/jsonpb"
	"github.com/steebchen/prisma-client-go/runtime/types"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	proto.UnimplementedDocumentServiceServer
}

func documentToProto(record *db.DocumentModel) *proto.Document {
	image, _ := record.Image()
	svg, _ := record.Svg()
	paths, _ := record.Paths()

	pathsStr := "{}"

	if paths != nil {
		pathsStr = string(paths)
	}

	pathsStruct := &structpb.Struct{}

	err := jsonpb.UnmarshalString(pathsStr, pathsStruct)
	if err != nil {
		panic(err)
	}

	return &proto.Document{
		Id:        record.ID,
		UserId:    uint32(record.UserID),
		Title:     record.Title,
		CreatedAt: timestamppb.New(record.CreatedAt),
		UpdatedAt: timestamppb.New(record.UpdatedAt),
		Image:     &image,
		Svg:       &svg,
		Paths:     pathsStruct,
	}
}

func (s *Server) GetDocumentList(ctx context.Context, req *proto.GetDocumentListRequest) (*proto.GetDocumentListReply, error) {

	client := prisma.GetPrismaClient()

	var documents = []*proto.Document{}

	var filter []db.DocumentWhereParam
	if req.UserId != nil {
		filter = append(filter, db.Document.UserID.Equals(int(*req.UserId)))
	}

	var query = client.Document.FindMany(filter...)
	if req.Limit != nil {
		query = query.Take(int(*req.Limit))
	}
	if req.Offset != nil {
		query = query.Skip(int(*req.Offset))
	}

	records, err := query.Exec(ctx)
	if err != nil {
		panic(err)
	}

	for _, record := range records {

		documents = append(documents, documentToProto(&record))
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

	return &proto.GetDocumentReply{Document: documentToProto(record)}, nil
}

func (s *Server) SaveDocument(ctx context.Context, req *proto.SaveDocumentRequest) (*proto.SaveDocumentReply, error) {

	client := prisma.GetPrismaClient()

	var err error
	paths := "{}"
	var description string
	var image string
	var svg string

	pathsJSON := structpb.NewStructValue(req.Document.Paths)

	if pathsJSON != nil {
		marshaler := &jsonpb.Marshaler{}
		paths, err = marshaler.MarshalToString(pathsJSON)
		if err != nil {
			panic(err)
		}
	}

	if req.Document.Description != nil {
		description = *req.Document.Description
	}
	if req.Document.Image != nil {
		image = *req.Document.Image
	}
	if req.Document.Svg != nil {
		svg = *req.Document.Svg
	}

	record, err := client.Document.CreateOne(
		db.Document.UserID.Set(int(req.UserId)),
		db.Document.Title.Set(req.Document.Title),
		db.Document.CreatedAt.Set(time.Now()),
		db.Document.UpdatedAt.Set(time.Now()),
		db.Document.Description.Set(description),
		db.Document.Image.Set(image),
		db.Document.Svg.Set(svg),
		db.Document.Paths.Set(types.JSON(paths)),
	).Exec(ctx)

	if err != nil {
		panic(err)
	}

	return &proto.SaveDocumentReply{Id: record.ID}, nil
}

func (s *Server) UpdateDocument(ctx context.Context, req *proto.UpdateDocumentRequest) (*proto.UpdateDocumentReply, error) {
	// Implement the UpdateDocument method here
	return &proto.UpdateDocumentReply{Id: req.Id}, nil
}

func (s *Server) DeleteDocument(ctx context.Context, req *proto.DeleteDocumentRequest) (*proto.DeleteDocumentReply, error) {
	// Implement the DeleteDocument method here
	return &proto.DeleteDocumentReply{Id: req.Id}, nil
}
