package file

import "fp_kata/internal/model"

type UserStorage interface {
	Create(user model.User) int
	Read(id int) (model.User, bool)
	Update(id int, user model.User) bool
	Delete(id int)
	GetByOrderId(id int) (model.User, bool)
}

type InMemoryStorage struct {
	store  map[int]model.User
	lastID int
}

func NewUserStorage() UserStorage {
	return &InMemoryStorage{
		store: make(map[int]model.User),
	}
}

func (s *InMemoryStorage) Create(user model.User) int {
	s.lastID++
	s.store[s.lastID] = user
	return s.lastID
}

func (s *InMemoryStorage) Read(id int) (model.User, bool) {
	user, exists := s.store[id]
	return user, exists
}

func (s *InMemoryStorage) Update(id int, user model.User) bool {
	_, exists := s.store[id]
	if exists {
		s.store[id] = user
	}
	return exists
}

func (s *InMemoryStorage) Delete(id int) {
	delete(s.store, id)
}

func (s *InMemoryStorage) GetByOrderId(orderId int) (model.User, bool) {

	for _, user := range s.store {
		if user.Orders != nil {
			for _, id := range user.Orders {
				if id == orderId {
					return user, true
				}
			}
		}
	}
	return model.User{}, false

}
