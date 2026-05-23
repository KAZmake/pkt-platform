package handler

import (
	"net/http"

	"github.com/KAZmake/pkt-platform/apps/api/core-api/internal/service"
	"github.com/KAZmake/pkt-platform/apps/api/core-api/pkg/response"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type NotificationHandler struct {
	svc *service.NotificationService
}

func NewNotificationHandler(svc *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{svc: svc}
}

// ListNotifications returns the current user's notifications.
// ?unread=true filters to unread only.
func (h *NotificationHandler) ListNotifications(w http.ResponseWriter, r *http.Request) {
	u, ok := userFromCtx(r.Context())
	if !ok {
		response.Unauthorized(w)
		return
	}

	unreadOnly := r.URL.Query().Get("unread") == "true"
	notifs, err := h.svc.ListForUser(r.Context(), u.ID, unreadOnly)
	if err != nil {
		response.InternalError(w)
		return
	}
	response.OK(w, notifs)
}

// UnreadCount returns the count of unread notifications.
func (h *NotificationHandler) UnreadCount(w http.ResponseWriter, r *http.Request) {
	u, ok := userFromCtx(r.Context())
	if !ok {
		response.Unauthorized(w)
		return
	}

	count, err := h.svc.UnreadCount(r.Context(), u.ID)
	if err != nil {
		response.InternalError(w)
		return
	}
	response.OK(w, map[string]int{"count": count})
}

// MarkRead marks a single notification as read.
func (h *NotificationHandler) MarkRead(w http.ResponseWriter, r *http.Request) {
	u, ok := userFromCtx(r.Context())
	if !ok {
		response.Unauthorized(w)
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "invalid notification id")
		return
	}

	if err := h.svc.MarkRead(r.Context(), id, u.ID); err != nil {
		response.InternalError(w)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// MarkAllRead marks all notifications for the current user as read.
func (h *NotificationHandler) MarkAllRead(w http.ResponseWriter, r *http.Request) {
	u, ok := userFromCtx(r.Context())
	if !ok {
		response.Unauthorized(w)
		return
	}

	if err := h.svc.MarkAllRead(r.Context(), u.ID); err != nil {
		response.InternalError(w)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
