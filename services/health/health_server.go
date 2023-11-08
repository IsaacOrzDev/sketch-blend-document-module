package health

import (
	"context"
	"sketch-blend-document-module/proto"
)

type Server struct {
	proto.UnimplementedHealthServer
}

func (s *Server) Check(ctx context.Context, req *proto.HealthCheckRequest) (*proto.HealthCheckResponse, error) {
	return &proto.HealthCheckResponse{
		Status: proto.HealthCheckResponse_SERVING,
	}, nil
}
