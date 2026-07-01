// Package server builds the Fiber application from its dependencies. Extracting
// this from main lets tests exercise the full HTTP stack (via app.Test) with
// mocked collaborators, without binding a port or a real database.
package server

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/xangel1111/qr-stats-fullstack/go-api/internal/handler"
	"github.com/xangel1111/qr-stats-fullstack/go-api/internal/middleware"
	"github.com/xangel1111/qr-stats-fullstack/go-api/internal/repository"
	"github.com/xangel1111/qr-stats-fullstack/go-api/internal/service"
)

// Deps holds everything the HTTP layer needs.
type Deps struct {
	Logger       *slog.Logger
	JWTSecret    string
	JWTExpiry    time.Duration
	AuthUsername string
	AuthPassword string
	QRService    *service.QRService
	Repo         repository.ComputationRepository
	DB           handler.Pinger
}

// New wires the handlers, middlewares and routes into a Fiber app.
func New(d Deps) *fiber.App {
	authHandler := handler.NewAuthHandler(d.JWTSecret, d.JWTExpiry, d.AuthUsername, d.AuthPassword)
	qrHandler := handler.NewQRHandler(d.QRService)
	compHandler := handler.NewComputationHandler(d.Repo)
	healthHandler := handler.NewHealthHandler(d.DB)

	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler(d.Logger),
	})
	app.Use(recover.New())
	app.Use(cors.New())
	app.Use(middleware.RequestLogger(d.Logger))

	// Public routes.
	app.Get("/health", healthHandler.Health)
	app.Post("/auth/login", authHandler.Login)

	// Protected routes.
	api := app.Group("/api/v1", middleware.JWTAuth(d.JWTSecret))
	api.Post("/qr", qrHandler.Compute)
	api.Get("/computations", compHandler.List)
	api.Get("/computations/:id", compHandler.Get)

	return app
}
