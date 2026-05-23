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

type ApplicationRepository struct {
	db *pgxpool.Pool
}

func NewApplicationRepository(db *pgxpool.Pool) *ApplicationRepository {
	return &ApplicationRepository{db: db}
}

const appCols = `id, borrower_id, program_id, assignee_id, status,
	amount, term_months, payment_type, created_at, updated_at`

type CreateApplicationInput struct {
	BorrowerID  uuid.UUID `json:"borrower_id"`
	ProgramID   uuid.UUID `json:"program_id"`
	Amount      float64   `json:"amount"`
	TermMonths  int       `json:"term_months"`
	PaymentType string    `json:"payment_type"`
}

func (r *ApplicationRepository) Create(ctx context.Context, inp CreateApplicationInput) (*model.Application, error) {
	row := r.db.QueryRow(ctx, `
		INSERT INTO applications (borrower_id, program_id, amount, term_months, payment_type)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING `+appCols,
		inp.BorrowerID, inp.ProgramID, inp.Amount, inp.TermMonths, inp.PaymentType,
	)
	return scanApplication(row)
}

func (r *ApplicationRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Application, error) {
	row := r.db.QueryRow(ctx,
		`SELECT `+appCols+` FROM applications WHERE id = $1`, id)
	return scanApplication(row)
}

// ListByBorrower returns applications for a specific borrower.
func (r *ApplicationRepository) ListByBorrower(ctx context.Context, borrowerID uuid.UUID) ([]*model.Application, error) {
	return r.queryList(ctx,
		`SELECT `+appCols+` FROM applications WHERE borrower_id = $1 ORDER BY created_at DESC`,
		borrowerID)
}

// ListAll returns all applications (for employee/expert/admin), optionally filtered by status.
func (r *ApplicationRepository) ListAll(ctx context.Context, status string) ([]*model.Application, error) {
	if status != "" {
		return r.queryList(ctx,
			`SELECT `+appCols+` FROM applications WHERE status = $1 ORDER BY created_at DESC`,
			status)
	}
	return r.queryList(ctx,
		`SELECT `+appCols+` FROM applications ORDER BY created_at DESC`)
}

// UpdateStatus changes the application status and records history.
// Uses a transaction so both updates are atomic.
func (r *ApplicationRepository) UpdateStatus(ctx context.Context, appID, actorID uuid.UUID, toStatus string, comment *string) (*model.Application, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	// Fetch current status
	var fromStatus string
	if err := tx.QueryRow(ctx,
		`SELECT status FROM applications WHERE id = $1`, appID).Scan(&fromStatus); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("fetch status: %w", err)
	}

	// Update application status
	row := tx.QueryRow(ctx, `
		UPDATE applications SET status = $2, updated_at = NOW()
		WHERE id = $1
		RETURNING `+appCols, appID, toStatus)

	app, err := scanApplication(row)
	if err != nil {
		return nil, err
	}

	// Insert history record (INSERT-only table)
	_, err = tx.Exec(ctx, `
		INSERT INTO application_history
		  (application_id, from_status, to_status, actor_id, comment)
		VALUES ($1, $2, $3, $4, $5)`,
		appID, fromStatus, toStatus, actorID, comment)
	if err != nil {
		return nil, fmt.Errorf("insert history: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}
	return app, nil
}

func (r *ApplicationRepository) GetHistory(ctx context.Context, appID uuid.UUID) ([]*model.ApplicationHistory, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, application_id, from_status, to_status, actor_id, comment, created_at
		FROM application_history
		WHERE application_id = $1
		ORDER BY created_at ASC`, appID)
	if err != nil {
		return nil, fmt.Errorf("get history: %w", err)
	}
	defer rows.Close()

	var history []*model.ApplicationHistory
	for rows.Next() {
		h := &model.ApplicationHistory{}
		if err := rows.Scan(
			&h.ID, &h.ApplicationID, &h.FromStatus, &h.ToStatus,
			&h.ActorID, &h.Comment, &h.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan history: %w", err)
		}
		history = append(history, h)
	}
	return history, nil
}

// ── helpers ──────────────────────────────────────────────────────────────────

func (r *ApplicationRepository) queryList(ctx context.Context, q string, args ...any) ([]*model.Application, error) {
	rows, err := r.db.Query(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("applications query: %w", err)
	}
	defer rows.Close()

	var apps []*model.Application
	for rows.Next() {
		a, err := scanApplication(rows)
		if err != nil {
			return nil, err
		}
		apps = append(apps, a)
	}
	return apps, nil
}

func scanApplication(s scanner) (*model.Application, error) {
	a := &model.Application{}
	err := s.Scan(
		&a.ID, &a.BorrowerID, &a.ProgramID, &a.AssigneeID, &a.Status,
		&a.Amount, &a.TermMonths, &a.PaymentType, &a.CreatedAt, &a.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("scan application: %w", err)
	}
	return a, nil
}
