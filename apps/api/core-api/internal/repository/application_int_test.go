package repository_test

import (
	"context"
	"testing"

	"github.com/KAZmake/pkt-platform/apps/api/core-api/internal/model"
	"github.com/KAZmake/pkt-platform/apps/api/core-api/internal/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ── fixture ───────────────────────────────────────────────────────────────────

func newApp(t *testing.T) *model.Application {
	t.Helper()
	ctx := context.Background()
	user := insertUser(t)
	borrower := insertBorrower(t, user.ID)
	program := insertProgram(t)
	programID := uuid.MustParse(program.ID)

	repo := repository.NewApplicationRepository(testPool)
	app, err := repo.Create(ctx, repository.CreateApplicationInput{
		BorrowerID:  borrower.ID,
		ProgramID:   programID,
		Amount:      500_000,
		TermMonths:  12,
		PaymentType: "annuity",
	})
	require.NoError(t, err)
	require.NotNil(t, app)
	return app
}

// ── tests ─────────────────────────────────────────────────────────────────────

func TestApplicationRepository_Create(t *testing.T) {
	ctx := context.Background()
	user := insertUser(t)
	borrower := insertBorrower(t, user.ID)
	program := insertProgram(t)
	repo := repository.NewApplicationRepository(testPool)

	got, err := repo.Create(ctx, repository.CreateApplicationInput{
		BorrowerID:  borrower.ID,
		ProgramID:   uuid.MustParse(program.ID),
		Amount:      750_000,
		TermMonths:  24,
		PaymentType: "differentiated",
	})
	require.NoError(t, err)
	require.NotNil(t, got)

	assert.NotEmpty(t, got.ID)
	assert.Equal(t, borrower.ID, got.BorrowerID)
	assert.Equal(t, uuid.MustParse(program.ID), got.ProgramID)
	assert.Equal(t, model.StatusReceived, got.Status)
	assert.InDelta(t, 750_000.0, got.Amount, 0.01)
	assert.Equal(t, 24, got.TermMonths)
	assert.Equal(t, "differentiated", got.PaymentType)
}

func TestApplicationRepository_GetByID(t *testing.T) {
	ctx := context.Background()
	app := newApp(t)
	repo := repository.NewApplicationRepository(testPool)

	got, err := repo.GetByID(ctx, app.ID)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, app.ID, got.ID)
	assert.Equal(t, model.StatusReceived, got.Status)
}

func TestApplicationRepository_GetByID_NotFound(t *testing.T) {
	ctx := context.Background()
	repo := repository.NewApplicationRepository(testPool)

	got, err := repo.GetByID(ctx, uuid.New())
	require.NoError(t, err)
	assert.Nil(t, got)
}

func TestApplicationRepository_ListByBorrower(t *testing.T) {
	ctx := context.Background()
	user := insertUser(t)
	borrower := insertBorrower(t, user.ID)
	program := insertProgram(t)
	repo := repository.NewApplicationRepository(testPool)

	programID := uuid.MustParse(program.ID)
	for i := 0; i < 3; i++ {
		_, err := repo.Create(ctx, repository.CreateApplicationInput{
			BorrowerID: borrower.ID, ProgramID: programID,
			Amount: 100_000, TermMonths: 6, PaymentType: "annuity",
		})
		require.NoError(t, err)
	}

	list, err := repo.ListByBorrower(ctx, borrower.ID)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(list), 3)
	for _, a := range list {
		assert.Equal(t, borrower.ID, a.BorrowerID)
	}
}

func TestApplicationRepository_ListAll_NoFilter(t *testing.T) {
	ctx := context.Background()
	app := newApp(t)
	repo := repository.NewApplicationRepository(testPool)

	all, err := repo.ListAll(ctx, "")
	require.NoError(t, err)

	var found bool
	for _, a := range all {
		if a.ID == app.ID {
			found = true
			break
		}
	}
	assert.True(t, found, "created application must appear in ListAll with no filter")
}

func TestApplicationRepository_ListAll_WithStatusFilter(t *testing.T) {
	ctx := context.Background()
	app := newApp(t)
	actor := insertUser(t)
	repo := repository.NewApplicationRepository(testPool)

	_, err := repo.UpdateStatus(ctx, app.ID, actor.ID, model.StatusPrimaryScoring, nil)
	require.NoError(t, err)

	received, err := repo.ListAll(ctx, model.StatusReceived)
	require.NoError(t, err)

	scoring, err := repo.ListAll(ctx, model.StatusPrimaryScoring)
	require.NoError(t, err)

	var inReceived, inScoring bool
	for _, a := range received {
		if a.ID == app.ID {
			inReceived = true
		}
	}
	for _, a := range scoring {
		if a.ID == app.ID {
			inScoring = true
		}
	}

	assert.False(t, inReceived, "app should no longer be in received list after advancing")
	assert.True(t, inScoring, "app should be in primary_scoring list")
}

func TestApplicationRepository_UpdateStatus(t *testing.T) {
	ctx := context.Background()
	app := newApp(t)
	actor := insertUser(t)
	repo := repository.NewApplicationRepository(testPool)

	comment := "Прошёл первичный скоринг"
	updated, err := repo.UpdateStatus(ctx, app.ID, actor.ID, model.StatusPrimaryScoring, &comment)
	require.NoError(t, err)
	require.NotNil(t, updated)
	assert.Equal(t, model.StatusPrimaryScoring, updated.Status)

	got, err := repo.GetByID(ctx, app.ID)
	require.NoError(t, err)
	assert.Equal(t, model.StatusPrimaryScoring, got.Status)
}

func TestApplicationRepository_GetHistory(t *testing.T) {
	ctx := context.Background()
	app := newApp(t)
	actor := insertUser(t)
	repo := repository.NewApplicationRepository(testPool)

	comment := "Первый переход"
	_, err := repo.UpdateStatus(ctx, app.ID, actor.ID, model.StatusPrimaryScoring, &comment)
	require.NoError(t, err)

	_, err = repo.UpdateStatus(ctx, app.ID, actor.ID, model.StatusSecurityCheck, nil)
	require.NoError(t, err)

	history, err := repo.GetHistory(ctx, app.ID)
	require.NoError(t, err)
	require.Len(t, history, 2)

	h0 := history[0]
	require.NotNil(t, h0.FromStatus)
	assert.Equal(t, model.StatusReceived, *h0.FromStatus)
	assert.Equal(t, model.StatusPrimaryScoring, h0.ToStatus)
	assert.Equal(t, actor.ID, h0.ActorID)
	require.NotNil(t, h0.Comment)
	assert.Equal(t, "Первый переход", *h0.Comment)

	h1 := history[1]
	require.NotNil(t, h1.FromStatus)
	assert.Equal(t, model.StatusPrimaryScoring, *h1.FromStatus)
	assert.Equal(t, model.StatusSecurityCheck, h1.ToStatus)
	assert.Nil(t, h1.Comment)
}

func TestApplicationRepository_History_IsInsertOnly(t *testing.T) {
	ctx := context.Background()
	app := newApp(t)
	actor := insertUser(t)
	repo := repository.NewApplicationRepository(testPool)

	_, err := repo.UpdateStatus(ctx, app.ID, actor.ID, model.StatusPrimaryScoring, nil)
	require.NoError(t, err)

	history, err := repo.GetHistory(ctx, app.ID)
	require.NoError(t, err)
	require.Len(t, history, 1)
	histID := history[0].ID

	// PostgreSQL rule silently swallows UPDATE on application_history
	_, execErr := testPool.Exec(ctx,
		`UPDATE application_history SET comment = 'hacked' WHERE id = $1`, histID)
	require.NoError(t, execErr)

	after, err := repo.GetHistory(ctx, app.ID)
	require.NoError(t, err)
	assert.Nil(t, after[0].Comment, "UPDATE must be silently ignored by DB rule")
}
