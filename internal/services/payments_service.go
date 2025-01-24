package services

import (
	"context"
	"fp_kata/common/utils"
	"fp_kata/internal/datasources"
	"fp_kata/internal/datasources/dsmodels"
	"fp_kata/internal/models"
)

const compPaymentsService = "PaymentsService"

type PaymentsService interface {
	StorePayment(ctx context.Context, payment models.Payment) (*models.Payment, error)
	GetPaymentsByOrder(ctx context.Context, orderId int) ([]*models.Payment, error)
	GetPaymentByID(ctx context.Context, id int) (*models.Payment, error)
}

type paymentsService struct {
	storage datasources.PaymentsDatasource
}

func NewPaymentsService(storage datasources.PaymentsDatasource) PaymentsService {
	return &paymentsService{storage: storage}
}

func (service *paymentsService) StorePayment(ctx context.Context, payment models.Payment) (*models.Payment, error) {
	utils.LogAction(ctx, compPaymentsService, "StorePayment")

	createdDsPayment := dsmodels.Payment{}
	var err error
	if payment.Id == 0 {
		createdDsPayment, err = service.storage.Create(ctx, *payment.ToDSModel())
	} else {
		createdDsPayment, err = service.storage.Update(ctx, *payment.ToDSModel())
	}
	if err != nil {
		return nil, err
	}
	createdPayment := models.MapToPayment(createdDsPayment, payment.User, payment.Order)
	return createdPayment, nil
}

func (service *paymentsService) GetPaymentByID(ctx context.Context, id int) (*models.Payment, error) {
	utils.LogAction(ctx, compPaymentsService, "GetPaymentByID")

	dsPayment, err := service.storage.Read(ctx, id)
	if err != nil {
		return nil, err
	}
	payment := models.MapToPayment(dsPayment, &models.User{ID: dsPayment.UserId}, &models.Order{ID: dsPayment.OrderId})
	return payment, nil
}

func (service *paymentsService) GetPaymentsByOrder(ctx context.Context, orderId int) ([]*models.Payment, error) {
	utils.LogAction(ctx, compPaymentsService, "GetPaymentsByOrder")

	dsPayments, err := service.storage.AllByOrderId(ctx, orderId)
	if err != nil {
		return nil, err
	}
	payments := make([]*models.Payment, len(dsPayments))
	for i, dsPayment := range dsPayments {
		payments[i] = models.MapToPayment(dsPayment, &models.User{ID: dsPayment.UserId}, &models.Order{ID: dsPayment.OrderId})
	}
	return payments, nil
}
