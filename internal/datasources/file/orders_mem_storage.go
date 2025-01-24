package file

import (
	"context"
	"errors"
	"fp_kata/common/utils"
	"fp_kata/internal/datasources"
	"fp_kata/internal/datasources/dsmodels"
)

const compOrdersStorage = "OrdersDatasource"

type inMemoryOrdersStorage struct {
	orders map[int]dsmodels.Order
}

func (s *inMemoryOrdersStorage) GetOrder(ctx context.Context, orderID int) (*dsmodels.Order, error) {
	utils.LogAction(ctx, compOrdersStorage, "GetOrder")

	order, exists := s.orders[orderID]
	if !exists {
		return nil, errors.New("order not found")
	}
	return &order, nil
}

func (s *inMemoryOrdersStorage) GetAllOrdersForUser(ctx context.Context, userID int) ([]dsmodels.Order, error) {
	utils.LogAction(ctx, compOrdersStorage, "GetAllOrdersForUser")

	userOrders := make([]dsmodels.Order, 0)
	for _, order := range s.orders {
		if order.UserId == userID {
			userOrders = append(userOrders, order)
		}
	}
	return userOrders, nil
}

func (s *inMemoryOrdersStorage) DeleteOrder(ctx context.Context, orderID int) error {
	utils.LogAction(ctx, compOrdersStorage, "DeleteOrder")

	if _, exists := s.orders[orderID]; !exists {
		return errors.New("order not found")
	}
	delete(s.orders, orderID)
	return nil
}

func (s *inMemoryOrdersStorage) UpdateOrder(ctx context.Context, order dsmodels.Order) (*dsmodels.Order, error) {
	utils.LogAction(ctx, compOrdersStorage, "UpdateOrder")

	_, exists := s.orders[order.ID]
	if !exists {
		return nil, errors.New("order not found")
	}
	s.orders[order.ID] = order
	return &order, nil
}

func (s *inMemoryOrdersStorage) InsertOrder(ctx context.Context, order dsmodels.Order) (*dsmodels.Order, error) {
	utils.LogAction(ctx, compOrdersStorage, "InsertOrder")

	if _, exists := s.orders[order.ID]; exists {
		return nil, errors.New("order already exists")
	}
	s.orders[order.ID] = order
	return &order, nil
}

func NewOrdersStorage() datasources.OrdersDatasource {
	return &inMemoryOrdersStorage{
		orders: make(map[int]dsmodels.Order),
	}
}
