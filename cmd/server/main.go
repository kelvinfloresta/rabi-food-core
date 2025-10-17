package main

import (
	"log"
	"rabi-food-core/config"
	"rabi-food-core/libs/database/gorm_adapter"
	"rabi-food-core/libs/http/fiber_adapter"
	"time"
)

func main() {
	time.Local = time.UTC

	db := gorm_adapter.New(config.ProductionDatabase)

	if err := db.Start(); err != nil {
		panic(err)
	}

	httpServer := fiber_adapter.New(db)

	log.Fatal(httpServer.Start(config.Port))
}
