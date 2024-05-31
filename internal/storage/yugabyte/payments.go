package yugabyte

import (
	"fmt"
	"fp_kata/internal/model"
)

type PaymentStorage interface {
	Create(model.Payment) (int, error)
	Read(int) (model.Payment, error)
	Update(int, model.Payment) error
	Delete(int) error
	AllByOrderId(int) ([]model.Payment, error)
}

type inMemoryPaymentStorage map[int]model.Payment

func NewPaymentStorage() PaymentStorage {
	return make(inMemoryPaymentStorage)
}

func (s inMemoryPaymentStorage) Create(p model.Payment) (int, error) {
	id := len(s) + 1
	p.PaymentID = id
	s[id] = p
	return id, nil
}

func (s inMemoryPaymentStorage) Read(id int) (model.Payment, error) {
	if p, exists := s[id]; exists {
		return p, nil
	}
	return model.Payment{}, fmt.Errorf("payment with id %d not found", id)
}

func (s inMemoryPaymentStorage) Update(id int, p model.Payment) error {
	if _, exists := s[id]; exists {
		p.PaymentID = id
		s[id] = p
		return nil
	}
	return fmt.Errorf("payment with id %d not found", id)
}

func (s inMemoryPaymentStorage) Delete(id int) error {
	if _, exists := s[id]; exists {
		delete(s, id)
		return nil
	}
	return fmt.Errorf("payment with id %d not found", id)
}

func (s inMemoryPaymentStorage) AllByOrderId(orderId int) ([]model.Payment, error) {

	var payments []model.Payment
	for _, payment := range s {
		if payment.OrderId == orderId {
			payments = append(payments, payment)
		}
	}
	if len(payments) == 0 {
		return nil, fmt.Errorf("no payments found with orderId %d", orderId)
	}
	return payments, nil
}
