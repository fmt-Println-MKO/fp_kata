package services

import (
	"fp_kata/internal/models"
	"fp_kata/pkg/log"
	zlog "github.com/rs/zerolog/log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthorizationService_IsAuthorized(t *testing.T) {
	log.InitLogger()
	ctx := log.NewBackgroundContext(&zlog.Logger)

	type testCase struct {
		name       string
		userId     int
		order      models.Order
		assertFunc func(t *testing.T, res bool, err error)
	}

	testCases := []testCase{
		{
			name:   "valid user ID matches order user",
			userId: 1,
			order:  models.Order{User: &models.User{ID: 1}},
			assertFunc: func(t *testing.T, res bool, err error) {

				assert.NoError(t, err, "Expected no error but got one")
				assert.Equal(t, true, res, "Expected result did not match")

			},
		},
		{
			name:   "valid user ID does not match order user",
			userId: 1,
			order:  models.Order{User: &models.User{ID: 2}},
			assertFunc: func(t *testing.T, res bool, err error) {

				assert.NoError(t, err, "Expected no error but got one")
				assert.Equal(t, false, res, "Expected result did not match")

			},
		},
		{
			name:   "user ID is zero",
			userId: 0,
			order:  models.Order{User: &models.User{ID: 1}},
			assertFunc: func(t *testing.T, res bool, err error) {

				assert.EqualError(t, err, "userId is required", "Expected no error but got one")
				assert.Equal(t, false, res, "Expected result did not match")

			},
		},
		{
			name:   "order has no user",
			userId: 1,
			order:  models.Order{User: nil},
			assertFunc: func(t *testing.T, res bool, err error) {

				assert.EqualError(t, err, "missing user on order", "Expected no error but got one")
				assert.Equal(t, false, res, "Expected result did not match")

			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc := NewAuthorizationService()
			res, err := svc.IsAuthorized(ctx, tc.userId, &tc.order)
			tc.assertFunc(t, res, err)
		})
	}
}
