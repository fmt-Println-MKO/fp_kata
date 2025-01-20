package models

import (
	"fp_kata/common"
	"fp_kata/internal/datasources/dsmodels"
)

type Payment struct {
	Id     int
	Amount float64
	Method common.PaymentMethod
	User   *User
	Order  *Order
}

func (p Payment) ToDSModel() *dsmodels.Payment {
	if p.User == nil || p.Order == nil {
		return nil
	}
	return &dsmodels.Payment{
		Id:      p.Id,
		Amount:  p.Amount,
		Method:  p.Method,
		UserId:  p.User.ID,
		OrderId: p.Order.ID,
	}
}
func MapToPayment(dsPayment dsmodels.Payment, user *User, order *Order) *Payment {
	if user == nil || order == nil {
		return nil
	}
	return &Payment{
		Id:     dsPayment.Id,
		Amount: dsPayment.Amount,
		Method: dsPayment.Method,
		User:   user,
		Order:  order,
	}
}
