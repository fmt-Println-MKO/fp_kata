package transports

type PaymentMethod string

const (
	CreditCard   PaymentMethod = "CreditCard"
	DebitCard    PaymentMethod = "DebitCard"
	PayPal       PaymentMethod = "PayPal"
	BankTransfer PaymentMethod = "BankTransfer"
)

type PaymentResponse struct {
	PaymentID     int
	PaymentAmount float64
	PaymentMethod PaymentMethod
	User          UserResponse
	Order         OrderResponse
}

func NewPaymentResponse(id int, amount float64, method PaymentMethod, u UserResponse, o OrderResponse) *PaymentResponse {
	return &PaymentResponse{
		PaymentID:     id,
		PaymentAmount: amount,
		PaymentMethod: method,
		User:          u,
		Order:         o,
	}
}
