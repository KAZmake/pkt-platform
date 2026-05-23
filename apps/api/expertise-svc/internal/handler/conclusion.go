package handler

import (
	"encoding/json"
	"net/http"

	"github.com/KAZmake/pkt-platform/apps/api/expertise-svc/internal/service"
	"github.com/KAZmake/pkt-platform/apps/api/expertise-svc/pkg/response"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type ConclusionHandler struct {
	svc *service.ConclusionService
}

func NewConclusionHandler(svc *service.ConclusionService) *ConclusionHandler {
	return &ConclusionHandler{svc: svc}
}

// SubmitConclusion submits an expert conclusion for a stage (2.2.3).
func (h *ConclusionHandler) SubmitConclusion(w http.ResponseWriter, r *http.Request) {
	appID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "invalid application id")
		return
	}

	actorID, _, ok := actorFromCtx(r.Context())
	if !ok {
		response.Unauthorized(w)
		return
	}

	var inp service.SubmitConclusionInput
	if err := json.NewDecoder(r.Body).Decode(&inp); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	conclusion, err := h.svc.Submit(r.Context(), appID, actorID, inp)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}
	response.Created(w, conclusion)
}

// ListConclusions returns all expert conclusions for an application.
func (h *ConclusionHandler) ListConclusions(w http.ResponseWriter, r *http.Request) {
	appID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "invalid application id")
		return
	}
	conclusions, err := h.svc.ListByApplication(r.Context(), appID)
	if err != nil {
		response.InternalError(w)
		return
	}
	response.OK(w, conclusions)
}

// AddVote adds a credit committee vote (2.2.5).
func (h *ConclusionHandler) AddVote(w http.ResponseWriter, r *http.Request) {
	appID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "invalid application id")
		return
	}

	actorID, _, ok := actorFromCtx(r.Context())
	if !ok {
		response.Unauthorized(w)
		return
	}

	var inp service.AddVoteInput
	if err := json.NewDecoder(r.Body).Decode(&inp); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	vote, err := h.svc.AddVote(r.Context(), appID, actorID, inp)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}
	response.Created(w, vote)
}

// ListVotes returns all committee votes for an application.
func (h *ConclusionHandler) ListVotes(w http.ResponseWriter, r *http.Request) {
	appID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "invalid application id")
		return
	}
	votes, err := h.svc.GetVotes(r.Context(), appID)
	if err != nil {
		response.InternalError(w)
		return
	}
	response.OK(w, votes)
}
