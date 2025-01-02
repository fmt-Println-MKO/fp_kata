package services

import (
	"errors"
	"fp_kata/internal/datasources"
	"fp_kata/internal/datasources/file"
	"fp_kata/internal/models"
)

type UserService interface {
	GetUserByID(id int) (*models.User, error)
}

type userService struct {
	storage datasources.UsersDatasource
}

func NewUserService() UserService {
	return &userService{
		storage: file.NewUserStorage(),
	}
}

func (us *userService) GetUserByID(id int) (*models.User, error) {
	user, exists := us.storage.Read(id)
	if !exists {
		return nil, errors.New("no user found for id")
	}
	return &user, nil
}
