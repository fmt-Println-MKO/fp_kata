package transports

import (
	"time"
)

type OrderResponse struct {
	ID             int
	ProductID      int
	Quantity       int
	Price          float64
	OrderDate      time.Time
	Payments       []PaymentResponse
	User           UserResponse
	hasWeightables bool
}
