package app

import (
	"auth/internal/config"
	grpcv1 "auth/internal/controller/grpc/v1"
	httpv1 "auth/internal/controller/http/v1"
	"auth/internal/repository"
	"auth/internal/repository/pgdb"
	"auth/internal/repository/rdb"
	"auth/internal/services"
	"auth/internal/services/jwt"
	"auth/pkg/hasher"
	"auth/pkg/httpserver"
	"auth/pkg/logger"
	"auth/pkg/postgres"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"auth/internal/app/grpc"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-redis/redis"
)

func Run(cfg *config.Config) {
	const op = "app - Run"
	log := logger.SetupLogger(cfg.Env)

	log.Info("Connecting to postgres")
	pg, err := postgres.New(cfg.PG.URL)
	if err != nil {
		log.Error(fmt.Sprintf("%s - postgres.New: %v", op, err))
		return
	}
	defer pg.Close()

	log.Info("Connecting to redis")
	log.Info("redis password: " + cfg.RDB.Password)
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.RDB.Host, cfg.RDB.Port),
		Password: cfg.RDB.Password,
	})
	if err := client.Ping().Err(); err != nil {
		log.Error(fmt.Sprintf("%s - redis.NewClient: %v", op, err))
		return
	}
	defer client.Close()

	log.Info("Initializing services")
	service := services.New(
		log,
		jwt.New(
			cfg.JWT.Secret,
			cfg.JWT.TokenTTL,
			cfg.JWT.RefreshTime,
		),
		hasher.New(cfg.Hasher.Salt),
		repository.New(
			pgdb.NewUserRepo(pg),
			rdb.NewRefreshRepo(client, cfg.JWT.RefreshTime),
		),
	)

	log.Info("Initializing HTTP server")
	log.Info("Initializing handlers & routes")

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Mount("/api/v1", httpv1.New(service).Routes())

	log.Info("Starting http server...")
	httpServer := httpserver.New(
		r,
		httpserver.Port(cfg.HTTP.Port),
		httpserver.ReadTimeout(cfg.HTTP.Timeout),
		httpserver.WriteTimeout(cfg.HTTP.Timeout),
	)

	log.Info("Initializing gRPC server")
	grpcServer := grpc.New(log, cfg.GRPC.Port, grpcv1.NewAuth(service))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	select {
	case sig := <-quit:
		log.Info("Application catch:", slog.String("signal", sig.String()))
	case err := <-httpServer.Notify():
		log.Error("Application HTTP server: ", err)
	case err := <-grpcServer.Notify():
		log.Error("Application gRPC server: ", err)
	}

	log.Info("Shutting down...")
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		if err := httpServer.Shutdown(); err != nil {
			log.Error("Application HTTP shutdown: ", err)
		}
	}()

	go func() {
		defer wg.Done()
		grpcServer.Stop()
	}()

	wg.Wait()
}
