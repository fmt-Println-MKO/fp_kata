package services

import (
	"fp_kata/internal/model"
	"fp_kata/internal/storage/yugabyte"
)

type PaymentService interface {
	GetPaymentsByOrder(order int) ([]*model.Payment, error)
	GetPaymentByID(id int) (*model.Payment, error)
}

type paymentService struct {
	storage yugabyte.PaymentStorage
}

func NewPaymentService() PaymentService {
	return &paymentService{storage: yugabyte.NewPaymentStorage()}
}

func (service *paymentService) GetPaymentByID(id int) (*model.Payment, error) {
	payment, err := service.storage.Read(id)
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (service *paymentService) GetPaymentsByOrder(order int) ([]*model.Payment, error) {

	payments, err := service.storage.AllByOrderId(order)
	if err != nil {
		return nil, err
	}
	var result []*model.Payment
	for _, p := range payments {
		result = append(result, &p)
	}
	return result, nil
}
