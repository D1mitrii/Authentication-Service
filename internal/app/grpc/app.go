package grpc

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"

	"auth/internal/controller/grpc/interceptors"
	grpcv1 "auth/internal/controller/grpc/v1"
	desc "auth/pkg/auth/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
	notify     chan error
}

func New(log *slog.Logger, port int, authService *grpcv1.Auth) *App {
	logOpts := []logging.Option{
		logging.WithLogOnEvents(
			logging.PayloadReceived,
			logging.PayloadSent,
		),
	}
	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandler(func(p any) (err error) {
			log.Error("gRPC server recovered from panic", slog.Any("panic", p))
			return status.Error(codes.Internal, "internal server error")
		}),
	}

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			recovery.UnaryServerInterceptor(recoveryOpts...),
			logging.UnaryServerInterceptor(InterceptorLogger(log), logOpts...),
			interceptors.MetricsInterceptor,
		),
	)
	desc.RegisterAuthV1Server(s, authService)
	app := &App{
		log:        log,
		gRPCServer: s,
		port:       port,
		notify:     make(chan error),
	}
	go app.Run()
	return app
}

// InterceptorLogger adapts slog logger to interceptor logger.
// This code is simple enough to be copied and not imported.
func InterceptorLogger(l *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}

func (a *App) Run() {
	defer close(a.notify)
	const op = "grpc.App.Run"

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		a.notify <- fmt.Errorf("%s: %w", op, err)
		return
	}

	a.log.Info("gRPC server started")
	a.notify <- a.gRPCServer.Serve(l)
}

func (a *App) Notify() chan error {
	return a.notify
}

func (a *App) Stop() {
	const op = "grpc.App.Stop"
	a.log.Info(fmt.Sprintf("%s - stopping gRPC server", op))
	a.gRPCServer.GracefulStop()
}
