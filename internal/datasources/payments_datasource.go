package datasources

import (
	"context"
	"fp_kata/internal/datasources/dsmodels"
)

type PaymentsDatasource interface {
	Create(ctx context.Context, payment dsmodels.Payment) (dsmodels.Payment, error)
	Read(ctx context.Context, paymentId int) (dsmodels.Payment, error)
	Update(ctx context.Context, payment dsmodels.Payment) (dsmodels.Payment, error)
	Delete(ctx context.Context, paymentId int) error
	AllByOrderId(ctx context.Context, paymentId int) ([]dsmodels.Payment, error)
}
