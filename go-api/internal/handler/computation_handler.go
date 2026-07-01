package handler

import (
	"errors"

	"github.com/gofiber/fiber/v2"

	"github.com/xangel1111/qr-stats-fullstack/go-api/internal/apperr"
	"github.com/xangel1111/qr-stats-fullstack/go-api/internal/dto"
	"github.com/xangel1111/qr-stats-fullstack/go-api/internal/repository"
)

const (
	defaultLimit = 20
	maxLimit     = 100
)

// ComputationHandler exposes the computation history (audit) endpoints.
type ComputationHandler struct {
	repo repository.ComputationRepository
}

func NewComputationHandler(repo repository.ComputationRepository) *ComputationHandler {
	return &ComputationHandler{repo: repo}
}

// Get handles GET /api/v1/computations/:id.
func (h *ComputationHandler) Get(c *fiber.Ctx) error {
	id := c.Params("id")

	res, err := h.repo.GetByID(c.Context(), id)
	if errors.Is(err, repository.ErrNotFound) {
		return apperr.NotFound("computation not found")
	}
	if err != nil {
		return apperr.Internal("could not fetch computation")
	}

	return c.JSON(res)
}

// List handles GET /api/v1/computations with pagination.
func (h *ComputationHandler) List(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", defaultLimit)
	offset := c.QueryInt("offset", 0)
	if limit <= 0 || limit > maxLimit {
		limit = defaultLimit
	}
	if offset < 0 {
		offset = 0
	}

	items, total, err := h.repo.List(c.Context(), limit, offset)
	if err != nil {
		return apperr.Internal("could not list computations")
	}

	return c.JSON(dto.ComputationList{
		Items:  items,
		Limit:  limit,
		Offset: offset,
		Total:  total,
	})
}
