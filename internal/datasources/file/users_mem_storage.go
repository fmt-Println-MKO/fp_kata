package file

import (
	"context"
	"fp_kata/internal/datasources"
	"fp_kata/internal/datasources/dsmodels"
	"fp_kata/pkg/log"
)

type inMemoryUsersStorage struct {
	store  map[int]dsmodels.User
	lastID int
}

func NewUsersStorage() datasources.UsersDatasource {
	return &inMemoryUsersStorage{
		store:  make(map[int]dsmodels.User),
		lastID: 0,
	}
}

func (s *inMemoryUsersStorage) Create(ctx context.Context, user dsmodels.User) (dsmodels.User, bool) {
	logger := log.GetLogger(ctx)
	logger.Debug().Str("comp", "UsersDatasource").Str("func", "Create").Send()
	s.lastID++

	if s.lastID > 10 {
		return dsmodels.User{}, false
	}

	user.ID = s.lastID
	s.store[s.lastID] = user
	return user, true
}

func (s *inMemoryUsersStorage) Read(ctx context.Context, id int) (dsmodels.User, bool) {
	logger := log.GetLogger(ctx)
	logger.Debug().Str("comp", "UsersDatasource").Str("func", "Read").Send()
	user, exists := s.store[id]
	if !exists {
		return dsmodels.User{}, false
	}
	return user, true
}

func (s *inMemoryUsersStorage) Update(ctx context.Context, id int, user dsmodels.User) bool {
	logger := log.GetLogger(ctx)
	logger.Debug().Str("comp", "UsersDatasource").Str("func", "Update").Send()

	_, exists := s.store[id]
	if exists {
		s.store[id] = user
	}
	return exists
}

func (s *inMemoryUsersStorage) Delete(ctx context.Context, id int) bool {
	logger := log.GetLogger(ctx)
	logger.Debug().Str("comp", "UsersDatasource").Str("func", "Delete").Send()

	if _, exists := s.store[id]; !exists {
		return false
	}

	delete(s.store, id)
	_, exists := s.store[id]
	return !exists
}
