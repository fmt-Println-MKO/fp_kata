package model

type PaymentMethod string

const (
	CreditCard   PaymentMethod = "CreditCard"
	DebitCard    PaymentMethod = "DebitCard"
	PayPal       PaymentMethod = "PayPal"
	BankTransfer PaymentMethod = "BankTransfer"
)

type Payment struct {
	PaymentID     int
	PaymentAmount float64
	PaymentMethod PaymentMethod
	UserId        int
	OrderId       int
}

func NewPayment(id int, amount float64, method PaymentMethod, u int, o int) *Payment {
	return &Payment{
		PaymentID:     id,
		PaymentAmount: amount,
		PaymentMethod: method,
		UserId:        u,
		OrderId:       o,
	}
}
