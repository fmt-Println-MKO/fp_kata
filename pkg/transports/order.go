package transports

import (
	"fp_kata/internal/models"
	"time"
)

type OrderResponse struct {
	ID             int                `json:"id,omitempty"`
	ProductID      int                `json:"product_id,omitempty"`
	Quantity       int                `json:"quantity,omitempty"`
	Price          float64            `json:"price,omitempty"`
	OrderDate      time.Time          `json:"order_date,omitempty"`
	Payments       []*PaymentResponse `json:"payments,omitempty"`
	User           *UserResponse      `json:"user,omitempty"`
	HasWeightables bool               `json:"has_weightables"`
}

type OrderCreateRequest struct {
	ProductID      int               `json:"product_id,omitempty" validate:"required" binding:"required"`
	Quantity       int               `json:"quantity,omitempty" validate:"required" binding:"required"`
	Price          float64           `json:"price,omitempty" validate:"required" binding:"required"`
	OrderDate      time.Time         `json:"order_date,omitempty" validate:"required" binding:"required"`
	Payments       []*PaymentRequest `json:"payments,omitempty" validate:"required" binding:"required"`
	HasWeightables bool              `json:"has_weightables,omitempty" binding:"required"`
}

func (orderRequest *OrderCreateRequest) ToOrder(user models.User) *models.Order {

	payments := make([]*models.Payment, len(orderRequest.Payments))
	for i, paymentReq := range orderRequest.Payments {
		payment := paymentReq.ToPayment(user)
		payments[i] = payment
	}

	return &models.Order{
		ProductID:      orderRequest.ProductID,
		Quantity:       orderRequest.Quantity,
		Price:          orderRequest.Price,
		OrderDate:      orderRequest.OrderDate,
		Payments:       payments,
		User:           &user,
		HasWeightables: orderRequest.HasWeightables,
	}
}

// MapToOrderResponse creates an OrderResponse from a models.Order.
func MapToOrderResponse(order models.Order) *OrderResponse {

	//user := UserResponse{}
	var user *UserResponse
	if order.User != nil {
		user = MapToUserResponse(*order.User)
	}

	return &OrderResponse{
		ID:             order.ID,
		ProductID:      order.ProductID,
		Quantity:       order.Quantity,
		Price:          order.Price,
		OrderDate:      order.OrderDate,
		Payments:       convertPayments(order.Payments),
		User:           user,
		HasWeightables: order.HasWeightables,
	}
}

// Helper function to convert a slice of models.Payment to []*PaymentResponse.
func convertPayments(payments []*models.Payment) []*PaymentResponse {
	paymentResponses := make([]*PaymentResponse, len(payments))
	for i, payment := range payments {
		paymentResponses[i] = MapToPaymentResponse(*payment)
	}
	return paymentResponses
}
