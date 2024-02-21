package main

import (
	"auth/internal/config"
	"auth/pkg/logger"
)

func main() {
	config := config.MustLoad()
	log := logger.SetupLogger(config.Env)
	log.Info("Success startup")
}
