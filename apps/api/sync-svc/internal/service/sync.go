// Package service contains the sync orchestration logic (2.3.1).
package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/KAZmake/pkt-platform/apps/api/sync-svc/internal/cache"
	"github.com/KAZmake/pkt-platform/apps/api/sync-svc/internal/model"
	"github.com/KAZmake/pkt-platform/apps/api/sync-svc/internal/onec"
	"github.com/KAZmake/pkt-platform/apps/api/sync-svc/internal/repository"
	"github.com/google/uuid"
)

// SyncService pulls data from 1С, writes to PostgreSQL, and warms Redis cache.
type SyncService struct {
	client   onec.Client
	loanRepo *repository.LoanRepository
	cache    *cache.Client
}

func NewSyncService(client onec.Client, loanRepo *repository.LoanRepository, cache *cache.Client) *SyncService {
	return &SyncService{client: client, loanRepo: loanRepo, cache: cache}
}

// Run performs a full sync cycle: 1С → DB → Redis (2.3.1).
// Called by the cron scheduler every 20 minutes.
func (s *SyncService) Run(ctx context.Context) {
	start := time.Now()
	slog.Info("sync: starting 1С sync cycle")

	loans, err := s.client.GetLoans()
	if err != nil {
		slog.Error("sync: failed to get loans from 1С", "error", err)
		return
	}
	slog.Info("sync: received loans from 1С", "count", len(loans))

	syncedByBorrower := map[string][]*model.Loan{}

	for _, raw := range loans {
		saved, err := s.upsertLoan(ctx, raw)
		if err != nil {
			slog.Warn("sync: upsert loan failed", "one_c_id", raw.OneCID, "error", err)
			continue
		}

		// Sync schedule
		if err := s.syncSchedule(ctx, raw.OneCID, saved.ID); err != nil {
			slog.Warn("sync: schedule sync failed", "loan", raw.OneCID, "error", err)
		}

		// Sync debts
		if err := s.syncDebts(ctx, raw.OneCID, saved.ID); err != nil {
			slog.Warn("sync: debts sync failed", "loan", raw.OneCID, "error", err)
		}

		key := saved.BorrowerID.String()
		syncedByBorrower[key] = append(syncedByBorrower[key], saved)
	}

	// Warm Redis cache per borrower
	for borrowerID, borrowerLoans := range syncedByBorrower {
		if err := s.cache.Set(ctx, cache.LoanKey(borrowerID), borrowerLoans, cache.DefaultTTL); err != nil {
			slog.Warn("sync: cache warm failed", "borrower_id", borrowerID, "error", err)
		}
	}
	// Warm all-loans cache
	if err := s.cache.Set(ctx, cache.AllLoansKey(), loans, cache.DefaultTTL); err != nil {
		slog.Warn("sync: all-loans cache warm failed", "error", err)
	}

	slog.Info("sync: cycle complete", "duration", time.Since(start).Round(time.Millisecond), "loans", len(loans))
}

func (s *SyncService) upsertLoan(ctx context.Context, raw onec.Loan) (*model.Loan, error) {
	l := model.Loan{
		OneCID:     raw.OneCID,
		Amount:     raw.Amount,
		Rate:       raw.Rate,
		TermMonths: raw.TermMonths,
		IssuedAt:   raw.IssuedAt,
		ExpiresAt:  raw.ExpiresAt,
		Status:     raw.Status,
	}

	if id, err := uuid.Parse(raw.BorrowerID); err == nil {
		l.BorrowerID = id
	}
	if raw.ProgramID != "" {
		if id, err := uuid.Parse(raw.ProgramID); err == nil {
			l.ProgramID = &id
		}
	}

	return s.loanRepo.UpsertLoan(ctx, l)
}

func (s *SyncService) syncSchedule(ctx context.Context, oneCID string, loanID uuid.UUID) error {
	rawItems, err := s.client.GetSchedule(oneCID)
	if err != nil {
		return err
	}

	items := make([]model.ScheduleItem, 0, len(rawItems))
	for _, ri := range rawItems {
		items = append(items, model.ScheduleItem{
			LoanID:    loanID,
			DueDate:   ri.DueDate,
			Principal: ri.Principal,
			Interest:  ri.Interest,
			Total:     ri.Total,
			IsPaid:    ri.IsPaid,
			PaidAt:    ri.PaidAt,
		})
	}
	if err := s.loanRepo.ReplaceSchedule(ctx, loanID, items); err != nil {
		return err
	}
	// Invalidate schedule cache
	_ = s.cache.Del(ctx, cache.ScheduleKey(loanID.String()))
	return nil
}

func (s *SyncService) syncDebts(ctx context.Context, oneCID string, loanID uuid.UUID) error {
	rawDebts, err := s.client.GetDebts(oneCID)
	if err != nil {
		return err
	}

	debts := make([]model.LoanDebt, 0, len(rawDebts))
	for _, rd := range rawDebts {
		debts = append(debts, model.LoanDebt{
			LoanID:      loanID,
			Type:        rd.Type,
			Amount:      rd.Amount,
			DaysOverdue: rd.DaysOverdue,
		})
	}
	if err := s.loanRepo.ReplaceDebts(ctx, loanID, debts); err != nil {
		return err
	}
	// Invalidate debts cache
	_ = s.cache.Del(ctx, cache.DebtsKey(loanID.String()))
	return nil
}
