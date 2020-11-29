package main

import (
	pb "awesomeProject3/greet"
	"context"
	"crypto/tls"
	"fmt"
	"github.com/google/martian"
	"github.com/google/martian/fifo"
	"github.com/google/martian/mitm"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"time"
)

// protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative paradox.proto

type Server struct {
	pb.UnimplementedGreetServiceServer
	history History
}

type History map[string][]byte

func (server *Server) ModifyRequest(req *http.Request) error {
	return nil
}

func (server *Server) ModifyResponse(res *http.Response) error {
	ctx := martian.NewContext(res.Request)

	raw, err := httputil.DumpResponse(res, true)
	if err != nil {
		log.Fatalf("Failed to dump response: %v", err)
		return err
	}

	id := ctx.ID()
	server.history[id] = raw

	return nil
}

func (server *Server) StartProxy() {
	proxy := martian.NewProxy()

	tr := &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: time.Second,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		DisableCompression: true,
	}
	proxy.SetRoundTripper(tr)

	x509c, priv, _ := mitm.NewAuthority("martian.proxy", "Martian Authority", 30*24*time.Hour)
	mc, _ := mitm.NewConfig(x509c, priv)

	mc.SkipTLSVerify(true)

	proxy.SetMITM(mc)

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen proxy: %v", err)
	}

	group := fifo.NewGroup()

	group.AddRequestModifier(server)
	group.AddResponseModifier(server)

	proxy.SetRequestModifier(group)
	proxy.SetResponseModifier(group)

	go proxy.Serve(listener)
}

func (server *Server) Greet(ctx context.Context, in *pb.GreetRequest) (*pb.GreetResponse, error) {
	log.Printf("Received: %v", in.GetResult())
	for key, _ := range server.history {
		fmt.Printf("Key: %v\n", key)
	}
	return &pb.GreetResponse{Result: "Complete:\n"}, nil
}

func (server *Server) History(in *pb.HistoryRequest, stream pb.GreetService_HistoryServer) error {
	fmt.Printf("GreetManyTimes function was invoked with %v\n", in.GetResult())
	for key, _ := range server.history {
		res := &pb.HistoryResponse{Result: key}
		stream.Send(res)
	}
	fmt.Printf("Completed\n")
	return nil
}

func main() {
	server := grpc.NewServer()

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen : %v", err)
	}

	modifier := &Server{history: make(map[string][]byte)}
	pb.RegisterGreetServiceServer(server, modifier)

	modifier.StartProxy()
	server.Serve(listener)
}
