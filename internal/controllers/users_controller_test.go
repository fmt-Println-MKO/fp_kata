package controllers

import (
	"bytes"
	"encoding/json"
	"fp_kata/internal/models"
	"fp_kata/internal/services"
	"fp_kata/mocks"
	"fp_kata/pkg/transports"
	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func createUsersTestApp(usersService services.UsersService, contextData *map[any]any) *fiber.App {
	app := fiber.New()

	mockData := make(map[any]any)
	if contextData != nil {
		mockData = *contextData
	}

	ctx := &mocks.CustomCtx{
		DefaultCtx: *fiber.NewDefaultCtx(app),
		MockLocals: mockData,
	}
	app.NewCtxFunc(func(app *fiber.App) fiber.CustomCtx {

		return ctx
	})

	controller := &UsersController{userService: usersService}
	app.Post("/users", controller.SignUp)
	app.Get("/users/me", controller.GetUser)

	return app
}

func TestSignUp(t *testing.T) {
	tests := []struct {
		name                 string
		inputBody            transports.UserCreateRequest
		mockSetup            func(service *mocks.UsersService)
		expectedStatus       int
		expectedResponseBody map[string]interface{}
	}{
		{
			name: "valid user data",
			inputBody: transports.UserCreateRequest{
				Email:    "testuser@example.com",
				Password: "password123",
			},
			mockSetup: func(service *mocks.UsersService) {
				service.On(
					"SignUp",
					mock.Anything,
					mock.AnythingOfType("models.User"),
				).Return(&models.User{
					ID:       1,
					Username: "testuser",
					Email:    "testuser@example.com",
					Password: "password123",
				}, nil)
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
			mockSetup:      func(service *mocks.UsersService) {},
			expectedStatus: fiber.StatusBadRequest,
			expectedResponseBody: map[string]interface{}{
				"error": "Invalid input",
			},
		},
		{
			name: "service error",
			inputBody: transports.UserCreateRequest{
				Email:    "testuser@example.com",
				Password: "password123",
			},
			mockSetup: func(service *mocks.UsersService) {
				service.On(
					"SignUp",
					mock.Anything,
					mock.AnythingOfType("models.User"),
				).Return(nil, assert.AnError)
			},
			expectedStatus: fiber.StatusInternalServerError,
			expectedResponseBody: map[string]interface{}{
				"error": "Could not create user",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &mocks.UsersService{}
			tt.mockSetup(mockService)

			mockContextData := mocks.ProvideBaseMockContextData(nil)
			app := createUsersTestApp(mockService, mockContextData)

			body, _ := json.Marshal(tt.inputBody)
			req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, _ := app.Test(req)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode, "unexpected status code")

			var buf bytes.Buffer
			buf.ReadFrom(resp.Body)
			responseBody := buf.String()

			expectedResponseBody, _ := json.Marshal(tt.expectedResponseBody)

			// Compare the responseBody with the encoded expectedResponseBody string.
			assert.JSONEq(t, string(expectedResponseBody), responseBody, "unexpected response body")

			mockService.AssertExpectations(t)
		})
	}
}

func TestGetUser(t *testing.T) {
	tests := []struct {
		name                 string
		authenticatedUser    *models.User
		mockSetup            func(service *mocks.UsersService)
		expectedStatus       int
		expectedResponseBody map[string]interface{}
	}{
		{
			name: "valid user",
			authenticatedUser: &models.User{
				ID: 1,
			},
			mockSetup: func(usersService *mocks.UsersService) {
				usersService.On(
					"GetUserByID",
					mock.Anything,
					1,
				).Return(&models.User{
					ID:       1,
					Username: "testuser",
					Email:    "testuser@example.com",
					Password: "password123",
				}, nil)

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
			name: "service error",
			authenticatedUser: &models.User{
				ID: 1,
			},
			mockSetup: func(usersService *mocks.UsersService) {
				usersService.On("GetUserByID", mock.Anything, 1).Return(nil, assert.AnError).Once()
			},
			expectedStatus: fiber.StatusInternalServerError,
			expectedResponseBody: map[string]interface{}{
				"error": "Could not retrieve user",
			},
		},
		{
			name:              "missing user ID in context",
			authenticatedUser: nil,
			mockSetup:         func(service *mocks.UsersService) {},
			expectedStatus:    fiber.StatusUnauthorized,
			expectedResponseBody: map[string]interface{}{
				"error": "Authorization token is missing",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockUsersService := &mocks.UsersService{}

			tc.mockSetup(mockUsersService)

			mockContextData := mocks.ProvideBaseMockContextData(tc.authenticatedUser)
			app := createUsersTestApp(mockUsersService, mockContextData)

			req := httptest.NewRequest(http.MethodGet, "/users/me", nil)

			resp, _ := app.Test(req)

			assert.Equal(t, tc.expectedStatus, resp.StatusCode, "unexpected status code")

			var buf bytes.Buffer
			_, _ = buf.ReadFrom(resp.Body)
			responseBody := buf.String()

			expectedResponseBody, _ := json.Marshal(tc.expectedResponseBody)
			assert.JSONEq(t, string(expectedResponseBody), responseBody, "unexpected response body")
			mockUsersService.AssertExpectations(t)

		})
	}
}
