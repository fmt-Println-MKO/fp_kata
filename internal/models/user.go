package models

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Orders   []int  `json:"orders"`
	Payments []int  `json:"payments"`
}
