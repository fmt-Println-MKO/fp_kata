package datasources

import (
	"context"
	"fp_kata/internal/datasources/dsmodels"
	"github.com/samber/mo"
)

type OrdersDatasource interface {
	GetOrder(ctx context.Context, orderID int) mo.Result[dsmodels.Order]
	GetAllOrdersForUser(ctx context.Context, userID int) ([]dsmodels.Order, error)
	DeleteOrder(ctx context.Context, orderID int) error
	UpdateOrder(ctx context.Context, order dsmodels.Order) (*dsmodels.Order, error)
	InsertOrder(ctx context.Context, order dsmodels.Order) (*dsmodels.Order, error)
}
