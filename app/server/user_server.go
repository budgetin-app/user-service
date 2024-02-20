package server

import (
	"context"

	"github.com/budgetin-app/user-service/app/controller"
	"github.com/budgetin-app/user-service/app/pkg/validator"
	pb "github.com/budgetin-app/user-service/app/proto"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserServerImpl struct {
	authController controller.AuthController
	pb.UnimplementedUserServer
}

func NewUserServer(authController controller.AuthController) *UserServerImpl {
	return &UserServerImpl{authController: authController}
}

func (s *UserServerImpl) RegisterUser(ctx context.Context, r *pb.AuthenticationRequest) (*pb.RegisterResponse, error) {
	// Request validation
	if !validator.IsValidUsername(r.Username) {
		return nil, status.Error(codes.InvalidArgument, "invalid username")
	}
	if !validator.IsValidEmail(r.Email) {
		return nil, status.Error(codes.InvalidArgument, "invalid email")
	}
	if !validator.IsValidPassword(r.Password) {
		return nil, status.Error(codes.InvalidArgument, "invalid password")
	}

	// Begin to register new user
	credential, err := s.authController.Register(r.Username, r.Email, r.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to register user: %v", err)
	}

	return &pb.RegisterResponse{UserId: *proto.Uint32(uint32(credential.ID))}, nil
}

func (s *UserServerImpl) LoginUser(ctx context.Context, r *pb.AuthenticationRequest) (*pb.LoginResponse, error) {
	// Request validation
	if !validator.IsValidPassword(r.Password) {
		return nil, status.Error(codes.InvalidArgument, "invalid password")
	}

	// Identifying the received identifier is email or username
	var isEmail bool
	var identifier string

	switch {
	case len(r.Username) > 0 && !validator.IsValidUsername(r.Username):
		return nil, status.Error(codes.InvalidArgument, "invalid username")
	case len(r.Email) > 0 && !validator.IsValidEmail(r.Email):
		return nil, status.Error(codes.InvalidArgument, "invalid email")
	case len(r.Username) > 0:
		isEmail, identifier = false, r.Username
	case len(r.Email) > 0:
		isEmail, identifier = true, r.Email
	default:
		return nil, status.Error(codes.InvalidArgument, "either username or email must be provided")
	}

	// Begin to authenticate user
	session, err := s.authController.Login(isEmail, identifier, r.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to login user: %v", err)
	}

	return &pb.LoginResponse{AuthToken: session.Token}, nil
}

func (s *UserServerImpl) LogoutUser(ctx context.Context, r *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	// Request validation
	if len(r.AuthToken) == 0 {
		return nil, status.Error(codes.InvalidArgument, "authentication token must be provided")
	}

	// Begin to logout the user
	success, err := s.authController.Logout(r.AuthToken)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to logout user: %v", err)
	}

	return &pb.LogoutResponse{Success: success}, nil
}

func (s *UserServerImpl) VerifyEmailAddress(ctx context.Context, r *pb.VerifyEmailRequest) (*pb.VerifyEmailResponse, error) {
	// Request validation
	if !validator.IsValidEmail(r.Email) {
		return nil, status.Error(codes.InvalidArgument, "invalid email")
	}

	// Begin to verify the email address
	verified, err := s.authController.VerifyEmail(r.Email)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to verify email: %v", err)
	}

	return &pb.VerifyEmailResponse{Verified: verified}, nil
}
