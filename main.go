package main

import (
	pb "TimeParadox/paradox"
	"context"
	"fmt"
	"github.com/google/martian"
	"github.com/google/martian/fifo"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"sync"
)

type Server struct {
	pb.UnimplementedParadoxServer
	notification chan string
	latest       string
}

func (s *Server) ModifyRequest(req *http.Request) error {
	return nil
}

func (s *Server) HelloWorld(ctx context.Context, in *pb.Request) (*pb.Reply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.Reply{Message: "Hello " + s.latest}, nil
}

func (s *Server) StartProxy() {
	proxy := martian.NewProxy()

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen proxy: %v", err)
	}

	top := fifo.NewGroup()

	top.AddRequestModifier(s)
	proxy.SetRequestModifier(top)

	go proxy.Serve(listener)
}

func (s *Server) StartRPC() {
	server := grpc.NewServer()

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen : %v", err)
	}

	pb.RegisterParadoxServer(server, s)

	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func main() {

	var mutex sync.Mutex
	notification := make(chan string)

	mod := &Server{
		mu:           mutex,
		notification: notification,
		latest:       "",
	}

	mod.StartProxy()
	mod.StartRPC()
}
