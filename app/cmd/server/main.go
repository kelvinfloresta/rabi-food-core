package main

import (
	"rabi-food-core/config"
	"rabi-food-core/libs/database/gorm_adapter"
	"rabi-food-core/libs/http/fiber_adapter"
	"time"

	"github.com/rs/zerolog/log"
)

func main() {
	log.Info().Msg("Starting Rabi Food Core Server...")
	log.Info().Str("env", config.Env).Msg("Environment")

	time.Local = time.UTC
	db := gorm_adapter.New(config.ProductionDatabase)

	if err := db.Start(); err != nil {
		log.Fatal().Err(err).Msg("Failed to start database")
	}

	httpServer := fiber_adapter.New(db)

	if err := httpServer.Start(config.Port); err != nil {
		log.Fatal().Err(err).Msg("Failed to start HTTP server")
	}
}
