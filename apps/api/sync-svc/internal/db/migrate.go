package db

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(databaseURL, migrationsDir, migrationsTable string) error {
	pgxURL := strings.Replace(databaseURL, "postgres://", "pgx5://", 1)
	pgxURL = strings.Replace(pgxURL, "postgresql://", "pgx5://", 1)
	sep := "?"
	if strings.Contains(pgxURL, "?") {
		sep = "&"
	}
	pgxURL += sep + "x-migrations-table=" + migrationsTable

	m, err := migrate.New("file://"+migrationsDir, pgxURL)
	if err != nil {
		return fmt.Errorf("migrate.New: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			slog.Info("migrations: no changes")
			return nil
		}
		return fmt.Errorf("migrate.Up: %w", err)
	}
	slog.Info("migrations applied")
	return nil
}
