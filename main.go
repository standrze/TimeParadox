package main

import (
	pb "TimeParadox/paradox"
	"context"
	"fmt"
	"github.com/google/martian"
	"github.com/google/martian/fifo"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
)

type Server struct {
	pb.UnimplementedParadoxServer
	notification chan string
	history      History
}

type History map[string][]byte

func (s *Server) ModifyRequest(req *http.Request) error {
	fmt.Println("===Request===")
	r, _ := httputil.DumpRequestOut(req, true)
	fmt.Println(string(r))
	return nil
}

func (s *Server) ModifyResponse(res *http.Response) error {
	fmt.Println("===Response===")
	r, _ := httputil.DumpRequestOut(res.Request, true)
	fmt.Println(string(r))
	r, _ = httputil.DumpResponse(res, true)
	fmt.Println(string(r))
	b, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(b))
	return nil
}

func (s *Server) HelloWorld(ctx context.Context, in *pb.Request) (*pb.Reply, error) {
	log.Printf("Received: %v", in.GetName())
	fmt.Println(len(s.history))
	return &pb.Reply{Message: "Complete"}, nil
}

func (s *Server) StartProxy() {
	proxy := martian.NewProxy()

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen proxy: %v", err)
	}

	top := fifo.NewGroup()

	top.AddRequestModifier(s)
	top.AddResponseModifier(s)

	proxy.SetRequestModifier(top)
	proxy.SetResponseModifier(top)

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
	notification := make(chan string)

	mod := &Server{
		notification: notification,
		history:      make(map[string][]byte),
	}

	mod.StartProxy()
	mod.StartRPC()
}
