package services

import (
	"errors"
	"fp_kata/internal/model"
	"fp_kata/internal/storage/file"
)

type UserService interface {
	GetUserByID(id int) (*model.User, error)
}

type userService struct {
	storage file.UserStorage
}

func NewUserService() UserService {
	return &userService{
		storage: file.NewUserStorage(),
	}
}

func (us *userService) GetUserByID(id int) (*model.User, error) {
	user, exists := us.storage.Read(id)
	if !exists {
		return nil, errors.New("no user found for id")
	}
	return &user, nil
}
