package file

import (
	"errors"
	"fp_kata/internal/model"
)

type OrderStorage interface {
	GetOrder(orderID int) (*model.Order, error)
	GetAllOrders(userID int) ([]*model.Order, error)
	DeleteOrder(orderID int) error
	UpdateOrder(order *model.Order) error
	InsertOrder(order *model.Order) error
}

type inMemoryOrderStorage struct {
	orders map[int]*model.Order
}

func (s *inMemoryOrderStorage) GetOrder(orderID int) (*model.Order, error) {
	order, exists := s.orders[orderID]
	if !exists {
		return nil, errors.New("order not found")
	}
	return order, nil
}

func (s *inMemoryOrderStorage) GetAllOrders(userID int) ([]*model.Order, error) {
	userOrders := make([]*model.Order, 0)
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

func (s *inMemoryOrderStorage) UpdateOrder(order *model.Order) error {
	_, exists := s.orders[order.ID]
	if !exists {
		return errors.New("order not found")
	}
	s.orders[order.ID] = order
	return nil
}

func (s *inMemoryOrderStorage) InsertOrder(order *model.Order) error {
	if _, exists := s.orders[order.ID]; exists {
		return errors.New("order already exists")
	}
	s.orders[order.ID] = order
	return nil
}

func NewOrderStorage() OrderStorage {
	return &inMemoryOrderStorage{
		orders: make(map[int]*model.Order),
	}
}
