package file

import (
	"errors"
	"fp_kata/internal/datasources"
	"fp_kata/internal/models"
)

type inMemoryOrderStorage struct {
	orders map[int]*models.Order
}

func (s *inMemoryOrderStorage) GetOrder(orderID int) (*models.Order, error) {
	order, exists := s.orders[orderID]
	if !exists {
		return nil, errors.New("order not found")
	}
	return order, nil
}

func (s *inMemoryOrderStorage) GetAllOrders(userID int) ([]*models.Order, error) {
	userOrders := make([]*models.Order, 0)
	for _, order := range s.orders {
		if order.UserId == userID {
			userOrders = append(userOrders, order)
		}
	}
	return userOrders, nil
}

func (s *inMemoryOrderStorage) DeleteOrder(orderID int) error {
	if _, exists := s.orders[orderID]; !exists {
		return errors.New("order not found")
	}
	delete(s.orders, orderID)
	return nil
}

func (s *inMemoryOrderStorage) UpdateOrder(order *models.Order) error {
	_, exists := s.orders[order.ID]
	if !exists {
		return errors.New("order not found")
	}
	s.orders[order.ID] = order
	return nil
}

func (s *inMemoryOrderStorage) InsertOrder(order *models.Order) error {
	if _, exists := s.orders[order.ID]; exists {
		return errors.New("order already exists")
	}
	s.orders[order.ID] = order
	return nil
}

func NewOrderStorage() datasources.OrdersDatasource {
	return &inMemoryOrderStorage{
		orders: make(map[int]*models.Order),
	}
}
