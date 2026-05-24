package repository_test

import (
	"context"
	"testing"

	"github.com/KAZmake/pkt-platform/apps/api/core-api/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNotificationRepository_Create(t *testing.T) {
	ctx := context.Background()
	user := insertUser(t)
	repo := repository.NewNotificationRepository(testPool)

	got, err := repo.Create(ctx, repository.CreateNotificationInput{
		UserID: user.ID,
		Type:   "status",
		Title:  "Заявка обновлена",
		Body:   "Ваша заявка перешла на этап скоринга",
	})
	require.NoError(t, err)
	require.NotNil(t, got)

	assert.NotEmpty(t, got.ID)
	assert.Equal(t, user.ID, got.UserID)
	assert.Equal(t, "status", got.Type)
	assert.Equal(t, "Заявка обновлена", got.Title)
	assert.Equal(t, "Ваша заявка перешла на этап скоринга", got.Body)
	assert.False(t, got.IsRead)
}

func TestNotificationRepository_ListByUser(t *testing.T) {
	ctx := context.Background()
	user := insertUser(t)
	other := insertUser(t)
	repo := repository.NewNotificationRepository(testPool)

	for i := 0; i < 3; i++ {
		_, err := repo.Create(ctx, repository.CreateNotificationInput{
			UserID: user.ID, Type: "system",
			Title: "Уведомление", Body: "Тело",
		})
		require.NoError(t, err)
	}
	// noise: notification for different user
	_, err := repo.Create(ctx, repository.CreateNotificationInput{
		UserID: other.ID, Type: "system", Title: "Другой", Body: "Другой",
	})
	require.NoError(t, err)

	list, err := repo.ListByUser(ctx, user.ID, false)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(list), 3, "should return at least 3 notifications for user")

	for _, n := range list {
		assert.Equal(t, user.ID, n.UserID)
	}
}

func TestNotificationRepository_ListByUser_UnreadOnly(t *testing.T) {
	ctx := context.Background()
	user := insertUser(t)
	repo := repository.NewNotificationRepository(testPool)

	// 2 notifications; we'll mark one as read
	n1, err := repo.Create(ctx, repository.CreateNotificationInput{
		UserID: user.ID, Type: "status", Title: "Первое", Body: "Тело",
	})
	require.NoError(t, err)

	n2, err := repo.Create(ctx, repository.CreateNotificationInput{
		UserID: user.ID, Type: "status", Title: "Второе", Body: "Тело",
	})
	require.NoError(t, err)

	err = repo.MarkRead(ctx, n1.ID, user.ID)
	require.NoError(t, err)

	unread, err := repo.ListByUser(ctx, user.ID, true)
	require.NoError(t, err)

	var foundRead, foundUnread bool
	for _, n := range unread {
		if n.ID == n1.ID {
			foundRead = true
		}
		if n.ID == n2.ID {
			foundUnread = true
		}
	}
	assert.False(t, foundRead, "read notification must not appear in unread list")
	assert.True(t, foundUnread, "unread notification must appear in unread list")
}

func TestNotificationRepository_UnreadCount(t *testing.T) {
	ctx := context.Background()
	user := insertUser(t)
	repo := repository.NewNotificationRepository(testPool)

	count0, err := repo.UnreadCount(ctx, user.ID)
	require.NoError(t, err)
	assert.Equal(t, 0, count0)

	for i := 0; i < 4; i++ {
		_, err := repo.Create(ctx, repository.CreateNotificationInput{
			UserID: user.ID, Type: "ticket", Title: "T", Body: "B",
		})
		require.NoError(t, err)
	}

	count4, err := repo.UnreadCount(ctx, user.ID)
	require.NoError(t, err)
	assert.Equal(t, 4, count4)
}

func TestNotificationRepository_MarkRead(t *testing.T) {
	ctx := context.Background()
	user := insertUser(t)
	repo := repository.NewNotificationRepository(testPool)

	notif, err := repo.Create(ctx, repository.CreateNotificationInput{
		UserID: user.ID, Type: "payment", Title: "Платёж", Body: "Напоминание",
	})
	require.NoError(t, err)
	assert.False(t, notif.IsRead)

	err = repo.MarkRead(ctx, notif.ID, user.ID)
	require.NoError(t, err)

	count, err := repo.UnreadCount(ctx, user.ID)
	require.NoError(t, err)
	assert.Equal(t, 0, count)
}

func TestNotificationRepository_MarkAllRead(t *testing.T) {
	ctx := context.Background()
	user := insertUser(t)
	repo := repository.NewNotificationRepository(testPool)

	for i := 0; i < 3; i++ {
		_, err := repo.Create(ctx, repository.CreateNotificationInput{
			UserID: user.ID, Type: "system", Title: "Сист.", Body: "B",
		})
		require.NoError(t, err)
	}

	before, err := repo.UnreadCount(ctx, user.ID)
	require.NoError(t, err)
	assert.Equal(t, 3, before)

	err = repo.MarkAllRead(ctx, user.ID)
	require.NoError(t, err)

	after, err := repo.UnreadCount(ctx, user.ID)
	require.NoError(t, err)
	assert.Equal(t, 0, after)
}
