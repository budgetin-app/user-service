package interceptor

import (
	"context"
	"os"

	log "github.com/sirupsen/logrus"
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
		log.WithFields(log.Fields{"method": info.FullMethod}).Infof("Request -> %v", req)
	}

	// Call the actual handler to process the request
	resp, err := handler(ctx, req)

	// Log outgoing response if LOG_LEVEL is set to DEBUG
	if logLevel == "DEBUG" {
		if err != nil {
			log.WithFields(log.Fields{"method": info.FullMethod}).Errorf("Response -> %v", resp)
		} else {
			log.WithFields(log.Fields{"method": info.FullMethod}).Infof("Response -> %v", resp)
		}
	}
	return resp, err
}
