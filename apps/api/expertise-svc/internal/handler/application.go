package handler

import (
	"context"
	"encoding/json"
	"net/http"

	apimw "github.com/KAZmake/pkt-platform/apps/api/expertise-svc/internal/middleware"
	"github.com/KAZmake/pkt-platform/apps/api/expertise-svc/internal/model"
	"github.com/KAZmake/pkt-platform/apps/api/expertise-svc/internal/service"
	"github.com/KAZmake/pkt-platform/apps/api/expertise-svc/pkg/response"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type ctxKey string

const syncedUserIDKey ctxKey = "actor_uuid"
const syncedRoleKey ctxKey = "actor_role"

// storeActor middleware extracts actor uuid + highest role from JWT and stores in ctx.
func storeActor(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := apimw.ClaimsFromCtx(r.Context())
		if !ok {
			response.Unauthorized(w)
			return
		}
		actorID, err := uuid.Parse(claims.Subject)
		if err != nil {
			response.Unauthorized(w)
			return
		}
		role := highestRole(claims.RealmAccess.Roles)
		ctx := context.WithValue(r.Context(), syncedUserIDKey, actorID)
		ctx = context.WithValue(ctx, syncedRoleKey, role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func actorFromCtx(ctx context.Context) (uuid.UUID, string, bool) {
	id, ok1 := ctx.Value(syncedUserIDKey).(uuid.UUID)
	role, ok2 := ctx.Value(syncedRoleKey).(string)
	return id, role, ok1 && ok2
}

func highestRole(roles []string) string {
	priority := map[string]int{model.RoleAdmin: 4, model.RoleExpert: 3, model.RoleEmployee: 2, model.RoleBorrower: 1}
	best := "public"
	for _, r := range roles {
		if priority[r] > priority[best] {
			best = r
		}
	}
	return best
}

// ── ApplicationHandler ────────────────────────────────────────────────────────

type ApplicationHandler struct {
	svc *service.ApplicationService
}

func NewApplicationHandler(svc *service.ApplicationService) *ApplicationHandler {
	return &ApplicationHandler{svc: svc}
}

// StoreActor exposes the middleware for use in router setup.
var StoreActor = storeActor

// ListApplications returns filtered list for the expert workstation (2.2.3).
func (h *ApplicationHandler) ListApplications(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	assigneeID := r.URL.Query().Get("assignee_id")

	apps, err := h.svc.List(r.Context(), status, assigneeID)
	if err != nil {
		response.InternalError(w)
		return
	}
	response.OK(w, apps)
}

// GetApplication returns the full application card (2.2.3).
func (h *ApplicationHandler) GetApplication(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "invalid application id")
		return
	}
	full, err := h.svc.GetFull(r.Context(), id)
	if err != nil {
		response.InternalError(w)
		return
	}
	if full == nil {
		response.NotFound(w)
		return
	}
	response.OK(w, full)
}

// ChangeStatus advances the FSM with role check (2.2.1, 2.2.2).
func (h *ApplicationHandler) ChangeStatus(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "invalid application id")
		return
	}

	actorID, role, ok := actorFromCtx(r.Context())
	if !ok {
		response.Unauthorized(w)
		return
	}

	var inp service.ChangeStatusInput
	if err := json.NewDecoder(r.Body).Decode(&inp); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	app, err := h.svc.ChangeStatus(r.Context(), id, actorID, role, inp)
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

// Assign sets the assignee for an application.
func (h *ApplicationHandler) Assign(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "invalid application id")
		return
	}

	var inp struct {
		AssigneeID string `json:"assignee_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&inp); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	assigneeID, err := uuid.Parse(inp.AssigneeID)
	if err != nil {
		response.BadRequest(w, "invalid assignee_id")
		return
	}

	if err := h.svc.Assign(r.Context(), id, assigneeID); err != nil {
		response.InternalError(w)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
