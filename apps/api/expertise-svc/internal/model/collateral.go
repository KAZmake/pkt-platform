package model

import (
	"time"

	"github.com/google/uuid"
)

// Collateral types
const (
	CollateralLand       = "land"
	CollateralEquipment  = "equipment"
	CollateralLivestock  = "livestock"
	CollateralRealEstate = "real_estate"
	CollateralOther      = "other"
)

type Collateral struct {
	ID                uuid.UUID  `json:"id"`
	Type              string     `json:"type"`
	Description       *string    `json:"description,omitempty"`
	EstimatedValue    *float64   `json:"estimated_value,omitempty"`
	CadastralNumber   *string    `json:"cadastral_number,omitempty"`
	InsuranceExpiry   *time.Time `json:"insurance_expiry,omitempty"`
	LastInventoryDate *time.Time `json:"last_inventory_date,omitempty"`
	IsReleased        bool       `json:"is_released"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

// ApplicationCollateral is the join record linking collaterals to applications.
type ApplicationCollateral struct {
	ApplicationID uuid.UUID  `json:"application_id"`
	CollateralID  uuid.UUID  `json:"collateral_id"`
	AttachedAt    time.Time  `json:"attached_at"`
	ReleasedAt    *time.Time `json:"released_at,omitempty"`
}
