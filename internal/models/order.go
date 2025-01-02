package models

import (
	"time"
)

type Order struct {
	ID             int
	ProductID      int
	Quantity       int
	Price          float64
	OrderDate      time.Time
	Payments       []int
	UserId         int
	hasWeightables bool
}
