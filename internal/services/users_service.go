package services

import (
	"context"
	"errors"
	"fp_kata/internal/datasources"
	"fp_kata/internal/models"
	"fp_kata/pkg/log"
)

type UsersService interface {
	GetUserByID(ctx context.Context, id int) (*models.User, error)
	SignUp(ctx context.Context, user models.User) (*models.User, error)
}

type usersService struct {
	storage     datasources.UsersDatasource
	authService AuthService
}

func NewUsersService(storage datasources.UsersDatasource, authService AuthService) UsersService {
	return &usersService{
		storage:     storage,
		authService: authService,
	}
}

func (us *usersService) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	logger := log.GetLogger(ctx)
	logger.Debug().Str("comp", "UsersService").Str("func", "GetUserByID").Send()
	dsUser, exists := us.storage.Read(ctx, id)
	if !exists {
		return nil, errors.New("no user found for id")
	}
	user := models.MapToUser(dsUser)
	return user, nil
}

func (us *usersService) SignUp(ctx context.Context, user models.User) (*models.User, error) {
	logger := log.GetLogger(ctx)
	logger.Debug().Str("comp", "UsersService").Str("func", "SignUp").Send()
	dsUser := user.ToDSModel()
	createdDsUser, created := us.storage.Create(ctx, *dsUser)
	if !created {
		return nil, errors.New("user storage is full")
	}
	createdUser := models.MapToUser(createdDsUser)
	_, err := us.authService.GenerateAuthToken(ctx, *createdUser)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}
