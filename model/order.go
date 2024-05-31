package model

import (
	"time"
)

type Order struct {
	ID             int
	ProductID      int
	Quantity       int
	Price          float64
	OrderDate      time.Time
	Payments       []Payment
	User           User
	hasWeightables bool
}
