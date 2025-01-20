package dsmodels

import "fp_kata/common"

type Payment struct {
	Id      int
	Amount  float64
	Method  common.PaymentMethod
	UserId  int
	OrderId int
}
