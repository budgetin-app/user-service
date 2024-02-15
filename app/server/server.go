package server

import (
	"github.com/Budgetin-Project/user-management-service/config"
	pb "github.com/Budgetin-Project/user-service/app/proto"
	"github.com/Budgetin-Project/user-service/app/server/interceptor"
	"google.golang.org/grpc"
)

func InitServer(config *config.Configuration) *grpc.Server {
	// Create a new gRPC server
	server := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.LoggingInterceptor),
	)

	// Register the "service implementation (gRPC server methods) with the gRPC server
	pb.RegisterUserServer(server, NewUserServer(config.AuthController))

	return server
}
