package controllers

import (
	"encoding/json"
	"fp_kata/internal/models"
	"fp_kata/internal/services"
	"fp_kata/pkg/transports"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
)

type OrderController struct {
	orderService   services.OrderService
	userService    services.UserService
	paymentService services.PaymentService
}

func NewOrderController() *OrderController {
	return &OrderController{
		orderService:   services.NewOrderService(),
		userService:    services.NewUserService(),
		paymentService: services.NewPaymentService(),
	}
}

func (c *OrderController) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/orders", c.GetOrders).Methods(http.MethodGet)
	router.HandleFunc("/orders/{id}", c.GetOrder).Methods(http.MethodGet)
	// You can add other handler functions based on the OrderService methods
}

// GetOrders handles "/orders" with method "GET"
func (c *OrderController) GetOrders(w http.ResponseWriter, r *http.Request) {

	log.Info().Msg("GetOrders")

	queryParams := r.URL.Query()
	userId := queryParams.Get("userId")

	log.Info().Str("func", "GetOrders").Str("userId", userId)
	log.Info().Msgf("GetOrders -  userId: %s", userId)

	uid, err := strconv.Atoi(userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	filter := func(order models.Order) bool {
		return order.Price > 20
	}

	orders, err := c.orderService.FilterOrders(uid, filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result := make([]transports.OrderResponse, len(orders))
	for i, order := range orders {

		payments := make([]transports.PaymentResponse, len(order.Payments))
		for i, id := range order.Payments {
			payment, err := c.paymentService.GetPaymentByID(id)
			if err != nil {
				payments[i] = transports.PaymentResponse{}
			} else {
				payments[i] = mapPayments(*payment)
			}
		}
		aUser, err := c.userService.GetUserByID(order.UserId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		user := mapUser(*aUser)

		result[i] = mapToModelOrder(order, payments, user)

	}

	resp, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

// GetOrder handles "/orders/{id}" with method "GET"
func (c *OrderController) GetOrder(w http.ResponseWriter, r *http.Request) {

	log.Info().Msg("GetOrder")

	queryParams := r.URL.Query()
	orderId := queryParams.Get("orderId")

	oid, err := strconv.Atoi(orderId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	order, err := c.orderService.GetOrder(oid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	payments := make([]transports.PaymentResponse, len(order.Payments))
	for i, id := range order.Payments {
		payment, err := c.paymentService.GetPaymentByID(id)
		if err != nil {
			payments[i] = transports.PaymentResponse{}
		} else {
			payments[i] = mapPayments(*payment)
		}
	}
	aUser, err := c.userService.GetUserByID(order.UserId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := mapUser(*aUser)

	result := mapToModelOrder(order, payments, user)

	resp, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

func mapToModelOrder(order models.Order, payments []transports.PaymentResponse, user transports.UserResponse) transports.OrderResponse {
	return transports.OrderResponse{
		ID:        order.ID,
		ProductID: order.ProductID,
		Quantity:  order.Quantity,
		Price:     order.Price,
		OrderDate: order.OrderDate,
		Payments:  payments,
		User:      user,
	}
}

func mapUser(user models.User) transports.UserResponse {
	return transports.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}
}

func mapPayments(payment models.Payment) transports.PaymentResponse {
	return transports.PaymentResponse{
		PaymentID:     payment.PaymentID,
		PaymentAmount: payment.PaymentAmount,
		PaymentMethod: transports.PaymentMethod(payment.PaymentMethod),
	}
}
