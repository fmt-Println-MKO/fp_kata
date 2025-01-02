package main

import (
	"fp_kata/internal/controllers"
	fpLog "fp_kata/pkg/log"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"net/http"
)

func main() {
	fpLog.InitLogger()

	log.Info().Msg("Application started")
	log.Info().Msg("Hello functional go")

	r := mux.NewRouter()

	ordersController := controllers.NewOrderController()

	ordersController.RegisterRoutes(r)

	http.ListenAndServe(":8000", r)
}
