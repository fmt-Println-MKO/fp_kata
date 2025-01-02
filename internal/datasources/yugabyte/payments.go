package yugabyte

import (
	"fmt"
	"fp_kata/internal/datasources"
	"fp_kata/internal/models"
)

type inMemoryPaymentStorage map[int]models.Payment

func NewPaymentStorage() datasources.PaymentsDatasource {
	return make(inMemoryPaymentStorage)
}

func (s inMemoryPaymentStorage) Create(p models.Payment) (int, error) {
	id := len(s) + 1
	p.PaymentID = id
	s[id] = p
	return id, nil
}

func (s inMemoryPaymentStorage) Read(id int) (models.Payment, error) {
	if p, exists := s[id]; exists {
		return p, nil
	}
	return models.Payment{}, fmt.Errorf("payment with id %d not found", id)
}

func (s inMemoryPaymentStorage) Update(id int, p models.Payment) error {
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

func (s inMemoryPaymentStorage) AllByOrderId(orderId int) ([]models.Payment, error) {

	var payments []models.Payment
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
