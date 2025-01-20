package mocks

import (
	"fp_kata/common/constants"
	"fp_kata/internal/models"
	"fp_kata/pkg/log"
	"github.com/gofiber/fiber/v3"
	zlog "github.com/rs/zerolog/log"
)

type CustomCtx struct {
	fiber.DefaultCtx
	MockLocals map[any]any
}

func (c *CustomCtx) Locals(key any, value ...any) any {

	if len(value) == 0 {
		return c.MockLocals[key]
	}
	c.MockLocals[key] = value[0]
	return value[0]
}

func ProvideBaseMockContextData(user *models.User) *map[any]any {
	contextData := make(map[any]any)
	contextData[log.LogFiberContextKey] = &zlog.Logger
	if user != nil {
		contextData[constants.AuthenticatedUserIdKey] = user.ID
		contextData[constants.AuthenticatedUserKey] = *user
	}
	return &contextData
}
