package grpc

import (
	"context"
	"fmt"
	authgrpc "github.com/HappyProgger/gRPC_auth/internal/gRPC/auth"
	"github.com/HappyProgger/gRPC_auth/internal/services/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"net"
)

type App struct {
	log        *slog.Logger
	gPPCServer grpc.Server
	port       int
	gRPCServer *grpc.Server
}

func New(log *slog.Logger, authService *auth.Auth, port int) *App {
	loggingOpts := []logging.Option{logging.WithLogOnEvents(
		logging.PayloadReceived, logging.PayloadSent,
	),
	}
	recoveryOpts := []recovery.Option{recovery.WithRecoveryHandler(func(p interface{}) (err error) {
		log.Error("Recovered from panic", slog.Any("panic", p))
		return status.Errorf(codes.Internal, "internal error")
	}),
	}
	gRPCServer := grpc.NewServer(grpc.ChainUnaryInterceptor(recovery.UnaryServerInterceptor(recoveryOpts...), logging.UnaryServerInterceptor(InterceptorLogger(log), loggingOpts...)))
	authgrpc.Register(gRPCServer, authService)
	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

// MustRun runs gRPC server and panics if any error occurs
func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

// Run runs gRPC server
func (a *App) Run() error {
	const op = "grpcapp.Run"

	// Создаем listener, который будет слушать TCP-сообщения, адресованные
	// Нашему gRPC-серверу
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	a.log.Info("grpc server started", slog.String("addr", l.Addr().String()))

	// Запускаем обработчик gRPC-сообщений
	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func InterceptorLogger(l *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}

func (a *App) Stop() {
	const op = "grpcapp.Stop"

	a.log.With(slog.String("op", op)).
		Info("stopping gRPC server", slog.Int("port", a.port))

	// Используем встроенный в gRPCServer механизм graceful shutdown
	a.gRPCServer.GracefulStop()
}
