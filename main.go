package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"sketch-blend-document-module/prisma"
	"sketch-blend-document-module/proto"
	"sketch-blend-document-module/services/document"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {

	server := grpc.NewServer()
	proto.RegisterDocumentServiceServer(server, &document.Server{})

	reflection.Register(server)

	listener, err := net.Listen("tcp", ":5003")
	if err != nil {
		panic(err)
	}

	go func() {
		if err := server.Serve(listener); err != nil {
			panic(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	fmt.Println("Shutting down server...")

	prisma.Disconnect()
	server.GracefulStop()

	fmt.Println("Server stopped.")

}
