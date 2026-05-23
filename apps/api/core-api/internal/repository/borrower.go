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

type BorrowerRepository struct {
	db *pgxpool.Pool
}

func NewBorrowerRepository(db *pgxpool.Pool) *BorrowerRepository {
	return &BorrowerRepository{db: db}
}

func (r *BorrowerRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*model.Borrower, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, user_id, inn, bin, org_name, activity_type, farm_id, created_at
		FROM borrowers WHERE user_id = $1`, userID)
	return scanBorrower(row)
}

func (r *BorrowerRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Borrower, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, user_id, inn, bin, org_name, activity_type, farm_id, created_at
		FROM borrowers WHERE id = $1`, id)
	return scanBorrower(row)
}

type UpdateBorrowerInput struct {
	OrgName      *string
	ActivityType *string
	Phone        *string // stored on users table, passed for convenience
}

func (r *BorrowerRepository) Update(ctx context.Context, id uuid.UUID, inp UpdateBorrowerInput) (*model.Borrower, error) {
	row := r.db.QueryRow(ctx, `
		UPDATE borrowers SET
		    org_name      = COALESCE($2, org_name),
		    activity_type = COALESCE($3, activity_type)
		WHERE id = $1
		RETURNING id, user_id, inn, bin, org_name, activity_type, farm_id, created_at`,
		id, inp.OrgName, inp.ActivityType,
	)
	return scanBorrower(row)
}

// ── helpers ──────────────────────────────────────────────────────────────────

func scanBorrower(s scanner) (*model.Borrower, error) {
	b := &model.Borrower{}
	err := s.Scan(
		&b.ID, &b.UserID, &b.INN, &b.BIN,
		&b.OrgName, &b.ActivityType, &b.FarmID, &b.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("scan borrower: %w", err)
	}
	return b, nil
}
