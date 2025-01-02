package services

import (
	"fp_kata/internal/datasources"
	"fp_kata/internal/datasources/yugabyte"
	"fp_kata/internal/models"
)

type PaymentService interface {
	GetPaymentsByOrder(order int) ([]*models.Payment, error)
	GetPaymentByID(id int) (*models.Payment, error)
}

type paymentService struct {
	storage datasources.PaymentsDatasource
}

func NewPaymentService() PaymentService {
	return &paymentService{storage: yugabyte.NewPaymentStorage()}
}

func (service *paymentService) GetPaymentByID(id int) (*models.Payment, error) {
	payment, err := service.storage.Read(id)
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (service *paymentService) GetPaymentsByOrder(order int) ([]*models.Payment, error) {

	payments, err := service.storage.AllByOrderId(order)
	if err != nil {
		return nil, err
	}
	var result []*models.Payment
	for _, p := range payments {
		result = append(result, &p)
	}
	return result, nil
}
