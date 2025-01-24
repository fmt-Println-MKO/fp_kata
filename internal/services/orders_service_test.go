package services

import (
	"errors"
	"fp_kata/internal/datasources/dsmodels"
	"fp_kata/pkg/log"
	zlog "github.com/rs/zerolog/log"
	"github.com/stretchr/testify/mock"
	"testing"

	"fp_kata/internal/models"
	"fp_kata/mocks"

	"github.com/stretchr/testify/assert"
)

func TestOrderService_StoreOrder(t *testing.T) {
	log.InitLogger()
	ctx := log.NewBackgroundContext(&zlog.Logger)

	tests := []struct {
		name       string
		userId     int
		order      models.Order
		mockSetup  func(storage *mocks.OrdersDatasource, paymentService *mocks.PaymentsService)
		assertFunc func(t *testing.T, err error, createdOrder *models.Order)
	}{
		{
			name:   "new order success",
			userId: 1,
			order: models.Order{
				User: &models.User{ID: 1},
				Payments: []*models.Payment{
					{Amount: 20.0},
				},
			},
			mockSetup: func(storage *mocks.OrdersDatasource, paymentService *mocks.PaymentsService) {
				paymentService.On("StorePayment", ctx, mock.MatchedBy(func(payment models.Payment) bool {
					return payment.Amount == 20.0
				})).Return(&models.Payment{Id: 1, Amount: 20.0}, nil)
				storage.On("InsertOrder", ctx, mock.Anything).Return(&dsmodels.Order{ID: 1, UserId: 1}, nil)
			},
			assertFunc: func(t *testing.T, err error, createdOrder *models.Order) {
				assert.NoError(t, err, "expected no error on storing new order")
				assert.NotNil(t, createdOrder, "expected a created order object")
				assert.Equal(t, 1, createdOrder.ID, "expected order ID to match")
			},
		},
		{
			name:   "missing user ID",
			userId: 0,
			order: models.Order{
				User: &models.User{ID: 1},
			},
			mockSetup: func(storage *mocks.OrdersDatasource, paymentService *mocks.PaymentsService) {
				// No mocks needed
			},
			assertFunc: func(t *testing.T, err error, createdOrder *models.Order) {
				assert.EqualError(t, err, "user id is required", "expected error for missing user ID")
				assert.Nil(t, createdOrder, "expected no created order for missing user ID")
			},
		},
		{
			name:   "user ID mismatch",
			userId: 2,
			order: models.Order{
				User: &models.User{ID: 1},
			},
			mockSetup: func(storage *mocks.OrdersDatasource, paymentService *mocks.PaymentsService) {
				// No mocks needed
			},
			assertFunc: func(t *testing.T, err error, createdOrder *models.Order) {
				assert.EqualError(t, err, "user id is required", "expected error for user ID mismatch")
				assert.Nil(t, createdOrder, "expected no created order for user ID mismatch")
			},
		},
		{
			name:   "payment service error",
			userId: 1,
			order: models.Order{
				User: &models.User{ID: 1},
				Payments: []*models.Payment{
					{Amount: 50.0},
				},
			},
			mockSetup: func(storage *mocks.OrdersDatasource, paymentService *mocks.PaymentsService) {
				paymentService.On("StorePayment", ctx, mock.Anything).Return(nil, errors.New("payment processing failed"))
			},
			assertFunc: func(t *testing.T, err error, createdOrder *models.Order) {
				assert.EqualError(t, err, "payment processing failed", "expected error for payment processing failure")
				assert.Nil(t, createdOrder, "expected no created order for payment processing error")
			},
		},
		{
			name:   "storage insert error",
			userId: 1,
			order: models.Order{
				User: &models.User{ID: 1},
				Payments: []*models.Payment{
					{Amount: 20.0},
				},
			},
			mockSetup: func(storage *mocks.OrdersDatasource, paymentService *mocks.PaymentsService) {
				paymentService.On("StorePayment", ctx, mock.Anything).Return(&models.Payment{Id: 1, Amount: 20.0}, nil)
				storage.On("InsertOrder", ctx, mock.Anything).Return(nil, errors.New("insert failed"))
			},
			assertFunc: func(t *testing.T, err error, createdOrder *models.Order) {
				assert.EqualError(t, err, "insert failed", "expected error for storage insert failure")
				assert.Nil(t, createdOrder, "expected no created order on storage insert failure")
			},
		},
		{
			name:   "update order success",
			userId: 1,
			order: models.Order{
				ID:   1,
				User: &models.User{ID: 1},
				Payments: []*models.Payment{
					{Id: 1, Amount: 30.0},
				},
			},
			mockSetup: func(storage *mocks.OrdersDatasource, paymentService *mocks.PaymentsService) {
				paymentService.On("StorePayment", ctx, mock.Anything).Return(&models.Payment{Id: 1, Amount: 30.0}, nil)
				storage.On("UpdateOrder", ctx, mock.Anything).Return(&dsmodels.Order{ID: 1, UserId: 1}, nil)
			},
			assertFunc: func(t *testing.T, err error, createdOrder *models.Order) {
				assert.NoError(t, err, "expected no error on updating order")
				assert.NotNil(t, createdOrder, "expected an updated order object")
				assert.Equal(t, 1, createdOrder.ID, "expected updated order ID to match")
			},
		},
		{
			name:   "storage update error",
			userId: 1,
			order: models.Order{
				ID:   1,
				User: &models.User{ID: 1},
				Payments: []*models.Payment{
					{Id: 1, Amount: 30.0},
				},
			},
			mockSetup: func(storage *mocks.OrdersDatasource, paymentService *mocks.PaymentsService) {
				paymentService.On("StorePayment", ctx, mock.Anything).Return(&models.Payment{Id: 1, Amount: 30.0}, nil)
				storage.On("UpdateOrder", ctx, mock.Anything).Return(nil, errors.New("update failed"))
			},
			assertFunc: func(t *testing.T, err error, createdOrder *models.Order) {
				assert.EqualError(t, err, "update failed", "expected error for storage update failure")
				assert.Nil(t, createdOrder, "expected no created order on storage update failure")
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			storage := mocks.NewOrdersDatasource(t)
			paymentService := mocks.NewPaymentsService(t)
			test.mockSetup(storage, paymentService)

			service := NewOrdersService(storage, paymentService)

			createdOrder, err := service.StoreOrder(ctx, test.userId, test.order)
			test.assertFunc(t, err, createdOrder)

			storage.AssertExpectations(t)
			paymentService.AssertExpectations(t)
		})
	}
}

func TestOrderService_processPayments(t *testing.T) {
	log.InitLogger()
	ctx := log.NewBackgroundContext(&zlog.Logger)

	tests := []struct {
		name       string
		order      models.Order
		mockSetup  func(paymentService *mocks.PaymentsService)
		assertFunc func(t *testing.T, payments []*models.Payment, err error)
	}{
		{
			name: "successfully process all payments",
			order: models.Order{
				ID: 1,
				Payments: []*models.Payment{
					{Amount: 50.0},
					{Amount: 25.0},
				},
			},
			mockSetup: func(paymentService *mocks.PaymentsService) {
				paymentService.On("StorePayment", ctx, mock.MatchedBy(func(p models.Payment) bool { return p.Amount == 50.0 })).Return(&models.Payment{Id: 1, Amount: 50.0}, nil)
				paymentService.On("StorePayment", ctx, mock.MatchedBy(func(p models.Payment) bool { return p.Amount == 25.0 })).Return(&models.Payment{Id: 2, Amount: 25.0}, nil)
			},
			assertFunc: func(t *testing.T, payments []*models.Payment, err error) {
				assert.NoError(t, err, "expected no error")
				assert.Len(t, payments, 2, "expected exactly 2 payments processed")
				assert.Equal(t, 1, payments[0].Id, "unexpected first payment ID")
				assert.Equal(t, 50.0, payments[0].Amount, "unexpected first payment amount")
				assert.Equal(t, 2, payments[1].Id, "unexpected second payment ID")
				assert.Equal(t, 25.0, payments[1].Amount, "unexpected second payment amount")
			},
		},
		{
			name: "handle payment processing failure",
			order: models.Order{
				ID: 1,
				Payments: []*models.Payment{
					{Amount: 100.0},
				},
			},
			mockSetup: func(paymentService *mocks.PaymentsService) {
				paymentService.On("StorePayment", ctx, mock.Anything).Return(nil, errors.New("payment service failed"))
			},
			assertFunc: func(t *testing.T, payments []*models.Payment, err error) {
				assert.Nil(t, payments, "expected no processed payments")
				assert.EqualError(t, err, "payment service failed", "unexpected error message")
			},
		},
		{
			name: "no payments to process",
			order: models.Order{
				ID:       1,
				Payments: []*models.Payment{},
			},
			mockSetup: func(paymentService *mocks.PaymentsService) {
				// No payments, no mocks needed
			},
			assertFunc: func(t *testing.T, payments []*models.Payment, err error) {
				assert.NoError(t, err, "expected no error for empty payments")
				assert.Empty(t, payments, "expected empty payment result")
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			paymentService := mocks.NewPaymentsService(t)
			test.mockSetup(paymentService)

			service := &ordersService{paymentService: paymentService}
			payments, err := service.processPayments(ctx, &test.order)

			test.assertFunc(t, payments, err)
			paymentService.AssertExpectations(t)
		})
	}
}

func TestOrderService_GetOrdersWithFilter(t *testing.T) {
	log.InitLogger()
	ctx := log.NewBackgroundContext(&zlog.Logger)

	tests := []struct {
		name       string
		userId     int
		filter     func(order *models.Order) bool
		mockSetup  func(storage *mocks.OrdersDatasource, paymentService *mocks.PaymentsService)
		assertFunc func(t *testing.T, err error, orders []*models.Order)
	}{
		{
			name:   "success case with filter",
			userId: 1,
			filter: func(order *models.Order) bool {
				return order.ID == 2
			},
			mockSetup: func(storage *mocks.OrdersDatasource, paymentService *mocks.PaymentsService) {
				storage.On("GetAllOrdersForUser", ctx, 1).Return(
					[]dsmodels.Order{
						{ID: 1, UserId: 1},
						{ID: 2, UserId: 1},
					}, nil)
				paymentService.On("GetPaymentsByOrder", ctx, 2).Return([]*models.Payment{}, nil)
			},
			assertFunc: func(t *testing.T, err error, orders []*models.Order) {
				assert.NoError(t, err, "expected no error")
				assert.Len(t, orders, 1, "expected 1 order to match the filter")

				expectedOrders := []*models.Order{
					{ID: 2, User: &models.User{ID: 1}, Payments: []*models.Payment{}},
				}
				assert.Equal(t, expectedOrders, orders, "orders do not match expected filtered results")
			},
		},
		{
			name:   "missing user id",
			userId: 0,
			filter: func(order *models.Order) bool {
				return true
			},
			mockSetup: func(storage *mocks.OrdersDatasource, paymentService *mocks.PaymentsService) {
				// No mocks needed
			},
			assertFunc: func(t *testing.T, err error, orders []*models.Order) {
				assert.EqualError(t, err, "user id is required", "expected error when user id is missing")
				assert.Nil(t, orders, "expected no orders when user id is missing")
			},
		},
		{
			name:   "storage error",
			userId: 1,
			filter: func(order *models.Order) bool {
				return true
			},
			mockSetup: func(storage *mocks.OrdersDatasource, paymentService *mocks.PaymentsService) {
				storage.On("GetAllOrdersForUser", ctx, 1).Return(nil, errors.New("storage error"))
			},
			assertFunc: func(t *testing.T, err error, orders []*models.Order) {
				assert.EqualError(t, err, "storage error", "expected storage error")
				assert.Nil(t, orders, "expected no orders when storage fails")
			},
		},
		{
			name:   "payment service error after filtering",
			userId: 1,
			filter: func(order *models.Order) bool {
				return order.ID == 1
			},
			mockSetup: func(storage *mocks.OrdersDatasource, paymentService *mocks.PaymentsService) {
				storage.On("GetAllOrdersForUser", ctx, 1).Return(
					[]dsmodels.Order{
						{ID: 1, UserId: 1},
					}, nil)
				paymentService.On("GetPaymentsByOrder", ctx, 1).Return(nil, errors.New("payment service error"))
			},
			assertFunc: func(t *testing.T, err error, orders []*models.Order) {
				assert.EqualError(t, err, "payment service error", "expected payment service error")
				assert.Nil(t, orders, "expected no orders when payment service fails after filtering")
			},
		},
		{
			name:   "no orders match the filter",
			userId: 1,
			filter: func(order *models.Order) bool {
				return false
			},
			mockSetup: func(storage *mocks.OrdersDatasource, paymentService *mocks.PaymentsService) {
				storage.On("GetAllOrdersForUser", ctx, 1).Return(
					[]dsmodels.Order{
						{ID: 1, UserId: 1},
					}, nil)
			},
			assertFunc: func(t *testing.T, err error, orders []*models.Order) {
				assert.NoError(t, err, "expected no error for unmatched filter")
				assert.Empty(t, orders, "expected no orders to match the filter")
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			storage := mocks.NewOrdersDatasource(t)
			paymentService := mocks.NewPaymentsService(t)
			test.mockSetup(storage, paymentService)

			service := NewOrdersService(storage, paymentService)
			orders, err := service.GetOrdersWithFilter(ctx, test.userId, test.filter)
			test.assertFunc(t, err, orders)

			storage.AssertExpectations(t)
			paymentService.AssertExpectations(t)
		})
	}
}

func assertError(t *testing.T, err error, expectedErr error) {
	assert.Nil(t, nil, "expected result to be nil")
	assert.EqualError(t, err, expectedErr.Error(), "unexpected error message")
}

func assertSuccess(t *testing.T, err error, expectedOrder *models.Order, actualOrder *models.Order) {
	assert.NoError(t, err, "expected no error")
	assert.Equal(t, expectedOrder, actualOrder, "unexpected order result")
}

func TestOrderService_GetOrder(t *testing.T) {
	log.InitLogger()
	ctx := log.NewBackgroundContext(&zlog.Logger)

	tests := []struct {
		name       string
		userId     int
		orderId    int
		mockSetup  func(storage *mocks.OrdersDatasource, paymentService *mocks.PaymentsService, authorizationService *mocks.AuthorizationService)
		assertFunc func(t *testing.T, err error, actualOrder *models.Order)
	}{
		{
			name:    "success case",
			userId:  1,
			orderId: 123,
			mockSetup: func(storage *mocks.OrdersDatasource, paymentService *mocks.PaymentsService, authorizationService *mocks.AuthorizationService) {
				storage.On("GetOrder", ctx, 123).Return(&dsmodels.Order{ID: 123, UserId: 1}, nil)
				paymentService.On("GetPaymentsByOrder", ctx, 123).Return([]*models.Payment{}, nil)
				authorizationService.On("IsAuthorized", ctx, 1, &models.Order{ID: 123, User: &models.User{ID: 1}, Payments: []*models.Payment{}}).Return(true, nil)
			},
			assertFunc: func(t *testing.T, err error, actualOrder *models.Order) {
				expectedOrder := &models.Order{ID: 123, User: &models.User{ID: 1}, Payments: []*models.Payment{}}
				assertSuccess(t, err, expectedOrder, actualOrder)
			},
		},
		{
			name:    "missing user id",
			userId:  0,
			orderId: 123,
			mockSetup: func(storage *mocks.OrdersDatasource, paymentService *mocks.PaymentsService, authorizationService *mocks.AuthorizationService) {
				// No mocks needed since the function returns at the beginning
			},
			assertFunc: func(t *testing.T, err error, actualOrder *models.Order) {
				assertError(t, err, errors.New("user id is required"))
			},
		},
		{
			name:    "order not found in storage",
			userId:  1,
			orderId: 123,
			mockSetup: func(storage *mocks.OrdersDatasource, paymentService *mocks.PaymentsService, authorizationService *mocks.AuthorizationService) {
				storage.On("GetOrder", ctx, 123).Return(nil, errors.New("order not found"))
			},
			assertFunc: func(t *testing.T, err error, actualOrder *models.Order) {
				assertError(t, err, errors.New("order not found"))
			},
		},
		{
			name:    "user not authorized",
			userId:  2,
			orderId: 123,
			mockSetup: func(storage *mocks.OrdersDatasource, paymentService *mocks.PaymentsService, authorizationService *mocks.AuthorizationService) {
				storage.On("GetOrder", ctx, 123).Return(&dsmodels.Order{ID: 123, UserId: 1}, nil)
				authorizationService.On("IsAuthorized", ctx, 2, &models.Order{ID: 123, User: &models.User{ID: 1}, Payments: []*models.Payment{}}).Return(false, nil)
			},
			assertFunc: func(t *testing.T, err error, actualOrder *models.Order) {
				assertError(t, err, errors.New("user is not authorized to access this order"))
			},
		},
		{
			name:    "error fetching payments",
			userId:  1,
			orderId: 123,
			mockSetup: func(storage *mocks.OrdersDatasource, paymentService *mocks.PaymentsService, authorizationService *mocks.AuthorizationService) {
				storage.On("GetOrder", ctx, 123).Return(&dsmodels.Order{ID: 123, UserId: 1}, nil)
				paymentService.On("GetPaymentsByOrder", ctx, 123).Return(nil, errors.New("payment fetch error"))
				authorizationService.On("IsAuthorized", ctx, 1, &models.Order{ID: 123, User: &models.User{ID: 1}, Payments: []*models.Payment{}}).Return(true, nil)
			},
			assertFunc: func(t *testing.T, err error, actualOrder *models.Order) {
				assertError(t, err, errors.New("payment fetch error"))
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			storage := mocks.NewOrdersDatasource(t)
			paymentService := mocks.NewPaymentsService(t)
			authorizationService := mocks.NewAuthorizationService(t)
			test.mockSetup(storage, paymentService, authorizationService)

			service := NewOrdersService(storage, paymentService, authorizationService)

			actualOrder, err := service.GetOrder(ctx, test.userId, test.orderId)
			test.assertFunc(t, err, actualOrder)

			storage.AssertExpectations(t)
			paymentService.AssertExpectations(t)
			authorizationService.AssertExpectations(t)
		})
	}
}

func TestOrderService_GetOrders(t *testing.T) {
	log.InitLogger()
	ctx := log.NewBackgroundContext(&zlog.Logger)

	tests := []struct {
		name       string
		userId     int
		mockSetup  func(storage *mocks.OrdersDatasource, paymentService *mocks.PaymentsService)
		assertFunc func(t *testing.T, err error, orders []*models.Order)
	}{
		{
			name:   "success case",
			userId: 1,
			mockSetup: func(storage *mocks.OrdersDatasource, paymentService *mocks.PaymentsService) {
				storage.On("GetAllOrdersForUser", ctx, 1).Return(
					[]dsmodels.Order{
						{ID: 1, UserId: 1},
						{ID: 2, UserId: 1},
					}, nil)
				paymentService.On("GetPaymentsByOrder", ctx, 1).Return([]*models.Payment{}, nil)
				paymentService.On("GetPaymentsByOrder", ctx, 2).Return([]*models.Payment{}, nil)
			},
			assertFunc: func(t *testing.T, err error, orders []*models.Order) {
				assert.NoError(t, err, "expected no error, got error")
				assert.Len(t, orders, 2, "expected 2 orders, but got a different count")

				expectedOrders := []*models.Order{
					{ID: 1, User: &models.User{ID: 1}, Payments: []*models.Payment{}},
					{ID: 2, User: &models.User{ID: 1}, Payments: []*models.Payment{}},
				}
				assert.Equal(t, expectedOrders, orders, "orders do not match expected output")
			},
		},
		{
			name:   "missing user id",
			userId: 0,
			mockSetup: func(storage *mocks.OrdersDatasource, paymentService *mocks.PaymentsService) {
				// No mocks needed
			},
			assertFunc: func(t *testing.T, err error, orders []*models.Order) {
				assert.EqualError(t, err, "user id is required", "expected error when user id is missing")
				assert.Nil(t, orders, "expected no orders when user id is missing")
			},
		},
		{
			name:   "storage error",
			userId: 1,
			mockSetup: func(storage *mocks.OrdersDatasource, paymentService *mocks.PaymentsService) {
				storage.On("GetAllOrdersForUser", ctx, 1).Return(nil, errors.New("storage error"))
			},
			assertFunc: func(t *testing.T, err error, orders []*models.Order) {
				assert.EqualError(t, err, "storage error", "expected storage error")
				assert.Nil(t, orders, "expected no orders when storage fails")
			},
		},
		{
			name:   "payment service error",
			userId: 1,
			mockSetup: func(storage *mocks.OrdersDatasource, paymentService *mocks.PaymentsService) {
				storage.On("GetAllOrdersForUser", ctx, 1).Return(
					[]dsmodels.Order{
						{ID: 1, UserId: 1},
					}, nil)
				paymentService.On("GetPaymentsByOrder", ctx, 1).Return(nil, errors.New("payment service error"))
			},
			assertFunc: func(t *testing.T, err error, orders []*models.Order) {
				assert.EqualError(t, err, "payment service error", "expected payment service error")
				assert.Nil(t, orders, "expected no orders when payment service fails")
			},
		},
		{
			name:   "no orders fetched",
			userId: 1,
			mockSetup: func(storage *mocks.OrdersDatasource, paymentService *mocks.PaymentsService) {
				storage.On("GetAllOrdersForUser", ctx, 1).Return([]dsmodels.Order{}, nil)
			},
			assertFunc: func(t *testing.T, err error, orders []*models.Order) {
				assert.NoError(t, err, "expected no error for empty results")
				assert.Empty(t, orders, "expected no orders fetched")
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			storage := mocks.NewOrdersDatasource(t)
			paymentService := mocks.NewPaymentsService(t)
			test.mockSetup(storage, paymentService)

			service := NewOrdersService(storage, paymentService)

			orders, err := service.GetOrders(ctx, test.userId)
			test.assertFunc(t, err, orders)

			storage.AssertExpectations(t)
			paymentService.AssertExpectations(t)
		})
	}
}
