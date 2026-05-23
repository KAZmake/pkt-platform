package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port            string
	Environment     string
	AnthropicAPIKey string
	AnthropicModel  string
	MaxTokens       int
	DirectusURL     string
	DirectusToken   string
	// RateLimitRPM is the max requests per minute per IP (0 = disabled).
	RateLimitRPM int
}

func Load() *Config {
	maxTokens, _ := strconv.Atoi(getEnv("MAX_TOKENS", "1024"))
	rateLimit, _ := strconv.Atoi(getEnv("RATE_LIMIT_RPM", "20"))
	return &Config{
		Port:            getEnv("PORT", "8083"),
		Environment:     getEnv("ENVIRONMENT", "development"),
		AnthropicAPIKey: os.Getenv("ANTHROPIC_API_KEY"),
		AnthropicModel:  getEnv("ANTHROPIC_MODEL", "claude-haiku-4-5-20251001"),
		MaxTokens:       maxTokens,
		DirectusURL:     getEnv("DIRECTUS_URL", "http://localhost:8055"),
		DirectusToken:   os.Getenv("DIRECTUS_TOKEN"),
		RateLimitRPM:    rateLimit,
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
