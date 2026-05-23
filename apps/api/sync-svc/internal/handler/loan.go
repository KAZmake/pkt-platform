package handler

import (
	"net/http"

	"github.com/KAZmake/pkt-platform/apps/api/sync-svc/internal/service"
	"github.com/KAZmake/pkt-platform/apps/api/sync-svc/pkg/response"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type LoanHandler struct {
	svc *service.LoanService
}

func NewLoanHandler(svc *service.LoanService) *LoanHandler {
	return &LoanHandler{svc: svc}
}

// ListLoans returns all loans, optionally filtered by ?status= (2.3.2).
func (h *LoanHandler) ListLoans(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	loans, err := h.svc.GetAll(r.Context(), status)
	if err != nil {
		response.InternalError(w)
		return
	}
	response.OK(w, loans)
}

// GetBorrowerLoans returns loans for a specific borrower (2.3.2).
func (h *LoanHandler) GetBorrowerLoans(w http.ResponseWriter, r *http.Request) {
	borrowerID, err := uuid.Parse(chi.URLParam(r, "borrower_id"))
	if err != nil {
		response.BadRequest(w, "invalid borrower_id")
		return
	}
	loans, err := h.svc.GetByBorrower(r.Context(), borrowerID)
	if err != nil {
		response.InternalError(w)
		return
	}
	response.OK(w, loans)
}

// GetLoan returns a single loan by ID (2.3.2).
func (h *LoanHandler) GetLoan(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "invalid loan id")
		return
	}
	loan, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		response.InternalError(w)
		return
	}
	if loan == nil {
		response.NotFound(w)
		return
	}
	response.OK(w, loan)
}

// GetSchedule returns the payment schedule for a loan (2.3.2).
func (h *LoanHandler) GetSchedule(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "invalid loan id")
		return
	}
	items, err := h.svc.GetSchedule(r.Context(), id)
	if err != nil {
		response.InternalError(w)
		return
	}
	response.OK(w, items)
}

// GetDebts returns debt records for a loan (2.3.2).
func (h *LoanHandler) GetDebts(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "invalid loan id")
		return
	}
	debts, err := h.svc.GetDebts(r.Context(), id)
	if err != nil {
		response.InternalError(w)
		return
	}
	response.OK(w, debts)
}
