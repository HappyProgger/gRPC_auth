package app

import (
	grpcapp "github.com/HappyProgger/gRPC_auth/internal/app/grpc"
	cfg "github.com/HappyProgger/gRPC_auth/internal/config"
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
	////////////////
	log.Info("starting url_shorter",
		slog.String("env", cfg.Cfg.Env),
		slog.String("version", "version_example"),
	)
	log.Debug("debug messages are enabled")

	ssoClient, err := ssogrpc
	////////////////
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
