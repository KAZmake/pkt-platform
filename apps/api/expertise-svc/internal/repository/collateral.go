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

type CollateralRepository struct {
	db *pgxpool.Pool
}

func NewCollateralRepository(db *pgxpool.Pool) *CollateralRepository {
	return &CollateralRepository{db: db}
}

const colCols = `id, type, description, estimated_value, cadastral_number,
	insurance_expiry, last_inventory_date, is_released, created_at, updated_at`

type CreateCollateralInput struct {
	Type              string
	Description       *string
	EstimatedValue    *float64
	CadastralNumber   *string
	InsuranceExpiry   *string // DATE as string
	LastInventoryDate *string // DATE as string
}

func (r *CollateralRepository) Create(ctx context.Context, inp CreateCollateralInput) (*model.Collateral, error) {
	row := r.db.QueryRow(ctx, `
		INSERT INTO collaterals (type, description, estimated_value, cadastral_number,
		                         insurance_expiry, last_inventory_date)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING `+colCols,
		inp.Type, inp.Description, inp.EstimatedValue, inp.CadastralNumber,
		inp.InsuranceExpiry, inp.LastInventoryDate,
	)
	return scanCollateral(row)
}

func (r *CollateralRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Collateral, error) {
	row := r.db.QueryRow(ctx, `SELECT `+colCols+` FROM collaterals WHERE id = $1`, id)
	return scanCollateral(row)
}

func (r *CollateralRepository) ListByApplication(ctx context.Context, appID uuid.UUID) ([]*model.Collateral, error) {
	rows, err := r.db.Query(ctx, `
		SELECT c.`+colCols+`
		FROM collaterals c
		JOIN application_collaterals ac ON ac.collateral_id = c.id
		WHERE ac.application_id = $1 AND ac.released_at IS NULL
		ORDER BY ac.attached_at DESC`, appID)
	if err != nil {
		return nil, fmt.Errorf("list collaterals: %w", err)
	}
	defer rows.Close()

	var cols []*model.Collateral
	for rows.Next() {
		c, err := scanCollateral(rows)
		if err != nil {
			return nil, err
		}
		cols = append(cols, c)
	}
	return cols, nil
}

// Update updates mutable fields of a collateral (2.2.4).
func (r *CollateralRepository) Update(ctx context.Context, id uuid.UUID, inp CreateCollateralInput) (*model.Collateral, error) {
	row := r.db.QueryRow(ctx, `
		UPDATE collaterals SET
			type                = COALESCE($2, type),
			description         = COALESCE($3, description),
			estimated_value     = COALESCE($4, estimated_value),
			cadastral_number    = COALESCE($5, cadastral_number),
			insurance_expiry    = COALESCE($6::date, insurance_expiry),
			last_inventory_date = COALESCE($7::date, last_inventory_date),
			updated_at          = NOW()
		WHERE id = $1
		RETURNING `+colCols,
		id, inp.Type, inp.Description, inp.EstimatedValue, inp.CadastralNumber,
		inp.InsuranceExpiry, inp.LastInventoryDate,
	)
	return scanCollateral(row)
}

// Attach links a collateral to an application.
func (r *CollateralRepository) Attach(ctx context.Context, appID, colID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO application_collaterals (application_id, collateral_id)
		VALUES ($1, $2)
		ON CONFLICT (application_id, collateral_id) DO NOTHING`,
		appID, colID)
	return err
}

// Release marks a collateral as released from an application and sets is_released.
func (r *CollateralRepository) Release(ctx context.Context, appID, colID uuid.UUID) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	_, err = tx.Exec(ctx, `
		UPDATE application_collaterals SET released_at = NOW()
		WHERE application_id = $1 AND collateral_id = $2 AND released_at IS NULL`,
		appID, colID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
		UPDATE collaterals SET is_released = TRUE, updated_at = NOW()
		WHERE id = $1`, colID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func scanCollateral(s scanner) (*model.Collateral, error) {
	c := &model.Collateral{}
	err := s.Scan(
		&c.ID, &c.Type, &c.Description, &c.EstimatedValue, &c.CadastralNumber,
		&c.InsuranceExpiry, &c.LastInventoryDate, &c.IsReleased, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("scan collateral: %w", err)
	}
	return c, nil
}
