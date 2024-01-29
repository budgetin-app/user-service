package server

import (
	pb "github.com/Budgetin-Project/user-service/app/proto/userservice"
	"google.golang.org/grpc"
)

func InitServer() *grpc.Server {
	// Create a new gRPC server
	server := grpc.NewServer()

	// Register the "service implementation (gRPC server methods) with the gRPC server
	pb.RegisterUserServer(server, &UserServerImpl{})

	return server
}
