package server

import (
	"github.com/Budgetin-Project/user-management-service/config"
	pb "github.com/Budgetin-Project/user-service/app/proto/userservice"
	"google.golang.org/grpc"
)

func InitServer(config *config.Configuration) *grpc.Server {
	// Create a new gRPC server
	server := grpc.NewServer()

	// Register the "service implementation (gRPC server methods) with the gRPC server
	pb.RegisterUserServer(server, NewUserServer(config.AuthController))

	return server
}
