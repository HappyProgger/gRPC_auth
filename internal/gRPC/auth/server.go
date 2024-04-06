package auth

import (
	"context"
	ssov1 "github.com/HappyProgger/gRPC_auth/protos/gen/go/happy.sso.v1"
	"github.com/golang/protobuf/protoc-gen-go/grpc"
	// Сгенерированный код
)

type serverAPI struct {
	ssov1.UnimplementedAuthServer // Хитрая штука, о ней ниже
	auth                          Auth
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
	) (userID int64, err error)
	IsAdmin(
		ctx context.Context,
		userId int64,
	) (isAdmin bool, err error)
}

func Register(gRPCServer *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServer(gRPCServer, &serverAPI{auth: auth})
}

func (s *serverAPI) Login(
	ctx context.Context,
	in *ssov1.LoginRequest,
) (*ssov1.LoginResponse, error) {
	// TODO
}

func (s *serverAPI) Register(
	ctx context.Context,
	in *ssov1.RegisterRequest,
) (*ssov1.RegisterResponse, error) {
	// TODO
}
