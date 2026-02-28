package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds all configuration for the BFF service loaded from environment variables.
type Config struct {
	Port        int
	PostgresDSN string
	RedisURL    string
	NatsURL     string
	MinIOURL    string
	MinIOUser   string
	MinIOPass   string
	LogLevel    string
}

// Load reads configuration from environment variables with sensible defaults.
func Load() (*Config, error) {
	port, err := getEnvInt("PORT", 8080)
	if err != nil {
		return nil, fmt.Errorf("invalid PORT: %w", err)
	}

	return &Config{
		Port:        port,
		PostgresDSN: getEnv("POSTGRES_DSN", "postgres://axle:axle@localhost:5432/axle?sslmode=disable"),
		RedisURL:    getEnv("REDIS_URL", "redis://localhost:6379"),
		NatsURL:     getEnv("NATS_URL", "nats://localhost:4222"),
		MinIOURL:    getEnv("MINIO_URL", "http://localhost:9000"),
		MinIOUser:   getEnv("MINIO_USER", "minioadmin"),
		MinIOPass:   getEnv("MINIO_PASS", "minioadmin"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
	}, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) (int, error) {
	v := os.Getenv(key)
	if v == "" {
		return fallback, nil
	}
	return strconv.Atoi(v)
}
