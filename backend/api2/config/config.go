package config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

// Config func to get env value
func Config(key string) string {
	// load .env file
	err := godotenv.Load(".env")
	if err != nil {
		slog.Info("failed to load .env, load os env directly")
	}
	return os.Getenv(key)
}

const (
	AccessTokenActiveDurationMinutes  = 15          // 15mins
	RefreshTokenActiveDurationMinutes = 60 * 24 * 7 // 7days
)
