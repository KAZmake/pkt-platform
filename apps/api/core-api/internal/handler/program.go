package handler

import (
	"encoding/json"
	"net/http"

	"github.com/KAZmake/pkt-platform/apps/api/core-api/internal/repository"
	"github.com/KAZmake/pkt-platform/apps/api/core-api/internal/service"
	"github.com/KAZmake/pkt-platform/apps/api/core-api/pkg/response"
	"github.com/go-chi/chi/v5"
)

type ProgramHandler struct {
	svc *service.ProgramService
}

func NewProgramHandler(svc *service.ProgramService) *ProgramHandler {
	return &ProgramHandler{svc: svc}
}

// ListPrograms — public. Returns active programs only.
// Employee/admin get all programs including inactive (?all=true).
func (h *ProgramHandler) ListPrograms(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var (
		programs interface{}
		err      error
	)
	if r.URL.Query().Get("all") == "true" {
		programs, err = h.svc.ListAll(ctx)
	} else {
		programs, err = h.svc.ListActive(ctx)
	}
	if err != nil {
		response.InternalError(w)
		return
	}
	response.OK(w, programs)
}

// GetProgram — public. Returns a single program by ID.
func (h *ProgramHandler) GetProgram(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	p, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		response.InternalError(w)
		return
	}
	if p == nil {
		response.NotFound(w)
		return
	}
	response.OK(w, p)
}

// CreateProgram — admin only.
func (h *ProgramHandler) CreateProgram(w http.ResponseWriter, r *http.Request) {
	var inp repository.CreateProgramInput
	if err := json.NewDecoder(r.Body).Decode(&inp); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	p, err := h.svc.Create(r.Context(), inp)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}
	response.Created(w, p)
}

// UpdateProgram — admin only.
func (h *ProgramHandler) UpdateProgram(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var inp repository.UpdateProgramInput
	if err := json.NewDecoder(r.Body).Decode(&inp); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	p, err := h.svc.Update(r.Context(), id, inp)
	if err != nil {
		response.InternalError(w)
		return
	}
	if p == nil {
		response.NotFound(w)
		return
	}
	response.OK(w, p)
}

// DeactivateProgram — admin only. Soft-deletes by setting is_active=false.
func (h *ProgramHandler) DeactivateProgram(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.svc.Deactivate(r.Context(), id); err != nil {
		if err.Error() == "not found" {
			response.NotFound(w)
			return
		}
		response.InternalError(w)
		return
	}
	response.OK(w, map[string]string{"message": "program deactivated"})
}
