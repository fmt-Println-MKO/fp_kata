package services

import (
	"context"
	"errors"
	"fp_kata/common/constants"
	"fp_kata/internal/datasources/dsmodels"
	"fp_kata/pkg/log"
	zlog "github.com/rs/zerolog/log"
	"github.com/samber/mo"
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
		name           string
		userId         int
		order          models.Order
		mockSetup      func(storage *mocks.OrdersDatasource, paymentService *mocks.PaymentsService)
		expectedResult mo.Result[*models.Order]
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
				storage.On("InsertOrder", ctx, mock.Anything).Return(mo.Ok(dsmodels.Order{ID: 1, UserId: 1}))
			},
			expectedResult: mo.Ok(&models.Order{ID: 1, User: &models.User{ID: 1}, Payments: []*models.Payment{{Id: 1, Amount: 20.0}}}),
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
			expectedResult: mo.Errf[*models.Order]("user id is required"),
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
			expectedResult: mo.Errf[*models.Order]("user id is required"),
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
			expectedResult: mo.Errf[*models.Order]("payment processing failed"),
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
				storage.On("InsertOrder", ctx, mock.Anything).Return(mo.Errf[dsmodels.Order]("insert failed"))
			},
			expectedResult: mo.Errf[*models.Order]("insert failed"),
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
				storage.On("UpdateOrder", ctx, mock.Anything).Return(mo.Ok(dsmodels.Order{ID: 1, UserId: 1}))
			},
			expectedResult: mo.Ok(&models.Order{ID: 1, User: &models.User{ID: 1}, Payments: []*models.Payment{{Id: 1, Amount: 30.0}}}),
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
				storage.On("UpdateOrder", ctx, mock.Anything).Return(mo.Errf[dsmodels.Order]("update failed"))
			},
			expectedResult: mo.Errf[*models.Order]("update failed"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			storage := mocks.NewOrdersDatasource(t)
			paymentService := mocks.NewPaymentsService(t)
			authorizationService := mocks.NewAuthorizationService(t)
			test.mockSetup(storage, paymentService)

			service := NewOrdersService(storage, paymentService, authorizationService)

			createdOrderResult := service.StoreOrder(ctx, test.userId, test.order)
			assert.Equal(t, test.expectedResult, createdOrderResult)

			storage.AssertExpectations(t)
			paymentService.AssertExpectations(t)
			authorizationService.AssertExpectations(t)
		})
	}
}

func TestOrderService_processPayments(t *testing.T) {
	log.InitLogger()
	ctx := log.NewBackgroundContext(&zlog.Logger)

	tests := []struct {
		name           string
		order          models.Order
		mockSetup      func(paymentService *mocks.PaymentsService)
		expectedResult mo.Result[[]*models.Payment]
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
			expectedResult: mo.Ok[[]*models.Payment]([]*models.Payment{{Id: 1, Amount: 50.0}, {Id: 2, Amount: 25.0}}),
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
			expectedResult: mo.Errf[[]*models.Payment]("payment service failed"),
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
			expectedResult: mo.Ok[[]*models.Payment]([]*models.Payment{}),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			paymentService := mocks.NewPaymentsService(t)
			test.mockSetup(paymentService)

			service := &ordersService{paymentService: paymentService}
			paymentsResult := service.processPayments(ctx, &test.order)
			assert.Equal(t, test.expectedResult, paymentsResult)
			paymentService.AssertExpectations(t)
		})
	}
}

func TestOrderService_GetOrdersWithFilter(t *testing.T) {
	log.InitLogger()
	ctx := log.NewBackgroundContext(&zlog.Logger)

	tests := []struct {
		name           string
		userId         int
		ctxUser        *models.User
		filter         func(order *models.Order) bool
		mockSetup      func(storage *mocks.OrdersDatasource, paymentService *mocks.PaymentsService, authorizationService *mocks.AuthorizationService)
		expectedResult mo.Result[[]*models.Order]
	}{
		{
			name:    "success case with filter",
			userId:  1,
			ctxUser: &models.User{ID: 1},
			filter: func(order *models.Order) bool {
				return order.ID == 2
			},
			mockSetup: func(storage *mocks.OrdersDatasource, paymentService *mocks.PaymentsService, authorizationService *mocks.AuthorizationService) {
				storage.On("GetAllOrdersForUser", mock.Anything, 1).Return(
					mo.Ok([]dsmodels.Order{
						{ID: 1, UserId: 1},
						{ID: 2, UserId: 1},
					}))
				paymentService.On("GetPaymentsByOrder", mock.Anything, 2).Return([]*models.Payment{}, nil)
				authorizationService.On("IsAuthorized", mock.Anything, 1, &models.Order{ID: 2, User: &models.User{ID: 1}, Payments: []*models.Payment{}}).Return(true, nil).Once()
			},
			expectedResult: mo.Ok([]*models.Order{
				{ID: 2, User: &models.User{ID: 1}, Payments: []*models.Payment{}},
			}),
		},
		{
			name:   "missing user id",
			userId: 0,
			filter: func(order *models.Order) bool {
				return true
			},
			mockSetup: func(storage *mocks.OrdersDatasource, paymentService *mocks.PaymentsService, authorizationService *mocks.AuthorizationService) {
				// No mocks needed
			},
			expectedResult: mo.Errf[[]*models.Order]("user id is required"),
		},
		{
			name:   "context missing authenticated user",
			userId: 1,
			filter: func(order *models.Order) bool {
				return true
			},
			mockSetup: func(storage *mocks.OrdersDatasource, paymentService *mocks.PaymentsService, authorizationService *mocks.AuthorizationService) {
				storage.On("GetAllOrdersForUser", mock.Anything, 1).Return(
					mo.Ok([]dsmodels.Order{
						{ID: 1, UserId: 1},
					}))
				paymentService.On("GetPaymentsByOrder", mock.Anything, 1).Return([]*models.Payment{}, nil)
				authorizationService.On("IsAuthorized", mock.Anything, 1, &models.Order{ID: 1, User: &models.User{ID: 1}, Payments: []*models.Payment{}}).Return(true, nil).Once()
			},
			expectedResult: mo.Errf[[]*models.Order]("user id is required"),
		},
		{
			name:   "storage error",
			userId: 1,
			filter: func(order *models.Order) bool {
				return true
			},
			mockSetup: func(storage *mocks.OrdersDatasource, paymentService *mocks.PaymentsService, authorizationService *mocks.AuthorizationService) {
				storage.On("GetAllOrdersForUser", mock.Anything, 1).Return(mo.Errf[[]dsmodels.Order]("storage error"))
			},
			expectedResult: mo.Errf[[]*models.Order]("storage error"),
		},
		{
			name:   "payment service error after filtering",
			userId: 1,
			filter: func(order *models.Order) bool {
				return order.ID == 1
			},
			mockSetup: func(storage *mocks.OrdersDatasource, paymentService *mocks.PaymentsService, authorizationService *mocks.AuthorizationService) {
				storage.On("GetAllOrdersForUser", mock.Anything, 1).Return(
					mo.Ok([]dsmodels.Order{
						{ID: 1, UserId: 1},
					}))
				paymentService.On("GetPaymentsByOrder", mock.Anything, 1).Return(nil, errors.New("payment service error"))
				authorizationService.On("IsAuthorized", mock.Anything, 1, &models.Order{ID: 1, User: &models.User{ID: 1}, Payments: []*models.Payment{}}).Return(true, nil).Once()
			},
			expectedResult: mo.Errf[[]*models.Order]("payment service error"),
		},
		{
			name:   "no orders match the filter",
			userId: 1,
			filter: func(order *models.Order) bool {
				return false
			},
			mockSetup: func(storage *mocks.OrdersDatasource, paymentService *mocks.PaymentsService, authorizationService *mocks.AuthorizationService) {
				storage.On("GetAllOrdersForUser", mock.Anything, 1).Return(
					mo.Ok([]dsmodels.Order{
						{ID: 1, UserId: 1},
					}))
			},
			expectedResult: mo.Ok[[]*models.Order]([]*models.Order{}),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			storage := mocks.NewOrdersDatasource(t)
			paymentService := mocks.NewPaymentsService(t)
			authorizationService := mocks.NewAuthorizationService(t)
			test.mockSetup(storage, paymentService, authorizationService)

			service := NewOrdersService(storage, paymentService, authorizationService)

			testCtx := context.WithValue(ctx, constants.AuthenticatedUserIdKey, test.userId)
			testCtx = context.WithValue(testCtx, constants.AuthenticatedUserKey, test.ctxUser)

			ordersResult := service.GetOrdersWithFilter(testCtx, test.userId, test.filter)
			assert.Equal(t, test.expectedResult, ordersResult)

			storage.AssertExpectations(t)
			paymentService.AssertExpectations(t)
			authorizationService.AssertExpectations(t)
		})
	}
}

func TestOrderService_GetOrder(t *testing.T) {
	log.InitLogger()
	ctx := log.NewBackgroundContext(&zlog.Logger)

	tests := []struct {
		name           string
		userId         int
		ctxUser        *models.User
		orderId        int
		mockSetup      func(storage *mocks.OrdersDatasource, paymentService *mocks.PaymentsService, authorizationService *mocks.AuthorizationService)
		expectedResult mo.Result[*models.Order]
	}{
		{
			name:    "success case",
			userId:  1,
			ctxUser: &models.User{ID: 1},
			orderId: 123,
			mockSetup: func(storage *mocks.OrdersDatasource, paymentService *mocks.PaymentsService, authorizationService *mocks.AuthorizationService) {
				storage.On("GetOrder", mock.Anything, 123).Return(mo.Ok(dsmodels.Order{ID: 123, UserId: 1}))
				paymentService.On("GetPaymentsByOrder", mock.Anything, 123).Return([]*models.Payment{}, nil)
				authorizationService.On("IsAuthorized", mock.Anything, 1, &models.Order{ID: 123, User: &models.User{ID: 1}, Payments: []*models.Payment{}}).Return(true, nil)
			},
			expectedResult: mo.Ok(&models.Order{ID: 123, User: &models.User{ID: 1}, Payments: []*models.Payment{}}),
		},
		{
			name:    "missing user id",
			userId:  0,
			orderId: 123,
			mockSetup: func(storage *mocks.OrdersDatasource, paymentService *mocks.PaymentsService, authorizationService *mocks.AuthorizationService) {
				// No mocks needed since the function returns at the beginning
			},
			expectedResult: mo.Errf[*models.Order]("user id is required"),
		},
		{
			name:    "order not found in storage",
			userId:  1,
			orderId: 123,
			mockSetup: func(storage *mocks.OrdersDatasource, paymentService *mocks.PaymentsService, authorizationService *mocks.AuthorizationService) {
				storage.On("GetOrder", mock.Anything, 123).Return(mo.Errf[dsmodels.Order]("order not found"))
			},
			expectedResult: mo.Errf[*models.Order]("order not found"),
		},
		{
			name:    "user not authorized",
			userId:  2,
			orderId: 123,
			mockSetup: func(storage *mocks.OrdersDatasource, paymentService *mocks.PaymentsService, authorizationService *mocks.AuthorizationService) {
				storage.On("GetOrder", mock.Anything, 123).Return(mo.Ok(dsmodels.Order{ID: 123, UserId: 1}))
				authorizationService.On("IsAuthorized", mock.Anything, 2, &models.Order{ID: 123, User: &models.User{ID: 1}, Payments: []*models.Payment{}}).Return(false, nil)
			},
			expectedResult: mo.Errf[*models.Order]("user is not authorized to access this order"),
		},
		{
			name:    "user authorization error",
			userId:  1,
			orderId: 123,
			mockSetup: func(storage *mocks.OrdersDatasource, paymentService *mocks.PaymentsService, authorizationService *mocks.AuthorizationService) {
				storage.On("GetOrder", mock.Anything, 123).Return(mo.Ok(dsmodels.Order{ID: 123, UserId: 1}))
				authorizationService.On("IsAuthorized", mock.Anything, 1, &models.Order{ID: 123, User: &models.User{ID: 1}, Payments: []*models.Payment{}}).Return(false, errors.New("userId is required"))
			},
			expectedResult: mo.Errf[*models.Order]("userId is required"),
		},
		{
			name:    "context missing authenticated user",
			userId:  1,
			orderId: 123,
			mockSetup: func(storage *mocks.OrdersDatasource, paymentService *mocks.PaymentsService, authorizationService *mocks.AuthorizationService) {
				storage.On("GetOrder", mock.Anything, 123).Return(mo.Ok(dsmodels.Order{ID: 123, UserId: 1}))
				paymentService.On("GetPaymentsByOrder", mock.Anything, 123).Return([]*models.Payment{}, nil)
				authorizationService.On("IsAuthorized", mock.Anything, 1, &models.Order{ID: 123, User: &models.User{ID: 1}, Payments: []*models.Payment{}}).Return(true, nil)
			},
			expectedResult: mo.Errf[*models.Order]("user id is required"),
		},
		{
			name:    "error fetching payments",
			userId:  1,
			ctxUser: &models.User{ID: 1},
			orderId: 123,
			mockSetup: func(storage *mocks.OrdersDatasource, paymentService *mocks.PaymentsService, authorizationService *mocks.AuthorizationService) {
				storage.On("GetOrder", mock.Anything, 123).Return(mo.Ok(dsmodels.Order{ID: 123, UserId: 1}))
				paymentService.On("GetPaymentsByOrder", mock.Anything, 123).Return(nil, errors.New("payment fetch error"))
				authorizationService.On("IsAuthorized", mock.Anything, 1, &models.Order{ID: 123, User: &models.User{ID: 1}, Payments: []*models.Payment{}}).Return(true, nil)
			},
			expectedResult: mo.Errf[*models.Order]("payment fetch error"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			storage := mocks.NewOrdersDatasource(t)
			paymentService := mocks.NewPaymentsService(t)
			authorizationService := mocks.NewAuthorizationService(t)
			test.mockSetup(storage, paymentService, authorizationService)

			service := NewOrdersService(storage, paymentService, authorizationService)

			testCtx := context.WithValue(ctx, constants.AuthenticatedUserIdKey, test.userId)
			testCtx = context.WithValue(testCtx, constants.AuthenticatedUserKey, test.ctxUser)

			actualOrderResult := service.GetOrder(testCtx, test.userId, test.orderId)
			assert.Equal(t, test.expectedResult, actualOrderResult)

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
		name           string
		userId         int
		ctxUser        *models.User
		mockSetup      func(storage *mocks.OrdersDatasource, paymentService *mocks.PaymentsService, authorizationService *mocks.AuthorizationService)
		expectedResult mo.Result[[]*models.Order]
	}{
		{
			name:    "success case",
			userId:  1,
			ctxUser: &models.User{ID: 1},
			mockSetup: func(storage *mocks.OrdersDatasource, paymentService *mocks.PaymentsService, authorizationService *mocks.AuthorizationService) {
				storage.On("GetAllOrdersForUser", mock.Anything, 1).Return(
					mo.Ok([]dsmodels.Order{
						{ID: 1, UserId: 1},
						{ID: 2, UserId: 1},
					}))
				paymentService.On("GetPaymentsByOrder", mock.Anything, 1).Return([]*models.Payment{}, nil)
				paymentService.On("GetPaymentsByOrder", mock.Anything, 2).Return([]*models.Payment{}, nil)
				authorizationService.On("IsAuthorized", mock.Anything, 1, &models.Order{ID: 1, User: &models.User{ID: 1}, Payments: []*models.Payment{}}).Return(true, nil)
				authorizationService.On("IsAuthorized", mock.Anything, 1, &models.Order{ID: 2, User: &models.User{ID: 1}, Payments: []*models.Payment{}}).Return(true, nil)
			},
			expectedResult: mo.Ok([]*models.Order{
				{ID: 1, User: &models.User{ID: 1}, Payments: []*models.Payment{}},
				{ID: 2, User: &models.User{ID: 1}, Payments: []*models.Payment{}},
			}),
		},
		{
			name:   "missing user id",
			userId: 0,
			mockSetup: func(storage *mocks.OrdersDatasource, paymentService *mocks.PaymentsService, authorizationService *mocks.AuthorizationService) {
				// No mocks needed
			},
			expectedResult: mo.Errf[[]*models.Order]("user id is required"),
		},
		{
			name:   "context missing user",
			userId: 1,
			mockSetup: func(storage *mocks.OrdersDatasource, paymentService *mocks.PaymentsService, authorizationService *mocks.AuthorizationService) {
				storage.On("GetAllOrdersForUser", mock.Anything, 1).Return(
					mo.Ok([]dsmodels.Order{
						{ID: 1, UserId: 1},
						{ID: 2, UserId: 1},
					}))
				paymentService.On("GetPaymentsByOrder", mock.Anything, 1).Return([]*models.Payment{}, nil)
				authorizationService.On("IsAuthorized", mock.Anything, 1, &models.Order{ID: 1, User: &models.User{ID: 1}, Payments: []*models.Payment{}}).Return(true, nil)
			},
			expectedResult: mo.Errf[[]*models.Order]("user id is required"),
		},
		{
			name:   "storage error",
			userId: 1,
			mockSetup: func(storage *mocks.OrdersDatasource, paymentService *mocks.PaymentsService, authorizationService *mocks.AuthorizationService) {
				storage.On("GetAllOrdersForUser", mock.Anything, 1).Return(mo.Errf[[]dsmodels.Order]("storage error"))
			},
			expectedResult: mo.Errf[[]*models.Order]("storage error"),
		},
		{
			name:   "payment service error",
			userId: 1,
			mockSetup: func(storage *mocks.OrdersDatasource, paymentService *mocks.PaymentsService, authorizationService *mocks.AuthorizationService) {
				storage.On("GetAllOrdersForUser", mock.Anything, 1).Return(
					mo.Ok([]dsmodels.Order{
						{ID: 1, UserId: 1},
					}))
				paymentService.On("GetPaymentsByOrder", mock.Anything, 1).Return(nil, errors.New("payment service error"))
				authorizationService.On("IsAuthorized", mock.Anything, 1, &models.Order{ID: 1, User: &models.User{ID: 1}, Payments: []*models.Payment{}}).Return(true, nil)
			},
			expectedResult: mo.Errf[[]*models.Order]("payment service error"),
		},
		{
			name:   "no orders fetched",
			userId: 1,
			mockSetup: func(storage *mocks.OrdersDatasource, paymentService *mocks.PaymentsService, authorizationService *mocks.AuthorizationService) {
				storage.On("GetAllOrdersForUser", mock.Anything, 1).Return(mo.Ok([]dsmodels.Order{}))
			},
			expectedResult: mo.Ok([]*models.Order{}),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			storage := mocks.NewOrdersDatasource(t)
			paymentService := mocks.NewPaymentsService(t)
			authorizationService := mocks.NewAuthorizationService(t)
			test.mockSetup(storage, paymentService, authorizationService)

			service := NewOrdersService(storage, paymentService, authorizationService)

			testCtx := context.WithValue(ctx, constants.AuthenticatedUserIdKey, test.userId)
			testCtx = context.WithValue(testCtx, constants.AuthenticatedUserKey, test.ctxUser)

			ordersResult := service.GetOrders(testCtx, test.userId)
			assert.Equal(t, test.expectedResult, ordersResult)

			storage.AssertExpectations(t)
			paymentService.AssertExpectations(t)
			authorizationService.AssertExpectations(t)
		})
	}
}
