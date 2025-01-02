package transports

type UserResponse struct {
	ID       int               `json:"id"`
	Username string            `json:"username"`
	Email    string            `json:"email"`
	Password string            `json:"password"`
	Orders   []OrderResponse   `json:"orders"`
	Payments []PaymentResponse `json:"payments"`
}
