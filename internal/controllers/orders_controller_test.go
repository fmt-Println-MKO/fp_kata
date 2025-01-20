package controllers

import (
	"bytes"
	"encoding/json"
	"fp_kata/common"
	"fp_kata/internal/models"
	"fp_kata/internal/services"
	"fp_kata/mocks"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"strconv"
	"time"

	"fp_kata/pkg/transports"
	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"testing"
)

func createTestOrdersController(mockOrdersService services.OrdersService, contextData *map[any]any) *fiber.App {
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
	controller := &OrdersController{orderService: mockOrdersService}
	app.Post("/orders", controller.CreateOrder)
	app.Get("/orders", controller.GetOrders)
	app.Get("/orders/:id", controller.GetOrder)

	return app

}

func TestCreateOrder(t *testing.T) {

	setupValidStoreOrderMock := func(mockOrdersService *mocks.OrdersService, body transports.OrderCreateRequest, user models.User, mockReturn *models.Order, mockError error) {
		order := body.ToOrder(user)
		mockOrdersService.
			On("StoreOrder", mock.Anything, user.ID, *order).
			Return(mockReturn, mockError)
	}

	tests := []struct {
		name string
		body transports.OrderCreateRequest

		user                   models.User
		mockReturn             *models.Order
		mockError              error
		setupOrdersServiceMock func(mockOrdersService *mocks.OrdersService, body transports.OrderCreateRequest, user models.User, mockReturn *models.Order, mockError error)

		expectedCode int
		expectedJSON map[string]interface{}
	}{
		{
			name: "success",
			body: transports.OrderCreateRequest{
				ProductID: 1,
				Quantity:  2,
				Price:     10.23,
				OrderDate: time.Date(2025, 1, 30, 10, 30, 0, 0, time.UTC),
				Payments: []*transports.PaymentRequest{
					{
						PaymentMethod: common.CreditCard,
						PaymentAmount: 10.23,
					},
				},
			},

			user: models.User{ID: 1, Username: "John Doe"},
			mockReturn: &models.Order{
				ID:        42,
				User:      &models.User{ID: 1, Username: "John Doe"},
				ProductID: 1,
				Quantity:  2,
				Price:     10.23,
				OrderDate: time.Date(2025, 1, 30, 10, 30, 0, 0, time.UTC),
				Payments: []*models.Payment{
					{
						Id:     1,
						Method: common.CreditCard,
						Amount: 10.23,
					},
				}},
			mockError:              nil,
			setupOrdersServiceMock: setupValidStoreOrderMock,

			expectedCode: fiber.StatusCreated,
			expectedJSON: map[string]interface{}{
				"has_weightables": false,
				"id":              42,
				"order_date":      "2025-01-30T10:30:00Z",
				"payments": []interface{}{map[string]interface{}{
					"amount": 10.23,
					"id":     1,
					"method": "CreditCard"}},
				"price":      10.23,
				"product_id": 1,
				"quantity":   2,
				"user": map[string]interface{}{
					"email":    "",
					"id":       1,
					"password": "",
					"username": "John Doe"}},
		},
		{
			name: "bad request",
			body: transports.OrderCreateRequest{}, // Missing required fields

			user:       models.User{ID: 1, Username: "John Doe"},
			mockReturn: nil,
			mockError:  nil,
			setupOrdersServiceMock: func(mockOrdersService *mocks.OrdersService, body transports.OrderCreateRequest, user models.User, mockReturn *models.Order, mockError error) {
			},
			expectedCode: fiber.StatusBadRequest,
			expectedJSON: map[string]interface{}{"details": "Key: 'OrderCreateRequest.ProductID' Error:Field validation for 'ProductID' failed on the 'required' tag\nKey: 'OrderCreateRequest.Quantity' Error:Field validation for 'Quantity' failed on the 'required' tag\nKey: 'OrderCreateRequest.Price' Error:Field validation for 'Price' failed on the 'required' tag\nKey: 'OrderCreateRequest.OrderDate' Error:Field validation for 'OrderDate' failed on the 'required' tag\nKey: 'OrderCreateRequest.Payments' Error:Field validation for 'Payments' failed on the 'required' tag", "error": "Validation failed"},
		},
		{
			name: "internal server error",
			body: transports.OrderCreateRequest{
				ProductID: 1,
				Quantity:  2,
				Price:     10.23,
				OrderDate: time.Date(2025, 1, 30, 10, 30, 0, 0, time.UTC),
				Payments: []*transports.PaymentRequest{
					{
						PaymentMethod: common.CreditCard,
						PaymentAmount: 10.23,
					},
				},
			},

			user:                   models.User{ID: 1, Username: "John Doe"},
			mockReturn:             nil,
			setupOrdersServiceMock: setupValidStoreOrderMock,
			mockError:              assert.AnError,
			expectedCode:           fiber.StatusInternalServerError,
			expectedJSON:           map[string]interface{}{"error": "Unable to create the order"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockOrdersService := new(mocks.OrdersService)
			tc.setupOrdersServiceMock(mockOrdersService, tc.body, tc.user, tc.mockReturn, tc.mockError)

			mockContextData := mocks.ProvideBaseMockContextData(&tc.user)
			app := createTestOrdersController(mockOrdersService, mockContextData)
			body, _ := json.Marshal(tc.body)
			req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)

			assert.Nil(t, err, "Handler should not return an error")
			assert.Equal(t, tc.expectedCode, resp.StatusCode, "Unexpected status code")

			var buf bytes.Buffer
			buf.ReadFrom(resp.Body)
			responseBody := buf.String()

			expectedResponseBody, _ := json.Marshal(tc.expectedJSON)
			assert.JSONEq(t, string(expectedResponseBody), responseBody, "Unexpected response JSON")

			mockOrdersService.AssertExpectations(t)
		})
	}
}

func TestGetOrders(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    string
		user           models.User
		mockReturn     []*models.Order
		mockError      error
		setServiceMock func(mockOrdersService *mocks.OrdersService, user models.User, queryParams string, mockReturn []*models.Order, mockError error)
		assertFunc     func(t *testing.T, responseBody string, responseCode int)
	}{
		{
			name:        "success - no filter",
			queryParams: "",
			user:        models.User{ID: 1, Username: "Jane Doe"},
			mockReturn: []*models.Order{
				{
					ID:        1,
					ProductID: 1,
					Price:     19.99,
					OrderDate: time.Date(2025, 1, 30, 10, 30, 0, 0, time.UTC),
				},
				{
					ID:        2,
					ProductID: 2,
					Price:     29.99,
					OrderDate: time.Date(2025, 2, 10, 12, 0, 0, 0, time.UTC),
				},
			},
			mockError: nil,
			setServiceMock: func(mockOrdersService *mocks.OrdersService, user models.User, queryParams string, mockReturn []*models.Order, mockError error) {
				mockOrdersService.On("GetOrders", mock.Anything, user.ID).Return(mockReturn, mockError)
			},
			assertFunc: func(t *testing.T, responseBody string, responseCode int) {
				assert.Equal(t, fiber.StatusOK, responseCode, "Unexpected status code")

				expectedResponseBody, _ := json.Marshal([]interface{}{
					map[string]interface{}{
						"has_weightables": false,
						"id":              1,
						"order_date":      "2025-01-30T10:30:00Z",
						"price":           19.99,
						"product_id":      1,
					},
					map[string]interface{}{
						"has_weightables": false,
						"id":              2,
						"order_date":      "2025-02-10T12:00:00Z",
						"price":           29.99,
						"product_id":      2},
				})
				assert.JSONEq(t, string(expectedResponseBody), responseBody, "Unexpected response JSON")
			},
		},
		{
			name:        "success - filter by price",
			queryParams: "?price=20",
			user:        models.User{ID: 1, Username: "Jane Doe"},
			mockReturn: []*models.Order{
				{
					ID:        2,
					ProductID: 2,
					Price:     29.99,
					OrderDate: time.Date(2025, 2, 10, 12, 0, 0, 0, time.UTC),
				},
				{
					ID:        3,
					ProductID: 3,
					Price:     19.99,
					OrderDate: time.Date(2025, 2, 10, 12, 0, 0, 0, time.UTC),
				},
			},
			mockError: nil,
			setServiceMock: func(mockOrdersService *mocks.OrdersService, user models.User, queryParams string, mockReturn []*models.Order, mockError error) {
				call := mockOrdersService.On("GetOrdersWithFilter", mock.Anything, user.ID, mock.AnythingOfType("func(*models.Order) bool"))
				call.Run(func(args mock.Arguments) {
					filter := args.Get(2).(func(order *models.Order) bool)
					assert.True(t, filter(mockReturn[0]), "Filter should match order with price 29.99")
					assert.False(t, filter(mockReturn[1]), "Filter should not match order with price 19.99")

					filteredOrders := make([]*models.Order, 0)
					for _, order := range mockReturn {
						if filter(order) {
							filteredOrders = append(filteredOrders, order)
						}
					}
					call.Return(filteredOrders, mockError)
				})
			},
			assertFunc: func(t *testing.T, responseBody string, responseCode int) {
				assert.Equal(t, fiber.StatusOK, responseCode, "Unexpected status code")

				expectedResponseBody, _ := json.Marshal([]interface{}{
					map[string]interface{}{
						"has_weightables": false,
						"id":              2,
						"order_date":      "2025-02-10T12:00:00Z",
						"price":           29.99,
						"product_id":      2},
				})
				assert.JSONEq(t, string(expectedResponseBody), responseBody, "Unexpected response JSON")
			},
		},
		{
			name:        "failure - invalid price filter",
			queryParams: "?price=abc",
			user:        models.User{ID: 1, Username: "Jane Doe"},
			mockReturn:  nil,
			mockError:   nil,
			setServiceMock: func(mockOrdersService *mocks.OrdersService, user models.User, queryParams string, mockReturn []*models.Order, mockError error) {
				// No service method is called for invalid price
			},
			assertFunc: func(t *testing.T, responseBody string, responseCode int) {
				assert.Equal(t, fiber.StatusBadRequest, responseCode, "Unexpected status code")
				assert.Contains(t, "Invalid price value", responseBody, "Unexpected response JSON")
			},
		},
		{
			name:        "failure - internal server error",
			queryParams: "",
			user:        models.User{ID: 1, Username: "Jane Doe"},
			mockReturn:  nil,
			mockError:   assert.AnError,
			setServiceMock: func(mockOrdersService *mocks.OrdersService, user models.User, queryParams string, mockReturn []*models.Order, mockError error) {
				mockOrdersService.On("GetOrders", mock.Anything, user.ID).Return(mockReturn, mockError)
			},
			assertFunc: func(t *testing.T, responseBody string, responseCode int) {
				assert.Equal(t, fiber.StatusInternalServerError, responseCode, "Unexpected status code")
				assert.Contains(t, "Error loading orders", responseBody, "Unexpected response JSON")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockOrdersService := new(mocks.OrdersService)
			tc.setServiceMock(mockOrdersService, tc.user, tc.queryParams, tc.mockReturn, tc.mockError)

			mockContextData := mocks.ProvideBaseMockContextData(&tc.user)
			app := createTestOrdersController(mockOrdersService, mockContextData)
			req := httptest.NewRequest(http.MethodGet, "/orders"+tc.queryParams, nil)
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)

			assert.Nil(t, err, "Handler should not return an error")

			var buf bytes.Buffer
			buf.ReadFrom(resp.Body)
			responseBody := buf.String()

			tc.assertFunc(t, responseBody, resp.StatusCode)

			mockOrdersService.AssertExpectations(t)
		})
	}
}

func TestGetOrder(t *testing.T) {
	tests := []struct {
		name             string
		orderID          string
		user             models.User
		mockReturn       *models.Order
		mockError        error
		setupServiceMock func(mockOrdersService *mocks.OrdersService, user models.User, orderID int, mockReturn *models.Order, mockError error)
		assertFunc       func(t *testing.T, responseBody string, responseCode int)
	}{
		{
			name:    "success - valid order ID",
			orderID: "1",
			user:    models.User{ID: 1, Username: "Jane Doe"},
			mockReturn: &models.Order{
				ID:        1,
				ProductID: 101,
				Price:     20.5,
				OrderDate: time.Date(2025, 1, 30, 10, 30, 0, 0, time.UTC),
			},
			mockError: nil,
			setupServiceMock: func(mockOrdersService *mocks.OrdersService, user models.User, orderID int, mockReturn *models.Order, mockError error) {
				mockOrdersService.On("GetOrder", mock.Anything, user.ID, orderID).Return(mockReturn, mockError)
			},
			assertFunc: func(t *testing.T, responseBody string, responseCode int) {
				assert.Equal(t, fiber.StatusOK, responseCode, "Unexpected status code")

				expectedResponseBody, _ := json.Marshal(map[string]interface{}{
					"id":              1,
					"product_id":      101,
					"price":           20.5,
					"order_date":      "2025-01-30T10:30:00Z",
					"has_weightables": false,
				})
				assert.JSONEq(t, string(expectedResponseBody), responseBody, "Unexpected response JSON")
			},
		},
		{
			name:       "failure - invalid order ID (non-numeric)",
			orderID:    "abc",
			user:       models.User{ID: 1, Username: "Jane Doe"},
			mockReturn: nil,
			mockError:  nil,
			setupServiceMock: func(mockOrdersService *mocks.OrdersService, user models.User, orderID int, mockReturn *models.Order, mockError error) {
			},
			assertFunc: func(t *testing.T, responseBody string, responseCode int) {
				assert.Equal(t, fiber.StatusBadRequest, responseCode, "Unexpected status code")
				assert.Contains(t, "strconv.Atoi: parsing \"abc\": invalid syntax", responseBody, "Unexpected response JSON")
			},
		},
		{
			name:       "failure - order not found",
			orderID:    "999",
			user:       models.User{ID: 1, Username: "John Doe"},
			mockReturn: nil,
			mockError:  assert.AnError,
			setupServiceMock: func(mockOrdersService *mocks.OrdersService, user models.User, orderID int, mockReturn *models.Order, mockError error) {
				mockOrdersService.On("GetOrder", mock.Anything, user.ID, orderID).Return(mockReturn, mockError)
			},
			assertFunc: func(t *testing.T, responseBody string, responseCode int) {
				assert.Equal(t, fiber.StatusInternalServerError, responseCode, "Unexpected status code")
				assert.Contains(t, "assert.AnError general error for testing", responseBody, "Unexpected response JSON")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockOrdersService := new(mocks.OrdersService)

			orderID, _ := strconv.Atoi(tc.orderID) // Handle invalid order ID
			tc.setupServiceMock(mockOrdersService, tc.user, orderID, tc.mockReturn, tc.mockError)

			mockContextData := mocks.ProvideBaseMockContextData(&tc.user)

			app := createTestOrdersController(mockOrdersService, mockContextData)
			req := httptest.NewRequest(http.MethodGet, "/orders/"+tc.orderID, nil)
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)

			assert.Nil(t, err, "Handler should not return an error")

			var buf bytes.Buffer
			buf.ReadFrom(resp.Body)
			responseBody := buf.String()

			tc.assertFunc(t, responseBody, resp.StatusCode)

			mockOrdersService.AssertExpectations(t)
		})
	}
}
