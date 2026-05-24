package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/KAZmake/pkt-platform/apps/api/core-api/internal/model"
	"github.com/KAZmake/pkt-platform/apps/api/core-api/internal/service"
	"github.com/google/uuid"
)

// ── mock ─────────────────────────────────────────────────────────────────────

type mockTicketSvc struct {
	createFn          func(context.Context, uuid.UUID, service.CreateTicketInput) (*model.Ticket, error)
	getWithMessagesFn func(context.Context, uuid.UUID) (*model.Ticket, []*model.TicketMessage, error)
	listForBorrowerFn func(context.Context, uuid.UUID) ([]*model.Ticket, error)
	listAllFn         func(context.Context, string) ([]*model.Ticket, error)
	changeStatusFn    func(context.Context, uuid.UUID, string) (*model.Ticket, error)
	addMessageFn      func(context.Context, uuid.UUID, uuid.UUID, service.AddMessageInput) (*model.TicketMessage, error)
}

func (m *mockTicketSvc) Create(ctx context.Context, borrowerID uuid.UUID, inp service.CreateTicketInput) (*model.Ticket, error) {
	if m.createFn != nil {
		return m.createFn(ctx, borrowerID, inp)
	}
	return nil, nil
}
func (m *mockTicketSvc) GetWithMessages(ctx context.Context, id uuid.UUID) (*model.Ticket, []*model.TicketMessage, error) {
	if m.getWithMessagesFn != nil {
		return m.getWithMessagesFn(ctx, id)
	}
	return nil, nil, nil
}
func (m *mockTicketSvc) ListForBorrower(ctx context.Context, id uuid.UUID) ([]*model.Ticket, error) {
	if m.listForBorrowerFn != nil {
		return m.listForBorrowerFn(ctx, id)
	}
	return nil, nil
}
func (m *mockTicketSvc) ListAll(ctx context.Context, status string) ([]*model.Ticket, error) {
	if m.listAllFn != nil {
		return m.listAllFn(ctx, status)
	}
	return nil, nil
}
func (m *mockTicketSvc) ChangeStatus(ctx context.Context, id uuid.UUID, status string) (*model.Ticket, error) {
	if m.changeStatusFn != nil {
		return m.changeStatusFn(ctx, id, status)
	}
	return nil, nil
}
func (m *mockTicketSvc) AddMessage(ctx context.Context, ticketID, authorID uuid.UUID, inp service.AddMessageInput) (*model.TicketMessage, error) {
	if m.addMessageFn != nil {
		return m.addMessageFn(ctx, ticketID, authorID, inp)
	}
	return nil, nil
}

// ── helpers ───────────────────────────────────────────────────────────────────

func withBorrower(r *http.Request, borrower *model.Borrower) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), syncedBorrowerKey, borrower))
}

func testTicket() *model.Ticket {
	return &model.Ticket{
		ID: uuid.New(), BorrowerID: uuid.New(), Type: "other",
		Subject: "Тестовый вопрос", Status: "open",
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
}

// ── ListTickets ───────────────────────────────────────────────────────────────

func TestTicketHandler_ListTickets_Employee(t *testing.T) {
	u := testUser()
	u.Role = model.RoleEmployee
	ticket := testTicket()

	h := NewTicketHandler(&mockTicketSvc{
		listAllFn: func(_ context.Context, _ string) ([]*model.Ticket, error) {
			return []*model.Ticket{ticket}, nil
		},
	})

	r := httptest.NewRequest(http.MethodGet, "/api/v1/tickets", nil)
	r = withUser(r, u)
	w := httptest.NewRecorder()
	h.ListTickets(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("want 200, got %d", w.Code)
	}
}

func TestTicketHandler_ListTickets_NoUser(t *testing.T) {
	h := NewTicketHandler(&mockTicketSvc{})
	r := httptest.NewRequest(http.MethodGet, "/api/v1/tickets", nil)
	w := httptest.NewRecorder()
	h.ListTickets(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("want 401, got %d", w.Code)
	}
}

func TestTicketHandler_ListTickets_BorrowerNoBorrowerCtx(t *testing.T) {
	u := testUser()
	u.Role = model.RoleBorrower // borrower but no borrower in context

	h := NewTicketHandler(&mockTicketSvc{})
	r := httptest.NewRequest(http.MethodGet, "/api/v1/tickets", nil)
	r = withUser(r, u)
	w := httptest.NewRecorder()
	h.ListTickets(w, r)

	if w.Code != http.StatusForbidden {
		t.Errorf("want 403 (no borrower in ctx), got %d", w.Code)
	}
}

func TestTicketHandler_ListTickets_BorrowerSuccess(t *testing.T) {
	u := testUser()
	u.Role = model.RoleBorrower
	borrower := &model.Borrower{ID: uuid.New(), UserID: u.ID}
	ticket := testTicket()

	h := NewTicketHandler(&mockTicketSvc{
		listForBorrowerFn: func(_ context.Context, _ uuid.UUID) ([]*model.Ticket, error) {
			return []*model.Ticket{ticket}, nil
		},
	})

	r := httptest.NewRequest(http.MethodGet, "/api/v1/tickets", nil)
	r = withUser(r, u)
	r = withBorrower(r, borrower)
	w := httptest.NewRecorder()
	h.ListTickets(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("want 200, got %d", w.Code)
	}
}

// ── GetTicket ─────────────────────────────────────────────────────────────────

func TestTicketHandler_GetTicket_Found(t *testing.T) {
	ticket := testTicket()
	h := NewTicketHandler(&mockTicketSvc{
		getWithMessagesFn: func(_ context.Context, id uuid.UUID) (*model.Ticket, []*model.TicketMessage, error) {
			if id == ticket.ID {
				return ticket, nil, nil
			}
			return nil, nil, nil
		},
	})

	r := httptest.NewRequest(http.MethodGet, "/api/v1/tickets/"+ticket.ID.String(), nil)
	r = withParam(r, "id", ticket.ID.String())
	w := httptest.NewRecorder()
	h.GetTicket(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("want 200, got %d", w.Code)
	}
}

func TestTicketHandler_GetTicket_NotFound(t *testing.T) {
	h := NewTicketHandler(&mockTicketSvc{
		getWithMessagesFn: func(_ context.Context, _ uuid.UUID) (*model.Ticket, []*model.TicketMessage, error) {
			return nil, nil, nil
		},
	})

	r := httptest.NewRequest(http.MethodGet, "/api/v1/tickets/"+uuid.New().String(), nil)
	r = withParam(r, "id", uuid.New().String())
	w := httptest.NewRecorder()
	h.GetTicket(w, r)

	if w.Code != http.StatusNotFound {
		t.Errorf("want 404, got %d", w.Code)
	}
}

func TestTicketHandler_GetTicket_InvalidID(t *testing.T) {
	h := NewTicketHandler(&mockTicketSvc{})
	r := httptest.NewRequest(http.MethodGet, "/api/v1/tickets/bad-id", nil)
	r = withParam(r, "id", "not-a-uuid")
	w := httptest.NewRecorder()
	h.GetTicket(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("want 400, got %d", w.Code)
	}
}

func TestTicketHandler_GetTicket_Error(t *testing.T) {
	h := NewTicketHandler(&mockTicketSvc{
		getWithMessagesFn: func(_ context.Context, _ uuid.UUID) (*model.Ticket, []*model.TicketMessage, error) {
			return nil, nil, errors.New("db error")
		},
	})

	r := httptest.NewRequest(http.MethodGet, "/api/v1/tickets/"+uuid.New().String(), nil)
	r = withParam(r, "id", uuid.New().String())
	w := httptest.NewRecorder()
	h.GetTicket(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("want 500, got %d", w.Code)
	}
}

// ── CreateTicket ──────────────────────────────────────────────────────────────

func TestTicketHandler_CreateTicket_Success(t *testing.T) {
	ticket := testTicket()
	borrower := &model.Borrower{ID: uuid.New(), UserID: uuid.New()}

	h := NewTicketHandler(&mockTicketSvc{
		createFn: func(_ context.Context, _ uuid.UUID, _ service.CreateTicketInput) (*model.Ticket, error) {
			return ticket, nil
		},
	})

	body, _ := json.Marshal(service.CreateTicketInput{Type: "other", Subject: "Тема"})
	r := httptest.NewRequest(http.MethodPost, "/api/v1/tickets", bytes.NewReader(body))
	r = withBorrower(r, borrower)
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.CreateTicket(w, r)

	if w.Code != http.StatusCreated {
		t.Errorf("want 201, got %d", w.Code)
	}
}

func TestTicketHandler_CreateTicket_NoBorrower(t *testing.T) {
	h := NewTicketHandler(&mockTicketSvc{})
	r := httptest.NewRequest(http.MethodPost, "/api/v1/tickets", nil)
	w := httptest.NewRecorder()
	h.CreateTicket(w, r)

	if w.Code != http.StatusForbidden {
		t.Errorf("want 403, got %d", w.Code)
	}
}

func TestTicketHandler_CreateTicket_BadJSON(t *testing.T) {
	borrower := &model.Borrower{ID: uuid.New()}
	h := NewTicketHandler(&mockTicketSvc{})

	r := httptest.NewRequest(http.MethodPost, "/api/v1/tickets", bytes.NewReader([]byte(`bad`)))
	r = withBorrower(r, borrower)
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.CreateTicket(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("want 400, got %d", w.Code)
	}
}

func TestTicketHandler_CreateTicket_ServiceError(t *testing.T) {
	borrower := &model.Borrower{ID: uuid.New()}
	h := NewTicketHandler(&mockTicketSvc{
		createFn: func(_ context.Context, _ uuid.UUID, _ service.CreateTicketInput) (*model.Ticket, error) {
			return nil, errors.New("invalid type")
		},
	})

	body, _ := json.Marshal(service.CreateTicketInput{Type: "bad", Subject: "X"})
	r := httptest.NewRequest(http.MethodPost, "/api/v1/tickets", bytes.NewReader(body))
	r = withBorrower(r, borrower)
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.CreateTicket(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("want 400, got %d", w.Code)
	}
}

// ── ChangeTicketStatus ────────────────────────────────────────────────────────

func TestTicketHandler_ChangeTicketStatus_Success(t *testing.T) {
	ticket := testTicket()
	ticket.Status = "in_progress"

	h := NewTicketHandler(&mockTicketSvc{
		changeStatusFn: func(_ context.Context, _ uuid.UUID, _ string) (*model.Ticket, error) {
			return ticket, nil
		},
	})

	body, _ := json.Marshal(map[string]string{"status": "in_progress"})
	r := httptest.NewRequest(http.MethodPatch, "/api/v1/tickets/"+ticket.ID.String()+"/status", bytes.NewReader(body))
	r = withParam(r, "id", ticket.ID.String())
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.ChangeTicketStatus(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("want 200, got %d", w.Code)
	}
}

func TestTicketHandler_ChangeTicketStatus_NotFound(t *testing.T) {
	h := NewTicketHandler(&mockTicketSvc{
		changeStatusFn: func(_ context.Context, _ uuid.UUID, _ string) (*model.Ticket, error) {
			return nil, nil
		},
	})

	body, _ := json.Marshal(map[string]string{"status": "in_progress"})
	r := httptest.NewRequest(http.MethodPatch, "/api/v1/tickets/"+uuid.New().String()+"/status", bytes.NewReader(body))
	r = withParam(r, "id", uuid.New().String())
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.ChangeTicketStatus(w, r)

	if w.Code != http.StatusNotFound {
		t.Errorf("want 404, got %d", w.Code)
	}
}

func TestTicketHandler_ChangeTicketStatus_ServiceError(t *testing.T) {
	h := NewTicketHandler(&mockTicketSvc{
		changeStatusFn: func(_ context.Context, _ uuid.UUID, _ string) (*model.Ticket, error) {
			return nil, errors.New("invalid status")
		},
	})

	body, _ := json.Marshal(map[string]string{"status": "bad"})
	r := httptest.NewRequest(http.MethodPatch, "/api/v1/tickets/"+uuid.New().String()+"/status", bytes.NewReader(body))
	r = withParam(r, "id", uuid.New().String())
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.ChangeTicketStatus(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("want 400, got %d", w.Code)
	}
}
