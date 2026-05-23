package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/KAZmake/pkt-platform/apps/api/core-api/pkg/response"
	"github.com/jackc/pgx/v5/pgxpool"
)

type HealthHandler struct {
	db *pgxpool.Pool
}

func NewHealthHandler(db *pgxpool.Pool) *HealthHandler {
	return &HealthHandler{db: db}
}

func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	dbStatus := "ok"
	if err := h.db.Ping(ctx); err != nil {
		dbStatus = "unavailable"
	}

	status := "ok"
	httpStatus := http.StatusOK
	if dbStatus != "ok" {
		status = "degraded"
		httpStatus = http.StatusServiceUnavailable
	}

	response.JSON(w, httpStatus, response.Envelope{
		"status":  status,
		"service": "core-api",
		"checks": map[string]string{
			"database": dbStatus,
		},
	})
}
