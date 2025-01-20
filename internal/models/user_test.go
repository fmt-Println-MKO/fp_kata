package models

import (
	"testing"

	"fp_kata/internal/datasources/dsmodels"
	"github.com/stretchr/testify/assert"
)

func TestUserToDSModel(t *testing.T) {
	type testCase struct {
		name     string
		user     User
		expected *dsmodels.User
	}

	tests := []testCase{
		{
			name: "valid complete user",
			user: User{
				ID:       1,
				Username: "testuser",
				Email:    "test@example.com",
				Password: "securepassword",
				Orders:   nil,
			},
			expected: &dsmodels.User{
				ID:       1,
				Username: "testuser",
				Email:    "test@example.com",
				Password: "securepassword",
			},
		},
		{
			name: "empty user",
			user: User{},
			expected: &dsmodels.User{
				ID:       0,
				Username: "",
				Email:    "",
				Password: "",
			},
		},
		{
			name: "user with some fields empty",
			user: User{
				ID:       2,
				Username: "",
				Email:    "incomplete@example.com",
				Password: "",
				Orders:   nil,
			},
			expected: &dsmodels.User{
				ID:       2,
				Username: "",
				Email:    "incomplete@example.com",
				Password: "",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.user.ToDSModel()
			assert.NotNil(t, result, "The result should not be nil")
			assert.Equal(t, tc.expected.ID, result.ID, "ID mismatch")
			assert.Equal(t, tc.expected.Username, result.Username, "Username mismatch")
			assert.Equal(t, tc.expected.Email, result.Email, "Email mismatch")
			assert.Equal(t, tc.expected.Password, result.Password, "Password mismatch")
		})
	}
}

func TestMapToUser(t *testing.T) {
	type testCase struct {
		name     string
		input    dsmodels.User
		expected *User
	}

	tests := []testCase{
		{
			name: "valid user",
			input: dsmodels.User{
				ID:       1,
				Username: "testuser",
				Email:    "test@example.com",
				Password: "securepassword",
			},
			expected: &User{
				ID:       1,
				Username: "testuser",
				Email:    "test@example.com",
				Password: "securepassword",
			},
		},
		{
			name:  "empty user",
			input: dsmodels.User{},
			expected: &User{
				ID:       0,
				Username: "",
				Email:    "",
				Password: "",
			},
		},
		{
			name: "user with missing fields",
			input: dsmodels.User{
				ID:    2,
				Email: "incomplete@example.com",
			},
			expected: &User{
				ID:       2,
				Username: "",
				Email:    "incomplete@example.com",
				Password: "",
			},
		},
		{
			name: "user with extra whitespace in fields",
			input: dsmodels.User{
				ID:       3,
				Username: "  spaceduser  ",
				Email:    " spaced@space.com ",
				Password: "   ",
			},
			expected: &User{
				ID:       3,
				Username: "  spaceduser  ",
				Email:    " spaced@space.com ",
				Password: "   ",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := MapToUser(tc.input)
			assert.NotNil(t, result, "Result should not be nil")
			assert.Equal(t, tc.expected.ID, result.ID, "ID mismatch for test case: %s", tc.name)
			assert.Equal(t, tc.expected.Username, result.Username, "Username mismatch for test case: %s", tc.name)
			assert.Equal(t, tc.expected.Email, result.Email, "Email mismatch for test case: %s", tc.name)
			assert.Equal(t, tc.expected.Password, result.Password, "Password mismatch for test case: %s", tc.name)
		})
	}
}
