package datasources

import (
	"context"
	"fp_kata/internal/datasources/dsmodels"
	"github.com/samber/mo"
)

type OrdersDatasource interface {
	GetOrder(ctx context.Context, orderID int) mo.Result[dsmodels.Order]
	GetAllOrdersForUser(ctx context.Context, userID int) mo.Result[[]dsmodels.Order]
	DeleteOrder(ctx context.Context, orderID int) error
	UpdateOrder(ctx context.Context, order dsmodels.Order) mo.Result[dsmodels.Order]
	InsertOrder(ctx context.Context, order dsmodels.Order) mo.Result[dsmodels.Order]
}
