package file

import (
	"context"
	"errors"
	"fp_kata/common/monads"
	"fp_kata/common/utils"
	"fp_kata/internal/datasources"
	"fp_kata/internal/datasources/dsmodels"
)

const compOrdersStorage = "OrdersDatasource"

type inMemoryOrdersStorage struct {
	orders map[int]dsmodels.Order
}

func (s *inMemoryOrdersStorage) GetOrder(ctx context.Context, orderID int) monads.Result[dsmodels.Order] {
	utils.LogAction(ctx, compOrdersStorage, "GetOrder")

	order, exists := s.orders[orderID]
	if !exists {
		return monads.Errf[dsmodels.Order]("order not found")
	}
	return monads.Ok(order)
}

func (s *inMemoryOrdersStorage) GetAllOrdersForUser(ctx context.Context, userID int) monads.Result[[]dsmodels.Order] {
	utils.LogAction(ctx, compOrdersStorage, "GetAllOrdersForUser")

	userOrders := make([]dsmodels.Order, 0)
	for _, order := range s.orders {
		if order.UserId == userID {
			userOrders = append(userOrders, order)
		}
	}
	return monads.Ok(userOrders)
}

func (s *inMemoryOrdersStorage) DeleteOrder(ctx context.Context, orderID int) error {
	utils.LogAction(ctx, compOrdersStorage, "DeleteOrder")

	if _, exists := s.orders[orderID]; !exists {
		return errors.New("order not found")
	}
	delete(s.orders, orderID)
	return nil
}

func (s *inMemoryOrdersStorage) UpdateOrder(ctx context.Context, order dsmodels.Order) monads.Result[dsmodels.Order] {
	utils.LogAction(ctx, compOrdersStorage, "UpdateOrder")

	_, exists := s.orders[order.ID]
	if !exists {
		return monads.Errf[dsmodels.Order]("order not found")
	}
	s.orders[order.ID] = order
	return monads.Ok(order)
}

func (s *inMemoryOrdersStorage) InsertOrder(ctx context.Context, order dsmodels.Order) monads.Result[dsmodels.Order] {
	utils.LogAction(ctx, compOrdersStorage, "InsertOrder")

	if _, exists := s.orders[order.ID]; exists {
		return monads.Errf[dsmodels.Order]("order already exists")
	}
	s.orders[order.ID] = order
	return monads.Ok(order)
}

func NewOrdersStorage() datasources.OrdersDatasource {
	return &inMemoryOrdersStorage{
		orders: make(map[int]dsmodels.Order),
	}
}

