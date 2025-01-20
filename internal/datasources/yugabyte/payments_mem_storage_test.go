package yugabyte

import (
	"context"
	"fmt"
	"fp_kata/common"
	zlog "github.com/rs/zerolog/log"
	"testing"

	"fp_kata/internal/datasources/dsmodels"
	"fp_kata/pkg/log"
	"github.com/stretchr/testify/assert"
)

func initTestPaymentsStorage(store map[int]dsmodels.Payment) (*inMemoryPaymentsStorage, context.Context) {
	log.InitLogger()
	ctx := log.NewBackgroundContext(&zlog.Logger)
	return &inMemoryPaymentsStorage{
		payments: store,
	}, ctx
}

// Helper function to quickly create payment objects.
func createPayment(id int, amount float64, method common.PaymentMethod, userID, orderID int) dsmodels.Payment {
	return dsmodels.Payment{
		Id:      id,
		Amount:  amount,
		Method:  method,
		UserId:  userID,
		OrderId: orderID,
	}
}

// Helper function to initialize the map of payments.
func createPaymentsMap(payments ...dsmodels.Payment) map[int]dsmodels.Payment {
	result := make(map[int]dsmodels.Payment)
	for _, payment := range payments {
		result[payment.Id] = payment
	}
	return result
}

func TestInMemoryPaymentsStorage_Create(t *testing.T) {
	const initialPaymentID = 1 // Common starting point for new payment IDs.

	type PaymentTestCase struct {
		name            string
		payment         dsmodels.Payment
		expected        dsmodels.Payment
		initialPayments map[int]dsmodels.Payment
		assert          func(t *testing.T, result dsmodels.Payment, err error, storage *inMemoryPaymentsStorage)
	}

	tests := []PaymentTestCase{
		{
			name:            "valid payment creation",
			payment:         createPayment(0, 100.0, common.CreditCard, 1, 123), // ID is 0 initially before assignment.
			expected:        createPayment(initialPaymentID, 100.0, common.CreditCard, 1, 123),
			initialPayments: createPaymentsMap(),
			assert: func(t *testing.T, result dsmodels.Payment, err error, storage *inMemoryPaymentsStorage) {
				assert.NoError(t, err, "unexpected error during valid payment creation")
				assert.Equal(t, initialPaymentID, result.Id, "unexpected id assigned to payment")
				assert.Equal(t, 1, len(storage.payments), "payment storage should contain exactly one payment")
				assert.Equal(t, result, storage.payments[result.Id], "stored payment does not match expected payment")
			},
		},
		{
			name:            "valid creation for second payment",
			payment:         createPayment(0, 200.0, common.PayPal, 2, 456),
			expected:        createPayment(initialPaymentID+1, 200.0, common.PayPal, 2, 456),
			initialPayments: createPaymentsMap(createPayment(initialPaymentID, 100.0, common.CreditCard, 1, 123)),
			assert: func(t *testing.T, result dsmodels.Payment, err error, storage *inMemoryPaymentsStorage) {
				assert.NoError(t, err, "unexpected error during valid payment creation")
				assert.Equal(t, initialPaymentID+1, result.Id, "unexpected id assigned to second payment")
				assert.Equal(t, 2, len(storage.payments), "payment storage should contain exactly two payments")
				assert.Equal(t, result, storage.payments[result.Id], "stored payment does not match expected payment")
			},
		},
		{
			name:     "creation with existing payment results in new ID",
			payment:  createPayment(0, 300.0, common.BankTransfer, 3, 789),
			expected: createPayment(initialPaymentID+2, 300.0, common.BankTransfer, 3, 789),
			initialPayments: createPaymentsMap(
				createPayment(initialPaymentID, 100.0, common.CreditCard, 1, 123),
				createPayment(initialPaymentID+1, 200.0, common.PayPal, 2, 456),
			),
			assert: func(t *testing.T, result dsmodels.Payment, err error, storage *inMemoryPaymentsStorage) {
				assert.NoError(t, err, "unexpected error during valid payment creation")
				assert.Equal(t, initialPaymentID+2, result.Id, "unexpected id assigned to third payment")
				assert.Equal(t, 3, len(storage.payments), "payment storage should contain exactly three payments")
				assert.Equal(t, result, storage.payments[result.Id], "stored payment does not match expected payment")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			storage, ctx := initTestPaymentsStorage(tc.initialPayments)
			result, err := storage.Create(ctx, tc.payment)
			tc.assert(t, result, err, storage)
		})
	}
}

func TestInMemoryPaymentsStorage_Read(t *testing.T) {
	type ReadPaymentTestCase struct {
		name            string
		paymentID       int
		initialPayments map[int]dsmodels.Payment
		expected        dsmodels.Payment
		expectedErr     error
		assert          func(t *testing.T, result dsmodels.Payment, err error)
	}

	tests := []ReadPaymentTestCase{
		{
			name:      "valid payment read",
			paymentID: 1,
			initialPayments: createPaymentsMap(
				createPayment(1, 100.0, common.CreditCard, 1, 101),
				createPayment(2, 200.0, common.PayPal, 2, 102),
			),
			expected:    createPayment(1, 100.0, common.CreditCard, 1, 101),
			expectedErr: nil,
			assert: func(t *testing.T, result dsmodels.Payment, err error) {
				assert.NoError(t, err, "unexpected error when reading existing payment")
				assert.Equal(t, 1, result.Id, "unexpected payment ID")
				assert.Equal(t, common.CreditCard, result.Method, "unexpected payment method")
			},
		},
		{
			name:            "payment not found",
			paymentID:       99,
			initialPayments: createPaymentsMap(createPayment(1, 100.0, common.CreditCard, 1, 101)),
			expected:        dsmodels.Payment{},
			expectedErr:     fmt.Errorf("payment with id 99 not found"),
			assert: func(t *testing.T, result dsmodels.Payment, err error) {
				assert.Error(t, err, "expected an error when reading non-existent payment")
				assert.Equal(t, "payment with id 99 not found", err.Error(), "unexpected error message")
				assert.Equal(t, dsmodels.Payment{}, result, "expected empty payment object for non-existent ID")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			storage, ctx := initTestPaymentsStorage(tc.initialPayments)
			result, err := storage.Read(ctx, tc.paymentID)
			tc.assert(t, result, err)
		})
	}
}

func TestInMemoryPaymentsStorage_Update(t *testing.T) {
	type UpdatePaymentTestCase struct {
		name            string
		updatePayment   dsmodels.Payment
		initialPayments map[int]dsmodels.Payment
		expected        dsmodels.Payment
		expectedErr     error
		assert          func(t *testing.T, result dsmodels.Payment, err error, storage *inMemoryPaymentsStorage)
	}

	tests := []UpdatePaymentTestCase{
		{
			name:          "update existing payment",
			updatePayment: createPayment(1, 150.0, common.CreditCard, 1, 101),
			initialPayments: createPaymentsMap(
				createPayment(1, 100.0, common.CreditCard, 1, 101),
			),
			expected:    createPayment(1, 150.0, common.CreditCard, 1, 101),
			expectedErr: nil,
			assert: func(t *testing.T, result dsmodels.Payment, err error, storage *inMemoryPaymentsStorage) {
				assert.NoError(t, err, "error should be nil for valid update")
				assert.Equal(t, result, storage.payments[result.Id], "updated payment does not match expected payment")
				assert.Equal(t, result.Amount, 150.0, "updated payment amount mismatch")
			},
		},
		{
			name:          "update non-existent payment",
			updatePayment: createPayment(99, 150.0, common.CreditCard, 1, 101),
			initialPayments: createPaymentsMap(
				createPayment(1, 100.0, common.CreditCard, 1, 101),
			),
			expected:    dsmodels.Payment{},
			expectedErr: fmt.Errorf("payment with id 99 not found"),
			assert: func(t *testing.T, result dsmodels.Payment, err error, storage *inMemoryPaymentsStorage) {
				assert.Error(t, err, "expected error for updating non-existent payment")
				assert.Equal(t, "payment with id 99 not found", err.Error(), "error message mismatch")
				assert.Equal(t, dsmodels.Payment{}, result, "result should be an empty payment object")
			},
		},
		{
			name:          "update payment with modified attributes",
			updatePayment: createPayment(1, 200.0, common.PayPal, 1, 101),
			initialPayments: createPaymentsMap(
				createPayment(1, 150.0, common.CreditCard, 1, 101),
			),
			expected:    createPayment(1, 200.0, common.PayPal, 1, 101),
			expectedErr: nil,
			assert: func(t *testing.T, result dsmodels.Payment, err error, storage *inMemoryPaymentsStorage) {
				assert.NoError(t, err, "error should be nil for valid update with modified attributes")
				assert.Equal(t, result, storage.payments[result.Id], "updated payment does not match expected payment")
				assert.Equal(t, result.Method, common.PayPal, "updated payment method mismatch")
				assert.Equal(t, result.Amount, 200.0, "updated payment amount mismatch")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			storage, ctx := initTestPaymentsStorage(tc.initialPayments)
			result, err := storage.Update(ctx, tc.updatePayment)
			tc.assert(t, result, err, storage)
		})
	}
}

func TestInMemoryPaymentsStorage_Delete(t *testing.T) {
	type DeletePaymentTestCase struct {
		name            string
		paymentID       int
		initialPayments map[int]dsmodels.Payment
		expectedErr     error
		assert          func(t *testing.T, storage *inMemoryPaymentsStorage, err error)
	}

	tests := []DeletePaymentTestCase{
		{
			name:      "delete existing payment",
			paymentID: 1,
			initialPayments: createPaymentsMap(
				createPayment(1, 100.0, common.CreditCard, 1, 101),
				createPayment(2, 200.0, common.PayPal, 2, 102),
			),
			expectedErr: nil,
			assert: func(t *testing.T, storage *inMemoryPaymentsStorage, err error) {
				assert.NoError(t, err, "unexpected error while deleting an existing payment")
				_, exists := storage.payments[1]
				assert.False(t, exists, "payment with ID 1 should have been deleted")
				assert.Equal(t, 1, len(storage.payments), "remaining payments count mismatch after delete")
			},
		},
		{
			name:      "delete non-existent payment",
			paymentID: 99,
			initialPayments: createPaymentsMap(
				createPayment(1, 100.0, common.CreditCard, 1, 101),
			),
			expectedErr: fmt.Errorf("payment with id 99 not found"),
			assert: func(t *testing.T, storage *inMemoryPaymentsStorage, err error) {
				assert.Error(t, err, "expected error while deleting a non-existent payment")
				assert.Equal(t, "payment with id 99 not found", err.Error(), "unexpected error message")
				assert.Equal(t, 1, len(storage.payments), "payments count mismatch when deleting non-existent payment")
			},
		},
		{
			name:      "delete last remaining payment",
			paymentID: 1,
			initialPayments: createPaymentsMap(
				createPayment(1, 100.0, common.CreditCard, 1, 101),
			),
			expectedErr: nil,
			assert: func(t *testing.T, storage *inMemoryPaymentsStorage, err error) {
				assert.NoError(t, err, "unexpected error while deleting the last remaining payment")
				assert.Empty(t, storage.payments, "storage should be empty after last payment is deleted")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			storage, ctx := initTestPaymentsStorage(tc.initialPayments)
			err := storage.Delete(ctx, tc.paymentID)
			tc.assert(t, storage, err)
		})
	}
}

func TestInMemoryPaymentsStorage_AllByOrderId(t *testing.T) {
	type AllByOrderIdTestCase struct {
		name            string
		orderID         int
		initialPayments map[int]dsmodels.Payment
		expected        []dsmodels.Payment
		expectedErr     error
		assert          func(t *testing.T, result []dsmodels.Payment, err error)
	}

	tests := []AllByOrderIdTestCase{
		{
			name:    "payments found for valid orderId",
			orderID: 101,
			initialPayments: createPaymentsMap(
				createPayment(1, 100.0, common.CreditCard, 1, 101),
				createPayment(2, 200.0, common.PayPal, 2, 101),
				createPayment(3, 300.0, common.BankTransfer, 3, 102),
			),
			expected: []dsmodels.Payment{
				createPayment(1, 100.0, common.CreditCard, 1, 101),
				createPayment(2, 200.0, common.PayPal, 2, 101),
			},
			expectedErr: nil,
			assert: func(t *testing.T, result []dsmodels.Payment, err error) {
				assert.NoError(t, err, "unexpected error when retrieving payments by orderId")
				assert.Len(t, result, 2, "expected exactly 2 payments for orderId 101")
				assert.Equal(t, 1, result[0].Id, "first payment ID mismatch")
				assert.Equal(t, 2, result[1].Id, "second payment ID mismatch")
			},
		},
		{
			name:    "no payments found for orderId with no matches",
			orderID: 999,
			initialPayments: createPaymentsMap(
				createPayment(1, 100.0, common.CreditCard, 1, 101),
			),
			expected:    nil,
			expectedErr: fmt.Errorf("no payments found with orderId 999"),
			assert: func(t *testing.T, result []dsmodels.Payment, err error) {
				assert.Error(t, err, "expected an error when no payments are found")
				assert.Equal(t, "no payments found with orderId 999", err.Error(), "unexpected error message")
				assert.Nil(t, result, "result should be nil when no payments are found")
			},
		},
		{
			name:            "no payments in storage",
			orderID:         101,
			initialPayments: createPaymentsMap(
			// No payments initialized
			),
			expected:    nil,
			expectedErr: fmt.Errorf("no payments found with orderId 101"),
			assert: func(t *testing.T, result []dsmodels.Payment, err error) {
				assert.Error(t, err, "expected an error when storage is empty")
				assert.Equal(t, "no payments found with orderId 101", err.Error(), "unexpected error message")
				assert.Nil(t, result, "result should be nil when storage is empty")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			storage, ctx := initTestPaymentsStorage(tc.initialPayments)
			result, err := storage.AllByOrderId(ctx, tc.orderID)
			tc.assert(t, result, err)
		})
	}
}

func TestNewPaymentsStorage(t *testing.T) {
	type NewPaymentsStorageTestCase struct {
		name string
	}

	tests := []NewPaymentsStorageTestCase{
		{
			name: "storage initialization",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := NewPaymentsStorage()
			storage, ok := result.(*inMemoryPaymentsStorage)
			assert.True(t, ok, "expected result to be of type *inMemoryPaymentsStorage")
			assert.NotNil(t, storage, "storage instance should not be nil")
			assert.NotNil(t, storage.payments, "payments map should be initialized")
			assert.Empty(t, storage.payments, "new storage should be empty initially")
		})
	}
}
