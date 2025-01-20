package models

import "fp_kata/internal/datasources/dsmodels"

type User struct {
	ID       int
	Username string
	Email    string
	Password string
	Orders   []*Order
}

func (u User) ToDSModel() *dsmodels.User {
	return &dsmodels.User{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
		Password: u.Password,
	}
}
func MapToUser(dsUser dsmodels.User) *User {

	return &User{
		ID:       dsUser.ID,
		Username: dsUser.Username,
		Email:    dsUser.Email,
		Password: dsUser.Password,
	}
}
