package db

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// RunMigrations applies all pending UP migrations from the given directory.
// migrationsTable is the postgres table used to track applied versions
// (e.g. "schema_migrations_core").
func RunMigrations(databaseURL, migrationsDir, migrationsTable string) error {
	// golang-migrate expects pgx5:// scheme for pgx/v5 driver
	pgxURL := "pgx5://" + databaseURL[len("postgres://"):]
	if len(databaseURL) >= 11 && databaseURL[:11] == "postgresql:" {
		pgxURL = "pgx5://" + databaseURL[len("postgresql://"):]
	}

	// Append x-migrations-table to URL so each service tracks its own versions
	sep := "?"
	for _, c := range pgxURL {
		if c == '?' {
			sep = "&"
			break
		}
	}
	pgxURL = fmt.Sprintf("%s%sx-migrations-table=%s", pgxURL, sep, migrationsTable)

	m, err := migrate.New("file://"+migrationsDir, pgxURL)
	if err != nil {
		return fmt.Errorf("migrate.New: %w", err)
	}
	defer func() {
		srcErr, dbErr := m.Close()
		if srcErr != nil {
			slog.Warn("migrate source close error", "error", srcErr)
		}
		if dbErr != nil {
			slog.Warn("migrate db close error", "error", dbErr)
		}
	}()

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			slog.Info("migrations: no changes", "table", migrationsTable)
			return nil
		}
		return fmt.Errorf("migrate.Up: %w", err)
	}

	version, dirty, _ := m.Version()
	slog.Info("migrations applied", "table", migrationsTable, "version", version, "dirty", dirty)
	return nil
}
