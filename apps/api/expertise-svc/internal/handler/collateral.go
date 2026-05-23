package handler

import (
	"encoding/json"
	"net/http"

	"github.com/KAZmake/pkt-platform/apps/api/expertise-svc/internal/service"
	"github.com/KAZmake/pkt-platform/apps/api/expertise-svc/pkg/response"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type CollateralHandler struct {
	svc *service.CollateralService
}

func NewCollateralHandler(svc *service.CollateralService) *CollateralHandler {
	return &CollateralHandler{svc: svc}
}

// CreateCollateral creates a new collateral card (2.2.4).
func (h *CollateralHandler) CreateCollateral(w http.ResponseWriter, r *http.Request) {
	var inp service.CollateralInput
	if err := json.NewDecoder(r.Body).Decode(&inp); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}
	col, err := h.svc.Create(r.Context(), inp)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}
	response.Created(w, col)
}

// GetCollateral returns a single collateral by ID.
func (h *CollateralHandler) GetCollateral(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "invalid collateral id")
		return
	}
	col, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		response.InternalError(w)
		return
	}
	if col == nil {
		response.NotFound(w)
		return
	}
	response.OK(w, col)
}

// ListApplicationCollaterals lists collaterals linked to an application.
func (h *CollateralHandler) ListApplicationCollaterals(w http.ResponseWriter, r *http.Request) {
	appID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "invalid application id")
		return
	}
	cols, err := h.svc.ListByApplication(r.Context(), appID)
	if err != nil {
		response.InternalError(w)
		return
	}
	response.OK(w, cols)
}

// UpdateCollateral updates a collateral card (2.2.4).
func (h *CollateralHandler) UpdateCollateral(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "invalid collateral id")
		return
	}
	var inp service.CollateralInput
	if err := json.NewDecoder(r.Body).Decode(&inp); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}
	col, err := h.svc.Update(r.Context(), id, inp)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}
	if col == nil {
		response.NotFound(w)
		return
	}
	response.OK(w, col)
}

// AttachCollateral links a collateral to an application.
func (h *CollateralHandler) AttachCollateral(w http.ResponseWriter, r *http.Request) {
	appID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "invalid application id")
		return
	}
	var inp struct {
		CollateralID string `json:"collateral_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&inp); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}
	colID, err := uuid.Parse(inp.CollateralID)
	if err != nil {
		response.BadRequest(w, "invalid collateral_id")
		return
	}
	if err := h.svc.Attach(r.Context(), appID, colID); err != nil {
		response.InternalError(w)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ReleaseCollateral releases a collateral from an application (2.2.4).
func (h *CollateralHandler) ReleaseCollateral(w http.ResponseWriter, r *http.Request) {
	appID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "invalid application id")
		return
	}
	colID, err := uuid.Parse(chi.URLParam(r, "col_id"))
	if err != nil {
		response.BadRequest(w, "invalid collateral id")
		return
	}
	if err := h.svc.Release(r.Context(), appID, colID); err != nil {
		response.InternalError(w)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
