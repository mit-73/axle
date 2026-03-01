package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds all configuration for the Gateway service.
type Config struct {
	Port     int
	NatsURL  string
	RedisURL string
	LogLevel string
}

// Load reads configuration from environment variables with sensible defaults.
func Load() (*Config, error) {
	port, err := getEnvInt("PORT", 9002)
	if err != nil {
		return nil, fmt.Errorf("invalid PORT: %w", err)
	}

	return &Config{
		Port:     port,
		NatsURL:  getEnv("NATS_URL", "nats://localhost:4222"),
		RedisURL: getEnv("REDIS_URL", "redis://localhost:6379"),
		LogLevel: getEnv("LOG_LEVEL", "info"),
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
