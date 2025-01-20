package controllers

import (
	"bytes"
	"encoding/json"
	"fp_kata/internal/app"
	"fp_kata/pkg/transports"
	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSignUp(t *testing.T) {
	tests := []struct {
		name                 string
		inputBody            transports.UserCreateRequest
		expectedStatus       int
		expectedResponseBody map[string]interface{}
	}{
		{
			name: "valid user data",
			inputBody: transports.UserCreateRequest{
				Email:    "testuser@example.com",
				Password: "password123",
			},
			expectedStatus: fiber.StatusCreated,
			expectedResponseBody: map[string]interface{}{
				"id":       1,
				"username": "testuser",
				"email":    "testuser@example.com",
				"password": "password123",
			},
		},
		{
			name: "invalid request body",
			inputBody: transports.UserCreateRequest{
				Email:    "",
				Password: "",
			},
			expectedStatus: fiber.StatusBadRequest,
			expectedResponseBody: map[string]interface{}{
				"error": "Invalid input",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := app.InitApp()
			body, _ := json.Marshal(tt.inputBody)
			req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, _ := app.Test(req)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode, "unexpected status code")

			var buf bytes.Buffer
			buf.ReadFrom(resp.Body)
			responseBody := buf.String()

			expectedResponseBody, _ := json.Marshal(tt.expectedResponseBody)
			assert.JSONEq(t, string(expectedResponseBody), responseBody, "unexpected response body")
		})
	}
}

func TestGetUser(t *testing.T) {
	tests := []struct {
		name                     string
		setupAuthorizationHeader func(req *http.Request)
		expectedStatus           int
		expectedResponseBody     map[string]interface{}
	}{
		{
			name: "valid user",

			setupAuthorizationHeader: func(req *http.Request) {
				req.Header.Set("Authorization", "token_1")
			},
			expectedStatus: fiber.StatusOK,
			expectedResponseBody: map[string]interface{}{
				"id":       1,
				"username": "testuser",
				"email":    "testuser@example.com",
				"password": "password123",
			},
		},
		{
			name: "missing user ID in context",
			setupAuthorizationHeader: func(req *http.Request) {

			},
			expectedStatus: fiber.StatusUnauthorized,
			expectedResponseBody: map[string]interface{}{
				"error": "Authorization token is missing",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			app := app.InitApp()
			PrepareUser(t, app)

			req := httptest.NewRequest(http.MethodGet, "/users/me", nil)
			tc.setupAuthorizationHeader(req)

			resp, _ := app.Test(req)

			assert.Equal(t, tc.expectedStatus, resp.StatusCode, "unexpected status code")

			var buf bytes.Buffer
			_, _ = buf.ReadFrom(resp.Body)
			responseBody := buf.String()

			expectedResponseBody, _ := json.Marshal(tc.expectedResponseBody)
			assert.JSONEq(t, string(expectedResponseBody), responseBody, "unexpected response body")
		})
	}
}

func PrepareUser(t *testing.T, app *fiber.App) {
	signUpRequest := transports.UserCreateRequest{
		Email:    "testuser@example.com",
		Password: "password123",
	}
	signUpBody, _ := json.Marshal(signUpRequest)
	signUpReq := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(signUpBody))
	signUpReq.Header.Set("Content-Type", "application/json")

	signUpResp, err := app.Test(signUpReq)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, signUpResp.StatusCode)

}
