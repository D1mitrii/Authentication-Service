package main

import (
	"auth/internal/config"
	"auth/pkg/logger"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	config := config.MustLoad()
	log := logger.SetupLogger(config.Env)
	log.Info("Success startup")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	sig := <-quit
	log.Info("Shutdown application", slog.String("signal", sig.String()))
}
