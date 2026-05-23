package handler

import (
	"encoding/json"
	"net/http"

	"github.com/KAZmake/pkt-platform/apps/api/core-api/internal/model"
	"github.com/KAZmake/pkt-platform/apps/api/core-api/internal/service"
	"github.com/KAZmake/pkt-platform/apps/api/core-api/pkg/response"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type TicketHandler struct {
	svc *service.TicketService
}

func NewTicketHandler(svc *service.TicketService) *TicketHandler {
	return &TicketHandler{svc: svc}
}

// ListTickets returns tickets scoped to role:
//   - borrower → only their own tickets
//   - employee/expert/admin → all (?status= filter)
func (h *TicketHandler) ListTickets(w http.ResponseWriter, r *http.Request) {
	u, ok := userFromCtx(r.Context())
	if !ok {
		response.Unauthorized(w)
		return
	}

	if u.Role == model.RoleBorrower {
		borrowerID, err := borrowerIDFromCtx(r.Context())
		if err != nil {
			response.Forbidden(w)
			return
		}
		tickets, err := h.svc.ListForBorrower(r.Context(), borrowerID)
		if err != nil {
			response.InternalError(w)
			return
		}
		response.OK(w, tickets)
		return
	}

	status := r.URL.Query().Get("status")
	tickets, err := h.svc.ListAll(r.Context(), status)
	if err != nil {
		response.InternalError(w)
		return
	}
	response.OK(w, tickets)
}

// CreateTicket creates a new ticket for the authenticated borrower.
func (h *TicketHandler) CreateTicket(w http.ResponseWriter, r *http.Request) {
	borrowerID, err := borrowerIDFromCtx(r.Context())
	if err != nil {
		response.Forbidden(w)
		return
	}

	var inp service.CreateTicketInput
	if err := json.NewDecoder(r.Body).Decode(&inp); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	ticket, err := h.svc.Create(r.Context(), borrowerID, inp)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}
	response.Created(w, ticket)
}

// GetTicket returns a ticket and its messages.
func (h *TicketHandler) GetTicket(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "invalid ticket id")
		return
	}

	ticket, messages, err := h.svc.GetWithMessages(r.Context(), id)
	if err != nil {
		response.InternalError(w)
		return
	}
	if ticket == nil {
		response.NotFound(w)
		return
	}
	response.OK(w, map[string]any{
		"ticket":   ticket,
		"messages": messages,
	})
}

// AddMessage appends a message to a ticket.
func (h *TicketHandler) AddMessage(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "invalid ticket id")
		return
	}

	u, ok := userFromCtx(r.Context())
	if !ok {
		response.Unauthorized(w)
		return
	}

	var inp service.AddMessageInput
	if err := json.NewDecoder(r.Body).Decode(&inp); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	msg, err := h.svc.AddMessage(r.Context(), id, u.ID, inp)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}
	if msg == nil {
		response.NotFound(w)
		return
	}
	response.Created(w, msg)
}

// ChangeTicketStatus updates the ticket status (employee/expert/admin only).
func (h *TicketHandler) ChangeTicketStatus(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "invalid ticket id")
		return
	}

	var inp struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&inp); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	ticket, err := h.svc.ChangeStatus(r.Context(), id, inp.Status)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}
	if ticket == nil {
		response.NotFound(w)
		return
	}
	response.OK(w, ticket)
}
