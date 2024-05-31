package services

import (
	"fp_kata/internal/model"
	"fp_kata/internal/storage/file"
)

type OrderService interface {
	StoreOrder(order model.Order) error
	GetOrder(id int) (model.Order, error)
	FilterOrders(userId int, filter func(order model.Order) bool) ([]model.Order, error)
}

type orderService struct {
	Storage file.OrderStorage
}

func NewOrderService() OrderService {
	return &orderService{Storage: file.NewOrderStorage()}
}

func (service *orderService) StoreOrder(order model.Order) error {

	if order.ID == 0 {
		err := service.Storage.InsertOrder(&order)
		if err != nil {
			return err
		}
	} else {
		err := service.Storage.UpdateOrder(&order)
		if err != nil {
			return err
		}
	}
	return nil
}

func (service *orderService) GetOrder(id int) (model.Order, error) {
	order, err := service.Storage.GetOrder(id)
	if err != nil {
		return model.Order{}, err
	}
	return *order, nil
}

func (service *orderService) FilterOrders(userId int, filter func(order model.Order) bool) ([]model.Order, error) {
	allOrders, err := service.Storage.GetAllOrders(userId)
	if err != nil {
		return nil, err
	}

	var filteredOrders []model.Order
	for _, order := range allOrders {
		if filter(*order) {
			filteredOrders = append(filteredOrders, *order)
		}
	}

	return filteredOrders, nil
}
