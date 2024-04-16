package auth

import (
	"context"
	"errors"
	"github.com/HappyProgger/gRPC_auth/internal/services/auth"
	ssov1 "github.com/HappyProgger/gRPC_auth/protos/proto/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverAPI struct {
	ssov1.UnimplementedAuthServer
	auth Auth
}
type Auth interface {
	Login(
		ctx context.Context,
		email string,
		password string,
		appID int,
	) (token string, err error)
	RegisterNewUser(
		ctx context.Context,
		email string,
		password string,
	) (userID int64, err error)
	IsAdmin(
		ctx context.Context,
		userID int64,
	) (isAdmin bool, err error)
}

func Register(gRPCServer *grpc.Server, auth *auth.Auth) {
	ssov1.RegisterAuthServer(gRPCServer, &serverAPI{auth: auth})
}

func (s *serverAPI) Login(
	ctx context.Context,
	in *ssov1.LoginRequest,
) (*ssov1.LoginResponse, error) {
	if in.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}
	if in.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}
	if in.GetServiceId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "app_id is required")
	}
	token, err := s.auth.Login(ctx, in.GetEmail(), in.GetPassword(), int(in.GetServiceId()))
	if err != nil {
		// Ошибку auth.ErrInvalidCredentials создадим ниже
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid email or password")
		}
		return nil, status.Error(codes.Internal, "failed to login")
	}
	return &ssov1.LoginResponse{JwtToken: token}, nil
}

func (s *serverAPI) Register(
	ctx context.Context,
	in *ssov1.RegisterRequest,
) (*ssov1.RegisterResponse, error) {
	if in.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}
	if in.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}
	uid, err := s.auth.RegisterNewUser(ctx, in.GetEmail(), in.GetPassword())
	if err != nil {
		// Ошибку auth.ErrInvalidCredentials мы создадим ниже
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.AlreadyExists, "user already exist")
		}
		return nil, status.Error(codes.Internal, "failed to register user")
	}
	return &ssov1.RegisterResponse{UserId: int64(uid)}, nil
}

func (s *serverAPI) IsAdmin(
	ctx context.Context,
	in *ssov1.IsAdminRequest,
) (*ssov1.IsAdminResponse, error) {
	if in.UserId == 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	isAdmin, err := s.auth.IsAdmin(ctx, in.GetUserId())
	if err != nil {
		return nil, status.Error(codes.Internal, "user_id not exist")
	}
	return &ssov1.IsAdminResponse{IsAdmin: bool(isAdmin)}, nil
}

//func (s *serverAPI) IsAdmin(
//	ctx context.Context,
//	in *ssov1.IsAdminRequest,
//) (*ssov1.IsAdminResponse, error) {
//	// TODO dsafdsaf
//
//}
