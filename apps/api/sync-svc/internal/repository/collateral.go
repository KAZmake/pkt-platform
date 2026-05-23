package repository

import (
	"context"
	"fmt"

	"github.com/KAZmake/pkt-platform/apps/api/sync-svc/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

// CollateralRepository reads from expertise-svc's collaterals table (shared DB, read-only).
type CollateralRepository struct {
	db *pgxpool.Pool
}

func NewCollateralRepository(db *pgxpool.Pool) *CollateralRepository {
	return &CollateralRepository{db: db}
}

// ListActive returns all non-released collaterals for monitoring (2.3.3).
func (r *CollateralRepository) ListActive(ctx context.Context) ([]*model.CollateralRef, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, type, description, estimated_value, insurance_expiry, last_inventory_date, is_released
		FROM collaterals
		WHERE is_released = FALSE
		ORDER BY insurance_expiry ASC NULLS LAST`)
	if err != nil {
		return nil, fmt.Errorf("list collaterals: %w", err)
	}
	defer rows.Close()

	var cols []*model.CollateralRef
	for rows.Next() {
		c := &model.CollateralRef{}
		if err := rows.Scan(&c.ID, &c.Type, &c.Description, &c.EstimatedValue,
			&c.InsuranceExpiry, &c.LastInventoryDate, &c.IsReleased); err != nil {
			return nil, fmt.Errorf("scan collateral: %w", err)
		}
		cols = append(cols, c)
	}
	return cols, nil
}
