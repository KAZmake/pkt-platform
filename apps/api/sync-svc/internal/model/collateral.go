package model

import (
	"time"

	"github.com/google/uuid"
)

// CollateralAlert is a monitoring alert produced by the collateral monitor (2.3.3).
type CollateralAlert struct {
	CollateralID uuid.UUID `json:"collateral_id"`
	Type         string    `json:"type"`     // insurance_expiring | inventory_overdue
	Severity     string    `json:"severity"` // warning | critical
	Message      string    `json:"message"`
	DaysUntil    int       `json:"days_until,omitempty"` // for insurance_expiring
	DaysSince    int       `json:"days_since,omitempty"` // for inventory_overdue
}

// CollateralRef is the read-only view of expertise-svc's collaterals table.
type CollateralRef struct {
	ID                uuid.UUID  `json:"id"`
	Type              string     `json:"type"`
	Description       *string    `json:"description,omitempty"`
	EstimatedValue    *float64   `json:"estimated_value,omitempty"`
	InsuranceExpiry   *time.Time `json:"insurance_expiry,omitempty"`
	LastInventoryDate *time.Time `json:"last_inventory_date,omitempty"`
	IsReleased        bool       `json:"is_released"`
}
