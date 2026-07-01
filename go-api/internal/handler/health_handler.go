package handler

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Pinger is the minimal database dependency the health check needs. Both
// *pgxpool.Pool and test fakes satisfy it.
type Pinger interface {
	Ping(ctx context.Context) error
}

// HealthHandler reports service readiness, including database connectivity.
type HealthHandler struct {
	db Pinger
}

func NewHealthHandler(db Pinger) *HealthHandler {
	return &HealthHandler{db: db}
}

// Health handles GET /health.
func (h *HealthHandler) Health(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 2*time.Second)
	defer cancel()

	dbStatus := "ok"
	status := "ok"
	code := fiber.StatusOK

	if err := h.db.Ping(ctx); err != nil {
		dbStatus = "error"
		status = "degraded"
		code = fiber.StatusServiceUnavailable
	}

	return c.Status(code).JSON(fiber.Map{
		"status": status,
		"checks": fiber.Map{"database": dbStatus},
	})
}
