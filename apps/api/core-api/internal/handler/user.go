package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	apimw "github.com/KAZmake/pkt-platform/apps/api/core-api/internal/middleware"
	"github.com/KAZmake/pkt-platform/apps/api/core-api/internal/model"
	"github.com/KAZmake/pkt-platform/apps/api/core-api/internal/service"
	"github.com/KAZmake/pkt-platform/apps/api/core-api/pkg/response"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// ── Context helpers ──────────────────────────────────────────────────────────

type ctxKey string

const (
	syncedUserKey     ctxKey = "synced_user"
	syncedBorrowerKey ctxKey = "synced_borrower"
)

func userFromCtx(ctx context.Context) (*model.User, bool) {
	u, ok := ctx.Value(syncedUserKey).(*model.User)
	return u, ok && u != nil
}

func borrowerIDFromCtx(ctx context.Context) (uuid.UUID, error) {
	b, ok := ctx.Value(syncedBorrowerKey).(*model.Borrower)
	if !ok || b == nil {
		return uuid.Nil, fmt.Errorf("no borrower in context")
	}
	return b.ID, nil
}

// ── SyncUser middleware ───────────────────────────────────────────────────────

// SyncUser upserts the Keycloak user into our DB on every authenticated request
// and stores the result in context. Must run after Authenticate.
func SyncUser(svc *service.UserService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := apimw.ClaimsFromCtx(r.Context())
			if !ok {
				response.Unauthorized(w)
				return
			}

			u, err := svc.SyncFromClaims(r.Context(), claims)
			if err != nil || u == nil {
				response.InternalError(w)
				return
			}

			ctx := context.WithValue(r.Context(), syncedUserKey, u)

			// For borrowers, eagerly load borrower profile into context.
			if u.Role == model.RoleBorrower {
				b, err := svc.GetBorrowerProfile(r.Context(), u.ID)
				if err == nil && b != nil {
					ctx = context.WithValue(ctx, syncedBorrowerKey, b)
				}
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// ── UserHandler ───────────────────────────────────────────────────────────────

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

// GetMe returns the current authenticated user's profile.
func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	u, ok := userFromCtx(r.Context())
	if !ok {
		response.InternalError(w)
		return
	}
	response.OK(w, u)
}

// UpdateMe updates first_name, last_name, phone of the current user.
func (h *UserHandler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	u, ok := userFromCtx(r.Context())
	if !ok {
		response.InternalError(w)
		return
	}

	var inp service.UpdateProfileInput
	if err := json.NewDecoder(r.Body).Decode(&inp); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	updated, err := h.svc.UpdateProfile(r.Context(), u.ID, inp)
	if err != nil || updated == nil {
		response.InternalError(w)
		return
	}
	response.OK(w, updated)
}

// GetMyBorrower returns the borrower profile for the current user.
func (h *UserHandler) GetMyBorrower(w http.ResponseWriter, r *http.Request) {
	u, ok := userFromCtx(r.Context())
	if !ok {
		response.InternalError(w)
		return
	}

	b, err := h.svc.GetBorrowerProfile(r.Context(), u.ID)
	if err != nil {
		response.InternalError(w)
		return
	}
	if b == nil {
		response.NotFound(w)
		return
	}
	response.OK(w, b)
}

// UpdateMyBorrower updates org_name and activity_type of the current borrower.
func (h *UserHandler) UpdateMyBorrower(w http.ResponseWriter, r *http.Request) {
	u, ok := userFromCtx(r.Context())
	if !ok {
		response.InternalError(w)
		return
	}

	var inp service.UpdateBorrowerInput
	if err := json.NewDecoder(r.Body).Decode(&inp); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	b, err := h.svc.UpdateBorrowerProfile(r.Context(), u.ID, inp)
	if err != nil || b == nil {
		response.InternalError(w)
		return
	}
	response.OK(w, b)
}

// ListUsers returns a paginated list of users (admin only).
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	users, total, err := h.svc.ListUsers(r.Context(), limit, offset)
	if err != nil {
		response.InternalError(w)
		return
	}
	response.OK(w, map[string]any{
		"items":  users,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

// GetUser returns a single user by ID (employee/expert/admin only).
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "invalid user id")
		return
	}

	u, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		response.InternalError(w)
		return
	}
	if u == nil {
		response.NotFound(w)
		return
	}
	response.OK(w, u)
}
