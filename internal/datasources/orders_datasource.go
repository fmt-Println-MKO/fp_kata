package datasources

import (
	"context"
	"fp_kata/internal/datasources/dsmodels"
)

type OrdersDatasource interface {
	GetOrder(ctx context.Context, orderID int) (*dsmodels.Order, error)
	GetAllOrdersForUser(ctx context.Context, userID int) ([]dsmodels.Order, error)
	DeleteOrder(ctx context.Context, orderID int) error
	UpdateOrder(ctx context.Context, order dsmodels.Order) (*dsmodels.Order, error)
	InsertOrder(ctx context.Context, order dsmodels.Order) (*dsmodels.Order, error)
}
