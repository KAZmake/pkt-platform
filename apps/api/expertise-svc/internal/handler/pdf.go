package handler

import (
	"net/http"
	"strconv"

	"github.com/KAZmake/pkt-platform/apps/api/expertise-svc/internal/pdf"
	"github.com/KAZmake/pkt-platform/apps/api/expertise-svc/internal/service"
	"github.com/KAZmake/pkt-platform/apps/api/expertise-svc/pkg/response"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type PDFHandler struct {
	appSvc  *service.ApplicationService
	concSvc *service.ConclusionService
}

func NewPDFHandler(appSvc *service.ApplicationService, concSvc *service.ConclusionService) *PDFHandler {
	return &PDFHandler{appSvc: appSvc, concSvc: concSvc}
}

// CommitteeProtocol generates and returns the CC protocol PDF (2.2.5).
func (h *PDFHandler) CommitteeProtocol(w http.ResponseWriter, r *http.Request) {
	appID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "invalid application id")
		return
	}

	full, err := h.appSvc.GetFull(r.Context(), appID)
	if err != nil || full == nil {
		if full == nil {
			response.NotFound(w)
		} else {
			response.InternalError(w)
		}
		return
	}

	data, err := pdf.CommitteeProtocol(full.Application, full.Votes)
	if err != nil {
		response.InternalError(w)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=\"committee_protocol_"+appID.String()+".pdf\"")
	_, _ = w.Write(data)
}

// LoanAgreement generates and returns the loan agreement PDF (2.2.6).
func (h *PDFHandler) LoanAgreement(w http.ResponseWriter, r *http.Request) {
	appID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "invalid application id")
		return
	}

	full, err := h.appSvc.GetFull(r.Context(), appID)
	if err != nil || full == nil {
		if full == nil {
			response.NotFound(w)
		} else {
			response.InternalError(w)
		}
		return
	}

	// Rate from query param (provided by caller from loan program); default 7%
	rate := 7.0
	if v := r.URL.Query().Get("rate"); v != "" {
		if parsed, err := strconv.ParseFloat(v, 64); err == nil {
			rate = parsed
		}
	}

	data, err := pdf.LoanAgreement(full.Application, rate)
	if err != nil {
		response.InternalError(w)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=\"loan_agreement_"+appID.String()+".pdf\"")
	_, _ = w.Write(data)
}
