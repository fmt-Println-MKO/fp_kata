package datasources

import "fp_kata/internal/models"

type UsersDatasource interface {
	Create(user models.User) int
	Read(id int) (models.User, bool)
	Update(id int, user models.User) bool
	Delete(id int)
	GetByOrderId(id int) (models.User, bool)
}
