package config

import (
	"os"
	"time"
)

// Config holds the application configuration
type Config struct {
	Port          string
	JWTSecret     string
	JWTExpiration time.Duration
}

// Load reads configuration from environment variables with sensible defaults
func Load() *Config {
	return &Config{
		Port:          getEnv("PORT", "8080"),
		JWTSecret:     getEnv("JWT_SECRET", "default-secret-key"),
		JWTExpiration: parseDuration(getEnv("JWT_EXPIRATION", "24h")),
	}
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// parseDuration parses a duration string, returning a default on error
func parseDuration(durationStr string) time.Duration {
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 24 * time.Hour // Default to 24 hours
	}
	return duration
}
