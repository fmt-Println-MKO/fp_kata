package datasources

import (
	"context"
	"fp_kata/common/monads"
	"fp_kata/internal/datasources/dsmodels"
)

type OrdersDatasource interface {
	GetOrder(ctx context.Context, orderID int) monads.Result[dsmodels.Order]
	GetAllOrdersForUser(ctx context.Context, userID int) monads.Result[[]dsmodels.Order]
	DeleteOrder(ctx context.Context, orderID int) error
	UpdateOrder(ctx context.Context, order dsmodels.Order) monads.Result[dsmodels.Order]
	InsertOrder(ctx context.Context, order dsmodels.Order) monads.Result[dsmodels.Order]
}
