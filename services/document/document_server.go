package document

import (
	"context"
	"sketch-blend-document-module/db"
	"sketch-blend-document-module/prisma"
	"sketch-blend-document-module/proto"
	"time"

	"github.com/golang/protobuf/jsonpb"
	"github.com/steebchen/prisma-client-go/runtime/types"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	proto.UnimplementedDocumentServiceServer
}

func documentToProtoDetail(record *db.DocumentModel) *proto.DocumentDetail {
	image, _ := record.Image()
	svg, _ := record.Svg()
	paths, _ := record.Paths()
	description, _ := record.Description()

	pathsStr := "{}"

	if paths != nil {
		pathsStr = string(paths)
	}

	pathsStruct := &structpb.Struct{}

	err := jsonpb.UnmarshalString(pathsStr, pathsStruct)
	if err != nil {
		panic(err)
	}

	return &proto.DocumentDetail{
		Id:          record.ID,
		UserId:      uint32(record.UserID),
		Title:       record.Title,
		Description: &description,
		CreatedAt:   timestamppb.New(record.CreatedAt),
		UpdatedAt:   timestamppb.New(record.UpdatedAt),
		Image:       &image,
		Svg:         &svg,
		Paths:       pathsStruct,
	}
}

func documentToProto(record *db.DocumentModel) *proto.Document {
	image := ""
	svg, _ := record.Svg()
	description, _ := record.Description()

	return &proto.Document{
		Id:          record.ID,
		UserId:      uint32(record.UserID),
		Title:       record.Title,
		CreatedAt:   timestamppb.New(record.CreatedAt),
		UpdatedAt:   timestamppb.New(record.UpdatedAt),
		Description: &description,
		Svg:         &svg,
		Image:       &image,
	}
}

func (s *Server) GetDocumentList(ctx context.Context, req *proto.GetDocumentListRequest) (*proto.GetDocumentListReply, error) {

	client := prisma.GetPrismaClient()

	var documents = []*proto.Document{}

	var filter []db.DocumentWhereParam
	if req.UserId != nil {
		filter = append(filter, db.Document.UserID.Equals(int(*req.UserId)))
	}

	var query = client.Document.FindMany(filter...).OrderBy(db.Document.CreatedAt.Order(db.SortOrderDesc))
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

	return &proto.GetDocumentListReply{Records: documents}, nil
}

func (s *Server) GetDocument(ctx context.Context, req *proto.GetDocumentRequest) (*proto.GetDocumentReply, error) {

	client := prisma.GetPrismaClient()

	record, err := client.Document.FindUnique(
		db.Document.ID.Equals(req.Id),
	).Exec(ctx)

	if err != nil {
		return &proto.GetDocumentReply{}, nil
	}

	return &proto.GetDocumentReply{Record: documentToProtoDetail(record)}, nil
}

func (s *Server) SaveDocument(ctx context.Context, req *proto.SaveDocumentRequest) (*proto.SaveDocumentReply, error) {

	client := prisma.GetPrismaClient()

	var err error
	paths := "{}"
	var description string
	var image string
	var svg string

	pathsJSON := structpb.NewStructValue(req.Data.Paths)

	if pathsJSON != nil {
		marshaler := &jsonpb.Marshaler{}
		paths, err = marshaler.MarshalToString(pathsJSON)
		if err != nil {
			panic(err)
		}
	}

	if req.Data.Description != nil {
		description = *req.Data.Description
	}
	if req.Data.Image != nil {
		image = *req.Data.Image
	}
	if req.Data.Svg != nil {
		svg = *req.Data.Svg
	}

	record, err := client.Document.CreateOne(
		db.Document.UserID.Set(int(req.UserId)),
		db.Document.Title.Set(req.Data.Title),
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

	client := prisma.GetPrismaClient()

	var err error
	paths := "{}"

	pathsJSON := structpb.NewStructValue(req.Data.Paths)

	if pathsJSON != nil {
		marshaler := &jsonpb.Marshaler{}
		paths, err = marshaler.MarshalToString(pathsJSON)
		if err != nil {
			panic(err)
		}
	}

	var updates []db.DocumentSetParam

	if req.Data.Title != "" {
		updates = append(updates, db.Document.Title.Set(req.Data.Title))
	}

	if req.Data.Description != nil {
		updates = append(updates, db.Document.Description.Set(*req.Data.Description))
	}
	if req.Data.Image != nil {
		updates = append(updates, db.Document.Image.Set(*req.Data.Image))
	}
	if req.Data.Svg != nil {
		updates = append(updates, db.Document.Svg.Set(*req.Data.Svg))
	}

	if req.Data.Paths != nil {
		updates = append(updates, db.Document.Paths.Set(types.JSON(paths)))
	}

	updates = append(updates, db.Document.UpdatedAt.Set(time.Now()))

	_, err = client.Document.FindUnique(
		db.Document.ID.Equals(req.Id),
	).Update(
		updates...,
	).Exec(ctx)

	if err != nil {
		panic(err)
	}

	return &proto.UpdateDocumentReply{Id: req.Id}, nil
}

func (s *Server) DeleteDocument(ctx context.Context, req *proto.DeleteDocumentRequest) (*proto.DeleteDocumentReply, error) {

	client := prisma.GetPrismaClient()

	_, err := client.Document.FindUnique(
		db.Document.ID.Equals(req.Id),
	).Delete().Exec(ctx)

	if err != nil {
		panic(err)
	}

	return &proto.DeleteDocumentReply{Id: req.Id}, nil
}
