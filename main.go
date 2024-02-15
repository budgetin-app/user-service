package main

import (
	"fmt"
	"net"

	log "github.com/sirupsen/logrus"

	"github.com/budgetin-app/user-management-service/config"
	"github.com/budgetin-app/user-service/app/pkg/helper/env"
	"github.com/budgetin-app/user-service/app/pkg/logger"
	"github.com/budgetin-app/user-service/app/server"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()

	// Initialize logger
	logger.InitLogger()
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
		log.WithFields(log.Fields{"port": port}).Fatal("Failed to listen")
	}

	// Log the server address where it's listening
	log.Infof("Server listening: %v", listen.Addr())

	// Start serving incoming gRPC requests
	if err := server.Serve(listen); err != nil {
		log.WithFields(log.Fields{
			"address": listen.Addr(),
		}).Fatal("Server failed to serve")
	}
}
