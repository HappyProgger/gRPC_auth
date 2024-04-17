package grpc

import (
	"context"
	"fmt"
	ssov1 "github.com/HappyProgger/gRPC_auth/protos/proto/gen/go/sso"
	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	slog "log/slog"
	"time"
)

//

func main() {

}

type Client struct {
	api ssov1.AuthClient
	log *slog.Logger
}

func New(
	ctx context.Context,
	log *slog.Logger,
	addr string,
	timeout time.Duration,
	retriesCount int,
) (*Client, error) {
	const op = "grpc.New"
	retriesOpts := []grpcretry.CallOption{
		grpcretry.WithCodes(codes.NotFound, codes.Aborted, codes.DeadlineExceeded),
		grpcretry.WithMax(uint(retriesCount)),
		grpcretry.WithPerRetryTimeout(timeout),
	}

	logOpts := []grpclog.Option{
		grpclog.WithLogOnEvents(grpclog.PayloadReceived, grpclog.PayloadSent),
	}

	//todo передать другие insecure в функцию ниже
	cc, err := grpc.DialContext(ctx, addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			grpclog.UnaryClientInterceptor(InterceptorLogger(log), logOpts...),
			grpcretry.UnaryClientInterceptor(retriesOpts...),
		))

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Client{
		api: ssov1.NewAuthClient(cc),
	}, nil
}

func (c *Client) IsAdmin(ctx context.Context, UserID int64) (bool, error) {
	op := "grpc.IsAdmin"
	resp, err := c.api.IsAdmin(ctx, &ssov1.IsAdminRequest{
		UserId: UserID,
	})

	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return resp.IsAdmin, nil
}

func InterceptorLogger(l *slog.Logger) grpclog.Logger {
	return grpclog.LoggerFunc(func(
		ctx context.Context,
		lvl grpclog.Level,
		msg string,
		fields ...any,
	) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}
