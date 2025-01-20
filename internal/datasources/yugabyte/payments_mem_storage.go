package yugabyte

import (
	"context"
	"fmt"
	"fp_kata/internal/datasources"
	"fp_kata/internal/datasources/dsmodels"
	"fp_kata/pkg/log"
	"sort"
)

type inMemoryPaymentsStorage struct {
	payments map[int]dsmodels.Payment
}

func NewPaymentsStorage() datasources.PaymentsDatasource {
	return &inMemoryPaymentsStorage{make(map[int]dsmodels.Payment)}
}

func (s inMemoryPaymentsStorage) Create(ctx context.Context, p dsmodels.Payment) (dsmodels.Payment, error) {
	logger := log.GetLogger(ctx)
	logger.Debug().Str("comp", "PaymentsDatasource").Str("func", "Create").Send()
	id := len(s.payments) + 1
	p.Id = id
	s.payments[id] = p
	return p, nil
}

func (s inMemoryPaymentsStorage) Read(ctx context.Context, id int) (dsmodels.Payment, error) {
	logger := log.GetLogger(ctx)
	logger.Debug().Str("comp", "PaymentsDatasource").Str("func", "Read").Send()
	if p, exists := s.payments[id]; exists {
		return p, nil
	}
	return dsmodels.Payment{}, fmt.Errorf("payment with id %d not found", id)
}

func (s inMemoryPaymentsStorage) Update(ctx context.Context, p dsmodels.Payment) (dsmodels.Payment, error) {
	logger := log.GetLogger(ctx)
	logger.Debug().Str("comp", "PaymentsDatasource").Str("func", "Update").Send()
	if _, exists := s.payments[p.Id]; exists {
		s.payments[p.Id] = p
		return p, nil
	}
	return dsmodels.Payment{}, fmt.Errorf("payment with id %d not found", p.Id)
}

func (s inMemoryPaymentsStorage) Delete(ctx context.Context, id int) error {
	logger := log.GetLogger(ctx)
	logger.Debug().Str("comp", "PaymentsDatasource").Str("func", "Delete").Send()
	if _, exists := s.payments[id]; exists {
		delete(s.payments, id)
		return nil
	}
	return fmt.Errorf("payment with id %d not found", id)
}

func (s inMemoryPaymentsStorage) AllByOrderId(ctx context.Context, orderId int) ([]dsmodels.Payment, error) {
	logger := log.GetLogger(ctx)
	logger.Debug().Str("comp", "PaymentsDatasource").Str("func", "AllByOrderId").Send()

	var payments []dsmodels.Payment
	for _, payment := range s.payments {
		if payment.OrderId == orderId {
			payments = append(payments, payment)
		}
	}

	// Sort the payments by Id in ascending order
	sort.Slice(payments, func(i, j int) bool {
		return payments[i].Id < payments[j].Id
	})
	if len(payments) == 0 {
		return nil, fmt.Errorf("no payments found with orderId %d", orderId)
	}
	return payments, nil
}
