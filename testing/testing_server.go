package testing

import (
	"context"
	"demo-system-document-module/prisma"
	"demo-system-document-module/proto"
)

type Server struct {
	proto.UnimplementedTestingServiceServer
}

func (s *Server) Test(ctx context.Context, req *proto.TestRequest) (*proto.TestResponse, error) {

	client := prisma.GetPrismaClient()

	var names []string

	users, err := client.User.FindMany().Exec(ctx)
	if err != nil {
		panic(err)
	}

	for _, user := range users {
		names = append(names, user.Name)
	}

	return &proto.TestResponse{
		Names: names,
	}, nil

}
