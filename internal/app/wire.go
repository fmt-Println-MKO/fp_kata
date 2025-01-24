//go:build wireinject
// +build wireinject

package app

import (
	"fp_kata/common/middleware"
	"fp_kata/internal/controllers"
	"fp_kata/internal/datasources/file"
	"fp_kata/internal/datasources/yugabyte"
	"fp_kata/internal/services"
	"github.com/gofiber/fiber/v3"
	"github.com/google/wire"
)

type AppModules struct {
	AuthMiddleware   fiber.Handler
	UsersController  controllers.UsersController
	OrdersController controllers.OrdersController
}

// Define a ProviderSet that provides AuthService once.
var AppModulesSet = wire.NewSet(
	// Dependencies used across multiple parts of the app.
	file.NewOrdersStorage,
	file.NewUsersStorage,
	yugabyte.NewPaymentsStorage,

	// Services
	services.NewAuthService,
	services.NewUsersService,
	services.NewPaymentsService,
	services.NewOrdersService,
	services.NewAuthorizationService,

	// Controllers
	controllers.NewUsersController,
	controllers.NewOrdersController,

	// Middleware
	middleware.AuthMiddleware,

	// Use our “newAppModules” provider function below to tie it all together.
	newAppModules,
)

// newAppModules ties together all the pieces into a single struct.
func newAppModules(
	authMW fiber.Handler,
	usersCtrl controllers.UsersController,
	ordersCtrl controllers.OrdersController,
) *AppModules {
	return &AppModules{
		AuthMiddleware:   authMW,
		UsersController:  usersCtrl,
		OrdersController: ordersCtrl,
	}
}

// InitializeAppModules wires up the entire application in one go.
func InitializeAppModules() *AppModules {
	wire.Build(AppModulesSet)
	return &AppModules{} // This return is never reached; Wire will generate the code.
}
