// Code generated by mockery v2.33.3. DO NOT EDIT.

package mocks

import (
	context "context"
	models "fp_kata/internal/models"

	mock "github.com/stretchr/testify/mock"
)

// UsersService is an autogenerated mock type for the UsersService type
type UsersService struct {
	mock.Mock
}

// GetUserByID provides a mock function with given fields: ctx, id
func (_m *UsersService) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	ret := _m.Called(ctx, id)

	var r0 *models.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int) (*models.User, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int) *models.User); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SignUp provides a mock function with given fields: ctx, user
func (_m *UsersService) SignUp(ctx context.Context, user models.User) (*models.User, error) {
	ret := _m.Called(ctx, user)

	var r0 *models.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, models.User) (*models.User, error)); ok {
		return rf(ctx, user)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.User) *models.User); ok {
		r0 = rf(ctx, user)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, models.User) error); ok {
		r1 = rf(ctx, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewUsersService creates a new instance of UsersService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUsersService(t interface {
	mock.TestingT
	Cleanup(func())
}) *UsersService {
	mock := &UsersService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
