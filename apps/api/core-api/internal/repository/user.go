package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/KAZmake/pkt-platform/apps/api/core-api/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetByKeycloakID(ctx context.Context, keycloakID string) (*model.User, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, keycloak_id, email, role,
		       first_name, last_name, phone,
		       created_at, updated_at
		FROM users WHERE keycloak_id = $1`, keycloakID)
	return scanUser(row)
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, keycloak_id, email, role,
		       first_name, last_name, phone,
		       created_at, updated_at
		FROM users WHERE id = $1`, id)
	return scanUser(row)
}

// Upsert inserts or updates a user by keycloak_id.
// Used for Keycloak sync on every authenticated request.
func (r *UserRepository) Upsert(ctx context.Context, u *model.User) (*model.User, error) {
	row := r.db.QueryRow(ctx, `
		INSERT INTO users (keycloak_id, email, role, first_name, last_name, phone)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (keycloak_id) DO UPDATE SET
		    email      = EXCLUDED.email,
		    role       = EXCLUDED.role,
		    first_name = COALESCE(EXCLUDED.first_name, users.first_name),
		    last_name  = COALESCE(EXCLUDED.last_name,  users.last_name),
		    updated_at = NOW()
		RETURNING id, keycloak_id, email, role,
		          first_name, last_name, phone,
		          created_at, updated_at`,
		u.KeycloakID, u.Email, u.Role, u.FirstName, u.LastName, u.Phone,
	)
	return scanUser(row)
}

type UpdateUserInput struct {
	FirstName *string
	LastName  *string
	Phone     *string
}

func (r *UserRepository) Update(ctx context.Context, id uuid.UUID, inp UpdateUserInput) (*model.User, error) {
	row := r.db.QueryRow(ctx, `
		UPDATE users SET
		    first_name = COALESCE($2, first_name),
		    last_name  = COALESCE($3, last_name),
		    phone      = COALESCE($4, phone),
		    updated_at = NOW()
		WHERE id = $1
		RETURNING id, keycloak_id, email, role,
		          first_name, last_name, phone,
		          created_at, updated_at`,
		id, inp.FirstName, inp.LastName, inp.Phone,
	)
	return scanUser(row)
}

func (r *UserRepository) List(ctx context.Context, limit, offset int) ([]*model.User, int, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, keycloak_id, email, role,
		       first_name, last_name, phone,
		       created_at, updated_at
		FROM users ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("users list query: %w", err)
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		u, err := scanUser(rows)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, u)
	}

	var total int
	r.db.QueryRow(ctx, `SELECT COUNT(*) FROM users`).Scan(&total) //nolint:errcheck

	return users, total, nil
}

// ── helpers ──────────────────────────────────────────────────────────────────

type scanner interface {
	Scan(dest ...any) error
}

func scanUser(s scanner) (*model.User, error) {
	u := &model.User{}
	err := s.Scan(
		&u.ID, &u.KeycloakID, &u.Email, &u.Role,
		&u.FirstName, &u.LastName, &u.Phone,
		&u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("scan user: %w", err)
	}
	return u, nil
}
