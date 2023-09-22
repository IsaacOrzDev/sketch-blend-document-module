package testing

import (
	"context"
)

type Server struct {
	UnimplementedTestingServiceServer
}

func (s *Server) Test(ctx context.Context, req *TestRequest) (*TestResponse, error) {

	return &TestResponse{Message: "Testing"}, nil

}
