package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/KAZmake/pkt-platform/apps/api/core-api/internal/model"
	"github.com/KAZmake/pkt-platform/apps/api/core-api/internal/repository"
	"github.com/google/uuid"
)

// ── mock ─────────────────────────────────────────────────────────────────────

type mockNotificationRepo struct {
	createFn      func(context.Context, repository.CreateNotificationInput) (*model.Notification, error)
	listByUserFn  func(context.Context, uuid.UUID, bool) ([]*model.Notification, error)
	unreadCountFn func(context.Context, uuid.UUID) (int, error)
	markReadFn    func(context.Context, uuid.UUID, uuid.UUID) error
	markAllReadFn func(context.Context, uuid.UUID) error
}

func (m *mockNotificationRepo) Create(ctx context.Context, inp repository.CreateNotificationInput) (*model.Notification, error) {
	if m.createFn != nil {
		return m.createFn(ctx, inp)
	}
	return nil, nil
}
func (m *mockNotificationRepo) ListByUser(ctx context.Context, userID uuid.UUID, unreadOnly bool) ([]*model.Notification, error) {
	if m.listByUserFn != nil {
		return m.listByUserFn(ctx, userID, unreadOnly)
	}
	return nil, nil
}
func (m *mockNotificationRepo) UnreadCount(ctx context.Context, userID uuid.UUID) (int, error) {
	if m.unreadCountFn != nil {
		return m.unreadCountFn(ctx, userID)
	}
	return 0, nil
}
func (m *mockNotificationRepo) MarkRead(ctx context.Context, id, userID uuid.UUID) error {
	if m.markReadFn != nil {
		return m.markReadFn(ctx, id, userID)
	}
	return nil
}
func (m *mockNotificationRepo) MarkAllRead(ctx context.Context, userID uuid.UUID) error {
	if m.markAllReadFn != nil {
		return m.markAllReadFn(ctx, userID)
	}
	return nil
}

func stubNotification(userID uuid.UUID) *model.Notification {
	return &model.Notification{
		ID: uuid.New(), UserID: userID,
		Type: "status", Title: "Заявка обновлена", Body: "Ваша заявка перешла на следующий этап",
		IsRead: false, CreatedAt: time.Now(),
	}
}

// ── Create ────────────────────────────────────────────────────────────────────

func TestNotificationService_Create(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	want := stubNotification(userID)

	svc := NewNotificationService(&mockNotificationRepo{
		createFn: func(_ context.Context, inp repository.CreateNotificationInput) (*model.Notification, error) {
			if inp.UserID != userID {
				return nil, errors.New("wrong user ID")
			}
			return want, nil
		},
	})
	got, err := svc.Create(ctx, repository.CreateNotificationInput{
		UserID: userID, Type: "status", Title: "Test", Body: "Body",
	})
	if err != nil {
		t.Fatal(err)
	}
	if got == nil || got.ID != want.ID {
		t.Errorf("want ID %v, got %v", want.ID, got)
	}
}

func TestNotificationService_Create_Error(t *testing.T) {
	ctx := context.Background()
	dbErr := errors.New("db error")
	svc := NewNotificationService(&mockNotificationRepo{
		createFn: func(_ context.Context, _ repository.CreateNotificationInput) (*model.Notification, error) {
			return nil, dbErr
		},
	})
	_, err := svc.Create(ctx, repository.CreateNotificationInput{UserID: uuid.New()})
	if !errors.Is(err, dbErr) {
		t.Errorf("want %v, got %v", dbErr, err)
	}
}

// ── ListForUser ───────────────────────────────────────────────────────────────

func TestNotificationService_ListForUser(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	notifs := []*model.Notification{stubNotification(userID), stubNotification(userID)}

	svc := NewNotificationService(&mockNotificationRepo{
		listByUserFn: func(_ context.Context, id uuid.UUID, unreadOnly bool) ([]*model.Notification, error) {
			if id != userID || unreadOnly {
				return nil, nil
			}
			return notifs, nil
		},
	})
	got, err := svc.ListForUser(ctx, userID, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 {
		t.Errorf("want 2 notifications, got %d", len(got))
	}
}

func TestNotificationService_ListForUser_UnreadOnly(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	var capturedUnread bool

	svc := NewNotificationService(&mockNotificationRepo{
		listByUserFn: func(_ context.Context, _ uuid.UUID, unreadOnly bool) ([]*model.Notification, error) {
			capturedUnread = unreadOnly
			return nil, nil
		},
	})
	svc.ListForUser(ctx, userID, true) //nolint:errcheck
	if !capturedUnread {
		t.Error("expected unreadOnly=true to be passed through")
	}
}

// ── UnreadCount ───────────────────────────────────────────────────────────────

func TestNotificationService_UnreadCount(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	svc := NewNotificationService(&mockNotificationRepo{
		unreadCountFn: func(_ context.Context, id uuid.UUID) (int, error) {
			if id == userID {
				return 5, nil
			}
			return 0, nil
		},
	})
	count, err := svc.UnreadCount(ctx, userID)
	if err != nil {
		t.Fatal(err)
	}
	if count != 5 {
		t.Errorf("want 5, got %d", count)
	}
}

// ── MarkRead / MarkAllRead ────────────────────────────────────────────────────

func TestNotificationService_MarkRead(t *testing.T) {
	ctx := context.Background()
	notifID := uuid.New()
	userID := uuid.New()
	var capturedID, capturedUserID uuid.UUID

	svc := NewNotificationService(&mockNotificationRepo{
		markReadFn: func(_ context.Context, id, uid uuid.UUID) error {
			capturedID = id
			capturedUserID = uid
			return nil
		},
	})
	if err := svc.MarkRead(ctx, notifID, userID); err != nil {
		t.Fatal(err)
	}
	if capturedID != notifID || capturedUserID != userID {
		t.Error("MarkRead called with wrong IDs")
	}
}

func TestNotificationService_MarkAllRead(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	var capturedUserID uuid.UUID

	svc := NewNotificationService(&mockNotificationRepo{
		markAllReadFn: func(_ context.Context, uid uuid.UUID) error {
			capturedUserID = uid
			return nil
		},
	})
	if err := svc.MarkAllRead(ctx, userID); err != nil {
		t.Fatal(err)
	}
	if capturedUserID != userID {
		t.Errorf("want userID %v, got %v", userID, capturedUserID)
	}
}

func TestNotificationService_MarkAllRead_Error(t *testing.T) {
	ctx := context.Background()
	dbErr := errors.New("db error")
	svc := NewNotificationService(&mockNotificationRepo{
		markAllReadFn: func(_ context.Context, _ uuid.UUID) error { return dbErr },
	})
	err := svc.MarkAllRead(ctx, uuid.New())
	if !errors.Is(err, dbErr) {
		t.Errorf("want %v, got %v", dbErr, err)
	}
}
