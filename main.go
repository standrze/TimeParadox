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
	pb "TimeParadox/paradox"
	"context"
	"github.com/google/martian"
	"github.com/google/martian/fifo"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"sync"
)

const (
	port        = ":50051"
	defaultPort = "8080"
)

type Server struct {
	pb.UnimplementedParadoxServer
	mu           sync.Mutex
	notification chan string
	latest       string
}

func (m *Server) ModifyRequest(req *http.Request) error {
	m.latest = req.RequestURI
	return nil
}

func (m *Server) HelloWorld(ctx context.Context, in *pb.Request) (*pb.Reply, error) {
	//s.mu.Lock()
	//st := <-s.notification
	//log.Printf("Received: %v", in.GetName())
	//s.mu.Unlock()
	return &pb.Reply{Message: "Hello " + m.latest}, nil
}

func (m *Server) StartProxy() {
	proxy := martian.NewProxy()
	listener, err := net.Listen("tcp", defaultPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	top := fifo.NewGroup()

	top.AddRequestModifier(m)
	proxy.SetRequestModifier(top)

	go proxy.Serve(listener)
}

func (m *Server) StartRPC() {
	server := grpc.NewServer()

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	pb.RegisterParadoxServer(server, m)

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
