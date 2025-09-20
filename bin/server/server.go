package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	pb "github.com/smukherj1/windows-agent/grpc/server"
	"golang.org/x/crypto/ssh"
	"google.golang.org/grpc"
)

var (
	port           = flag.Int("port", 8000, "The server port")
	remoteServer   = flag.String("remote.server", "", "The address of the remote ssh server.")
	remoteUser     = flag.String("remote.user", "", "The user on the remote ssh server.")
	remotePassword = flag.String("remote.password", "", "The ssh password for the remote ssh server.")
)

type server struct {
	pb.UnimplementedServiceServer
}

func (s *server) Hello(_ context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetMessage())
	return &pb.HelloReply{Message: "Hello " + in.GetMessage()}, nil
}

func startReverseSshTunnel(port int, remoteServer, remoteUser, remotePassword string) {
	config := &ssh.ClientConfig{
		User: remoteUser,
		Auth: []ssh.AuthMethod{
			ssh.Password(remotePassword),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	log.Printf("Starting reverse SSH tunnel from %v@%v:%v (password length=%v) -> localhost:%v.", remoteUser, remoteServer, port, len(remotePassword), port)
	// Connect to the remote SSH server.
	client, err := ssh.Dial("tcp", remoteServer+":22", config)
	if err != nil {
		log.Fatalf("unable to connect to remote server: %v", err)
	}
	defer client.Close()

	// Listen on the remote server.
	remotePort := fmt.Sprintf("localhost:%d", port)
	l, err := client.Listen("tcp", remotePort)
	if err != nil {
		log.Fatalf("unable to listen on remote server: %v", err)
	}
	defer l.Close()

	localPort := fmt.Sprintf("localhost:%d", port)
	for {
		remote, err := l.Accept()
		if err != nil {
			log.Printf("failed to accept remote connection: %v", err)
			continue
		}
		log.Println("Got remote connection from", remote.RemoteAddr(), "to", remote.LocalAddr())

		local, err := net.Dial("tcp", localPort)
		if err != nil {
			log.Printf("failed to dial local server: %v", err)
			continue
		}

		go func() {
			defer local.Close()
			defer remote.Close()
			io.Copy(local, remote)
		}()

		go func() {
			defer local.Close()
			defer remote.Close()
			io.Copy(remote, local)
		}()
	}
}

func main() {
	flag.Parse()

	if *remoteServer != "" {
		go startReverseSshTunnel(*port, *remoteServer, *remoteUser, *remotePassword)
	}

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
