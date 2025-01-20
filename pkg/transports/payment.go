package transports

import (
	"fp_kata/common"
	"fp_kata/internal/models"
)

type PaymentResponse struct {
	Id     int                  `json:"id"`
	Amount float64              `json:"amount"`
	Method common.PaymentMethod `json:"method"`
}

func MapToPaymentResponse(payment models.Payment) *PaymentResponse {
	return &PaymentResponse{
		Id:     payment.Id,
		Amount: payment.Amount,
		Method: payment.Method,
	}
}

type PaymentRequest struct {
	PaymentAmount float64              `json:"payment_amount"`
	PaymentMethod common.PaymentMethod `json:"payment_method"`
}

func (p PaymentRequest) ToPayment(user models.User) *models.Payment {
	return &models.Payment{
		Amount: p.PaymentAmount,
		Method: p.PaymentMethod,
		User:   &user,
	}
}
