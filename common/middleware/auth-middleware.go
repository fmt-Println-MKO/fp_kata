package middleware

import (
	"fp_kata/common/constants"
	"fp_kata/internal/services"
	"fp_kata/pkg/log"
	"github.com/gofiber/fiber/v3"
)

func AuthMiddleware(authService services.AuthService, userService services.UsersService) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		logger := log.GetFiberLogger(ctx).With().Logger()
		context := log.NewBackgroundContext(&logger)
		// Get the token from the Authorization header
		authHeader := ctx.Get("Authorization")
		if authHeader == "" {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authorization token is missing",
			})
		}

		// Extract the token
		token := authHeader

		// Use the UsersService to load the user
		userId, err := authService.GetUserIDByToken(context, token)
		if err != nil {
			logger.Error().Err(err).Msg("Error getting user ID from token")
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		user, err := userService.GetUserByID(context, userId)
		if err != nil {
			logger.Error().Err(err).Msg("Error getting user from token")
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}
		logger = logger.With().Int("userId", user.ID).Logger()
		log.SetFiberLogger(ctx, &logger)

		// Add the user to the Fiber context
		ctx.Locals(constants.AuthenticatedUserKey, *user)
		ctx.Locals(constants.AuthenticatedUserIdKey, user.ID)

		// Proceed to the next handler
		return ctx.Next()
	}
}
