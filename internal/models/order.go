package models

import (
	"fp_kata/internal/datasources/dsmodels"
	"time"
)

type Order struct {
	ID             int
	ProductID      int
	Quantity       int
	Price          float64
	OrderDate      time.Time
	Payments       []*Payment
	User           *User
	HasWeightables bool
}

// ToDSModel converts the Order struct to the dsmodels.Order struct
func (o *Order) ToDSModel() *dsmodels.Order {

	if o.User == nil {
		return nil
	}

	dsPayments := make([]int, len(o.Payments))
	for i, payment := range o.Payments {
		dsPayments[i] = payment.Id
	}

	return &dsmodels.Order{
		ID:             o.ID,
		ProductID:      o.ProductID,
		Quantity:       o.Quantity,
		Price:          o.Price,
		OrderDate:      o.OrderDate,
		Payments:       dsPayments,
		UserId:         o.User.ID,
		HasWeightables: o.HasWeightables,
	}
}

// MapToOrder maps the fields from a dsmodels.Order struct to the Order struct
func MapToOrder(dso dsmodels.Order, payments []*Payment, user *User) *Order {
	return &Order{
		ID:             dso.ID,
		ProductID:      dso.ProductID,
		Quantity:       dso.Quantity,
		Price:          dso.Price,
		OrderDate:      dso.OrderDate,
		Payments:       payments,
		User:           user,
		HasWeightables: dso.HasWeightables,
	}

}
