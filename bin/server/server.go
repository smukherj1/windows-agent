package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	pb "github.com/smukherj1/windows-agent/grpc/server"
	"golang.org/x/crypto/ssh"
	"google.golang.org/grpc"
)

var (
	port         = flag.Int("port", 8000, "The server port")
	remoteServer = flag.String("remote.server", "", "The address of the remote ssh server.")
	remoteUser   = flag.String("remote.user", "", "The user on the remote ssh server.")
	privateKey   = flag.String("private.key", "", "The path to the private key for the remote ssh server.")
)

type server struct {
	pb.UnimplementedServiceServer
}

func (s *server) Hello(_ context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetMessage())
	return &pb.HelloReply{Message: "Hello " + in.GetMessage()}, nil
}

func startReverseSshTunnel(port int, remoteServer, remoteUser, privateKey string) {
	key, err := os.ReadFile(privateKey)
	if err != nil {
		log.Fatalf("unable to read private key: %v", err)
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("unable to parse private key: %v", err)
	}

	config := &ssh.ClientConfig{
		User: remoteUser,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Connect to the remote SSH server.
	client, err := ssh.Dial("tcp", remoteServer, config)
	if err != nil {
		log.Fatalf("unable to connect to remote server: %v", err)
	}
	defer client.Close()

	// Listen on the remote server.
	remotePort := fmt.Sprintf(":%d", port)
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
		go startReverseSshTunnel(*port, *remoteServer, *remoteUser, *privateKey)
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
