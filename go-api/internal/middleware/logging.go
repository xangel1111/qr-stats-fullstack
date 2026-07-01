package middleware

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
)

// RequestLogger emits one structured log line per request.
func RequestLogger(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		logger.Info("request",
			"method", c.Method(),
			"path", c.Path(),
			"status", c.Response().StatusCode(),
			"duration_ms", time.Since(start).Milliseconds(),
		)
		return err
	}
}
