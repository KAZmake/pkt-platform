package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/KAZmake/pkt-platform/apps/api/core-api/internal/model"
	"github.com/KAZmake/pkt-platform/apps/api/core-api/internal/repository"
	"github.com/google/uuid"
)

// ── mock ─────────────────────────────────────────────────────────────────────

type mockTicketRepo struct {
	createFn         func(context.Context, repository.CreateTicketInput) (*model.Ticket, error)
	getByIDFn        func(context.Context, uuid.UUID) (*model.Ticket, error)
	getMessagesFn    func(context.Context, uuid.UUID) ([]*model.TicketMessage, error)
	listByBorrowerFn func(context.Context, uuid.UUID) ([]*model.Ticket, error)
	listAllFn        func(context.Context, string) ([]*model.Ticket, error)
	updateStatusFn   func(context.Context, uuid.UUID, string) (*model.Ticket, error)
	addMessageFn     func(context.Context, uuid.UUID, uuid.UUID, string) (*model.TicketMessage, error)
}

func (m *mockTicketRepo) Create(ctx context.Context, inp repository.CreateTicketInput) (*model.Ticket, error) {
	if m.createFn != nil {
		return m.createFn(ctx, inp)
	}
	return nil, nil
}
func (m *mockTicketRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.Ticket, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, nil
}
func (m *mockTicketRepo) GetMessages(ctx context.Context, id uuid.UUID) ([]*model.TicketMessage, error) {
	if m.getMessagesFn != nil {
		return m.getMessagesFn(ctx, id)
	}
	return nil, nil
}
func (m *mockTicketRepo) ListByBorrower(ctx context.Context, id uuid.UUID) ([]*model.Ticket, error) {
	if m.listByBorrowerFn != nil {
		return m.listByBorrowerFn(ctx, id)
	}
	return nil, nil
}
func (m *mockTicketRepo) ListAll(ctx context.Context, status string) ([]*model.Ticket, error) {
	if m.listAllFn != nil {
		return m.listAllFn(ctx, status)
	}
	return nil, nil
}
func (m *mockTicketRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status string) (*model.Ticket, error) {
	if m.updateStatusFn != nil {
		return m.updateStatusFn(ctx, id, status)
	}
	return nil, nil
}
func (m *mockTicketRepo) AddMessage(ctx context.Context, ticketID, authorID uuid.UUID, body string) (*model.TicketMessage, error) {
	if m.addMessageFn != nil {
		return m.addMessageFn(ctx, ticketID, authorID, body)
	}
	return nil, nil
}

func stubTicket() *model.Ticket {
	return &model.Ticket{
		ID:         uuid.New(),
		BorrowerID: uuid.New(),
		Type:       "other",
		Subject:    "Вопрос по займу",
		Status:     "open",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

// ── Create ────────────────────────────────────────────────────────────────────

func TestTicketService_Create_InvalidType(t *testing.T) {
	ctx := context.Background()
	svc := NewTicketService(&mockTicketRepo{})
	_, err := svc.Create(ctx, uuid.New(), CreateTicketInput{Type: "invalid", Subject: "Тема"})
	if err == nil {
		t.Fatal("expected error for invalid type")
	}
}

func TestTicketService_Create_ValidTypes(t *testing.T) {
	ctx := context.Background()
	validTypes := []string{"early_repayment", "restructuring", "prolongation", "other"}
	ticket := stubTicket()

	for _, tp := range validTypes {
		t.Run(tp, func(t *testing.T) {
			mock := &mockTicketRepo{
				createFn: func(_ context.Context, _ repository.CreateTicketInput) (*model.Ticket, error) {
					return ticket, nil
				},
			}
			svc := NewTicketService(mock)
			got, err := svc.Create(ctx, uuid.New(), CreateTicketInput{Type: tp, Subject: "Тема"})
			if err != nil {
				t.Fatalf("unexpected error for type %q: %v", tp, err)
			}
			if got == nil {
				t.Error("expected ticket, got nil")
			}
		})
	}
}

func TestTicketService_Create_MissingSubject(t *testing.T) {
	ctx := context.Background()
	svc := NewTicketService(&mockTicketRepo{})
	_, err := svc.Create(ctx, uuid.New(), CreateTicketInput{Type: "other", Subject: ""})
	if err == nil {
		t.Fatal("expected error for missing subject")
	}
	if err.Error() != "subject is required" {
		t.Errorf("want 'subject is required', got %q", err.Error())
	}
}

func TestTicketService_Create_Success(t *testing.T) {
	ctx := context.Background()
	want := stubTicket()
	mock := &mockTicketRepo{
		createFn: func(_ context.Context, _ repository.CreateTicketInput) (*model.Ticket, error) {
			return want, nil
		},
	}
	svc := NewTicketService(mock)
	got, err := svc.Create(ctx, uuid.New(), CreateTicketInput{Type: "other", Subject: "Тема"})
	if err != nil {
		t.Fatal(err)
	}
	if got == nil || got.ID != want.ID {
		t.Errorf("want ID %v, got %v", want.ID, got)
	}
}

// ── ChangeStatus ──────────────────────────────────────────────────────────────

func TestTicketService_ChangeStatus_InvalidStatus(t *testing.T) {
	ctx := context.Background()
	svc := NewTicketService(&mockTicketRepo{})
	_, err := svc.ChangeStatus(ctx, uuid.New(), "unknown_status")
	if err == nil {
		t.Fatal("expected error for invalid status")
	}
}

func TestTicketService_ChangeStatus_ValidStatuses(t *testing.T) {
	ctx := context.Background()
	validStatuses := []string{"open", "in_progress", "resolved", "closed"}
	ticket := stubTicket()

	for _, status := range validStatuses {
		t.Run(status, func(t *testing.T) {
			mock := &mockTicketRepo{
				updateStatusFn: func(_ context.Context, _ uuid.UUID, _ string) (*model.Ticket, error) {
					return ticket, nil
				},
			}
			svc := NewTicketService(mock)
			got, err := svc.ChangeStatus(ctx, uuid.New(), status)
			if err != nil {
				t.Fatalf("unexpected error for status %q: %v", status, err)
			}
			if got == nil {
				t.Error("expected ticket, got nil")
			}
		})
	}
}

func TestTicketService_ChangeStatus_RepoError(t *testing.T) {
	ctx := context.Background()
	dbErr := errors.New("db error")
	mock := &mockTicketRepo{
		updateStatusFn: func(_ context.Context, _ uuid.UUID, _ string) (*model.Ticket, error) {
			return nil, dbErr
		},
	}
	svc := NewTicketService(mock)
	_, err := svc.ChangeStatus(ctx, uuid.New(), "open")
	if !errors.Is(err, dbErr) {
		t.Errorf("want %v, got %v", dbErr, err)
	}
}

// ── AddMessage ────────────────────────────────────────────────────────────────

func TestTicketService_AddMessage_EmptyBody(t *testing.T) {
	ctx := context.Background()
	svc := NewTicketService(&mockTicketRepo{})
	_, err := svc.AddMessage(ctx, uuid.New(), uuid.New(), AddMessageInput{Body: ""})
	if err == nil {
		t.Fatal("expected error for empty body")
	}
	if err.Error() != "body is required" {
		t.Errorf("want 'body is required', got %q", err.Error())
	}
}

func TestTicketService_AddMessage_TicketNotFound(t *testing.T) {
	ctx := context.Background()
	mock := &mockTicketRepo{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (*model.Ticket, error) { return nil, nil },
	}
	svc := NewTicketService(mock)
	got, err := svc.AddMessage(ctx, uuid.New(), uuid.New(), AddMessageInput{Body: "привет"})
	if err != nil {
		t.Fatal(err)
	}
	if got != nil {
		t.Error("expected nil for not-found ticket")
	}
}

func TestTicketService_AddMessage_Success(t *testing.T) {
	ctx := context.Background()
	ticket := stubTicket()
	msg := &model.TicketMessage{ID: uuid.New(), TicketID: ticket.ID, Body: "привет", CreatedAt: time.Now()}

	mock := &mockTicketRepo{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (*model.Ticket, error) { return ticket, nil },
		addMessageFn: func(_ context.Context, _, _ uuid.UUID, body string) (*model.TicketMessage, error) {
			return msg, nil
		},
	}
	svc := NewTicketService(mock)
	got, err := svc.AddMessage(ctx, ticket.ID, uuid.New(), AddMessageInput{Body: "привет"})
	if err != nil {
		t.Fatal(err)
	}
	if got == nil || got.ID != msg.ID {
		t.Errorf("want message ID %v, got %v", msg.ID, got)
	}
}

// ── GetWithMessages ───────────────────────────────────────────────────────────

func TestTicketService_GetWithMessages_NotFound(t *testing.T) {
	ctx := context.Background()
	mock := &mockTicketRepo{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (*model.Ticket, error) { return nil, nil },
	}
	svc := NewTicketService(mock)
	t_, msgs, err := svc.GetWithMessages(ctx, uuid.New())
	if err != nil || t_ != nil || msgs != nil {
		t.Errorf("want nil,nil,nil for not found; got %v,%v,%v", t_, msgs, err)
	}
}

func TestTicketService_GetWithMessages_Success(t *testing.T) {
	ctx := context.Background()
	ticket := stubTicket()
	msgs := []*model.TicketMessage{{ID: uuid.New(), Body: "test", CreatedAt: time.Now()}}

	mock := &mockTicketRepo{
		getByIDFn:     func(_ context.Context, _ uuid.UUID) (*model.Ticket, error) { return ticket, nil },
		getMessagesFn: func(_ context.Context, _ uuid.UUID) ([]*model.TicketMessage, error) { return msgs, nil },
	}
	svc := NewTicketService(mock)
	gotTicket, gotMsgs, err := svc.GetWithMessages(ctx, ticket.ID)
	if err != nil {
		t.Fatal(err)
	}
	if gotTicket == nil || gotTicket.ID != ticket.ID {
		t.Errorf("want ticket ID %v, got %v", ticket.ID, gotTicket)
	}
	if len(gotMsgs) != 1 {
		t.Errorf("want 1 message, got %d", len(gotMsgs))
	}
}

// ── ListForBorrower / ListAll ─────────────────────────────────────────────────

func TestTicketService_ListForBorrower(t *testing.T) {
	ctx := context.Background()
	borrowerID := uuid.New()
	ticket := stubTicket()
	mock := &mockTicketRepo{
		listByBorrowerFn: func(_ context.Context, id uuid.UUID) ([]*model.Ticket, error) {
			if id != borrowerID {
				return nil, nil
			}
			return []*model.Ticket{ticket}, nil
		},
	}
	svc := NewTicketService(mock)
	tickets, err := svc.ListForBorrower(ctx, borrowerID)
	if err != nil {
		t.Fatal(err)
	}
	if len(tickets) != 1 {
		t.Errorf("want 1 ticket, got %d", len(tickets))
	}
}

func TestTicketService_ListAll_WithStatus(t *testing.T) {
	ctx := context.Background()
	var capturedStatus string
	mock := &mockTicketRepo{
		listAllFn: func(_ context.Context, status string) ([]*model.Ticket, error) {
			capturedStatus = status
			return nil, nil
		},
	}
	svc := NewTicketService(mock)
	svc.ListAll(ctx, "open") //nolint:errcheck
	if capturedStatus != "open" {
		t.Errorf("want status 'open', got %q", capturedStatus)
	}
}
