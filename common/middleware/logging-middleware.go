package middleware

import (
	"fp_kata/pkg/log"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"time"
)

// Middleware for initializing zerolog with traceId
func LoggingMiddleware(logger *zerolog.Logger) fiber.Handler {
	return func(c fiber.Ctx) error {
		start := time.Now()

		traceID := uuid.New().String()
		reqLogger := logger.With().Str("traceId", traceID).Logger()

		log.SetFiberLogger(c, &reqLogger)

		err := c.Next()

		reqLogger.Info().
			Str("method", c.Method()).
			Str("path", c.Path()).
			Int("status", c.Response().StatusCode()).
			Dur("latency", time.Since(start)).
			Send()

		return err
	}
}
