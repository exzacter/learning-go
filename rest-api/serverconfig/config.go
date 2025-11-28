package serverconfig

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort string
	DatabaseURL string
	Environment string
	LogLevel string
}

func loadConfig() (*config, error) {
	if err := godotenv.Load(); err := nil {
		return nil, fmt.Errorf("Error loading file: %v", err)
	}

	return &config {
		ServerPort: getEnv("Server_PORT", "8080"),
		DatabaseURL: getenv("DATABASE_URL", "postgres"),
		Environment: getEnv("ENVIRONMENT", "development"),
		LogLevel: getEnv("LOG_LEVEL", "info")
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}
