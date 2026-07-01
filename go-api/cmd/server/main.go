// Command server is the entrypoint of the Go QR API. It loads configuration,
// connects to PostgreSQL, runs migrations, wires the dependencies (manual DI)
// and starts the Fiber HTTP server.
package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/xangel1111/qr-stats-fullstack/go-api/db"
	"github.com/xangel1111/qr-stats-fullstack/go-api/internal/client"
	"github.com/xangel1111/qr-stats-fullstack/go-api/internal/config"
	"github.com/xangel1111/qr-stats-fullstack/go-api/internal/repository"
	"github.com/xangel1111/qr-stats-fullstack/go-api/internal/server"
	"github.com/xangel1111/qr-stats-fullstack/go-api/internal/service"
)

func main() {
	// Load .env for local development; ignored if the file is absent.
	_ = godotenv.Load()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	cfg, err := config.Load()
	if err != nil {
		logger.Error("configuration error", "error", err)
		os.Exit(1)
	}

	if cfg.AutoMigrate {
		if err := repository.Migrate(cfg.DatabaseURL, db.Migrations); err != nil {
			logger.Error("migration error", "error", err)
			os.Exit(1)
		}
		logger.Info("database migrations applied")
	}

	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		logger.Error("database connection error", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	// Manual dependency injection (constructor wiring).
	repo := repository.New(pool)
	statsClient := client.NewStatsClient(cfg.NodeStatsURL, cfg.NodeTimeout)
	qrService := service.NewQRService(statsClient, repo)

	app := server.New(server.Deps{
		Logger:       logger,
		JWTSecret:    cfg.JWTSecret,
		JWTExpiry:    cfg.JWTExpiry,
		AuthUsername: cfg.AuthUsername,
		AuthPassword: cfg.AuthPassword,
		QRService:    qrService,
		Repo:         repo,
		DB:           pool,
	})

	// Listen on all interfaces, dual-stack (IPv4 + IPv6). Some platforms route
	// to services over IPv6 (e.g. Railway private networking), so binding only
	// to 0.0.0.0 would make the service unreachable there.
	addr := "[::]:" + cfg.Port
	logger.Info("starting server", "addr", addr)
	if err := app.Listen(addr); err != nil {
		logger.Error("server error", "error", err)
		os.Exit(1)
	}
}
