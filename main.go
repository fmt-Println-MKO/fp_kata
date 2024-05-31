package main

import (
	"fp_kata/api/orders"
	"fp_kata/common/log"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	logger := log.GetLogger()

	logger.Info("Application started")
	logger.Info("Hello functional go")

	r := mux.NewRouter()

	ordersController := orders.NewOrderController()

	ordersController.RegisterRoutes(r)

	http.ListenAndServe(":8000", r)
}
