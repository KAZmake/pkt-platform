package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/KAZmake/pkt-platform/apps/api/expertise-svc/pkg/response"
	"github.com/jackc/pgx/v5/pgxpool"
)

type HealthHandler struct{ db *pgxpool.Pool }

func NewHealthHandler(db *pgxpool.Pool) *HealthHandler { return &HealthHandler{db: db} }

func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	dbStatus := "ok"
	if err := h.db.Ping(ctx); err != nil {
		dbStatus = "error"
	}
	response.OK(w, map[string]any{
		"status":  "ok",
		"service": "expertise-svc",
		"checks":  map[string]string{"database": dbStatus},
	})
}
