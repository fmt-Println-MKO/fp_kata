package services

import (
	"fp_kata/internal/datasources"
	"fp_kata/internal/datasources/file"
	"fp_kata/internal/models"
)

type OrderService interface {
	StoreOrder(order models.Order) error
	GetOrder(id int) (models.Order, error)
	FilterOrders(userId int, filter func(order models.Order) bool) ([]models.Order, error)
}

type orderService struct {
	Storage datasources.OrdersDatasource
}

func NewOrderService() OrderService {
	return &orderService{Storage: file.NewOrderStorage()}
}

func (service *orderService) StoreOrder(order models.Order) error {

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

func (service *orderService) GetOrder(id int) (models.Order, error) {
	order, err := service.Storage.GetOrder(id)
	if err != nil {
		return models.Order{}, err
	}
	return *order, nil
}

func (service *orderService) FilterOrders(userId int, filter func(order models.Order) bool) ([]models.Order, error) {
	allOrders, err := service.Storage.GetAllOrders(userId)
	if err != nil {
		return nil, err
	}

	var filteredOrders []models.Order
	for _, order := range allOrders {
		if filter(*order) {
			filteredOrders = append(filteredOrders, *order)
		}
	}

	return filteredOrders, nil
}
