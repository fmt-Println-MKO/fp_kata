package log

import (
	"context"
	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

const LogFiberContextKey = "logger"
const LogContextKey = "logger"

func InitLogger() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func GetFiberLogger(ctx fiber.Ctx) *zerolog.Logger {
	return ctx.Locals(LogFiberContextKey).(*zerolog.Logger)
}

func SetFiberLogger(ctx fiber.Ctx, logger *zerolog.Logger) {
	ctx.Locals(LogFiberContextKey, logger)
}

func GetLogger(ctx context.Context) *zerolog.Logger {
	return ctx.Value(LogContextKey).(*zerolog.Logger)
}

func NewBackgroundContext(logger *zerolog.Logger) context.Context {
	return context.WithValue(context.Background(), LogContextKey, logger)
}
