package env

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

func GetString(key string, defaultValue string) string {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		slog.Error("Error loading .env file", "error", err)
		os.Exit(1)
	}

	if val := os.Getenv(key); val != "" {
		return val
	}

	return defaultValue
}
