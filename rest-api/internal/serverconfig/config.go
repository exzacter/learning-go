package serverconfig

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort  string
	DatabaseURL string
	Environment string
	LogLevel    string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("Error loading file: %v", err)
	}

	return &Config{
		ServerPort:  GetEnv("SERVER_PORT", "8080"),
		DatabaseURL: GetEnv("DATABASE_URL", "postgres"),
		Environment: GetEnv("ENVIRONMENT", "development"),
		LogLevel:    GetEnv("LOG_LEVEL", "info"),
	}, nil
}

func GetEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}
