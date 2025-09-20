package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/smukherj1/windows-agent/grpc/server"
	"google.golang.org/grpc"
)

var port = flag.Int("port", 8000, "The server port")

type server struct {
	pb.UnimplementedServiceServer
}

func (s *server) Hello(_ context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetMessage())
	return &pb.HelloReply{Message: "Hello " + in.GetMessage()}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterServiceServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
