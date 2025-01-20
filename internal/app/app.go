package app

import (
	"fp_kata/common/middleware"
	fpLog "fp_kata/pkg/log"
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog/log"
)

func InitApp() *fiber.App {
	fpLog.InitLogger()
	log.Info().Msg("Hello non-functional go")
	app := fiber.New(fiber.Config{
		AppName:     "fp_kata",
		JSONEncoder: sonic.Marshal,
		JSONDecoder: sonic.Unmarshal,
	})

	app.Use(middleware.LoggingMiddleware(&log.Logger))

	appModules := InitializeAppModules()
	appModules.OrdersController.RegisterOrderRoutes(app, appModules.AuthMiddleware)
	appModules.UsersController.RegisterUserRoutes(app, appModules.AuthMiddleware)
	return app
}
