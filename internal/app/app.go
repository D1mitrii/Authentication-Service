package app

import (
	"auth/internal/config"
	"auth/pkg/httpserver"
	"auth/pkg/logger"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Run(cfg *config.Config) {
	log := logger.SetupLogger(cfg.Env)

	log.Info("Initializing application")
	log.Info("Initializing handlers & routes...")
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Alive!"))
	})
	log.Info("Starting http server...")
	s := httpserver.New(
		r,
		httpserver.Port(cfg.HTTP.Port),
		httpserver.ReadTimeout(cfg.HTTP.Timeout),
		httpserver.WriteTimeout(cfg.HTTP.Timeout),
	)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	select {
	case sig := <-quit:
		log.Info("Application catch:", slog.String("signal", sig.String()))
	case err := <-s.Notify():
		log.Error("Application", err)
	}

	log.Info("Shutting down...")
	if err := s.Shutdown(); err != nil {
		log.Error("Application shutdown:", err)
	}
}
