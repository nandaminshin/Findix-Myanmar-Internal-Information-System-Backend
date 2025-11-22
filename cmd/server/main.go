package main

import (
	"fmiis/internal/app"
	"fmiis/internal/config"
	"log"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config :", err)
	}

	application, err := app.NewApp(cfg)
	if err != nil {
		log.Fatal("Failed to initialize app:", err)
	}

	// Start server
	application.StartServer()
}
