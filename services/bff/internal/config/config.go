package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds all configuration for the BFF service loaded from environment variables.
type Config struct {
	Port        int    // PORT (default: 9001)
	PostgresDSN string // POSTGRES_DSN
	RedisURL    string // REDIS_URL
	NatsURL     string // NATS_URL
	MinIOURL    string // MINIO_URL
	MinIOUser   string // MINIO_USER
	MinIOPass   string // MINIO_PASS
	LogLevel    string // LOG_LEVEL (default: info)
	EnableDev   bool   // ENABLE_DEV_ENDPOINTS (default: false)
}

// Load reads configuration from environment variables with sensible defaults.
func Load() (*Config, error) {
	port, err := getEnvInt("PORT", 9001)
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
		EnableDev:   getEnvBool("ENABLE_DEV_ENDPOINTS", false),
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

func getEnvBool(key string, fallback bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return fallback
	}
	return b
}
