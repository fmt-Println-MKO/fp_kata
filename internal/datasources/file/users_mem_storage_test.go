package file

import (
	"context"
	"fp_kata/internal/datasources/dsmodels"
	"fp_kata/pkg/log"
	zlog "github.com/rs/zerolog/log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func initTestUsersStorage(store map[int]dsmodels.User, initialID int) (*inMemoryUsersStorage, context.Context) {
	log.InitLogger()
	ctx := log.NewBackgroundContext(&zlog.Logger)
	return &inMemoryUsersStorage{
		store:  store,
		lastID: initialID,
	}, ctx
}

func newUser(username, email, password string) dsmodels.User {
	return dsmodels.User{
		Username: username,
		Email:    email,
		Password: password,
	}
}

func TestInMemoryStorage_Create(t *testing.T) {
	tests := []struct {
		name        string
		initialID   int
		user        dsmodels.User
		wantSuccess bool
		wantID      int
		verify      func(*testing.T, *inMemoryUsersStorage, dsmodels.User, bool, bool, int)
	}{
		{
			name:        "valid creation at initial state",
			initialID:   0,
			user:        newUser("test1", "test1@example.com", "password1"),
			wantSuccess: true,
			wantID:      1,
			verify:      verifySuccessfulCreation,
		},
		{
			name:        "valid creation with existing entries",
			initialID:   5,
			user:        newUser("test2", "test2@example.com", "password2"),
			wantSuccess: true,
			wantID:      6,
			verify:      verifySuccessfulCreation,
		},
		{
			name:        "exceeding limit of 10 entries",
			initialID:   10,
			user:        newUser("test3", "test3@example.com", "password3"),
			wantSuccess: false,
			wantID:      0,
			verify:      verifyFailedCreation,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			storage, ctx := initTestUsersStorage(make(map[int]dsmodels.User), tc.initialID)
			createdUser, success := storage.Create(ctx, tc.user)
			tc.verify(t, storage, createdUser, success, tc.wantSuccess, tc.wantID)
		})
	}
}

func verifySuccessfulCreation(t *testing.T, s *inMemoryUsersStorage, createdUser dsmodels.User, success, wantSuccess bool, wantID int) {
	assert.Equal(t, wantSuccess, success, "success mismatch")
	assert.Equal(t, wantID, createdUser.ID, "user ID mismatch")
	assert.Equal(t, createdUser, s.store[wantID], "user in store mismatch")
}

func verifyFailedCreation(t *testing.T, _ *inMemoryUsersStorage, createdUser dsmodels.User, success, wantSuccess bool, _ int) {
	assert.Equal(t, wantSuccess, success, "success mismatch")
	assert.Equal(t, 0, createdUser.ID, "expected no user to be created, but got ID")
}

func TestInMemoryStorage_Read(t *testing.T) {

	tests := []struct {
		name     string
		initial  map[int]dsmodels.User
		id       int
		wantUser dsmodels.User
		wantOK   bool
		verify   func(*testing.T, dsmodels.User, bool, dsmodels.User, bool)
	}{
		{
			name: "read existing user",
			initial: map[int]dsmodels.User{
				1: {ID: 1, Username: "test1", Email: "test1@example.com"},
			},
			id:       1,
			wantUser: dsmodels.User{ID: 1, Username: "test1", Email: "test1@example.com"},
			wantOK:   true,
			verify:   verifyUserReadSuccess,
		},
		{
			name: "read non-existing user",
			initial: map[int]dsmodels.User{
				1: {ID: 1, Username: "test1", Email: "test1@example.com"},
			},
			id:       2,
			wantUser: dsmodels.User{},
			wantOK:   false,
			verify:   verifyUserReadFailure,
		},
	}

	for _, tc := range tests {
		storage, ctx := initTestUsersStorage(tc.initial, len(tc.initial)+1)
		t.Run(tc.name, func(t *testing.T) {
			gotUser, gotOK := storage.Read(ctx, tc.id)

			tc.verify(t, gotUser, gotOK, tc.wantUser, tc.wantOK)
		})
	}
}

func verifyUserReadSuccess(t *testing.T, gotUser dsmodels.User, gotOK bool, wantUser dsmodels.User, wantOK bool) {
	assert.Equal(t, wantOK, gotOK, "expected success mismatch")
	assert.Equal(t, wantUser, gotUser, "user data mismatch")
}

func verifyUserReadFailure(t *testing.T, gotUser dsmodels.User, gotOK bool, wantUser dsmodels.User, wantOK bool) {
	assert.Equal(t, wantOK, gotOK, "expected failure mismatch")
	assert.Equal(t, wantUser, gotUser, "non-existing user mismatch")
}

func TestInMemoryStorage_Update(t *testing.T) {

	tests := []struct {
		name       string
		initial    map[int]dsmodels.User
		id         int
		updateUser dsmodels.User
		wantOK     bool
		verify     func(*testing.T, *inMemoryUsersStorage, dsmodels.User, bool, map[int]dsmodels.User, bool)
	}{
		{
			name: "update existing user",
			initial: map[int]dsmodels.User{
				1: {ID: 1, Username: "test1", Email: "test1@example.com", Password: "password1"},
			},
			id:         1,
			updateUser: dsmodels.User{ID: 1, Username: "updatedUser", Email: "updated@example.com", Password: "newpassword"},
			wantOK:     true,
			verify:     verifyUpdateSuccess,
		},
		{
			name: "update non-existing user",
			initial: map[int]dsmodels.User{
				1: {ID: 1, Username: "test1", Email: "test1@example.com", Password: "password1"},
			},
			id:         2,
			updateUser: dsmodels.User{ID: 2, Username: "nonexistentUser", Email: "nonexistent@example.com", Password: "doesnotmatter"},
			wantOK:     false,
			verify:     verifyUpdateFailure,
		},
	}

	for _, tc := range tests {
		storage, ctx := initTestUsersStorage(tc.initial, len(tc.initial)+1)
		t.Run(tc.name, func(t *testing.T) {
			gotOK := storage.Update(ctx, tc.id, tc.updateUser)

			tc.verify(t, storage, tc.updateUser, gotOK, tc.initial, tc.wantOK)
		})
	}
}

func verifyUpdateSuccess(t *testing.T, s *inMemoryUsersStorage, updatedUser dsmodels.User, gotOK bool, initial map[int]dsmodels.User, wantOK bool) {
	assert.Equal(t, wantOK, gotOK, "expected success mismatch")
	assert.Equal(t, updatedUser, s.store[updatedUser.ID], "user data not correctly updated in store")
}

func verifyUpdateFailure(t *testing.T, s *inMemoryUsersStorage, updatedUser dsmodels.User, gotOK bool, initial map[int]dsmodels.User, wantOK bool) {
	assert.Equal(t, wantOK, gotOK, "expected failure mismatch")
	_, exists := initial[updatedUser.ID]
	assert.False(t, exists, "unexpected user modification for non-existing ID")
	assert.NotContains(t, s.store, updatedUser.ID, "user should not exist in store after failed update")
}

func TestInMemoryStorage_Delete(t *testing.T) {

	tests := []struct {
		name        string
		initial     map[int]dsmodels.User
		idToDelete  int
		wantSuccess bool
		verify      func(*testing.T, *inMemoryUsersStorage, int, bool, bool)
	}{
		{
			name: "delete existing user",
			initial: map[int]dsmodels.User{
				1: {ID: 1, Username: "test1", Email: "test1@example.com"},
				2: {ID: 2, Username: "test2", Email: "test2@example.com"},
			},
			idToDelete:  1,
			wantSuccess: true,
			verify:      verifySuccessfulDeletion,
		},
		{
			name: "delete non-existing user",
			initial: map[int]dsmodels.User{
				1: {ID: 1, Username: "test1", Email: "test1@example.com"},
			},
			idToDelete:  2,
			wantSuccess: false,
			verify:      verifyFailedDeletion,
		},
		{
			name:        "delete user in empty storage",
			initial:     map[int]dsmodels.User{},
			idToDelete:  1,
			wantSuccess: false,
			verify:      verifyFailedDeletion,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			storage, ctx := initTestUsersStorage(tc.initial, len(tc.initial)+1)

			success := storage.Delete(ctx, tc.idToDelete)

			tc.verify(t, storage, tc.idToDelete, success, tc.wantSuccess)
		})
	}
}

func verifySuccessfulDeletion(t *testing.T, s *inMemoryUsersStorage, id int, success, wantSuccess bool) {
	assert.Equal(t, wantSuccess, success, "deletion success mismatch")
	_, exists := s.store[id]
	assert.False(t, exists, "user ID should no longer exist in store")
}

func verifyFailedDeletion(t *testing.T, s *inMemoryUsersStorage, id int, success, wantSuccess bool) {
	assert.Equal(t, wantSuccess, success, "deletion success mismatch")
	_, exists := s.store[id]
	assert.False(t, success, "delete should fail but succeeded")
	assert.False(t, exists, "delete should fail but succeeded")
}

func TestNewUserStorage(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "new storage initialization",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := NewUsersStorage()
			memStorage, ok := storage.(*inMemoryUsersStorage)
			assert.True(t, ok, "expected storage to be of type *inMemoryUsersStorage")
			assert.NotNil(t, memStorage, "storage instance should not be nil")
			assert.NotNil(t, memStorage.store, "storage map should be initialized")
			assert.Empty(t, memStorage.store, "storage map should be empty initially")
			assert.Equal(t, 0, memStorage.lastID, "lastID should be initialized to 0")
		})
	}
}
