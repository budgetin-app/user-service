package interceptor

import (
	"context"
	"log"
	"os"

	"google.golang.org/grpc"
)

// LoggingInterceptor intercept the incoming request and outcoming response
// when the LOG_LEVEL set to DEBUG
func LoggingInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	// Check the value of the LOG_LEVEL environment variable
	logLevel := os.Getenv("LOG_LEVEL")

	// Log incoming request if LOG_LEVEL is set to DEBUG
	if logLevel == "DEBUG" {
		log.Printf("[LOG] gRPC method: %s, request: %v", info.FullMethod, req)
	}

	// Call the actual handler to process the request
	resp, err := handler(ctx, req)

	// Log outgoing response if LOG_LEVEL is set to DEBUG
	if logLevel == "DEBUG" {
		log.Printf("[LOG] gRPC method: %s, response: %v", info.FullMethod, resp)
	}

	return resp, err
}
