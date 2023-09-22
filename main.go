package main

import (
	testing "demo-system-document-module/proto"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {

	server := grpc.NewServer()
	testing.RegisterTestingServiceServer(server, &testing.Server{})
	reflection.Register(server)

	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		panic(err)
	}

	if err := server.Serve(listener); err != nil {
		panic(err)
	}

}
