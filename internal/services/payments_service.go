package services

import (
	"context"
	"fp_kata/internal/datasources"
	"fp_kata/internal/datasources/dsmodels"
	"fp_kata/internal/models"
	"fp_kata/pkg/log"
)

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
	logger := log.GetLogger(ctx)
	logger.Debug().Str("comp", "PaymentsService").Str("func", "StorePayment").Send()
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
	logger := log.GetLogger(ctx)
	logger.Debug().Str("comp", "PaymentsService").Str("func", "GetPaymentByID").Send()
	dsPayment, err := service.storage.Read(ctx, id)
	if err != nil {
		return nil, err
	}
	payment := models.MapToPayment(dsPayment, &models.User{ID: dsPayment.UserId}, &models.Order{ID: dsPayment.OrderId})
	return payment, nil
}

func (service *paymentsService) GetPaymentsByOrder(ctx context.Context, orderId int) ([]*models.Payment, error) {
	logger := log.GetLogger(ctx)
	logger.Debug().Str("comp", "PaymentsService").Str("func", "GetPaymentsByOrder").Send()

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
