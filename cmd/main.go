package main

import (
	"fp_kata/internal/app"
	"github.com/rs/zerolog/log"
)

func main() {

	app := app.InitApp()

	if err := app.Listen(":8000"); err != nil {
		log.Fatal().Err(err).Msg("Failed to start the application")
	}
}
