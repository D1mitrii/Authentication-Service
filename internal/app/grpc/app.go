package grpc

import (
	"fmt"
	"log/slog"
	"net"

	grpcv1 "auth/internal/controller/grpc/v1"
	desc "auth/pkg/auth/v1"

	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
	notify     chan error
}

func New(log *slog.Logger, port int, authService *grpcv1.Auth) *App {
	s := grpc.NewServer()
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
