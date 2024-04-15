package app

import (
	grpcapp "github.com/HappyProgger/gRPC_auth/internal/app/grpc"
	"github.com/HappyProgger/gRPC_auth/internal/services/auth"
	"github.com/HappyProgger/gRPC_auth/storage/postgres"
	"log/slog"
	"time"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	cfgPath string,
	tokenTTL time.Duration,
) *App {
	//todo вставить динамический путь до приложения 23 string

	storage, err := postgres.New(cfgPath)
	if err != nil {
		panic(err)
	}

	authService := auth.New(log, storage, storage, storage, tokenTTL)
	grpcApp := grpcapp.New(log, authService, grpcPort)
	return &App{
		GRPCServer: grpcApp,
	}
}
