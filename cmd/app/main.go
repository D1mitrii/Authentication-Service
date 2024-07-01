package main

import (
	"github.com/d1mitrii/authentication-service/internal/app"
	"github.com/d1mitrii/authentication-service/internal/config"
)

func main() {
	config := config.MustLoad()
	app.Run(config)
}
