package server

import (
	"context"

	pb "github.com/Budgetin-Project/user-service/app/proto/userservice"
)

type UserServerImpl struct {
	pb.UnimplementedUserServer
}

func (s *UserServerImpl) RegisterUser(context.Context, *pb.AuthenticationRequest) (*pb.RegisterResponse, error) {
	// TODO: Not yet implemented
	return &pb.RegisterResponse{UserId: 1}, nil
}
func (s *UserServerImpl) LoginUser(context.Context, *pb.AuthenticationRequest) (*pb.LoginResponse, error) {
	// TODO: Not yet implemented
	return &pb.LoginResponse{AuthToken: "dummy-token"}, nil
}
func (s *UserServerImpl) LogoutUser(context.Context, *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	// TODO: Not yet implemented
	return &pb.LogoutResponse{Success: true}, nil
}
func (s *UserServerImpl) VerifyEmailAddress(context.Context, *pb.VerifyEmailRequest) (*pb.VerifyEmailResponse, error) {
	// TODO: Not yet implemented
	return &pb.VerifyEmailResponse{Verified: true}, nil
}
