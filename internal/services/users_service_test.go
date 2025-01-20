package services

import (
	"errors"
	"fp_kata/internal/datasources/dsmodels"
	"fp_kata/internal/models"
	"fp_kata/mocks"
	"fp_kata/pkg/log"
	zlog "github.com/rs/zerolog/log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUserByID(t *testing.T) {
	log.InitLogger()
	ctx := log.NewBackgroundContext(&zlog.Logger)

	mockStorage := mocks.NewUsersDatasource(t)
	mockAuthService := mocks.NewAuthService(t)
	userSvc := NewUsersService(mockStorage, mockAuthService)

	dsUser := dsmodels.User{ID: 1, Username: "John Doe"}
	expectedUser := models.MapToUser(dsUser)

	testCases := []struct {
		name         string
		id           int
		mockSetup    func()
		expectedUser *models.User
		assertError  func(t *testing.T, err error)
	}{
		{
			name: "user exists",
			id:   1,
			mockSetup: func() {
				mockStorage.
					On("Read", ctx, 1).
					Return(dsUser, true).
					Once()
			},
			expectedUser: expectedUser,
			assertError: func(t *testing.T, err error) {
				assert.NoError(t, err, "unexpected error occurred")
			},
		},
		{
			name: "user does not exist",
			id:   2,
			mockSetup: func() {
				mockStorage.
					On("Read", ctx, 2).
					Return(dsmodels.User{}, false).
					Once()
			},
			expectedUser: nil,
			assertError: func(t *testing.T, err error) {
				assert.EqualError(t, err, "no user found for id", "expected error message does not match")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()

			user, err := userSvc.GetUserByID(ctx, tc.id)

			assert.Equal(t, tc.expectedUser, user, "expected user result does not match")
			tc.assertError(t, err)

			mockStorage.AssertExpectations(t)
		})
	}
}

func TestSignUp(t *testing.T) {
	log.InitLogger()
	ctx := log.NewBackgroundContext(&zlog.Logger)

	mockStorage := mocks.NewUsersDatasource(t)
	mockAuthService := mocks.NewAuthService(t)
	userSvc := NewUsersService(mockStorage, mockAuthService)

	inputUser := models.User{ID: 1, Username: "John Doe", Email: "john.doe@email.com", Password: "password123"}
	dsInputUser := inputUser.ToDSModel()
	createdDsUser := *dsInputUser
	expectedUser := models.MapToUser(createdDsUser)

	testCases := []struct {
		name         string
		input        models.User
		mockSetup    func()
		expectedUser *models.User
		assertError  func(t *testing.T, err error)
	}{
		{
			name:  "successful signup",
			input: inputUser,
			mockSetup: func() {
				mockStorage.
					On("Create", ctx, *dsInputUser).
					Return(createdDsUser, true).
					Once()
				mockAuthService.
					On("GenerateAuthToken", ctx, *expectedUser).
					Return("dummy_token", nil).
					Once()
			},
			expectedUser: expectedUser,
			assertError: func(t *testing.T, err error) {
				assert.NoError(t, err, "unexpected error occurred during signup")
			},
		},
		{
			name:  "failed signup - storage full",
			input: inputUser,
			mockSetup: func() {
				mockStorage.
					On("Create", ctx, *dsInputUser).
					Return(dsmodels.User{}, false).
					Once()
			},
			expectedUser: nil,
			assertError: func(t *testing.T, err error) {
				assert.EqualError(t, err, "user storage is full", "expected error message does not match")
			},
		},
		{
			name:  "failed signup - auth token error",
			input: inputUser,
			mockSetup: func() {
				mockStorage.
					On("Create", ctx, *dsInputUser).
					Return(createdDsUser, true).
					Once()
				mockAuthService.
					On("GenerateAuthToken", ctx, *expectedUser).
					Return("", errors.New("auth token generation failed")).
					Once()
			},
			expectedUser: nil,
			assertError: func(t *testing.T, err error) {
				assert.EqualError(t, err, "auth token generation failed", "expected error message does not match")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()

			user, err := userSvc.SignUp(ctx, tc.input)

			assert.Equal(t, tc.expectedUser, user, "expected user result does not match")
			tc.assertError(t, err)

			mockStorage.AssertExpectations(t)
			mockAuthService.AssertExpectations(t)
		})
	}
}
