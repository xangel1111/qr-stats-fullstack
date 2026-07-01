package handler

import (
	"github.com/gofiber/fiber/v2"

	"github.com/xangel1111/qr-stats-fullstack/go-api/internal/apperr"
	"github.com/xangel1111/qr-stats-fullstack/go-api/internal/dto"
	"github.com/xangel1111/qr-stats-fullstack/go-api/internal/middleware"
	"github.com/xangel1111/qr-stats-fullstack/go-api/internal/service"
)

// QRHandler exposes the QR computation endpoint.
type QRHandler struct {
	service *service.QRService
}

func NewQRHandler(s *service.QRService) *QRHandler {
	return &QRHandler{service: s}
}

// Compute handles POST /api/v1/qr.
func (h *QRHandler) Compute(c *fiber.Ctx) error {
	var req dto.QRRequest
	if err := c.BodyParser(&req); err != nil {
		return apperr.BadRequest("invalid request body")
	}

	username, _ := c.Locals(middleware.UsernameKey).(string)
	token, _ := c.Locals(middleware.TokenKey).(string)

	res, err := h.service.Compute(c.Context(), req.Matrix, username, token)
	if err != nil {
		return err // already an *apperr.APIError; mapped by the error handler
	}

	return c.Status(fiber.StatusOK).JSON(res)
}
