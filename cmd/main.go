package main

import (
	"backend-test-app/internal/app"
	"backend-test-app/internal/config"
)

func main() {
	cfg := config.Load()
	app.Run(cfg)
}
