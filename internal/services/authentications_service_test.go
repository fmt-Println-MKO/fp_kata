package services

import (
	"errors"
	"fp_kata/internal/models"
	"fp_kata/pkg/log"
	zlog "github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetUserIDByToken(t *testing.T) {
	log.InitLogger()
	ctx := log.NewBackgroundContext(&zlog.Logger)

	type testCase struct {
		name          string
		authToken     string
		tokenStorage  map[string]int
		expectedUser  int
		expectedError string
		validate      func(*testing.T, int, error)
	}

	testCases := []testCase{
		{
			name:      "valid_token",
			authToken: "validToken",
			tokenStorage: map[string]int{
				"validToken": 42,
			},
			expectedUser:  42,
			expectedError: "",
			validate: func(t *testing.T, userID int, err error) {
				assert.NoError(t, err, "unexpected error for valid token")
				assert.Equal(t, 42, userID, "expected user ID does not match actual user ID")
			},
		},
		{
			name:          "missing_token",
			authToken:     "",
			tokenStorage:  map[string]int{},
			expectedUser:  0,
			expectedError: "invalid auth token",
			validate: func(t *testing.T, userID int, err error) {
				assert.EqualError(t, err, "invalid auth token", "expected error does not match actual error")
				assert.Equal(t, 0, userID, "expected user ID to be zero for missing token")
			},
		},
		{
			name:      "invalid_token_not_found",
			authToken: "invalidToken",
			tokenStorage: map[string]int{
				"validToken": 42,
			},
			expectedUser:  0,
			expectedError: "auth token not found",
			validate: func(t *testing.T, userID int, err error) {
				assert.EqualError(t, err, "auth token not found", "expected error does not match actual error")
				assert.Equal(t, 0, userID, "expected user ID to be zero for invalid token")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			authSvc := &authService{
				tokenStorage: tc.tokenStorage,
			}

			userID, err := authSvc.GetUserIDByToken(ctx, tc.authToken)
			tc.validate(t, userID, err)
		})
	}
}

func TestGenerateAuthToken(t *testing.T) {
	type testCase struct {
		name          string
		user          models.User
		expectedError error
		validate      func(*testing.T, string, error)
	}

	authSvc := NewAuthService()

	testCases := []testCase{
		{
			name: "valid_user",
			user: models.User{
				ID:       1,
				Username: "testuser",
				Email:    "testuser@example.com",
			},
			expectedError: nil,
			validate: func(t *testing.T, token string, err error) {
				assert.NoError(t, err, "expected no error, but got an error")
				assert.NotEmpty(t, token, "expected a valid token, but got an empty token")
			},
		},
		{
			name: "invalid_user_zero_id",
			user: models.User{
				ID:       0,
				Username: "invaliduser",
				Email:    "invaliduser@example.com",
			},
			expectedError: errors.New("invalid user"),
			validate: func(t *testing.T, token string, err error) {
				assert.EqualError(t, err, "invalid user", "expected error does not match actual error")
				assert.Empty(t, token, "expected empty token on error, but got a token")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			log.InitLogger()
			ctx := log.NewBackgroundContext(&zlog.Logger)
			token, err := authSvc.GenerateAuthToken(ctx, tc.user)
			tc.validate(t, token, err)
		})
	}
}
