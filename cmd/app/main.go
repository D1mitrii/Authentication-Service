package main

import (
	"auth/internal/app"
	"auth/internal/config"
)

func main() {
	config := config.MustLoad()
	app.Run(config)
}
