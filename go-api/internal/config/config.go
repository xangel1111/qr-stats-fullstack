// Package config loads runtime configuration from environment variables with
// sensible defaults. Required secrets fail fast at startup.
package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds all runtime configuration for the Go API.
type Config struct {
	Port         string
	DatabaseURL  string
	AutoMigrate  bool
	JWTSecret    string
	JWTExpiry    time.Duration
	NodeStatsURL string
	NodeTimeout  time.Duration
	AuthUsername string
	AuthPassword string
}

// Load reads configuration from the environment. It returns an error if a
// required value (JWT secret, database URL) is missing.
func Load() (*Config, error) {
	cfg := &Config{
		Port:         getEnv("PORT", "8080"),
		DatabaseURL:  getEnv("DATABASE_URL", ""),
		AutoMigrate:  getBool("AUTO_MIGRATE", true),
		JWTSecret:    getEnv("JWT_SECRET", ""),
		JWTExpiry:    getDuration("JWT_EXPIRY", time.Hour),
		NodeStatsURL: getEnv("NODE_STATS_URL", "http://localhost:3000"),
		NodeTimeout:  getDuration("NODE_TIMEOUT", 5*time.Second),
		AuthUsername: getEnv("AUTH_USERNAME", "demo"),
		AuthPassword: getEnv("AUTH_PASSWORD", "demo123"),
	}

	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}
	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

func getBool(key string, fallback bool) bool {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			return b
		}
	}
	return fallback
}

// getDuration accepts either a Go duration ("500ms", "5s") or a plain integer
// interpreted as seconds.
func getDuration(key string, fallback time.Duration) time.Duration {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		return fallback
	}
	if d, err := time.ParseDuration(v); err == nil {
		return d
	}
	if secs, err := strconv.Atoi(v); err == nil {
		return time.Duration(secs) * time.Second
	}
	return fallback
}
