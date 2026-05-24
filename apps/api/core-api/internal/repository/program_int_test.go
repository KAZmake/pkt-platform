package repository_test

import (
	"context"
	"testing"

	"github.com/KAZmake/pkt-platform/apps/api/core-api/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProgramRepository_Create(t *testing.T) {
	ctx := context.Background()
	repo := repository.NewProgramRepository(testPool)

	got, err := repo.Create(ctx, repository.CreateProgramInput{
		Name:          "Интеграционный 2025",
		Rate:          8.5,
		MinAmount:     200_000,
		MaxAmount:     10_000_000,
		MinTermMonths: 12,
		MaxTermMonths: 84,
		ActivityTypes: []string{"crop_farming", "livestock"},
	})
	require.NoError(t, err)
	require.NotNil(t, got)

	assert.NotEmpty(t, got.ID)
	assert.Equal(t, "Интеграционный 2025", got.Name)
	assert.InDelta(t, 8.5, got.Rate, 0.001)
	assert.InDelta(t, 200_000.0, got.MinAmount, 0.01)
	assert.InDelta(t, 10_000_000.0, got.MaxAmount, 0.01)
	assert.Equal(t, 12, got.MinTermMonths)
	assert.Equal(t, 84, got.MaxTermMonths)
	assert.True(t, got.IsActive)
	assert.ElementsMatch(t, []string{"crop_farming", "livestock"}, got.ActivityTypes)
}

func TestProgramRepository_GetByID(t *testing.T) {
	ctx := context.Background()
	repo := repository.NewProgramRepository(testPool)

	created, err := repo.Create(ctx, repository.CreateProgramInput{
		Name: "GetByID Program", Rate: 5.0, MinAmount: 50_000,
		MaxAmount: 1_000_000, MinTermMonths: 3, MaxTermMonths: 36,
	})
	require.NoError(t, err)

	got, err := repo.GetByID(ctx, created.ID)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, created.ID, got.ID)
	assert.Equal(t, "GetByID Program", got.Name)
}

func TestProgramRepository_GetByID_NotFound(t *testing.T) {
	ctx := context.Background()
	repo := repository.NewProgramRepository(testPool)

	got, err := repo.GetByID(ctx, "00000000-0000-0000-0000-000000000000")
	require.NoError(t, err)
	assert.Nil(t, got)
}

func TestProgramRepository_List_ReturnsCreated(t *testing.T) {
	ctx := context.Background()
	repo := repository.NewProgramRepository(testPool)

	p, err := repo.Create(ctx, repository.CreateProgramInput{
		Name: "ListTest Program", Rate: 6.0, MinAmount: 100_000,
		MaxAmount: 2_000_000, MinTermMonths: 6, MaxTermMonths: 48,
	})
	require.NoError(t, err)

	all, err := repo.List(ctx, false)
	require.NoError(t, err)

	var found bool
	for _, prog := range all {
		if prog.ID == p.ID {
			found = true
			break
		}
	}
	assert.True(t, found, "newly created program should appear in List(false)")
}

func TestProgramRepository_List_ActiveOnly(t *testing.T) {
	ctx := context.Background()
	repo := repository.NewProgramRepository(testPool)

	active, err := repo.Create(ctx, repository.CreateProgramInput{
		Name: "Active Program", Rate: 7.0, MinAmount: 100_000,
		MaxAmount: 3_000_000, MinTermMonths: 6, MaxTermMonths: 60,
	})
	require.NoError(t, err)

	inactive, err := repo.Create(ctx, repository.CreateProgramInput{
		Name: "Inactive Program", Rate: 9.0, MinAmount: 100_000,
		MaxAmount: 3_000_000, MinTermMonths: 6, MaxTermMonths: 60,
	})
	require.NoError(t, err)

	err = repo.SetActive(ctx, inactive.ID, false)
	require.NoError(t, err)

	activeList, err := repo.List(ctx, true)
	require.NoError(t, err)

	var foundActive, foundInactive bool
	for _, p := range activeList {
		if p.ID == active.ID {
			foundActive = true
		}
		if p.ID == inactive.ID {
			foundInactive = true
		}
	}
	assert.True(t, foundActive, "active program must be in activeOnly list")
	assert.False(t, foundInactive, "inactive program must not be in activeOnly list")
}

func TestProgramRepository_Update(t *testing.T) {
	ctx := context.Background()
	repo := repository.NewProgramRepository(testPool)

	p, err := repo.Create(ctx, repository.CreateProgramInput{
		Name: "Before Update", Rate: 7.0, MinAmount: 100_000,
		MaxAmount: 2_000_000, MinTermMonths: 6, MaxTermMonths: 48,
	})
	require.NoError(t, err)

	newName := "After Update"
	newRate := 9.5
	updated, err := repo.Update(ctx, p.ID, repository.UpdateProgramInput{
		Name: &newName,
		Rate: &newRate,
	})
	require.NoError(t, err)
	require.NotNil(t, updated)

	assert.Equal(t, "After Update", updated.Name)
	assert.InDelta(t, 9.5, updated.Rate, 0.001)
	// unchanged fields preserved
	assert.Equal(t, p.MaxTermMonths, updated.MaxTermMonths)
}

func TestProgramRepository_SetActive(t *testing.T) {
	ctx := context.Background()
	repo := repository.NewProgramRepository(testPool)

	p, err := repo.Create(ctx, repository.CreateProgramInput{
		Name: "SetActive Test", Rate: 6.5, MinAmount: 100_000,
		MaxAmount: 2_000_000, MinTermMonths: 6, MaxTermMonths: 36,
	})
	require.NoError(t, err)
	assert.True(t, p.IsActive)

	err = repo.SetActive(ctx, p.ID, false)
	require.NoError(t, err)

	got, err := repo.GetByID(ctx, p.ID)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.False(t, got.IsActive)

	err = repo.SetActive(ctx, p.ID, true)
	require.NoError(t, err)

	got, err = repo.GetByID(ctx, p.ID)
	require.NoError(t, err)
	assert.True(t, got.IsActive)
}
