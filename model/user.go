package model

type User struct {
	ID       int       `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
	Orders   []Order   `json:"orders"`
	Payments []Payment `json:"payments"`
}
