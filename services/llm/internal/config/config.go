package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds all configuration for the LLM service.
type Config struct {
	Port        int    // PORT (default: 9003)
	PostgresDSN string // POSTGRES_DSN
	NatsURL     string // NATS_URL
	LogLevel    string // LOG_LEVEL (default: info)
	EnableDev   bool   // ENABLE_DEV_ENDPOINTS (default: false)

	// Bifrost provider settings (at least one must be set for LLM calls to work).
	OpenAIAPIKey    string // OPENAI_API_KEY
	AnthropicAPIKey string // ANTHROPIC_API_KEY
	DefaultModel    string // DEFAULT_MODEL (default: gpt-4o-mini)
	DefaultProvider string // DEFAULT_PROVIDER (default: openai)
}

// Load reads configuration from environment variables with sensible defaults.
func Load() (*Config, error) {
	port, err := getEnvInt("PORT", 9003)
	if err != nil {
		return nil, fmt.Errorf("invalid PORT: %w", err)
	}

	defaultModel := getEnv("DEFAULT_MODEL", "gpt-4o-mini")
	defaultProvider := getEnv("DEFAULT_PROVIDER", "openai")

	return &Config{
		Port:            port,
		PostgresDSN:     os.Getenv("POSTGRES_DSN"),
		NatsURL:         getEnv("NATS_URL", "nats://localhost:4222"),
		LogLevel:        getEnv("LOG_LEVEL", "info"),
		EnableDev:       getEnvBool("ENABLE_DEV_ENDPOINTS", false),
		OpenAIAPIKey:    os.Getenv("OPENAI_API_KEY"),
		AnthropicAPIKey: os.Getenv("ANTHROPIC_API_KEY"),
		DefaultModel:    defaultModel,
		DefaultProvider: defaultProvider,
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
