package config

import (
	"os"
	"strings"
)

type Config struct {
	Port           string
	DatabaseURL    string
	MigrationsDir  string
	KeycloakURL    string
	KeycloakRealm  string
	ValkeyURL      string
	NatsURL        string
	MinioEndpoint  string
	MinioAccessKey string
	MinioSecretKey string
	MinioUseSSL    bool
	ResendAPIKey   string
	EmailFrom      string
	CabinetURL     string
	Environment    string
}

func Load() *Config {
	return &Config{
		Port:           getEnv("PORT", "8080"),
		DatabaseURL:    getEnv("DATABASE_URL", "postgres://pkt:pkt_secret@localhost:5433/pkt_db?sslmode=disable"),
		MigrationsDir:  getEnv("MIGRATIONS_DIR", "migrations"),
		KeycloakURL:    getEnv("KEYCLOAK_URL", "http://localhost:8080"),
		KeycloakRealm:  getEnv("KEYCLOAK_REALM", "pkt"),
		ValkeyURL:      getEnv("VALKEY_URL", "redis://localhost:6380/0"),
		NatsURL:        getEnv("NATS_URL", "nats://localhost:4222"),
		MinioEndpoint:  getEnv("MINIO_ENDPOINT", "localhost:9000"),
		MinioAccessKey: getEnv("MINIO_ACCESS_KEY", "minioadmin"),
		MinioSecretKey: getEnv("MINIO_SECRET_KEY", "minioadmin123"),
		MinioUseSSL:    strings.ToLower(getEnv("MINIO_USE_SSL", "false")) == "true",
		ResendAPIKey:   getEnv("RESEND_API_KEY", ""),
		EmailFrom:      getEnv("EMAIL_FROM", "no-reply@pkt.kz"),
		CabinetURL:     getEnv("CABINET_URL", "http://localhost:3000/cabinet"),
		Environment:    getEnv("ENVIRONMENT", "development"),
	}
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
