// Code generated by mockery v2.33.3. DO NOT EDIT.

package mocks

import (
	context "context"

	dsmodels "fp_kata/internal/datasources/dsmodels"

	mo "github.com/samber/mo"

	mock "github.com/stretchr/testify/mock"
)

// OrdersDatasource is an autogenerated mock type for the OrdersDatasource type
type OrdersDatasource struct {
	mock.Mock
}

// DeleteOrder provides a mock function with given fields: ctx, orderID
func (_m *OrdersDatasource) DeleteOrder(ctx context.Context, orderID int) error {
	ret := _m.Called(ctx, orderID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int) error); ok {
		r0 = rf(ctx, orderID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAllOrdersForUser provides a mock function with given fields: ctx, userID
func (_m *OrdersDatasource) GetAllOrdersForUser(ctx context.Context, userID int) ([]dsmodels.Order, error) {
	ret := _m.Called(ctx, userID)

	var r0 []dsmodels.Order
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int) ([]dsmodels.Order, error)); ok {
		return rf(ctx, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int) []dsmodels.Order); ok {
		r0 = rf(ctx, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]dsmodels.Order)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetOrder provides a mock function with given fields: ctx, orderID
func (_m *OrdersDatasource) GetOrder(ctx context.Context, orderID int) mo.Result[dsmodels.Order] {
	ret := _m.Called(ctx, orderID)

	var r0 mo.Result[dsmodels.Order]
	if rf, ok := ret.Get(0).(func(context.Context, int) mo.Result[dsmodels.Order]); ok {
		r0 = rf(ctx, orderID)
	} else {
		r0 = ret.Get(0).(mo.Result[dsmodels.Order])
	}

	return r0
}

// InsertOrder provides a mock function with given fields: ctx, order
func (_m *OrdersDatasource) InsertOrder(ctx context.Context, order dsmodels.Order) (*dsmodels.Order, error) {
	ret := _m.Called(ctx, order)

	var r0 *dsmodels.Order
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, dsmodels.Order) (*dsmodels.Order, error)); ok {
		return rf(ctx, order)
	}
	if rf, ok := ret.Get(0).(func(context.Context, dsmodels.Order) *dsmodels.Order); ok {
		r0 = rf(ctx, order)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dsmodels.Order)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, dsmodels.Order) error); ok {
		r1 = rf(ctx, order)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateOrder provides a mock function with given fields: ctx, order
func (_m *OrdersDatasource) UpdateOrder(ctx context.Context, order dsmodels.Order) (*dsmodels.Order, error) {
	ret := _m.Called(ctx, order)

	var r0 *dsmodels.Order
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, dsmodels.Order) (*dsmodels.Order, error)); ok {
		return rf(ctx, order)
	}
	if rf, ok := ret.Get(0).(func(context.Context, dsmodels.Order) *dsmodels.Order); ok {
		r0 = rf(ctx, order)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dsmodels.Order)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, dsmodels.Order) error); ok {
		r1 = rf(ctx, order)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewOrdersDatasource creates a new instance of OrdersDatasource. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewOrdersDatasource(t interface {
	mock.TestingT
	Cleanup(func())
}) *OrdersDatasource {
	mock := &OrdersDatasource{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
