package datasources

import "fp_kata/internal/models"

type PaymentsDatasource interface {
	Create(models.Payment) (int, error)
	Read(int) (models.Payment, error)
	Update(int, models.Payment) error
	Delete(int) error
	AllByOrderId(int) ([]models.Payment, error)
}
