package config

import "os"

type Config struct {
	Port          string
	DatabaseURL   string
	MigrationsDir string
	ValkeyURL     string
	NatsURL       string
	OneCBaseURL   string // 1С HTTP-service base URL; empty = mock mode
	OneCUser      string
	OneCPassword  string
	CronSchedule  string // cron expression, default every 20 min
	Environment   string
}

func Load() *Config {
	return &Config{
		Port:          getEnv("PORT", "8082"),
		DatabaseURL:   getEnv("DATABASE_URL", "postgres://pkt:pkt_secret@localhost:5433/pkt_db?sslmode=disable"),
		MigrationsDir: getEnv("MIGRATIONS_DIR", "migrations"),
		ValkeyURL:     getEnv("VALKEY_URL", "redis://localhost:6380/0"),
		NatsURL:       getEnv("NATS_URL", "nats://localhost:4222"),
		OneCBaseURL:   getEnv("ONEC_BASE_URL", ""),
		OneCUser:      getEnv("ONEC_USER", ""),
		OneCPassword:  getEnv("ONEC_PASSWORD", ""),
		CronSchedule:  getEnv("CRON_SCHEDULE", "*/20 * * * *"),
		Environment:   getEnv("ENVIRONMENT", "development"),
	}
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
