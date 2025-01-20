package common

type PaymentMethod string

const (
	CreditCard   PaymentMethod = "CreditCard"
	DebitCard    PaymentMethod = "DebitCard"
	PayPal       PaymentMethod = "PayPal"
	BankTransfer PaymentMethod = "BankTransfer"
)
