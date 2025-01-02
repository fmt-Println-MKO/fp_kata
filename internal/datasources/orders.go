package datasources

import "fp_kata/internal/models"

type OrdersDatasource interface {
	GetOrder(orderID int) (*models.Order, error)
	GetAllOrders(userID int) ([]*models.Order, error)
	DeleteOrder(orderID int) error
	UpdateOrder(order *models.Order) error
	InsertOrder(order *models.Order) error
}
