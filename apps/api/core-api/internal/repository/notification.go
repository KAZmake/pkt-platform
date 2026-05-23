package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/KAZmake/pkt-platform/apps/api/core-api/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NotificationRepository struct {
	db *pgxpool.Pool
}

func NewNotificationRepository(db *pgxpool.Pool) *NotificationRepository {
	return &NotificationRepository{db: db}
}

const notifCols = `id, user_id, type, title, body, is_read, created_at`

type CreateNotificationInput struct {
	UserID uuid.UUID
	Type   string
	Title  string
	Body   string
}

func (r *NotificationRepository) Create(ctx context.Context, inp CreateNotificationInput) (*model.Notification, error) {
	row := r.db.QueryRow(ctx, `
		INSERT INTO notifications (user_id, type, title, body)
		VALUES ($1, $2, $3, $4)
		RETURNING `+notifCols,
		inp.UserID, inp.Type, inp.Title, inp.Body,
	)
	return scanNotification(row)
}

func (r *NotificationRepository) ListByUser(ctx context.Context, userID uuid.UUID, unreadOnly bool) ([]*model.Notification, error) {
	q := `SELECT ` + notifCols + ` FROM notifications WHERE user_id = $1`
	if unreadOnly {
		q += ` AND is_read = FALSE`
	}
	q += ` ORDER BY created_at DESC LIMIT 100`

	rows, err := r.db.Query(ctx, q, userID)
	if err != nil {
		return nil, fmt.Errorf("list notifications: %w", err)
	}
	defer rows.Close()

	var notifs []*model.Notification
	for rows.Next() {
		n, err := scanNotification(rows)
		if err != nil {
			return nil, err
		}
		notifs = append(notifs, n)
	}
	return notifs, nil
}

func (r *NotificationRepository) UnreadCount(ctx context.Context, userID uuid.UUID) (int, error) {
	var count int
	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM notifications WHERE user_id = $1 AND is_read = FALSE`, userID,
	).Scan(&count)
	return count, err
}

func (r *NotificationRepository) MarkRead(ctx context.Context, id, userID uuid.UUID) error {
	_, err := r.db.Exec(ctx,
		`UPDATE notifications SET is_read = TRUE WHERE id = $1 AND user_id = $2`, id, userID)
	return err
}

func (r *NotificationRepository) MarkAllRead(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.Exec(ctx,
		`UPDATE notifications SET is_read = TRUE WHERE user_id = $1 AND is_read = FALSE`, userID)
	return err
}

func scanNotification(s scanner) (*model.Notification, error) {
	n := &model.Notification{}
	err := s.Scan(&n.ID, &n.UserID, &n.Type, &n.Title, &n.Body, &n.IsRead, &n.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("scan notification: %w", err)
	}
	return n, nil
}
