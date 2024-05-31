package main

import (
	"fp_kata/api/orders"
	fpLog "fp_kata/common/log"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"net/http"
)

func main() {
	fpLog.InitLogger()

	log.Info().Msg("Application started")
	log.Info().Msg("Hello functional go")

	r := mux.NewRouter()

	ordersController := orders.NewOrderController()

	ordersController.RegisterRoutes(r)

	http.ListenAndServe(":8000", r)
}
