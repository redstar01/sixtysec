package main

import (
	"log"

	"github.com/redstar01/sixtysec/config"
	"github.com/redstar01/sixtysec/internal/app"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	app.Run(cfg)
}
