package transports

import (
	"fp_kata/common"
	"fp_kata/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapToPaymentResponse(t *testing.T) {
	tests := []struct {
		name   string
		input  models.Payment
		expect *PaymentResponse
	}{
		{
			name: "valid payment with user and order",
			input: models.Payment{
				Id:     1,
				Amount: 100.50,
				Method: common.CreditCard,
				User:   &models.User{ID: 10},
				Order:  &models.Order{ID: 20},
			},
			expect: &PaymentResponse{
				Id:     1,
				Amount: 100.50,
				Method: common.CreditCard,
			},
		},
		{
			name: "valid payment without user or order",
			input: models.Payment{
				Id:     2,
				Amount: 50.25,
				Method: common.PayPal,
				User:   nil,
				Order:  nil,
			},
			expect: &PaymentResponse{
				Id:     2,
				Amount: 50.25,
				Method: common.PayPal,
			},
		},
		{
			name: "valid payment with zero amount",
			input: models.Payment{
				Id:     3,
				Amount: 0,
				Method: common.BankTransfer,
			},
			expect: &PaymentResponse{
				Id:     3,
				Amount: 0,
				Method: common.BankTransfer,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := MapToPaymentResponse(test.input)

			assert.Equal(t, test.expect.Id, result.Id, "Id mismatch")
			assert.Equal(t, test.expect.Amount, result.Amount, "Amount mismatch")
			assert.Equal(t, test.expect.Method, result.Method, "Method mismatch")
		})
	}
}

func TestToPayment(t *testing.T) {
	tests := []struct {
		name          string
		request       PaymentRequest
		user          models.User
		expected      *models.Payment
		assertMessage string
	}{
		{
			name: "valid payment request",
			request: PaymentRequest{
				PaymentAmount: 100.50,
				PaymentMethod: common.CreditCard,
			},
			user: models.User{ID: 1},
			expected: &models.Payment{
				Amount: 100.50,
				Method: common.CreditCard,
				User:   &models.User{ID: 1},
			},
			assertMessage: "Valid payment conversion failed",
		},
		{
			name: "default values in request",
			request: PaymentRequest{
				PaymentAmount: 0,
				PaymentMethod: common.PayPal,
			},
			user: models.User{},
			expected: &models.Payment{
				Amount: 0,
				Method: common.PayPal,
				User:   &models.User{},
			},
			assertMessage: "Default values payment conversion failed",
		},
		{
			name: "zero payment amount",
			request: PaymentRequest{
				PaymentAmount: 0,
				PaymentMethod: common.DebitCard,
			},
			user: models.User{ID: 2},
			expected: &models.Payment{
				Amount: 0,
				Method: common.DebitCard,
				User:   &models.User{ID: 2},
			},
			assertMessage: "Zero payment amount conversion failed",
		},
		{
			name: "large payment amount",
			request: PaymentRequest{
				PaymentAmount: 1e9,
				PaymentMethod: common.BankTransfer,
			},
			user: models.User{ID: 3},
			expected: &models.Payment{
				Amount: 1e9,
				Method: common.BankTransfer,
				User:   &models.User{ID: 3},
			},
			assertMessage: "Large payment amount conversion failed",
		},
		{
			name: "empty user in request",
			request: PaymentRequest{
				PaymentAmount: 51.25,
				PaymentMethod: common.PayPal,
			},
			user: models.User{},
			expected: &models.Payment{
				Amount: 51.25,
				Method: common.PayPal,
				User:   &models.User{},
			},
			assertMessage: "Empty user conversion failed",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.request.ToPayment(test.user)

			assert.Equal(t, test.expected.Amount, result.Amount, test.assertMessage+": Amount mismatch")
			assert.Equal(t, test.expected.Method, result.Method, test.assertMessage+": Method mismatch")
			assert.Equal(t, test.expected.User, result.User, test.assertMessage+": User mismatch")
		})
	}
}
