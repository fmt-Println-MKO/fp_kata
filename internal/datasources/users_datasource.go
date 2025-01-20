package datasources

import (
	"context"
	"fp_kata/internal/datasources/dsmodels"
)

type UsersDatasource interface {
	Create(ctx context.Context, user dsmodels.User) (dsmodels.User, bool)
	Read(ctx context.Context, id int) (dsmodels.User, bool)
	Update(ctx context.Context, id int, user dsmodels.User) bool
	Delete(ctx context.Context, id int) bool
}
