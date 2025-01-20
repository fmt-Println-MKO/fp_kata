package models

import (
	"fp_kata/internal/datasources/dsmodels"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOrderToDSModel(t *testing.T) {
	tests := []struct {
		name            string
		order           func() *Order
		expectedModel   func() *dsmodels.Order
		expectedMessage string
	}{
		{
			name: "valid_order_single_payment",
			order: func() *Order {
				return &Order{
					ID:        1,
					ProductID: 101,
					Quantity:  2,
					Price:     15.5,
					OrderDate: time.Date(2023, 10, 10, 12, 0, 0, 0, time.UTC),
					Payments: []*Payment{
						{Id: 201},
					},
					User:           &User{ID: 301},
					HasWeightables: false,
				}
			},
			expectedModel: func() *dsmodels.Order {
				return &dsmodels.Order{
					ID:             1,
					ProductID:      101,
					Quantity:       2,
					Price:          15.5,
					OrderDate:      time.Date(2023, 10, 10, 12, 0, 0, 0, time.UTC),
					Payments:       []int{201},
					UserId:         301,
					HasWeightables: false,
				}
			},
			expectedMessage: "valid order with single payment failed",
		},
		{
			name: "valid_order_multiple_payments",
			order: func() *Order {
				return &Order{
					ID:        2,
					ProductID: 102,
					Quantity:  5,
					Price:     50.0,
					OrderDate: time.Date(2023, 8, 1, 8, 0, 0, 0, time.UTC),
					Payments: []*Payment{
						{Id: 202},
						{Id: 203},
					},
					User:           &User{ID: 302},
					HasWeightables: true,
				}
			},
			expectedModel: func() *dsmodels.Order {
				return &dsmodels.Order{
					ID:             2,
					ProductID:      102,
					Quantity:       5,
					Price:          50.0,
					OrderDate:      time.Date(2023, 8, 1, 8, 0, 0, 0, time.UTC),
					Payments:       []int{202, 203},
					UserId:         302,
					HasWeightables: true,
				}
			},
			expectedMessage: "valid order with multiple payments failed",
		},
		{
			name: "valid_order_no_payments",
			order: func() *Order {
				return &Order{
					ID:        3,
					ProductID: 103,
					Quantity:  1,
					Price:     25.0,
					OrderDate: time.Date(2023, 5, 15, 15, 0, 0, 0, time.UTC),
					Payments:  []*Payment{},
					User:      &User{ID: 303},
				}
			},
			expectedModel: func() *dsmodels.Order {
				return &dsmodels.Order{
					ID:             3,
					ProductID:      103,
					Quantity:       1,
					Price:          25.0,
					OrderDate:      time.Date(2023, 5, 15, 15, 0, 0, 0, time.UTC),
					Payments:       []int{},
					UserId:         303,
					HasWeightables: false,
				}
			},
			expectedMessage: "valid order with no payments failed",
		},
		{
			name: "nil_user",
			order: func() *Order {
				return &Order{
					ID:        4,
					ProductID: 104,
					Quantity:  3,
					Price:     75.0,
					OrderDate: time.Date(2023, 7, 20, 0, 0, 0, 0, time.UTC),
					Payments: []*Payment{
						{Id: 204},
					},
					User: nil,
				}
			},
			expectedModel: func() *dsmodels.Order {
				return nil
			},

			expectedMessage: "order with nil user failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			order := tt.order()
			expectedModel := tt.expectedModel()

			actualModel := order.ToDSModel()

			assert.Equal(t, expectedModel, actualModel, tt.expectedMessage)
		})
	}
}

func TestMapToOrder(t *testing.T) {
	tests := []struct {
		name            string
		inputOrder      func() dsmodels.Order
		inputPayments   func() []*Payment
		inputUser       func() *User
		expectedOutput  func() *Order
		expectedMessage string
	}{
		{
			name: "valid_order_with_all_fields",
			inputOrder: func() dsmodels.Order {
				return dsmodels.Order{
					ID:             1,
					ProductID:      101,
					Quantity:       2,
					Price:          15.5,
					OrderDate:      time.Date(2023, 10, 10, 12, 0, 0, 0, time.UTC),
					Payments:       []int{201},
					UserId:         301,
					HasWeightables: false,
				}
			},
			inputPayments: func() []*Payment {
				return []*Payment{{Id: 201}}
			},
			inputUser: func() *User {
				return &User{ID: 301}
			},
			expectedOutput: func() *Order {
				return &Order{
					ID:             1,
					ProductID:      101,
					Quantity:       2,
					Price:          15.5,
					OrderDate:      time.Date(2023, 10, 10, 12, 0, 0, 0, time.UTC),
					Payments:       []*Payment{{Id: 201}},
					User:           &User{ID: 301},
					HasWeightables: false,
				}
			},
			expectedMessage: "valid order with all fields failed",
		},
		{
			name: "nil_user",
			inputOrder: func() dsmodels.Order {
				return dsmodels.Order{
					ID:             2,
					ProductID:      102,
					Quantity:       3,
					Price:          20.0,
					OrderDate:      time.Date(2023, 11, 5, 0, 0, 0, 0, time.UTC),
					Payments:       []int{202},
					UserId:         302,
					HasWeightables: true,
				}
			},
			inputPayments: func() []*Payment {
				return []*Payment{{Id: 202}}
			},
			inputUser: func() *User {
				return nil
			},
			expectedOutput: func() *Order {
				return &Order{
					ID:             2,
					ProductID:      102,
					Quantity:       3,
					Price:          20.0,
					OrderDate:      time.Date(2023, 11, 5, 0, 0, 0, 0, time.UTC),
					Payments:       []*Payment{{Id: 202}},
					User:           nil,
					HasWeightables: true,
				}
			},
			expectedMessage: "order with nil user failed",
		},
		{
			name: "nil_payments",
			inputOrder: func() dsmodels.Order {
				return dsmodels.Order{
					ID:             3,
					ProductID:      103,
					Quantity:       1,
					Price:          30.0,
					OrderDate:      time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC),
					Payments:       nil,
					UserId:         303,
					HasWeightables: false,
				}
			},
			inputPayments: func() []*Payment {
				return nil
			},
			inputUser: func() *User {
				return &User{ID: 303}
			},
			expectedOutput: func() *Order {
				return &Order{
					ID:             3,
					ProductID:      103,
					Quantity:       1,
					Price:          30.0,
					OrderDate:      time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC),
					Payments:       nil,
					User:           &User{ID: 303},
					HasWeightables: false,
				}
			},
			expectedMessage: "order with nil payments failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputOrder := tt.inputOrder()
			inputPayments := tt.inputPayments()
			inputUser := tt.inputUser()
			expectedOutput := tt.expectedOutput()

			actualOutput := MapToOrder(inputOrder, inputPayments, inputUser)

			assert.Equal(t, expectedOutput, actualOutput, tt.expectedMessage)
		})
	}
}
