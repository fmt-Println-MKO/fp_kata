package transports

import (
	"fp_kata/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToUser(t *testing.T) {
	tests := []struct {
		name         string
		request      UserCreateRequest
		validateFunc func(t *testing.T, user *models.User)
	}{
		{
			name: "valid request",
			request: UserCreateRequest{
				Email:    "test@example.com",
				Password: "securepassword123",
			},
			validateFunc: func(t *testing.T, user *models.User) {
				assert.NotNil(t, user, "User should not be nil")
				assert.NotEmpty(t, user.Username, "Username should be generated")
				assert.Equal(t, "test@example.com", user.Email, "Email mismatch")
				assert.Equal(t, "securepassword123", user.Password, "Password mismatch")
			},
		},
		{
			name: "missing email",
			request: UserCreateRequest{
				Email:    "",
				Password: "securepassword123",
			},
			validateFunc: func(t *testing.T, user *models.User) {
				assert.Nil(t, user, "User should be nil")
			},
		},
		{
			name: "missing password",
			request: UserCreateRequest{
				Email:    "test@example.com",
				Password: "",
			},
			validateFunc: func(t *testing.T, user *models.User) {
				assert.Nil(t, user, "User should be nil")
			},
		},
		{
			name: "empty request",
			request: UserCreateRequest{
				Email:    "",
				Password: "",
			},
			validateFunc: func(t *testing.T, user *models.User) {
				assert.Nil(t, user, "User should be nil")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			user := tc.request.ToUser()

			tc.validateFunc(t, user) // Call validation helper
		})
	}
}
func TestExtractName(t *testing.T) {
	tests := []struct {
		name     string
		request  *UserCreateRequest
		validate func(t *testing.T, username string)
	}{
		{
			name: "valid email",
			request: &UserCreateRequest{
				Email: "test.user@example.com",
			},
			validate: func(t *testing.T, username string) {
				assert.Equal(t, "test.user", username, "Expected username to match email prefix")
			},
		},
		{
			name: "invalid email structure",
			request: &UserCreateRequest{
				Email: "invalid-email",
			},
			validate: func(t *testing.T, username string) {
				assert.Contains(t, username, "User", "Expected username to be randomly generated")
			},
		},
		{
			name: "empty email",
			request: &UserCreateRequest{
				Email: "",
			},
			validate: func(t *testing.T, username string) {
				assert.Contains(t, username, "User", "Expected username to be randomly generated for empty email")
			},
		},
		{
			name:    "nil request",
			request: nil,
			validate: func(t *testing.T, username string) {
				assert.Contains(t, username, "User", "Expected username to be randomly generated for nil request")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			username := extractName(tc.request)

			tc.validate(t, username)
		})
	}
}

func TestMapToUserResponse(t *testing.T) {
	tests := []struct {
		name     string
		user     models.User
		expected *UserResponse
	}{
		{
			name: "valid user",
			user: models.User{
				ID:       1,
				Username: "test.user",
				Email:    "test.user@example.com",
				Password: "securepassword123",
			},
			expected: &UserResponse{
				ID:       1,
				Username: "test.user",
				Email:    "test.user@example.com",
				Password: "securepassword123",
			},
		},
		{
			name: "empty user fields",
			user: models.User{
				ID:       0,
				Username: "",
				Email:    "",
				Password: "",
			},
			expected: &UserResponse{
				ID:       0,
				Username: "",
				Email:    "",
				Password: "",
			},
		},
		{
			name: "missing email",
			user: models.User{
				ID:       2,
				Username: "missing.email",
				Email:    "",
				Password: "password",
			},
			expected: &UserResponse{
				ID:       2,
				Username: "missing.email",
				Email:    "",
				Password: "password",
			},
		},
		{
			name: "missing username",
			user: models.User{
				ID:       3,
				Username: "",
				Email:    "email.only@example.com",
				Password: "password",
			},
			expected: &UserResponse{
				ID:       3,
				Username: "",
				Email:    "email.only@example.com",
				Password: "password",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := MapToUserResponse(tc.user)

			assert.NotNil(t, result, "Expected non-nil UserResponse")
			assert.Equal(t, tc.expected.ID, result.ID, "ID does not match")
			assert.Equal(t, tc.expected.Username, result.Username, "Username does not match")
			assert.Equal(t, tc.expected.Email, result.Email, "Email does not match")
			assert.Equal(t, tc.expected.Password, result.Password, "Password does not match")
		})
	}
}
