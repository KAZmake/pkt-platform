package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/KAZmake/pkt-platform/apps/api/core-api/internal/model"
	"github.com/KAZmake/pkt-platform/apps/api/core-api/internal/repository"
	"github.com/google/uuid"
)

// ── mock ─────────────────────────────────────────────────────────────────────

type mockNotifSvc struct {
	createFn      func(context.Context, repository.CreateNotificationInput) (*model.Notification, error)
	listForUserFn func(context.Context, uuid.UUID, bool) ([]*model.Notification, error)
	unreadCountFn func(context.Context, uuid.UUID) (int, error)
	markReadFn    func(context.Context, uuid.UUID, uuid.UUID) error
	markAllReadFn func(context.Context, uuid.UUID) error
}

func (m *mockNotifSvc) Create(ctx context.Context, inp repository.CreateNotificationInput) (*model.Notification, error) {
	if m.createFn != nil {
		return m.createFn(ctx, inp)
	}
	return nil, nil
}
func (m *mockNotifSvc) ListForUser(ctx context.Context, userID uuid.UUID, unreadOnly bool) ([]*model.Notification, error) {
	if m.listForUserFn != nil {
		return m.listForUserFn(ctx, userID, unreadOnly)
	}
	return nil, nil
}
func (m *mockNotifSvc) UnreadCount(ctx context.Context, userID uuid.UUID) (int, error) {
	if m.unreadCountFn != nil {
		return m.unreadCountFn(ctx, userID)
	}
	return 0, nil
}
func (m *mockNotifSvc) MarkRead(ctx context.Context, id, userID uuid.UUID) error {
	if m.markReadFn != nil {
		return m.markReadFn(ctx, id, userID)
	}
	return nil
}
func (m *mockNotifSvc) MarkAllRead(ctx context.Context, userID uuid.UUID) error {
	if m.markAllReadFn != nil {
		return m.markAllReadFn(ctx, userID)
	}
	return nil
}

// ── helpers ───────────────────────────────────────────────────────────────────

func withUser(r *http.Request, u *model.User) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), syncedUserKey, u))
}

func testUser() *model.User {
	return &model.User{
		ID:        uuid.New(),
		Role:      model.RoleBorrower,
		Email:     "test@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func testNotification(userID uuid.UUID) *model.Notification {
	return &model.Notification{
		ID: uuid.New(), UserID: userID, Type: "status",
		Title: "Обновление", Body: "Заявка обновлена", IsRead: false, CreatedAt: time.Now(),
	}
}

// ── ListNotifications ─────────────────────────────────────────────────────────

func TestNotificationHandler_ListNotifications_Success(t *testing.T) {
	u := testUser()
	notifs := []*model.Notification{testNotification(u.ID)}

	h := NewNotificationHandler(&mockNotifSvc{
		listForUserFn: func(_ context.Context, id uuid.UUID, _ bool) ([]*model.Notification, error) {
			if id != u.ID {
				return nil, nil
			}
			return notifs, nil
		},
	})

	r := httptest.NewRequest(http.MethodGet, "/api/v1/notifications", nil)
	r = withUser(r, u)
	w := httptest.NewRecorder()
	h.ListNotifications(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("want 200, got %d", w.Code)
	}
}

func TestNotificationHandler_ListNotifications_NoUser(t *testing.T) {
	h := NewNotificationHandler(&mockNotifSvc{})
	r := httptest.NewRequest(http.MethodGet, "/api/v1/notifications", nil)
	w := httptest.NewRecorder()
	h.ListNotifications(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("want 401, got %d", w.Code)
	}
}

func TestNotificationHandler_ListNotifications_Error(t *testing.T) {
	u := testUser()
	h := NewNotificationHandler(&mockNotifSvc{
		listForUserFn: func(_ context.Context, _ uuid.UUID, _ bool) ([]*model.Notification, error) {
			return nil, errors.New("db error")
		},
	})

	r := httptest.NewRequest(http.MethodGet, "/api/v1/notifications", nil)
	r = withUser(r, u)
	w := httptest.NewRecorder()
	h.ListNotifications(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("want 500, got %d", w.Code)
	}
}

// ── UnreadCount ───────────────────────────────────────────────────────────────

func TestNotificationHandler_UnreadCount_Success(t *testing.T) {
	u := testUser()
	h := NewNotificationHandler(&mockNotifSvc{
		unreadCountFn: func(_ context.Context, _ uuid.UUID) (int, error) { return 3, nil },
	})

	r := httptest.NewRequest(http.MethodGet, "/api/v1/notifications/unread-count", nil)
	r = withUser(r, u)
	w := httptest.NewRecorder()
	h.UnreadCount(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("want 200, got %d", w.Code)
	}
}

func TestNotificationHandler_UnreadCount_NoUser(t *testing.T) {
	h := NewNotificationHandler(&mockNotifSvc{})
	r := httptest.NewRequest(http.MethodGet, "/api/v1/notifications/unread-count", nil)
	w := httptest.NewRecorder()
	h.UnreadCount(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("want 401, got %d", w.Code)
	}
}

func TestNotificationHandler_UnreadCount_Error(t *testing.T) {
	u := testUser()
	h := NewNotificationHandler(&mockNotifSvc{
		unreadCountFn: func(_ context.Context, _ uuid.UUID) (int, error) { return 0, errors.New("db error") },
	})

	r := httptest.NewRequest(http.MethodGet, "/api/v1/notifications/unread-count", nil)
	r = withUser(r, u)
	w := httptest.NewRecorder()
	h.UnreadCount(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("want 500, got %d", w.Code)
	}
}

// ── MarkRead ──────────────────────────────────────────────────────────────────

func TestNotificationHandler_MarkRead_Success(t *testing.T) {
	u := testUser()
	notifID := uuid.New()
	h := NewNotificationHandler(&mockNotifSvc{
		markReadFn: func(_ context.Context, _, _ uuid.UUID) error { return nil },
	})

	r := httptest.NewRequest(http.MethodPatch, "/api/v1/notifications/"+notifID.String()+"/read", nil)
	r = withUser(r, u)
	r = withParam(r, "id", notifID.String())
	w := httptest.NewRecorder()
	h.MarkRead(w, r)

	if w.Code != http.StatusNoContent {
		t.Errorf("want 204, got %d", w.Code)
	}
}

func TestNotificationHandler_MarkRead_NoUser(t *testing.T) {
	h := NewNotificationHandler(&mockNotifSvc{})
	r := httptest.NewRequest(http.MethodPatch, "/api/v1/notifications/"+uuid.New().String()+"/read", nil)
	r = withParam(r, "id", uuid.New().String())
	w := httptest.NewRecorder()
	h.MarkRead(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("want 401, got %d", w.Code)
	}
}

func TestNotificationHandler_MarkRead_InvalidID(t *testing.T) {
	u := testUser()
	h := NewNotificationHandler(&mockNotifSvc{})

	r := httptest.NewRequest(http.MethodPatch, "/api/v1/notifications/bad-id/read", nil)
	r = withUser(r, u)
	r = withParam(r, "id", "not-a-uuid")
	w := httptest.NewRecorder()
	h.MarkRead(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("want 400, got %d", w.Code)
	}
}

// ── MarkAllRead ───────────────────────────────────────────────────────────────

func TestNotificationHandler_MarkAllRead_Success(t *testing.T) {
	u := testUser()
	h := NewNotificationHandler(&mockNotifSvc{
		markAllReadFn: func(_ context.Context, _ uuid.UUID) error { return nil },
	})

	r := httptest.NewRequest(http.MethodPatch, "/api/v1/notifications/read-all", nil)
	r = withUser(r, u)
	w := httptest.NewRecorder()
	h.MarkAllRead(w, r)

	if w.Code != http.StatusNoContent {
		t.Errorf("want 204, got %d", w.Code)
	}
}

func TestNotificationHandler_MarkAllRead_NoUser(t *testing.T) {
	h := NewNotificationHandler(&mockNotifSvc{})
	r := httptest.NewRequest(http.MethodPatch, "/api/v1/notifications/read-all", nil)
	w := httptest.NewRecorder()
	h.MarkAllRead(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("want 401, got %d", w.Code)
	}
}
