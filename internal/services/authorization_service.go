package services

import (
	"context"
	"errors"
	"fp_kata/common/utils"
	"fp_kata/internal/models"
)

const compAuthorizationService = "AuthorizationService"

type AuthorizationService interface {
	IsAuthorized(ctx context.Context, userId int, order models.Order) (bool, error)
}

type authorizationService struct{}

func NewAuthorizationService() AuthorizationService {
	return &authorizationService{}
}

func (a *authorizationService) IsAuthorized(ctx context.Context, userId int, order models.Order) (bool, error) {
	utils.LogAction(ctx, compAuthorizationService, "isAuthorized")

	if userId == 0 {
		return false, errors.New("userId is required")
	}
	if order.User == nil {
		return false, errors.New("missing user on order")
	}
	return order.User.ID == userId, nil
}
