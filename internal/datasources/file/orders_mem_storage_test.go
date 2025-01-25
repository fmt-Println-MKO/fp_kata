package file

import (
	"context"
	"fp_kata/common/monads"
	"fp_kata/internal/datasources"
	"fp_kata/internal/datasources/dsmodels"
	"fp_kata/pkg/log"
	zlog "github.com/rs/zerolog/log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func initTestOrdersStorage(store map[int]dsmodels.Order) (*inMemoryOrdersStorage, context.Context) {
	log.InitLogger()
	ctx := log.NewBackgroundContext(&zlog.Logger)
	return &inMemoryOrdersStorage{
		orders: store,
	}, ctx
}

func TestGetOrder(t *testing.T) {

	const errOrderNotFound = "order not found"

	tests := []struct {
		name           string
		initialOrders  map[int]dsmodels.Order
		orderID        int
		expectedResult monads.Result[dsmodels.Order]
	}{
		{
			name: "OrderExists",
			initialOrders: map[int]dsmodels.Order{
				1: {ID: 1, UserId: 123, Payments: []int{1, 2}},
			},
			orderID:        1,
			expectedResult: monads.Ok(dsmodels.Order{ID: 1, UserId: 123, Payments: []int{1, 2}}),
		},
		{
			name: "OrderDoesNotExist",
			initialOrders: map[int]dsmodels.Order{
				1: {ID: 1, UserId: 123, Payments: []int{1, 2}},
			},
			orderID:        2,
			expectedResult: monads.Errf[dsmodels.Order](errOrderNotFound),
		},
		{
			name:           "EmptyOrdersStorage",
			initialOrders:  map[int]dsmodels.Order{},
			orderID:        1,
			expectedResult: monads.Errf[dsmodels.Order](errOrderNotFound),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			storage, ctx := initTestOrdersStorage(tc.initialOrders)

			orderResult := storage.GetOrder(ctx, tc.orderID)

			assert.Equal(t, tc.expectedResult, orderResult, "unexpected result")
		})
	}
}

func TestGetAllOrdersForUser(t *testing.T) {

	tests := []struct {
		name           string
		initialOrders  map[int]dsmodels.Order
		userID         int
		expectedResult monads.Result[[]dsmodels.Order]
	}{
		{
			name: "UserHasOrders",
			initialOrders: map[int]dsmodels.Order{
				1: {ID: 1, UserId: 123},
				2: {ID: 2, UserId: 123},
				3: {ID: 3, UserId: 456},
			},
			userID: 123,
			expectedResult: monads.Ok([]dsmodels.Order{
				{ID: 1, UserId: 123},
				{ID: 2, UserId: 123},
			}),
		},
		{
			name: "UserHasNoOrders",
			initialOrders: map[int]dsmodels.Order{
				1: {ID: 1, UserId: 789},
			},
			userID:         123,
			expectedResult: monads.Ok([]dsmodels.Order{}),
		},
		{
			name:           "EmptyStorage",
			initialOrders:  map[int]dsmodels.Order{},
			userID:         123,
			expectedResult: monads.Ok([]dsmodels.Order{}),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			storage, ctx := initTestOrdersStorage(tc.initialOrders)
			ordersResult := storage.GetAllOrdersForUser(ctx, tc.userID)
			assert.Equal(t, tc.expectedResult, ordersResult, "unexpected result")
		})
	}
}

func TestDeleteOrder(t *testing.T) {
	tests := []struct {
		name          string
		initialOrders map[int]dsmodels.Order
		orderID       int
		validate      func(*testing.T, error, map[int]dsmodels.Order)
	}{
		{
			name: "OrderExists",
			initialOrders: map[int]dsmodels.Order{
				1: {ID: 1, UserId: 123},
				2: {ID: 2, UserId: 456},
			},
			orderID: 1,
			validate: func(t *testing.T, err error, orders map[int]dsmodels.Order) {
				assert.NoError(t, err, "unexpected error when deleting order")
				assert.Len(t, orders, 1, "expected only one order to remain")
				_, exists := orders[1]
				assert.False(t, exists, "order with ID 1 should have been deleted")
			},
		},
		{
			name: "OrderDoesNotExist",
			initialOrders: map[int]dsmodels.Order{
				1: {ID: 1, UserId: 123},
			},
			orderID: 2,
			validate: func(t *testing.T, err error, orders map[int]dsmodels.Order) {
				assert.EqualError(t, err, "order not found", "expected error mismatch")
				assert.Len(t, orders, 1, "expected orders storage to remain unchanged")
				order, exists := orders[1]
				assert.True(t, exists, "order with ID 1 should still be present")
				assert.Equal(t, 1, order.ID, "order ID mismatch in storage")
			},
		},
		{
			name:          "EmptyOrders",
			initialOrders: map[int]dsmodels.Order{},
			orderID:       1,
			validate: func(t *testing.T, err error, orders map[int]dsmodels.Order) {
				assert.EqualError(t, err, "order not found", "expected error mismatch")
				assert.Empty(t, orders, "expected no orders in storage")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			storage, ctx := initTestOrdersStorage(tc.initialOrders)

			err := storage.DeleteOrder(ctx, tc.orderID)

			tc.validate(t, err, storage.orders)
		})
	}
}

func TestUpdateOrder(t *testing.T) {

	tests := []struct {
		name           string
		initialOrders  map[int]dsmodels.Order
		updateOrder    dsmodels.Order
		expectedError  string
		expectedResult monads.Result[dsmodels.Order]
	}{
		{
			name: "UpdateExistingOrder",
			initialOrders: map[int]dsmodels.Order{
				1: {ID: 1, UserId: 123, Quantity: 2},
			},
			updateOrder:    dsmodels.Order{ID: 1, UserId: 123, Quantity: 3},
			expectedError:  "",
			expectedResult: monads.Ok(dsmodels.Order{ID: 1, UserId: 123, Quantity: 3}),
		},
		{
			name: "UpdateNonExistentOrder",
			initialOrders: map[int]dsmodels.Order{
				1: {ID: 1, UserId: 123, Quantity: 2},
			},
			updateOrder:    dsmodels.Order{ID: 2, UserId: 456, Quantity: 1},
			expectedError:  "order not found",
			expectedResult: monads.Errf[dsmodels.Order]("order not found"),
		},
		{
			name:           "UpdateOrderInEmptyStorage",
			initialOrders:  map[int]dsmodels.Order{},
			updateOrder:    dsmodels.Order{ID: 1, UserId: 123, Quantity: 1},
			expectedError:  "order not found",
			expectedResult: monads.Errf[dsmodels.Order]("order not found"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			storage, ctx := initTestOrdersStorage(tc.initialOrders)
			orderResult := storage.UpdateOrder(ctx, tc.updateOrder)
			assert.Equal(t, tc.expectedResult, orderResult, "unexpected result")
		})
	}
}

func TestInsertOrder(t *testing.T) {

	tests := []struct {
		name          string
		initialOrders map[int]dsmodels.Order
		orderToInsert dsmodels.Order
		expectedError string
		validate      func(*testing.T, *dsmodels.Order, map[int]dsmodels.Order, error)
	}{
		{
			name: "InsertNewOrder",
			initialOrders: map[int]dsmodels.Order{
				1: {ID: 1, UserId: 123, Quantity: 1},
			},
			orderToInsert: dsmodels.Order{ID: 2, UserId: 456, Quantity: 2},
			expectedError: "",
			validate: func(t *testing.T, order *dsmodels.Order, orders map[int]dsmodels.Order, err error) {
				assert.NoError(t, err, "unexpected error when inserting a new order")
				assert.Equal(t, &dsmodels.Order{ID: 2, UserId: 456, Quantity: 2}, order, "inserted order details mismatch")
				assert.Len(t, orders, 2, "expected orders map to contain two entries")
				_, exists := orders[2]
				assert.True(t, exists, "new order with ID 2 should exist in storage")
			},
		},
		{
			name: "InsertExistingOrder",
			initialOrders: map[int]dsmodels.Order{
				1: {ID: 1, UserId: 123, Quantity: 1},
			},
			orderToInsert: dsmodels.Order{ID: 1, UserId: 123, Quantity: 2},
			expectedError: "order already exists",
			validate: func(t *testing.T, order *dsmodels.Order, orders map[int]dsmodels.Order, err error) {
				assert.Nil(t, order, "expected no order to be returned when inserting existing order")
				assert.EqualError(t, err, "order already exists", "expected error mismatch")
				assert.Len(t, orders, 1, "orders map should remain unchanged with one entry")
				existingOrder, exists := orders[1]
				assert.True(t, exists, "existing order with ID 1 should still exist")
				assert.Equal(t, &dsmodels.Order{ID: 1, UserId: 123, Quantity: 1}, &existingOrder, "existing order details should remain unchanged")
			},
		},
		{
			name:          "InsertIntoEmptyStorage",
			initialOrders: map[int]dsmodels.Order{},
			orderToInsert: dsmodels.Order{ID: 1, UserId: 123, Quantity: 1},
			expectedError: "",
			validate: func(t *testing.T, order *dsmodels.Order, orders map[int]dsmodels.Order, err error) {
				assert.NoError(t, err, "unexpected error when inserting into empty storage")
				assert.Equal(t, &dsmodels.Order{ID: 1, UserId: 123, Quantity: 1}, order, "unexpected order details after insertion")
				assert.Len(t, orders, 1, "expected orders map to contain one entry")
				_, exists := orders[1]
				assert.True(t, exists, "new order with ID 1 should exist in storage")
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			storage, ctx := initTestOrdersStorage(tc.initialOrders)
			insertedOrder, err := storage.InsertOrder(ctx, tc.orderToInsert)
			tc.validate(t, insertedOrder, storage.orders, err)
		})
	}
}

func TestNewOrderStorage(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "CreateEmptyStorage",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			storage := NewOrdersStorage()

			inMemoryStorage, ok := storage.(*inMemoryOrdersStorage)
			assert.True(t, ok, "expected storage to be of type *inMemoryOrdersStorage")
			assert.NotNil(t, inMemoryStorage, "expected storage to be initialized")
			assert.Empty(t, inMemoryStorage.orders, "expected storage to be empty upon initialization")

			assert.Implements(t, (*datasources.OrdersDatasource)(nil), storage, "storage does not implement OrdersDatasource interface")
		})
	}
}
