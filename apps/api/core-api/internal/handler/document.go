package handler

import (
	"encoding/json"
	"net/http"

	"github.com/KAZmake/pkt-platform/apps/api/core-api/internal/service"
	"github.com/KAZmake/pkt-platform/apps/api/core-api/pkg/response"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type DocumentHandler struct {
	svc         *service.DocumentService
	scheduleSvc *service.ScheduleService
}

func NewDocumentHandler(svc *service.DocumentService, scheduleSvc *service.ScheduleService) *DocumentHandler {
	return &DocumentHandler{svc: svc, scheduleSvc: scheduleSvc}
}

// ListApplicationDocuments returns all documents for an application.
func (h *DocumentHandler) ListApplicationDocuments(w http.ResponseWriter, r *http.Request) {
	appID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "invalid application id")
		return
	}
	docs, err := h.svc.ListByOwner(r.Context(), appID, "application")
	if err != nil {
		response.InternalError(w)
		return
	}
	response.OK(w, docs)
}

// InitiateUpload registers a document and returns a presigned PUT URL for direct MinIO upload.
func (h *DocumentHandler) InitiateUpload(w http.ResponseWriter, r *http.Request) {
	appID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "invalid application id")
		return
	}

	var inp struct {
		Name     string  `json:"name"`
		MimeType *string `json:"mime_type"`
	}
	if err := json.NewDecoder(r.Body).Decode(&inp); err != nil || inp.Name == "" {
		response.BadRequest(w, "name is required")
		return
	}

	result, err := h.svc.InitiateUpload(r.Context(), appID, "application", inp.Name, inp.MimeType)
	if err != nil {
		response.InternalError(w)
		return
	}
	response.Created(w, result)
}

// GetDownloadURL returns a presigned GET URL for a document.
func (h *DocumentHandler) GetDownloadURL(w http.ResponseWriter, r *http.Request) {
	docID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "invalid document id")
		return
	}

	result, err := h.svc.GetDownloadURL(r.Context(), docID)
	if err != nil {
		response.InternalError(w)
		return
	}
	if result == nil {
		response.NotFound(w)
		return
	}
	response.OK(w, result)
}

// DeleteDocument removes a document from MinIO and the DB.
func (h *DocumentHandler) DeleteDocument(w http.ResponseWriter, r *http.Request) {
	docID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "invalid document id")
		return
	}
	if err := h.svc.Delete(r.Context(), docID); err != nil {
		response.InternalError(w)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// GetSchedule computes and returns the payment schedule for an application.
func (h *DocumentHandler) GetSchedule(w http.ResponseWriter, r *http.Request) {
	appID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "invalid application id")
		return
	}

	result, err := h.scheduleSvc.Calculate(r.Context(), appID)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}
	if result == nil {
		response.NotFound(w)
		return
	}
	response.OK(w, result)
}
