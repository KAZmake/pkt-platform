package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/KAZmake/pkt-platform/apps/api/core-api/internal/model"
	"github.com/KAZmake/pkt-platform/apps/api/core-api/internal/repository"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go/jetstream"
)

type ApplicationService struct {
	repo *repository.ApplicationRepository
	js   jetstream.JetStream // nil if NATS unavailable
}

func NewApplicationService(repo *repository.ApplicationRepository, js jetstream.JetStream) *ApplicationService {
	return &ApplicationService{repo: repo, js: js}
}

type CreateApplicationInput struct {
	ProgramID   uuid.UUID `json:"program_id"`
	Amount      float64   `json:"amount"`
	TermMonths  int       `json:"term_months"`
	PaymentType string    `json:"payment_type"`
}

func (s *ApplicationService) Create(ctx context.Context, borrowerID uuid.UUID, inp CreateApplicationInput) (*model.Application, error) {
	if inp.Amount <= 0 {
		return nil, fmt.Errorf("amount must be positive")
	}
	if inp.TermMonths <= 0 {
		return nil, fmt.Errorf("term_months must be positive")
	}
	if inp.PaymentType != "annuity" && inp.PaymentType != "differentiated" {
		return nil, fmt.Errorf("payment_type must be annuity or differentiated")
	}

	app, err := s.repo.Create(ctx, repository.CreateApplicationInput{
		BorrowerID:  borrowerID,
		ProgramID:   inp.ProgramID,
		Amount:      inp.Amount,
		TermMonths:  inp.TermMonths,
		PaymentType: inp.PaymentType,
	})
	if err != nil {
		return nil, err
	}

	s.publish("application.created", map[string]any{
		"application_id": app.ID,
		"borrower_id":    app.BorrowerID,
		"program_id":     app.ProgramID,
		"amount":         app.Amount,
		"created_at":     app.CreatedAt,
	})

	return app, nil
}

func (s *ApplicationService) GetByID(ctx context.Context, id uuid.UUID) (*model.Application, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ApplicationService) GetWithHistory(ctx context.Context, id uuid.UUID) (*model.Application, []*model.ApplicationHistory, error) {
	app, err := s.repo.GetByID(ctx, id)
	if err != nil || app == nil {
		return nil, nil, err
	}
	history, err := s.repo.GetHistory(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	return app, history, nil
}

func (s *ApplicationService) ListForBorrower(ctx context.Context, borrowerID uuid.UUID) ([]*model.Application, error) {
	return s.repo.ListByBorrower(ctx, borrowerID)
}

func (s *ApplicationService) ListAll(ctx context.Context, status string) ([]*model.Application, error) {
	return s.repo.ListAll(ctx, status)
}

type ChangeStatusInput struct {
	ToStatus string  `json:"to_status"`
	Comment  *string `json:"comment"`
}

func (s *ApplicationService) ChangeStatus(ctx context.Context, appID, actorID uuid.UUID, inp ChangeStatusInput) (*model.Application, error) {
	app, err := s.repo.GetByID(ctx, appID)
	if err != nil {
		return nil, err
	}
	if app == nil {
		return nil, fmt.Errorf("application not found")
	}

	if !model.IsValidTransition(app.Status, inp.ToStatus) {
		return nil, fmt.Errorf("invalid transition: %s → %s", app.Status, inp.ToStatus)
	}

	updated, err := s.repo.UpdateStatus(ctx, appID, actorID, inp.ToStatus, inp.Comment)
	if err != nil {
		return nil, err
	}

	s.publish("application.status_changed", map[string]any{
		"application_id": appID,
		"from_status":    app.Status,
		"to_status":      inp.ToStatus,
		"actor_id":       actorID,
		"changed_at":     time.Now(),
	})

	return updated, nil
}

// publish sends a JSON event to NATS JetStream. Non-fatal if NATS is down.
func (s *ApplicationService) publish(subject string, payload any) {
	if s.js == nil {
		return
	}
	data, err := json.Marshal(payload)
	if err != nil {
		slog.Warn("nats: marshal error", "subject", subject, "error", err)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if _, err := s.js.Publish(ctx, subject, data); err != nil {
		slog.Warn("nats: publish error", "subject", subject, "error", err)
	}
}
