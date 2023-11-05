package grpc

import (
	"context"
	"fmt"
	"github.com/mirshahriar/multiplexing-simple/grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct{}


func (s *grpcServer) EchoMessage(ctx context.Context, req *proto.EchoRequest) (*proto.EchoResponse, error) {
	fmt.Println("gRPC server request received")
	return &proto.EchoResponse{Message: fmt.Sprintf("echo %s from grpc!", req.Message)}, nil
}

func NewGRPCServer() *grpc.Server {
	server := grpc.NewServer()
	proto.RegisterEchoServiceServer(server, &grpcServer{})
	reflection.Register(server)
	return server
}

func NewGRPCHandler() proto.EchoServiceServer {
	return &grpcServer{}
}
