// Package middleware holds the Fiber middlewares: centralized error handling,
// JWT authentication and structured request logging.
package middleware

import (
	"errors"
	"log/slog"

	"github.com/gofiber/fiber/v2"

	"github.com/xangel1111/qr-stats-fullstack/go-api/internal/apperr"
)

type errorEnvelope struct {
	Error errorBody `json:"error"`
}

type errorBody struct {
	Code    string   `json:"code"`
	Message string   `json:"message"`
	Details []string `json:"details,omitempty"`
}

// ErrorHandler is Fiber's central error handler. It renders our typed
// APIError as the shared JSON envelope and falls back to 500 for anything
// unexpected (which it also logs).
func ErrorHandler(logger *slog.Logger) fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		var apiErr *apperr.APIError
		if errors.As(err, &apiErr) {
			return c.Status(apiErr.Status).JSON(errorEnvelope{
				Error: errorBody{Code: apiErr.Code, Message: apiErr.Message, Details: apiErr.Details},
			})
		}

		// Fiber's own errors (e.g. unmatched route -> 404).
		var fiberErr *fiber.Error
		if errors.As(err, &fiberErr) {
			return c.Status(fiberErr.Code).JSON(errorEnvelope{
				Error: errorBody{Code: "ERROR", Message: fiberErr.Message},
			})
		}

		logger.Error("unhandled error", "path", c.Path(), "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(errorEnvelope{
			Error: errorBody{Code: "INTERNAL_ERROR", Message: "internal server error"},
		})
	}
}
