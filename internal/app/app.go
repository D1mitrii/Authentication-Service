package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/d1mitrii/authentication-service/internal/config"
	grpcv1 "github.com/d1mitrii/authentication-service/internal/controller/grpc/v1"
	"github.com/d1mitrii/authentication-service/internal/controller/http/middlewares"
	httpv1 "github.com/d1mitrii/authentication-service/internal/controller/http/v1"
	"github.com/d1mitrii/authentication-service/internal/metrics"
	"github.com/d1mitrii/authentication-service/internal/repository"
	"github.com/d1mitrii/authentication-service/internal/repository/pgdb"
	"github.com/d1mitrii/authentication-service/internal/repository/rdb"
	"github.com/d1mitrii/authentication-service/internal/services"
	"github.com/d1mitrii/authentication-service/internal/services/jwt"
	"github.com/d1mitrii/authentication-service/pkg/hasher"
	"github.com/d1mitrii/authentication-service/pkg/httpserver"
	"github.com/d1mitrii/authentication-service/pkg/logger"
	"github.com/d1mitrii/authentication-service/pkg/postgres"

	"github.com/d1mitrii/authentication-service/internal/app/grpc"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
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
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.RDB.Host, cfg.RDB.Port),
		Password: cfg.RDB.Password,
	})
	if err := client.Ping(context.Background()).Err(); err != nil {
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

	log.Info("Initializing HTTP server for metrics")
	m := http.NewServeMux()
	reg := prometheus.NewRegistry()
	metrics.Init(reg)
	promHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
	m.Handle("/metrics", promHandler)
	metricsServer := httpserver.New(
		m,
		httpserver.Port(cfg.Prometheus.Port),
	)

	log.Info("Initializing HTTP server")
	log.Info("Initializing handlers & routes")

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middlewares.MetricsMiddleware)
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
		log.Error("Application HTTP server: ", slog.Any("err", err.Error()))
	case err := <-grpcServer.Notify():
		log.Error("Application gRPC server: ", slog.Any("err", err.Error()))
	case err := <-metricsServer.Notify():
		log.Error("Application Metrics server: ", slog.Any("err", err.Error()))
	}

	log.Info("Shutting down...")
	wg := sync.WaitGroup{}
	wg.Add(3)

	go func() {
		defer wg.Done()
		if err := httpServer.Shutdown(); err != nil {
			log.Error("Application HTTP shutdown: ", slog.Any("err", err.Error()))
		}
	}()

	go func() {
		defer wg.Done()
		if err := metricsServer.Shutdown(); err != nil {
			log.Error("Application Metrics server: ", slog.Any("err", err.Error()))
		}
	}()

	go func() {
		defer wg.Done()
		grpcServer.Stop()
	}()

	wg.Wait()
}
