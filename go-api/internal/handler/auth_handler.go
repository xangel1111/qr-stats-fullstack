// Package handler holds the Fiber HTTP handlers. Handlers only deal with
// transport concerns (parse, delegate, respond); business logic lives in the
// service and domain layers.
package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/xangel1111/qr-stats-fullstack/go-api/internal/apperr"
	"github.com/xangel1111/qr-stats-fullstack/go-api/internal/dto"
	"github.com/xangel1111/qr-stats-fullstack/go-api/internal/middleware"
)

// AuthHandler issues JWTs for a single configured demo user.
type AuthHandler struct {
	secret   string
	expiry   time.Duration
	username string
	password string
}

func NewAuthHandler(secret string, expiry time.Duration, username, password string) *AuthHandler {
	return &AuthHandler{secret: secret, expiry: expiry, username: username, password: password}
}

// Login validates credentials and returns a signed JWT.
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return apperr.BadRequest("invalid request body")
	}

	if req.Username != h.username || req.Password != h.password {
		return apperr.Unauthorized("invalid credentials")
	}

	token, err := middleware.GenerateToken(h.secret, req.Username, h.expiry)
	if err != nil {
		return apperr.Internal("could not generate token")
	}

	return c.JSON(dto.LoginResponse{
		TokenType: "Bearer",
		Token:     token,
		ExpiresIn: int(h.expiry.Seconds()),
	})
}
