package service

import (
	"context"
	"log/slog"

	"github.com/KAZmake/pkt-platform/apps/api/sync-svc/internal/cache"
	"github.com/KAZmake/pkt-platform/apps/api/sync-svc/internal/model"
	"github.com/KAZmake/pkt-platform/apps/api/sync-svc/internal/repository"
	"github.com/google/uuid"
)

// LoanService handles read requests for loans with Redis-first caching (2.3.2).
type LoanService struct {
	repo  *repository.LoanRepository
	cache *cache.Client
}

func NewLoanService(repo *repository.LoanRepository, cache *cache.Client) *LoanService {
	return &LoanService{repo: repo, cache: cache}
}

func (s *LoanService) GetByBorrower(ctx context.Context, borrowerID uuid.UUID) ([]*model.Loan, error) {
	var loans []*model.Loan
	hit, err := s.cache.Get(ctx, cache.LoanKey(borrowerID.String()), &loans)
	if err != nil {
		slog.Warn("cache get error", "key", cache.LoanKey(borrowerID.String()), "error", err)
	}
	if hit {
		return loans, nil
	}
	loans, err = s.repo.GetByBorrower(ctx, borrowerID)
	if err != nil {
		return nil, err
	}
	if len(loans) > 0 {
		_ = s.cache.Set(ctx, cache.LoanKey(borrowerID.String()), loans, cache.DefaultTTL)
	}
	return loans, nil
}

func (s *LoanService) GetAll(ctx context.Context, status string) ([]*model.Loan, error) {
	if status == "" {
		var loans []*model.Loan
		if hit, _ := s.cache.Get(ctx, cache.AllLoansKey(), &loans); hit {
			return loans, nil
		}
	}
	return s.repo.GetAll(ctx, status)
}

func (s *LoanService) GetByID(ctx context.Context, id uuid.UUID) (*model.Loan, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *LoanService) GetSchedule(ctx context.Context, loanID uuid.UUID) ([]*model.ScheduleItem, error) {
	var items []*model.ScheduleItem
	hit, err := s.cache.Get(ctx, cache.ScheduleKey(loanID.String()), &items)
	if err != nil {
		slog.Warn("cache get error", "key", cache.ScheduleKey(loanID.String()), "error", err)
	}
	if hit {
		return items, nil
	}
	items, err = s.repo.GetSchedule(ctx, loanID)
	if err != nil {
		return nil, err
	}
	if len(items) > 0 {
		_ = s.cache.Set(ctx, cache.ScheduleKey(loanID.String()), items, cache.DefaultTTL)
	}
	return items, nil
}

func (s *LoanService) GetDebts(ctx context.Context, loanID uuid.UUID) ([]*model.LoanDebt, error) {
	var debts []*model.LoanDebt
	hit, err := s.cache.Get(ctx, cache.DebtsKey(loanID.String()), &debts)
	if err != nil {
		slog.Warn("cache get error", "key", cache.DebtsKey(loanID.String()), "error", err)
	}
	if hit {
		return debts, nil
	}
	debts, err = s.repo.GetDebts(ctx, loanID)
	if err != nil {
		return nil, err
	}
	_ = s.cache.Set(ctx, cache.DebtsKey(loanID.String()), debts, cache.DefaultTTL)
	return debts, nil
}
