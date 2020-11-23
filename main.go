/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a server for Greeter service.
package main

import (
	"context"
	"github.com/google/martian"
	"github.com/google/martian/fifo"
	"log"
	"net"
	"net/http"

	"google.golang.org/grpc"
	pb
)

const (
	port = ":50051"
)

var m = &Modifier{"Hello"}

type Modifier struct {
	name		string
}

func (m *Modifier) ModifyRequest(req *http.Request) error {
	m.name = req.RequestURI
	return nil
}

func (m *Modifier) ServeHTTP(http.ResponseWriter, *http.Request) {
	//Header().Set("Content-Type", "application/json")
	//WriteHeader(http.StatusOK)
	//Write([]byte(`{"message": "yo" }`))
}

func (m *Modifier) GetName() string {
	return m.name;
}

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) HelloWorld(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + "fucker!"}, nil
}
// protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative paradox.proto
func (s *server) CreateChan() string {
	return ""
}

func main() {

	proxy := martian.NewProxy()
	defer proxy.Close()
	listener, _ := net.Listen("tcp", ":8080")
	top := fifo.NewGroup()
	top.AddRequestModifier(m)
	proxy.SetRequestModifier(top)
	go proxy.Serve(listener)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}