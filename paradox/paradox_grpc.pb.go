// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package greet

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// GreetServiceClient is the client API for GreetService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GreetServiceClient interface {
	// Unary
	Greet(ctx context.Context, in *GreetRequest, opts ...grpc.CallOption) (*GreetResponse, error)
	History(ctx context.Context, in *HistoryRequest, opts ...grpc.CallOption) (GreetService_HistoryClient, error)
}

type greetServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewGreetServiceClient(cc grpc.ClientConnInterface) GreetServiceClient {
	return &greetServiceClient{cc}
}

func (c *greetServiceClient) Greet(ctx context.Context, in *GreetRequest, opts ...grpc.CallOption) (*GreetResponse, error) {
	out := new(GreetResponse)
	err := c.cc.Invoke(ctx, "/greet.GreetService/Greet", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *greetServiceClient) History(ctx context.Context, in *HistoryRequest, opts ...grpc.CallOption) (GreetService_HistoryClient, error) {
	stream, err := c.cc.NewStream(ctx, &_GreetService_serviceDesc.Streams[0], "/greet.GreetService/History", opts...)
	if err != nil {
		return nil, err
	}
	x := &greetServiceHistoryClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type GreetService_HistoryClient interface {
	Recv() (*HistoryResponse, error)
	grpc.ClientStream
}

type greetServiceHistoryClient struct {
	grpc.ClientStream
}

func (x *greetServiceHistoryClient) Recv() (*HistoryResponse, error) {
	m := new(HistoryResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GreetServiceServer is the server API for GreetService service.
// All implementations must embed UnimplementedGreetServiceServer
// for forward compatibility
type GreetServiceServer interface {
	// Unary
	Greet(context.Context, *GreetRequest) (*GreetResponse, error)
	History(*HistoryRequest, GreetService_HistoryServer) error
	mustEmbedUnimplementedGreetServiceServer()
}

// UnimplementedGreetServiceServer must be embedded to have forward compatible implementations.
type UnimplementedGreetServiceServer struct {
}

func (UnimplementedGreetServiceServer) Greet(context.Context, *GreetRequest) (*GreetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Greet not implemented")
}
func (UnimplementedGreetServiceServer) History(*HistoryRequest, GreetService_HistoryServer) error {
	return status.Errorf(codes.Unimplemented, "method History not implemented")
}
func (UnimplementedGreetServiceServer) mustEmbedUnimplementedGreetServiceServer() {}

// UnsafeGreetServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GreetServiceServer will
// result in compilation errors.
type UnsafeGreetServiceServer interface {
	mustEmbedUnimplementedGreetServiceServer()
}

func RegisterGreetServiceServer(s grpc.ServiceRegistrar, srv GreetServiceServer) {
	s.RegisterService(&_GreetService_serviceDesc, srv)
}

func _GreetService_Greet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GreetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GreetServiceServer).Greet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/greet.GreetService/Greet",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GreetServiceServer).Greet(ctx, req.(*GreetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GreetService_History_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(HistoryRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(GreetServiceServer).History(m, &greetServiceHistoryServer{stream})
}

type GreetService_HistoryServer interface {
	Send(*HistoryResponse) error
	grpc.ServerStream
}

type greetServiceHistoryServer struct {
	grpc.ServerStream
}

func (x *greetServiceHistoryServer) Send(m *HistoryResponse) error {
	return x.ServerStream.SendMsg(m)
}

var _GreetService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "greet.GreetService",
	HandlerType: (*GreetServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Greet",
			Handler:    _GreetService_Greet_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "History",
			Handler:       _GreetService_History_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "paradox.proto",
}
