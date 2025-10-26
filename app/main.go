package main

import (
	"rabi-food-core/config"
	"rabi-food-core/libs/database"
	"rabi-food-core/libs/di"
	"rabi-food-core/libs/http"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/samber/do"
)

func main() {
	log.Info().Msg("Starting Rabi Food Core Server...")
	log.Info().Str("env", config.Env).Msg("Environment")

	time.Local = time.UTC

	injector := di.NewProduction()
	db := do.MustInvoke[database.Database](injector)

	err := db.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start database")
	}

	httpServer := do.MustInvoke[http.HTTPServer](injector)

	err = httpServer.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start HTTP server")
	}
}
