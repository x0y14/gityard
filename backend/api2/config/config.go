package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config func to get env value
func Config(key string) string {
	// load .env file
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Print("Error loading .env file")
	}
	return os.Getenv(key)
}

const (
	AccessTokenActiveDurationMinutes  = 15          // 15mins
	RefreshTokenActiveDurationMinutes = 60 * 24 * 7 // 7days
)
