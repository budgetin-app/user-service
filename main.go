package main

import (
	"fmt"
	"log"
	"net"

	"github.com/Budgetin-Project/user-management-service/config"
	"github.com/Budgetin-Project/user-service/app/pkg/helper/env"
	"github.com/Budgetin-Project/user-service/app/server"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}

func main() {
	// Initialize the configuration for dependency injection
	cfg := config.Configure()

	// Initialize grpc server
	server := server.InitServer(cfg)

	// Listener for incoming TCP connections on the specified ports
	port := env.GetenvOrDefault("SERVER_PORT", "50051")
	listen, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Log the server address where it's listening
	log.Printf("server listening at %v", listen.Addr())

	// Start serving incoming gRPC requests
	if err := server.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
