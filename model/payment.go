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
	User          User
	Order         Order
}

func NewPayment(id int, amount float64, method PaymentMethod, u User, o Order) *Payment {
	return &Payment{
		PaymentID:     id,
		PaymentAmount: amount,
		PaymentMethod: method,
		User:          u,
		Order:         o,
	}
}
