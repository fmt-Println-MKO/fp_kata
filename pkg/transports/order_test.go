package transports

import (
	"fp_kata/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestToOrder(t *testing.T) {
	tests := []struct {
		name         string
		inputRequest OrderCreateRequest
		inputUser    models.User
		expected     *models.Order
		errorMessage string
	}{
		{
			name: "valid_order",
			inputRequest: OrderCreateRequest{
				ProductID: 101,
				Quantity:  2,
				Price:     100.50,
				OrderDate: time.Date(2025, 1, 30, 10, 30, 0, 0, time.UTC),
				Payments: []*PaymentRequest{
					{
						PaymentAmount: 100.50,
						PaymentMethod: "CreditCard",
					},
				},
				HasWeightables: true,
			},
			inputUser: models.User{
				ID:       1,
				Username: "testuser",
				Email:    "test@example.com",
			},
			expected: &models.Order{
				ProductID: 101,
				Quantity:  2,
				Price:     100.50,
				OrderDate: time.Date(2025, 1, 30, 10, 30, 0, 0, time.UTC),
				Payments: []*models.Payment{
					{
						Id:     0, // Set dynamically elsewhere
						Amount: 100.50,
						Method: "CreditCard",
						User: &models.User{
							ID:       1,
							Username: "testuser",
							Email:    "test@example.com",
						},
					},
				},
				User: &models.User{
					ID:       1,
					Username: "testuser",
					Email:    "test@example.com",
				},
				HasWeightables: true,
			},
			errorMessage: "Valid input should map correctly to a valid order",
		},
		{
			name: "no_payments",
			inputRequest: OrderCreateRequest{
				ProductID:      102,
				Quantity:       1,
				Price:          50.00,
				OrderDate:      time.Date(2025, 2, 10, 15, 0, 0, 0, time.UTC),
				Payments:       nil,
				HasWeightables: false,
			},
			inputUser: models.User{
				ID:       2,
				Username: "nopaymentsuser",
				Email:    "nopayments@example.com",
			},
			expected: &models.Order{
				ProductID: 102,
				Quantity:  1,
				Price:     50.00,
				OrderDate: time.Date(2025, 2, 10, 15, 0, 0, 0, time.UTC),
				Payments:  []*models.Payment{},
				User: &models.User{
					ID:       2,
					Username: "nopaymentsuser",
					Email:    "nopayments@example.com",
				},
				HasWeightables: false,
			},
			errorMessage: "Input without payments should set an empty payments slice",
		},
		{
			name: "zero_values",
			inputRequest: OrderCreateRequest{
				ProductID:      0,
				Quantity:       0,
				Price:          0.0,
				OrderDate:      time.Time{},
				Payments:       nil,
				HasWeightables: false,
			},
			inputUser: models.User{
				ID:       0,
				Username: "",
				Email:    "",
			},
			expected: &models.Order{
				ProductID: 0,
				Quantity:  0,
				Price:     0.0,
				OrderDate: time.Time{},
				Payments:  []*models.Payment{},
				User: &models.User{
					ID:       0,
					Username: "",
					Email:    "",
				},
				HasWeightables: false,
			},
			errorMessage: "Input with zero values should map correctly to order with defaults",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.inputRequest.ToOrder(tt.inputUser)
			assert.Equal(t, tt.expected, actual, tt.errorMessage)
		})
	}
}

func TestMapToOrderResponse(t *testing.T) {
	tests := []struct {
		name         string
		input        models.Order
		expected     *OrderResponse
		errorMessage string
	}{
		{
			name: "valid_order",
			input: models.Order{
				ID:        1,
				ProductID: 101,
				Quantity:  3,
				Price:     200.50,
				OrderDate: time.Date(2025, 1, 30, 10, 30, 0, 0, time.UTC),
				Payments: []*models.Payment{
					{
						Id:     1,
						Amount: 200.50,
						Method: "CreditCard",
					},
				},
				User: &models.User{
					ID:       1,
					Email:    "john.doe@example.com",
					Username: "johndoe",
				},
				HasWeightables: true,
			},
			expected: &OrderResponse{
				ID:        1,
				ProductID: 101,
				Quantity:  3,
				Price:     200.50,
				OrderDate: time.Date(2025, 1, 30, 10, 30, 0, 0, time.UTC),
				Payments: []*PaymentResponse{
					{
						Id:     1,
						Amount: 200.50,
						Method: "CreditCard",
					},
				},
				User: &UserResponse{
					ID:       1,
					Username: "johndoe",
					Email:    "john.doe@example.com",
				},
				HasWeightables: true,
			},
			errorMessage: "Expected correct mapping with all fields populated, but result differs",
		},
		{
			name: "order_with_nil_user",
			input: models.Order{
				ID:        2,
				ProductID: 102,
				Quantity:  1,
				Price:     100.00,
				OrderDate: time.Date(2025, 2, 10, 15, 0, 0, 0, time.UTC),
				Payments: []*models.Payment{
					{
						Id:     2,
						Amount: 100.00,
						Method: "PayPal",
					},
				},
				User:           nil,
				HasWeightables: false,
			},
			expected: &OrderResponse{
				ID:        2,
				ProductID: 102,
				Quantity:  1,
				Price:     100.00,
				OrderDate: time.Date(2025, 2, 10, 15, 0, 0, 0, time.UTC),
				Payments: []*PaymentResponse{
					{
						Id:     2,
						Amount: 100.00,
						Method: "PayPal",
					},
				},
				User:           nil,
				HasWeightables: false,
			},
			errorMessage: "Expected correct mapping with a nil user, but result differs",
		},
		{
			name: "empty_order",
			input: models.Order{
				ID:             0,
				ProductID:      0,
				Quantity:       0,
				Price:          0.0,
				OrderDate:      time.Time{},
				Payments:       nil,
				User:           nil,
				HasWeightables: false,
			},
			expected: &OrderResponse{
				ID:             0,
				ProductID:      0,
				Quantity:       0,
				Price:          0.0,
				OrderDate:      time.Time{},
				Payments:       []*PaymentResponse{},
				User:           nil,
				HasWeightables: false,
			},
			errorMessage: "Expected correct mapping with empty order, but result differs",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := MapToOrderResponse(tt.input)
			assert.Equal(t, tt.expected, actual, tt.errorMessage)
		})
	}
}

func TestConvertPayments(t *testing.T) {
	tests := []struct {
		name         string
		input        []*models.Payment
		expected     []*PaymentResponse
		errorMessage string
	}{
		{
			name:         "empty_input",
			input:        []*models.Payment{},
			expected:     []*PaymentResponse{},
			errorMessage: "Expected an empty slice for empty input, but got a different result",
		},
		{
			name: "single_payment",
			input: []*models.Payment{
				{
					Id:     123,
					Amount: 100.50,
					Method: "credit_card",
				},
			},
			expected: []*PaymentResponse{
				{
					Id:     123,
					Amount: 100.50,
					Method: "credit_card",
				},
			},
			errorMessage: "Expected single converted payment, but got a different result",
		},
		{
			name: "multiple_payments",
			input: []*models.Payment{
				{
					Id:     123,
					Amount: 100.00,
					Method: "credit_card",
				},
				{
					Id:     456,
					Amount: 200.80,
					Method: "paypal",
				},
			},
			expected: []*PaymentResponse{
				{
					Id:     123,
					Amount: 100.00,
					Method: "credit_card",
				},
				{
					Id:     456,
					Amount: 200.80,
					Method: "paypal",
				},
			},
			errorMessage: "Expected multiple converted payments, but got a different result",
		},
		{
			name:         "nil_input",
			input:        nil,
			expected:     []*PaymentResponse{},
			errorMessage: "Expected an empty slice for nil input, but got a different result",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := convertPayments(tt.input)
			assert.Equal(t, tt.expected, actual, tt.errorMessage)
		})
	}
}
