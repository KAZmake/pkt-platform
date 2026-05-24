package service

import (
	"context"
	"fmt"

	"github.com/KAZmake/pkt-platform/apps/api/core-api/internal/model"
	"github.com/KAZmake/pkt-platform/apps/api/core-api/internal/repository"
)

type programRepo interface {
	List(ctx context.Context, activeOnly bool) ([]*model.LoanProgram, error)
	GetByID(ctx context.Context, id string) (*model.LoanProgram, error)
	Create(ctx context.Context, inp repository.CreateProgramInput) (*model.LoanProgram, error)
	Update(ctx context.Context, id string, inp repository.UpdateProgramInput) (*model.LoanProgram, error)
	SetActive(ctx context.Context, id string, active bool) error
}

type ProgramService struct {
	repo programRepo
}

func NewProgramService(repo programRepo) *ProgramService {
	return &ProgramService{repo: repo}
}

func (s *ProgramService) ListActive(ctx context.Context) ([]*model.LoanProgram, error) {
	return s.repo.List(ctx, true)
}

func (s *ProgramService) ListAll(ctx context.Context) ([]*model.LoanProgram, error) {
	return s.repo.List(ctx, false)
}

func (s *ProgramService) GetByID(ctx context.Context, id string) (*model.LoanProgram, error) {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, nil
	}
	return p, nil
}

func (s *ProgramService) Create(ctx context.Context, inp repository.CreateProgramInput) (*model.LoanProgram, error) {
	if inp.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if inp.Rate <= 0 {
		return nil, fmt.Errorf("rate must be positive")
	}
	if inp.MinAmount <= 0 || inp.MaxAmount <= inp.MinAmount {
		return nil, fmt.Errorf("invalid amount range")
	}
	if inp.MinTermMonths <= 0 || inp.MaxTermMonths < inp.MinTermMonths {
		return nil, fmt.Errorf("invalid term range")
	}
	return s.repo.Create(ctx, inp)
}

func (s *ProgramService) Update(ctx context.Context, id string, inp repository.UpdateProgramInput) (*model.LoanProgram, error) {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, nil
	}
	return s.repo.Update(ctx, id, inp)
}

func (s *ProgramService) Deactivate(ctx context.Context, id string) error {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("not found")
	}
	return s.repo.SetActive(ctx, id, false)
}
