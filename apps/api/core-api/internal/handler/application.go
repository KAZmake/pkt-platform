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

type ApplicationHandler struct {
	svc *service.ApplicationService
}

func NewApplicationHandler(svc *service.ApplicationService) *ApplicationHandler {
	return &ApplicationHandler{svc: svc}
}

// ListApplications returns applications scoped to the caller's role:
//   - borrower → only their own applications
//   - employee/expert/admin → all applications (filterable by ?status=)
func (h *ApplicationHandler) ListApplications(w http.ResponseWriter, r *http.Request) {
	u, ok := userFromCtx(r.Context())
	if !ok {
		response.Unauthorized(w)
		return
	}

	if u.Role == model.RoleBorrower {
		b, err := borrowerIDFromCtx(r.Context())
		if err != nil {
			response.Forbidden(w)
			return
		}
		apps, err := h.svc.ListForBorrower(r.Context(), b)
		if err != nil {
			response.InternalError(w)
			return
		}
		response.OK(w, apps)
		return
	}

	status := r.URL.Query().Get("status")
	apps, err := h.svc.ListAll(r.Context(), status)
	if err != nil {
		response.InternalError(w)
		return
	}
	response.OK(w, apps)
}

// GetApplication returns a single application with its FSM history.
func (h *ApplicationHandler) GetApplication(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "invalid application id")
		return
	}

	app, history, err := h.svc.GetWithHistory(r.Context(), id)
	if err != nil {
		response.InternalError(w)
		return
	}
	if app == nil {
		response.NotFound(w)
		return
	}
	response.OK(w, map[string]any{
		"application": app,
		"history":     history,
	})
}

// CreateApplication is called by a borrower to submit a new credit application.
func (h *ApplicationHandler) CreateApplication(w http.ResponseWriter, r *http.Request) {
	u, ok := userFromCtx(r.Context())
	if !ok {
		response.Unauthorized(w)
		return
	}

	borrowerID, err := borrowerIDFromCtx(r.Context())
	if err != nil {
		response.Forbidden(w)
		return
	}

	var inp service.CreateApplicationInput
	if err := json.NewDecoder(r.Body).Decode(&inp); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}
	_ = u

	app, err := h.svc.Create(r.Context(), borrowerID, inp)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}
	response.Created(w, app)
}

// ChangeStatus advances the FSM for an application.
func (h *ApplicationHandler) ChangeStatus(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "invalid application id")
		return
	}

	u, ok := userFromCtx(r.Context())
	if !ok {
		response.Unauthorized(w)
		return
	}

	var inp service.ChangeStatusInput
	if err := json.NewDecoder(r.Body).Decode(&inp); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	app, err := h.svc.ChangeStatus(r.Context(), id, u.ID, inp)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}
	if app == nil {
		response.NotFound(w)
		return
	}
	response.OK(w, app)
}
