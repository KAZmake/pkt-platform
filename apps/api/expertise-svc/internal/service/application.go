package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/KAZmake/pkt-platform/apps/api/expertise-svc/internal/model"
	"github.com/KAZmake/pkt-platform/apps/api/expertise-svc/internal/repository"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go/jetstream"
)

// FullApplication is the aggregated view of the expert workstation (2.2.3).
type FullApplication struct {
	Application *model.Application          `json:"application"`
	History     []*model.ApplicationHistory `json:"history"`
	Collaterals []*model.Collateral         `json:"collaterals"`
	Conclusions []*model.ExpertConclusion   `json:"conclusions"`
	Votes       []*model.CommitteeVote      `json:"votes"`
}

type ApplicationService struct {
	appRepo  *repository.ApplicationRepository
	colRepo  *repository.CollateralRepository
	concRepo *repository.ConclusionRepository
	js       jetstream.JetStream
}

func NewApplicationService(
	appRepo *repository.ApplicationRepository,
	colRepo *repository.CollateralRepository,
	concRepo *repository.ConclusionRepository,
	js jetstream.JetStream,
) *ApplicationService {
	return &ApplicationService{appRepo: appRepo, colRepo: colRepo, concRepo: concRepo, js: js}
}

// List returns applications, optionally filtered by status and/or assignee (2.2.3).
func (s *ApplicationService) List(ctx context.Context, status, assigneeID string) ([]*model.Application, error) {
	return s.appRepo.List(ctx, status, assigneeID)
}

// GetFull returns the full application card for the expert workstation (2.2.3).
func (s *ApplicationService) GetFull(ctx context.Context, id uuid.UUID) (*FullApplication, error) {
	app, err := s.appRepo.GetByID(ctx, id)
	if err != nil || app == nil {
		return nil, err
	}

	history, err := s.appRepo.GetHistory(ctx, id)
	if err != nil {
		return nil, err
	}

	collaterals, err := s.colRepo.ListByApplication(ctx, id)
	if err != nil {
		return nil, err
	}

	conclusions, err := s.concRepo.ListByApplication(ctx, id)
	if err != nil {
		return nil, err
	}

	votes, err := s.concRepo.GetVotes(ctx, id)
	if err != nil {
		return nil, err
	}

	return &FullApplication{
		Application: app,
		History:     history,
		Collaterals: collaterals,
		Conclusions: conclusions,
		Votes:       votes,
	}, nil
}

type ChangeStatusInput struct {
	ToStatus string  `json:"to_status"`
	Comment  *string `json:"comment"`
}

// ChangeStatus validates FSM + role (2.2.1, 2.2.2) and updates the application.
func (s *ApplicationService) ChangeStatus(ctx context.Context, appID, actorID uuid.UUID, role string, inp ChangeStatusInput) (*model.Application, error) {
	app, err := s.appRepo.GetByID(ctx, appID)
	if err != nil {
		return nil, err
	}
	if app == nil {
		return nil, nil
	}

	// 2.2.2 — role must own current stage
	if !model.CanActOnStage(role, app.Status) {
		return nil, fmt.Errorf("role %q cannot act on stage %q", role, app.Status)
	}

	// FSM graph check
	if !model.IsValidTransition(app.Status, inp.ToStatus) {
		return nil, fmt.Errorf("invalid transition: %s → %s", app.Status, inp.ToStatus)
	}

	updated, err := s.appRepo.UpdateStatus(ctx, appID, actorID, inp.ToStatus, inp.Comment)
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

	// 2.2.7 — if issued, publish issuance event for sync-svc
	if inp.ToStatus == model.StatusIssued {
		s.publish("application.issued", map[string]any{
			"application_id": appID,
			"borrower_id":    app.BorrowerID,
			"program_id":     app.ProgramID,
			"amount":         app.Amount,
			"term_months":    app.TermMonths,
			"payment_type":   app.PaymentType,
			"issued_at":      time.Now(),
			"issued_by":      actorID,
		})
		slog.Info("application issued — sync event published", "application_id", appID)
	}

	return updated, nil
}

// Assign assigns the application to an employee/expert (2.2.3).
func (s *ApplicationService) Assign(ctx context.Context, appID, assigneeID uuid.UUID) error {
	return s.appRepo.Assign(ctx, appID, assigneeID)
}

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
