package file

import (
	"fp_kata/internal/datasources"
	"fp_kata/internal/models"
)

type InMemoryStorage struct {
	store  map[int]models.User
	lastID int
}

func NewUserStorage() datasources.UsersDatasource {
	return &InMemoryStorage{
		store: make(map[int]models.User),
	}
}

func (s *InMemoryStorage) Create(user models.User) int {
	s.lastID++
	s.store[s.lastID] = user
	return s.lastID
}

func (s *InMemoryStorage) Read(id int) (models.User, bool) {
	user, exists := s.store[id]
	return user, exists
}

func (s *InMemoryStorage) Update(id int, user models.User) bool {
	_, exists := s.store[id]
	if exists {
		s.store[id] = user
	}
	return exists
}

func (s *InMemoryStorage) Delete(id int) {
	delete(s.store, id)
}

func (s *InMemoryStorage) GetByOrderId(orderId int) (models.User, bool) {

	for _, user := range s.store {
		if user.Orders != nil {
			for _, id := range user.Orders {
				if id == orderId {
					return user, true
				}
			}
		}
	}
	return models.User{}, false

}
