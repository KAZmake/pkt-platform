package service

import (
	"context"

	"github.com/KAZmake/pkt-platform/apps/api/core-api/internal/model"
	"github.com/KAZmake/pkt-platform/apps/api/core-api/internal/repository"
	"github.com/google/uuid"
)

type notificationRepo interface {
	Create(ctx context.Context, inp repository.CreateNotificationInput) (*model.Notification, error)
	ListByUser(ctx context.Context, userID uuid.UUID, unreadOnly bool) ([]*model.Notification, error)
	UnreadCount(ctx context.Context, userID uuid.UUID) (int, error)
	MarkRead(ctx context.Context, id, userID uuid.UUID) error
	MarkAllRead(ctx context.Context, userID uuid.UUID) error
}

type NotificationService struct {
	repo notificationRepo
}

func NewNotificationService(repo notificationRepo) *NotificationService {
	return &NotificationService{repo: repo}
}

func (s *NotificationService) Create(ctx context.Context, inp repository.CreateNotificationInput) (*model.Notification, error) {
	return s.repo.Create(ctx, inp)
}

func (s *NotificationService) ListForUser(ctx context.Context, userID uuid.UUID, unreadOnly bool) ([]*model.Notification, error) {
	return s.repo.ListByUser(ctx, userID, unreadOnly)
}

func (s *NotificationService) UnreadCount(ctx context.Context, userID uuid.UUID) (int, error) {
	return s.repo.UnreadCount(ctx, userID)
}

func (s *NotificationService) MarkRead(ctx context.Context, id, userID uuid.UUID) error {
	return s.repo.MarkRead(ctx, id, userID)
}

func (s *NotificationService) MarkAllRead(ctx context.Context, userID uuid.UUID) error {
	return s.repo.MarkAllRead(ctx, userID)
}
