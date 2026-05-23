package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/KAZmake/pkt-platform/apps/api/expertise-svc/internal/model"
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

func (r *ApplicationRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Application, error) {
	row := r.db.QueryRow(ctx, `SELECT `+appCols+` FROM applications WHERE id = $1`, id)
	return scanApplication(row)
}

func (r *ApplicationRepository) List(ctx context.Context, status, assigneeID string) ([]*model.Application, error) {
	q := `SELECT ` + appCols + ` FROM applications WHERE 1=1`
	args := []any{}
	i := 1
	if status != "" {
		q += fmt.Sprintf(` AND status = $%d`, i)
		args = append(args, status)
		i++
	}
	if assigneeID != "" {
		q += fmt.Sprintf(` AND assignee_id = $%d`, i)
		args = append(args, assigneeID)
		i++
	}
	q += ` ORDER BY updated_at DESC`

	rows, err := r.db.Query(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("list applications: %w", err)
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

// UpdateStatus atomically updates status and writes history (2.2.8 audit log).
func (r *ApplicationRepository) UpdateStatus(ctx context.Context, appID, actorID uuid.UUID, toStatus string, comment *string) (*model.Application, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	var fromStatus string
	if err := tx.QueryRow(ctx,
		`SELECT status FROM applications WHERE id = $1`, appID).Scan(&fromStatus); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("fetch status: %w", err)
	}

	row := tx.QueryRow(ctx, `
		UPDATE applications SET status = $2, updated_at = NOW()
		WHERE id = $1
		RETURNING `+appCols, appID, toStatus)

	app, err := scanApplication(row)
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO application_history (application_id, from_status, to_status, actor_id, comment)
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

// Assign sets the assignee_id for an application.
func (r *ApplicationRepository) Assign(ctx context.Context, appID, assigneeID uuid.UUID) error {
	_, err := r.db.Exec(ctx,
		`UPDATE applications SET assignee_id = $2, updated_at = NOW() WHERE id = $1`,
		appID, assigneeID)
	return err
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
