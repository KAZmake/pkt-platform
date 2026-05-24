package service

import (
	"context"
	"fmt"

	"github.com/KAZmake/pkt-platform/apps/api/core-api/internal/model"
	"github.com/KAZmake/pkt-platform/apps/api/core-api/internal/repository"
	"github.com/google/uuid"
)

var validTicketTypes = map[string]bool{
	"early_repayment": true,
	"restructuring":   true,
	"prolongation":    true,
	"other":           true,
}

var validTicketStatuses = map[string]bool{
	"open":        true,
	"in_progress": true,
	"resolved":    true,
	"closed":      true,
}

type ticketRepo interface {
	Create(ctx context.Context, inp repository.CreateTicketInput) (*model.Ticket, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Ticket, error)
	GetMessages(ctx context.Context, ticketID uuid.UUID) ([]*model.TicketMessage, error)
	ListByBorrower(ctx context.Context, borrowerID uuid.UUID) ([]*model.Ticket, error)
	ListAll(ctx context.Context, status string) ([]*model.Ticket, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) (*model.Ticket, error)
	AddMessage(ctx context.Context, ticketID, authorID uuid.UUID, body string) (*model.TicketMessage, error)
}

type TicketService struct {
	repo ticketRepo
}

func NewTicketService(repo ticketRepo) *TicketService {
	return &TicketService{repo: repo}
}

type CreateTicketInput struct {
	Type    string `json:"type"`
	Subject string `json:"subject"`
}

type AddMessageInput struct {
	Body string `json:"body"`
}

func (s *TicketService) Create(ctx context.Context, borrowerID uuid.UUID, inp CreateTicketInput) (*model.Ticket, error) {
	if !validTicketTypes[inp.Type] {
		return nil, fmt.Errorf("invalid type: must be early_repayment, restructuring, prolongation, or other")
	}
	if inp.Subject == "" {
		return nil, fmt.Errorf("subject is required")
	}
	return s.repo.Create(ctx, repository.CreateTicketInput{
		BorrowerID: borrowerID,
		Type:       inp.Type,
		Subject:    inp.Subject,
	})
}

func (s *TicketService) GetWithMessages(ctx context.Context, id uuid.UUID) (*model.Ticket, []*model.TicketMessage, error) {
	t, err := s.repo.GetByID(ctx, id)
	if err != nil || t == nil {
		return nil, nil, err
	}
	msgs, err := s.repo.GetMessages(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	return t, msgs, nil
}

func (s *TicketService) ListForBorrower(ctx context.Context, borrowerID uuid.UUID) ([]*model.Ticket, error) {
	return s.repo.ListByBorrower(ctx, borrowerID)
}

func (s *TicketService) ListAll(ctx context.Context, status string) ([]*model.Ticket, error) {
	return s.repo.ListAll(ctx, status)
}

func (s *TicketService) ChangeStatus(ctx context.Context, id uuid.UUID, status string) (*model.Ticket, error) {
	if !validTicketStatuses[status] {
		return nil, fmt.Errorf("invalid status: must be open, in_progress, resolved, or closed")
	}
	t, err := s.repo.UpdateStatus(ctx, id, status)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (s *TicketService) AddMessage(ctx context.Context, ticketID, authorID uuid.UUID, inp AddMessageInput) (*model.TicketMessage, error) {
	if inp.Body == "" {
		return nil, fmt.Errorf("body is required")
	}
	// Verify ticket exists
	t, err := s.repo.GetByID(ctx, ticketID)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, nil
	}
	return s.repo.AddMessage(ctx, ticketID, authorID, inp.Body)
}
