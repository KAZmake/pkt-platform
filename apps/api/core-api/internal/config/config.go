package config

import "os"

type Config struct {
	Port          string
	DatabaseURL   string
	KeycloakURL   string
	KeycloakRealm string
	ValkeyURL     string
	Environment   string
}

func Load() *Config {
	return &Config{
		Port:          getEnv("PORT", "8080"),
		DatabaseURL:   getEnv("DATABASE_URL", "postgres://pkt:pkt_secret@localhost:5433/pkt_db?sslmode=disable"),
		KeycloakURL:   getEnv("KEYCLOAK_URL", "http://localhost:8080"),
		KeycloakRealm: getEnv("KEYCLOAK_REALM", "pkt"),
		ValkeyURL:     getEnv("VALKEY_URL", "redis://localhost:6380/0"),
		Environment:   getEnv("ENVIRONMENT", "development"),
	}
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
