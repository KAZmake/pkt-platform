package service

import (
	"context"
	"fmt"

	"github.com/KAZmake/pkt-platform/apps/api/expertise-svc/internal/model"
	"github.com/KAZmake/pkt-platform/apps/api/expertise-svc/internal/repository"
	"github.com/google/uuid"
)

var validCollateralTypes = map[string]bool{
	model.CollateralLand:       true,
	model.CollateralEquipment:  true,
	model.CollateralLivestock:  true,
	model.CollateralRealEstate: true,
	model.CollateralOther:      true,
}

type CollateralService struct {
	repo *repository.CollateralRepository
}

func NewCollateralService(repo *repository.CollateralRepository) *CollateralService {
	return &CollateralService{repo: repo}
}

type CollateralInput struct {
	Type              string   `json:"type"`
	Description       *string  `json:"description"`
	EstimatedValue    *float64 `json:"estimated_value"`
	CadastralNumber   *string  `json:"cadastral_number"`
	InsuranceExpiry   *string  `json:"insurance_expiry"`    // YYYY-MM-DD
	LastInventoryDate *string  `json:"last_inventory_date"` // YYYY-MM-DD
}

func (s *CollateralService) Create(ctx context.Context, inp CollateralInput) (*model.Collateral, error) {
	if !validCollateralTypes[inp.Type] {
		return nil, fmt.Errorf("invalid type: must be land, equipment, livestock, real_estate, or other")
	}
	return s.repo.Create(ctx, repository.CreateCollateralInput{
		Type:              inp.Type,
		Description:       inp.Description,
		EstimatedValue:    inp.EstimatedValue,
		CadastralNumber:   inp.CadastralNumber,
		InsuranceExpiry:   inp.InsuranceExpiry,
		LastInventoryDate: inp.LastInventoryDate,
	})
}

func (s *CollateralService) GetByID(ctx context.Context, id uuid.UUID) (*model.Collateral, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *CollateralService) ListByApplication(ctx context.Context, appID uuid.UUID) ([]*model.Collateral, error) {
	return s.repo.ListByApplication(ctx, appID)
}

func (s *CollateralService) Update(ctx context.Context, id uuid.UUID, inp CollateralInput) (*model.Collateral, error) {
	if inp.Type != "" && !validCollateralTypes[inp.Type] {
		return nil, fmt.Errorf("invalid type")
	}
	return s.repo.Update(ctx, id, repository.CreateCollateralInput{
		Type:              inp.Type,
		Description:       inp.Description,
		EstimatedValue:    inp.EstimatedValue,
		CadastralNumber:   inp.CadastralNumber,
		InsuranceExpiry:   inp.InsuranceExpiry,
		LastInventoryDate: inp.LastInventoryDate,
	})
}

func (s *CollateralService) Attach(ctx context.Context, appID, colID uuid.UUID) error {
	return s.repo.Attach(ctx, appID, colID)
}

func (s *CollateralService) Release(ctx context.Context, appID, colID uuid.UUID) error {
	return s.repo.Release(ctx, appID, colID)
}
