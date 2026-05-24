package repository_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/KAZmake/pkt-platform/apps/api/core-api/internal/db"
	"github.com/KAZmake/pkt-platform/apps/api/core-api/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var testPool *pgxpool.Pool

func TestMain(m *testing.M) {
	ctx := context.Background()

	_, thisFile, _, _ := runtime.Caller(0)
	pkgDir := filepath.Dir(thisFile)

	pgc, err := tcpostgres.Run(ctx,
		"timescale/timescaledb-ha:pg16",
		tcpostgres.WithDatabase("testdb"),
		tcpostgres.WithUsername("test"),
		tcpostgres.WithPassword("test"),
		tcpostgres.WithInitScripts(filepath.Join(pkgDir, "testdata", "init.sql")),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(90*time.Second),
		),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "start postgres container: %v\n", err)
		os.Exit(1)
	}
	defer pgc.Terminate(ctx) //nolint:errcheck

	connStr, err := pgc.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		fmt.Fprintf(os.Stderr, "connection string: %v\n", err)
		os.Exit(1)
	}

	migrationsDir := filepath.Join(pkgDir, "..", "..", "migrations")
	if err := db.RunMigrations(connStr, migrationsDir, "schema_migrations_test"); err != nil {
		fmt.Fprintf(os.Stderr, "run migrations: %v\n", err)
		os.Exit(1)
	}

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "create pool: %v\n", err)
		os.Exit(1)
	}
	testPool = pool
	defer pool.Close()

	os.Exit(m.Run())
}

// ── seed helpers ──────────────────────────────────────────────────────────────

func insertUser(t *testing.T) *model.User {
	t.Helper()
	ctx := context.Background()
	u := &model.User{}
	err := testPool.QueryRow(ctx, `
		INSERT INTO users (keycloak_id, email, role, first_name, last_name)
		VALUES ($1, $2, 'borrower', 'Test', 'User')
		RETURNING id, keycloak_id, email, role, first_name, last_name, phone, created_at, updated_at`,
		"kc-"+uuid.New().String(),
		"test-"+uuid.New().String()+"@example.com",
	).Scan(&u.ID, &u.KeycloakID, &u.Email, &u.Role,
		&u.FirstName, &u.LastName, &u.Phone,
		&u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		t.Fatalf("insertUser: %v", err)
	}
	return u
}

func insertBorrower(t *testing.T, userID uuid.UUID) *model.Borrower {
	t.Helper()
	ctx := context.Background()
	b := &model.Borrower{}
	err := testPool.QueryRow(ctx, `
		INSERT INTO borrowers (user_id, inn)
		VALUES ($1, $2)
		RETURNING id, user_id, inn, bin, org_name, activity_type, farm_id, created_at`,
		userID,
		fmt.Sprintf("%012d", time.Now().UnixNano()%1_000_000_000_000),
	).Scan(&b.ID, &b.UserID, &b.INN, &b.BIN,
		&b.OrgName, &b.ActivityType, &b.FarmID, &b.CreatedAt)
	if err != nil {
		t.Fatalf("insertBorrower: %v", err)
	}
	return b
}

func insertProgram(t *testing.T) *model.LoanProgram {
	t.Helper()
	ctx := context.Background()
	p := &model.LoanProgram{}
	err := testPool.QueryRow(ctx, `
		INSERT INTO loan_programs
		  (name, rate, min_amount, max_amount, min_term_months, max_term_months)
		VALUES ($1, 7.5, 100000, 5000000, 6, 60)
		RETURNING id, name, name_kz, name_en, rate,
		          min_amount, max_amount, min_term_months, max_term_months,
		          activity_types, is_active, created_at, updated_at`,
		"Агро "+uuid.New().String()[:8],
	).Scan(&p.ID, &p.Name, &p.NameKZ, &p.NameEN, &p.Rate,
		&p.MinAmount, &p.MaxAmount, &p.MinTermMonths, &p.MaxTermMonths,
		&p.ActivityTypes, &p.IsActive, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		t.Fatalf("insertProgram: %v", err)
	}
	return p
}
