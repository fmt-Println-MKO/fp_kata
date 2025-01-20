package controllers

import (
	"fp_kata/common/constants"
	"fp_kata/internal/services"
	"fp_kata/pkg/log"
	"fp_kata/pkg/transports"
	"github.com/gofiber/fiber/v3"
)

// UsersController handles user-related operations
type UsersController struct {
	userService services.UsersService
}

// NewUsersController creates a new instance of UsersController
func NewUsersController(userService services.UsersService) UsersController {
	return UsersController{userService: userService}
}

// RegisterUserRoutes registers the routes for UsersController
func (c *UsersController) RegisterUserRoutes(app *fiber.App, authMiddleware fiber.Handler) {
	app.Post("/users", c.SignUp)
	app.Get("/users/me", c.GetUser, authMiddleware)
}

// SignUp creates a new user
func (c *UsersController) SignUp(ctx fiber.Ctx) error {
	logger := log.GetFiberLogger(ctx).With().Logger()

	logger.Debug().Str("comp", "UsersController").Str("func", "SignUp").Send()

	log.SetFiberLogger(ctx, &logger)
	userRequest := new(transports.UserCreateRequest)

	if err := ctx.Bind().Body(userRequest); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	user := userRequest.ToUser()
	if user == nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}
	context := log.NewBackgroundContext(&logger)

	newUser, err := c.userService.SignUp(context, *user)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not create user",
		})
	}
	return ctx.Status(fiber.StatusCreated).JSON(transports.MapToUserResponse(*newUser))
}

// GetUser retrieves a user by their ID
func (c *UsersController) GetUser(ctx fiber.Ctx) error {
	logger := log.GetFiberLogger(ctx).With().Logger()
	logger.Debug().Str("comp", "UsersController").Str("func", "GetUser").Send()

	userIdValue := ctx.Locals(constants.AuthenticatedUserIdKey)
	if userIdValue == nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token is missing",
		})
	}
	userId := userIdValue.(int)

	context := log.NewBackgroundContext(&logger)
	user, err := c.userService.GetUserByID(context, userId)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not retrieve user",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(transports.MapToUserResponse(*user))
}
