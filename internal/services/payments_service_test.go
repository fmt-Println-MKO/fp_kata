package services

import (
	"errors"
	"fp_kata/internal/datasources/dsmodels"
	"fp_kata/internal/models"
	"fp_kata/mocks"
	"fp_kata/pkg/log"
	zlog "github.com/rs/zerolog/log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func validateSuccess(expected *models.Payment) func(*testing.T, *models.Payment, error) {
	return func(t *testing.T, result *models.Payment, err error) {
		assert.NoError(t, err, "Expected no error but got one")
		assert.NotNil(t, result, "Expected a valid result but got nil")
		assert.Equal(t, expected.Id, result.Id, "Payment ID mismatch")
		assert.Equal(t, expected.Amount, result.Amount, "Payment amount mismatch")
		assert.Equal(t, expected.User.ID, result.User.ID, "User ID mismatch")
		assert.Equal(t, expected.Order.ID, result.Order.ID, "Order ID mismatch")
	}
}

func validateError(expectedError string) func(*testing.T, *models.Payment, error) {
	return func(t *testing.T, result *models.Payment, err error) {
		assert.Error(t, err, "Expected an error but got none")
		assert.EqualError(t, err, expectedError, "Error message mismatch")
		assert.Nil(t, result, "Expected result to be nil when an error occurs")
	}
}

func TestGetPaymentsByOrder(t *testing.T) {
	log.InitLogger()
	ctx := log.NewBackgroundContext(&zlog.Logger)

	tests := []struct {
		name      string
		orderID   int
		mockSetup func(*mocks.PaymentsDatasource)
		validate  func(*testing.T, []*models.Payment, error)
	}{
		{
			name:    "Payments Found for Order",
			orderID: 1,
			mockSetup: func(mockStorage *mocks.PaymentsDatasource) {
				mockStorage.On("AllByOrderId", mock.Anything, 1).Return([]dsmodels.Payment{
					{
						Id:      101,
						Amount:  10.50,
						UserId:  1,
						OrderId: 1,
					},
					{
						Id:      102,
						Amount:  20.00,
						UserId:  2,
						OrderId: 1,
					},
				}, nil)
			},
			validate: func(t *testing.T, result []*models.Payment, err error) {
				assert.NoError(t, err, "Expected no error but got one")
				assert.NotNil(t, result, "Expected a valid result but got nil")

				assert.Len(t, result, 2, "Expected 2 payments but got a different number")
				assert.Equal(t, 101, result[0].Id, "First payment ID mismatch")
				assert.Equal(t, 10.50, result[0].Amount, "First payment amount mismatch")
				assert.Equal(t, 102, result[1].Id, "Second payment ID mismatch")
				assert.Equal(t, 20.00, result[1].Amount, "Second payment amount mismatch")
			},
		},
		{
			name:    "No Payments for Order",
			orderID: 2,
			mockSetup: func(mockStorage *mocks.PaymentsDatasource) {
				mockStorage.On("AllByOrderId", mock.Anything, 2).Return([]dsmodels.Payment{}, nil)
			},
			validate: func(t *testing.T, result []*models.Payment, err error) {
				assert.NoError(t, err, "Expected no error but got one")
				assert.NotNil(t, result, "Expected a valid result but got nil")
				assert.Empty(t, result, "Expected no payments but got some")
			},
		},
		{
			name:    "Error Fetching Payments",
			orderID: 3,
			mockSetup: func(mockStorage *mocks.PaymentsDatasource) {
				mockStorage.On("AllByOrderId", mock.Anything, 3).Return(nil, errors.New("fetch error"))
			},
			validate: func(t *testing.T, result []*models.Payment, err error) {
				assert.Error(t, err, "Expected an error but got none")
				assert.EqualError(t, err, "fetch error", "Error message mismatch")
				assert.Nil(t, result, "Expected result to be nil on error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := mocks.NewPaymentsDatasource(t)
			if tt.mockSetup != nil {
				tt.mockSetup(mockStorage)
			}

			service := NewPaymentsService(mockStorage)
			result, err := service.GetPaymentsByOrder(ctx, tt.orderID)

			tt.validate(t, result, err)

			mockStorage.AssertExpectations(t)
		})
	}
}

func TestStorePayment(t *testing.T) {

	log.InitLogger()
	ctx := log.NewBackgroundContext(&zlog.Logger)

	tests := []struct {
		name      string
		payment   models.Payment
		mockSetup func(*mocks.PaymentsDatasource)
		validate  func(*testing.T, *models.Payment, error)
	}{
		{
			name: "Successful Payment Creation",
			payment: models.Payment{
				Id:     0,
				Amount: 100.50,
				User:   &models.User{ID: 1},
				Order:  &models.Order{ID: 1},
			},
			mockSetup: func(mockStorage *mocks.PaymentsDatasource) {
				mockStorage.On("Create", mock.Anything, mock.MatchedBy(func(p dsmodels.Payment) bool {
					return p.Amount == 100.50 && p.UserId == 1 && p.OrderId == 1
				})).Return(dsmodels.Payment{
					Id:      1,
					Amount:  100.50,
					UserId:  1,
					OrderId: 1,
				}, nil)
			},
			validate: validateSuccess(&models.Payment{
				Id:     1,
				Amount: 100.50,
				User:   &models.User{ID: 1},
				Order:  &models.Order{ID: 1},
			}),
		},
		{
			name: "Successful Payment Update",
			payment: models.Payment{
				Id:     1,
				Amount: 200.75,
				User:   &models.User{ID: 2},
				Order:  &models.Order{ID: 3},
			},
			mockSetup: func(mockStorage *mocks.PaymentsDatasource) {
				mockStorage.On("Update", mock.Anything, mock.MatchedBy(func(p dsmodels.Payment) bool {
					return p.Id == 1 && p.Amount == 200.75 && p.UserId == 2 && p.OrderId == 3
				})).Return(dsmodels.Payment{
					Id:      1,
					Amount:  200.75,
					UserId:  2,
					OrderId: 3,
				}, nil)
			},
			validate: validateSuccess(&models.Payment{
				Id:     1,
				Amount: 200.75,
				User:   &models.User{ID: 2},
				Order:  &models.Order{ID: 3},
			}),
		},
		{
			name: "Creation Fails",
			payment: models.Payment{
				Id:     0,
				Amount: 300.00,
				User:   &models.User{ID: 3},
				Order:  &models.Order{ID: 2},
			},
			mockSetup: func(mockStorage *mocks.PaymentsDatasource) {
				mockStorage.On("Create", mock.Anything, mock.Anything).Return(dsmodels.Payment{}, errors.New("creation error"))
			},
			validate: validateError("creation error"),
		},
		{
			name: "Update Fails",
			payment: models.Payment{
				Id:     5,
				Amount: 150.00,
				User:   &models.User{ID: 4},
				Order:  &models.Order{ID: 5},
			},
			mockSetup: func(mockStorage *mocks.PaymentsDatasource) {
				mockStorage.On("Update", mock.Anything, mock.Anything).Return(dsmodels.Payment{}, errors.New("update error"))
			},
			validate: validateError("update error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := mocks.NewPaymentsDatasource(t)
			if tt.mockSetup != nil {
				tt.mockSetup(mockStorage)
			}

			service := NewPaymentsService(mockStorage)
			result, err := service.StorePayment(ctx, tt.payment)

			tt.validate(t, result, err)

			mockStorage.AssertExpectations(t)
		})
	}
}

func TestGetPaymentByID(t *testing.T) {
	log.InitLogger()
	ctx := log.NewBackgroundContext(&zlog.Logger)

	tests := []struct {
		name      string
		id        int
		mockSetup func(*mocks.PaymentsDatasource)
		validate  func(*testing.T, *models.Payment, error)
	}{
		{
			name: "Payment Found",
			id:   1,
			mockSetup: func(mockStorage *mocks.PaymentsDatasource) {
				mockStorage.On("Read", mock.Anything, 1).Return(dsmodels.Payment{
					Id:      1,
					Amount:  100.50,
					UserId:  1,
					OrderId: 1,
				}, nil)
			},
			validate: validateSuccess(&models.Payment{
				Id:     1,
				Amount: 100.50,
				User:   &models.User{ID: 1},
				Order:  &models.Order{ID: 1},
			}),
		},
		{
			name: "Payment Not Found",
			id:   2,
			mockSetup: func(mockStorage *mocks.PaymentsDatasource) {
				mockStorage.On("Read", mock.Anything, 2).Return(dsmodels.Payment{}, errors.New("payment not found"))
			},
			validate: validateError("payment not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := mocks.NewPaymentsDatasource(t)
			if tt.mockSetup != nil {
				tt.mockSetup(mockStorage)
			}

			service := NewPaymentsService(mockStorage)
			result, err := service.GetPaymentByID(ctx, tt.id)

			tt.validate(t, result, err)
		})
	}
}
