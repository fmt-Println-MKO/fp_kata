package models

import (
	"fp_kata/common"
	"fp_kata/internal/datasources/dsmodels"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPaymentToDSModel(t *testing.T) {
	validateNilResult := func(t *testing.T, result *dsmodels.Payment, _ *dsmodels.Payment) {
		assert.Nil(t, result, "The result should be nil")
	}

	validateNonNilResult := func(t *testing.T, result *dsmodels.Payment, expected *dsmodels.Payment) {
		assert.NotNil(t, result, "The result should not be nil")
		assert.Equal(t, expected.Id, result.Id, "Payment ID mismatch")
		assert.Equal(t, expected.Amount, result.Amount, "Payment Amount mismatch")
		assert.Equal(t, expected.Method, result.Method, "Payment Method mismatch")
		assert.Equal(t, expected.UserId, result.UserId, "UserId mismatch")
		assert.Equal(t, expected.OrderId, result.OrderId, "OrderId mismatch")
	}

	tests := []struct {
		name           string
		payment        Payment
		expectedResult *dsmodels.Payment
		validate       func(t *testing.T, result *dsmodels.Payment, expected *dsmodels.Payment)
	}{
		{
			name: "valid_full_data",
			payment: Payment{
				Id:     1,
				Amount: 100.50,
				Method: common.CreditCard,
				User:   &User{ID: 42},
				Order:  &Order{ID: 84},
			},
			expectedResult: &dsmodels.Payment{
				Id:      1,
				Amount:  100.50,
				Method:  common.CreditCard,
				UserId:  42,
				OrderId: 84,
			},
			validate: validateNonNilResult,
		},
		{
			name: "user_nil",
			payment: Payment{
				Id:     2,
				Amount: 200.75,
				Method: common.PayPal,
				User:   nil,
				Order:  &Order{ID: 105},
			},
			expectedResult: nil,
			validate:       validateNilResult,
		},
		{
			name: "order_nil",
			payment: Payment{
				Id:     3,
				Amount: 300.00,
				Method: common.DebitCard,
				User:   &User{ID: 56},
				Order:  nil,
			},
			expectedResult: nil,
			validate:       validateNilResult,
		},
		{
			name: "negative_amount",
			payment: Payment{
				Id:     4,
				Amount: -50.00,
				Method: common.DebitCard,
				User:   &User{ID: 77},
				Order:  &Order{ID: 88},
			},
			expectedResult: &dsmodels.Payment{
				Id:      4,
				Amount:  -50.00,
				Method:  common.DebitCard,
				UserId:  77,
				OrderId: 88,
			},
			validate: validateNonNilResult,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			result := testCase.payment.ToDSModel()
			testCase.validate(t, result, testCase.expectedResult)
		})
	}
}

func TestMapToPayment(t *testing.T) {
	validateResult := func(t *testing.T, result *Payment, expected *Payment) {
		assert.NotNil(t, result, "The result should not be nil")
		assert.Equal(t, expected.Id, result.Id, "Payment ID mismatch")
		assert.Equal(t, expected.Amount, result.Amount, "Payment Amount mismatch")
		assert.Equal(t, expected.Method, result.Method, "Payment Method mismatch")
		assert.Equal(t, expected.User, result.User, "User mismatch")
		assert.Equal(t, expected.Order, result.Order, "Order mismatch")
	}

	validateNilResult := func(t *testing.T, result *Payment, _ *Payment) {
		assert.Nil(t, result, "The result should be nil")
	}

	tests := []struct {
		name           string
		dsPayment      dsmodels.Payment
		user           *User
		order          *Order
		expectedResult *Payment
		validate       func(t *testing.T, result *Payment, expected *Payment)
	}{
		{
			name: "valid_full_data",
			dsPayment: dsmodels.Payment{
				Id:      1,
				Amount:  100.50,
				Method:  common.CreditCard,
				UserId:  42,
				OrderId: 84,
			},
			user:  &User{ID: 42},
			order: &Order{ID: 84},
			expectedResult: &Payment{
				Id:     1,
				Amount: 100.50,
				Method: common.CreditCard,
				User:   &User{ID: 42},
				Order:  &Order{ID: 84},
			},
			validate: validateResult,
		},
		{
			name: "missing_user",
			dsPayment: dsmodels.Payment{
				Id:      2,
				Amount:  200.75,
				Method:  common.PayPal,
				UserId:  0,
				OrderId: 105,
			},
			user:           nil,
			order:          &Order{ID: 105},
			expectedResult: nil,
			validate:       validateNilResult,
		},
		{
			name: "missing_order",
			dsPayment: dsmodels.Payment{
				Id:      3,
				Amount:  300.00,
				Method:  common.DebitCard,
				UserId:  56,
				OrderId: 0,
			},
			user:           &User{ID: 56},
			order:          nil,
			expectedResult: nil,
			validate:       validateNilResult,
		},
		{
			name: "negative_amount",
			dsPayment: dsmodels.Payment{
				Id:      4,
				Amount:  -50.00,
				Method:  common.DebitCard,
				UserId:  77,
				OrderId: 88,
			},
			user:  &User{ID: 77},
			order: &Order{ID: 88},
			expectedResult: &Payment{
				Id:     4,
				Amount: -50.00,
				Method: common.DebitCard,
				User:   &User{ID: 77},
				Order:  &Order{ID: 88},
			},
			validate: validateResult,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			result := MapToPayment(testCase.dsPayment, testCase.user, testCase.order)
			testCase.validate(t, result, testCase.expectedResult)
		})
	}
}
