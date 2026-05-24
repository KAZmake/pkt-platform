package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/KAZmake/pkt-platform/apps/api/core-api/internal/model"
	"github.com/KAZmake/pkt-platform/apps/api/core-api/internal/repository"
)

// ── mock ─────────────────────────────────────────────────────────────────────

type mockProgramRepo struct {
	listFn      func(context.Context, bool) ([]*model.LoanProgram, error)
	getByIDFn   func(context.Context, string) (*model.LoanProgram, error)
	createFn    func(context.Context, repository.CreateProgramInput) (*model.LoanProgram, error)
	updateFn    func(context.Context, string, repository.UpdateProgramInput) (*model.LoanProgram, error)
	setActiveFn func(context.Context, string, bool) error
}

func (m *mockProgramRepo) List(ctx context.Context, active bool) ([]*model.LoanProgram, error) {
	if m.listFn != nil {
		return m.listFn(ctx, active)
	}
	return nil, nil
}
func (m *mockProgramRepo) GetByID(ctx context.Context, id string) (*model.LoanProgram, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, nil
}
func (m *mockProgramRepo) Create(ctx context.Context, inp repository.CreateProgramInput) (*model.LoanProgram, error) {
	if m.createFn != nil {
		return m.createFn(ctx, inp)
	}
	return nil, nil
}
func (m *mockProgramRepo) Update(ctx context.Context, id string, inp repository.UpdateProgramInput) (*model.LoanProgram, error) {
	if m.updateFn != nil {
		return m.updateFn(ctx, id, inp)
	}
	return nil, nil
}
func (m *mockProgramRepo) SetActive(ctx context.Context, id string, active bool) error {
	if m.setActiveFn != nil {
		return m.setActiveFn(ctx, id, active)
	}
	return nil
}

func stubProgram() *model.LoanProgram {
	return &model.LoanProgram{
		ID: "prog-1", Name: "Агро 2025", Rate: 7.5,
		MinAmount: 100_000, MaxAmount: 5_000_000,
		MinTermMonths: 6, MaxTermMonths: 60,
		IsActive: true, CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
}

func validCreateInput() repository.CreateProgramInput {
	return repository.CreateProgramInput{
		Name: "Агро 2025", Rate: 7.5,
		MinAmount: 100_000, MaxAmount: 5_000_000,
		MinTermMonths: 6, MaxTermMonths: 60,
	}
}

// ── Create validation ─────────────────────────────────────────────────────────

func TestProgramService_Create_Validation(t *testing.T) {
	ctx := context.Background()

	cases := []struct {
		name    string
		inp     repository.CreateProgramInput
		wantErr string
	}{
		{
			"missing name",
			repository.CreateProgramInput{Rate: 7, MinAmount: 100_000, MaxAmount: 1_000_000, MinTermMonths: 6, MaxTermMonths: 60},
			"name is required",
		},
		{
			"zero rate",
			repository.CreateProgramInput{Name: "X", Rate: 0, MinAmount: 100_000, MaxAmount: 1_000_000, MinTermMonths: 6, MaxTermMonths: 60},
			"rate must be positive",
		},
		{
			"negative rate",
			repository.CreateProgramInput{Name: "X", Rate: -0.1, MinAmount: 100_000, MaxAmount: 1_000_000, MinTermMonths: 6, MaxTermMonths: 60},
			"rate must be positive",
		},
		{
			"zero min amount",
			repository.CreateProgramInput{Name: "X", Rate: 7, MinAmount: 0, MaxAmount: 1_000_000, MinTermMonths: 6, MaxTermMonths: 60},
			"invalid amount range",
		},
		{
			"max amount equals min amount",
			repository.CreateProgramInput{Name: "X", Rate: 7, MinAmount: 500_000, MaxAmount: 500_000, MinTermMonths: 6, MaxTermMonths: 60},
			"invalid amount range",
		},
		{
			"max amount less than min amount",
			repository.CreateProgramInput{Name: "X", Rate: 7, MinAmount: 2_000_000, MaxAmount: 1_000_000, MinTermMonths: 6, MaxTermMonths: 60},
			"invalid amount range",
		},
		{
			"zero min term",
			repository.CreateProgramInput{Name: "X", Rate: 7, MinAmount: 100_000, MaxAmount: 1_000_000, MinTermMonths: 0, MaxTermMonths: 60},
			"invalid term range",
		},
		{
			"max term less than min term",
			repository.CreateProgramInput{Name: "X", Rate: 7, MinAmount: 100_000, MaxAmount: 1_000_000, MinTermMonths: 24, MaxTermMonths: 12},
			"invalid term range",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			svc := NewProgramService(&mockProgramRepo{})
			_, err := svc.Create(ctx, tc.inp)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if err.Error() != tc.wantErr {
				t.Errorf("want %q, got %q", tc.wantErr, err.Error())
			}
		})
	}
}

func TestProgramService_Create_Success(t *testing.T) {
	ctx := context.Background()
	want := stubProgram()

	mock := &mockProgramRepo{
		createFn: func(_ context.Context, _ repository.CreateProgramInput) (*model.LoanProgram, error) {
			return want, nil
		},
	}
	svc := NewProgramService(mock)

	got, err := svc.Create(ctx, validCreateInput())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil || got.ID != want.ID {
		t.Errorf("want ID %q, got %v", want.ID, got)
	}
}

func TestProgramService_Create_RepoError(t *testing.T) {
	ctx := context.Background()
	wantErr := errors.New("db error")
	mock := &mockProgramRepo{
		createFn: func(_ context.Context, _ repository.CreateProgramInput) (*model.LoanProgram, error) {
			return nil, wantErr
		},
	}
	svc := NewProgramService(mock)
	_, err := svc.Create(ctx, validCreateInput())
	if !errors.Is(err, wantErr) {
		t.Errorf("want %v, got %v", wantErr, err)
	}
}

// ── List ──────────────────────────────────────────────────────────────────────

func TestProgramService_ListActive(t *testing.T) {
	ctx := context.Background()
	var gotActive bool
	mock := &mockProgramRepo{
		listFn: func(_ context.Context, active bool) ([]*model.LoanProgram, error) {
			gotActive = active
			return []*model.LoanProgram{stubProgram()}, nil
		},
	}
	svc := NewProgramService(mock)
	progs, err := svc.ListActive(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if !gotActive {
		t.Error("expected activeOnly=true")
	}
	if len(progs) != 1 {
		t.Errorf("want 1 program, got %d", len(progs))
	}
}

func TestProgramService_ListAll(t *testing.T) {
	ctx := context.Background()
	var gotActive bool
	mock := &mockProgramRepo{
		listFn: func(_ context.Context, active bool) ([]*model.LoanProgram, error) {
			gotActive = active
			return nil, nil
		},
	}
	svc := NewProgramService(mock)
	svc.ListAll(ctx) //nolint:errcheck
	if gotActive {
		t.Error("expected activeOnly=false")
	}
}

// ── GetByID ───────────────────────────────────────────────────────────────────

func TestProgramService_GetByID_Found(t *testing.T) {
	ctx := context.Background()
	want := stubProgram()
	mock := &mockProgramRepo{
		getByIDFn: func(_ context.Context, id string) (*model.LoanProgram, error) {
			if id == want.ID {
				return want, nil
			}
			return nil, nil
		},
	}
	svc := NewProgramService(mock)
	got, err := svc.GetByID(ctx, want.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got == nil || got.ID != want.ID {
		t.Errorf("want %q, got %v", want.ID, got)
	}
}

func TestProgramService_GetByID_NotFound(t *testing.T) {
	ctx := context.Background()
	mock := &mockProgramRepo{
		getByIDFn: func(_ context.Context, _ string) (*model.LoanProgram, error) { return nil, nil },
	}
	svc := NewProgramService(mock)
	got, err := svc.GetByID(ctx, "missing")
	if err != nil {
		t.Fatal(err)
	}
	if got != nil {
		t.Error("expected nil, got program")
	}
}

// ── Update ────────────────────────────────────────────────────────────────────

func TestProgramService_Update_NotFound(t *testing.T) {
	ctx := context.Background()
	mock := &mockProgramRepo{
		getByIDFn: func(_ context.Context, _ string) (*model.LoanProgram, error) { return nil, nil },
	}
	svc := NewProgramService(mock)
	got, err := svc.Update(ctx, "missing", repository.UpdateProgramInput{})
	if err != nil {
		t.Fatal(err)
	}
	if got != nil {
		t.Error("expected nil for not-found update")
	}
}

func TestProgramService_Update_RepoError(t *testing.T) {
	ctx := context.Background()
	dbErr := errors.New("connection refused")
	mock := &mockProgramRepo{
		getByIDFn: func(_ context.Context, _ string) (*model.LoanProgram, error) { return nil, dbErr },
	}
	svc := NewProgramService(mock)
	_, err := svc.Update(ctx, "id", repository.UpdateProgramInput{})
	if !errors.Is(err, dbErr) {
		t.Errorf("want %v, got %v", dbErr, err)
	}
}

func TestProgramService_Update_Success(t *testing.T) {
	ctx := context.Background()
	orig := stubProgram()
	updated := stubProgram()
	updated.Rate = 9.0

	mock := &mockProgramRepo{
		getByIDFn: func(_ context.Context, _ string) (*model.LoanProgram, error) { return orig, nil },
		updateFn: func(_ context.Context, _ string, _ repository.UpdateProgramInput) (*model.LoanProgram, error) {
			return updated, nil
		},
	}
	svc := NewProgramService(mock)
	got, err := svc.Update(ctx, orig.ID, repository.UpdateProgramInput{})
	if err != nil {
		t.Fatal(err)
	}
	if got == nil || got.Rate != 9.0 {
		t.Errorf("unexpected result: %v", got)
	}
}

// ── Deactivate ────────────────────────────────────────────────────────────────

func TestProgramService_Deactivate_NotFound(t *testing.T) {
	ctx := context.Background()
	mock := &mockProgramRepo{
		getByIDFn: func(_ context.Context, _ string) (*model.LoanProgram, error) { return nil, nil },
	}
	svc := NewProgramService(mock)
	err := svc.Deactivate(ctx, "missing")
	if err == nil {
		t.Fatal("expected error")
	}
	if err.Error() != "not found" {
		t.Errorf("want 'not found', got %q", err.Error())
	}
}

func TestProgramService_Deactivate_Success(t *testing.T) {
	ctx := context.Background()
	var capturedID string
	var capturedActive bool

	mock := &mockProgramRepo{
		getByIDFn: func(_ context.Context, _ string) (*model.LoanProgram, error) { return stubProgram(), nil },
		setActiveFn: func(_ context.Context, id string, active bool) error {
			capturedID = id
			capturedActive = active
			return nil
		},
	}
	svc := NewProgramService(mock)
	if err := svc.Deactivate(ctx, "prog-1"); err != nil {
		t.Fatal(err)
	}
	if capturedID != "prog-1" {
		t.Errorf("want ID 'prog-1', got %q", capturedID)
	}
	if capturedActive {
		t.Error("SetActive should be called with false")
	}
}

func TestProgramService_Deactivate_RepoError(t *testing.T) {
	ctx := context.Background()
	dbErr := errors.New("connection refused")
	mock := &mockProgramRepo{
		getByIDFn: func(_ context.Context, _ string) (*model.LoanProgram, error) { return nil, dbErr },
	}
	svc := NewProgramService(mock)
	err := svc.Deactivate(ctx, "id")
	if !errors.Is(err, dbErr) {
		t.Errorf("want %v, got %v", dbErr, err)
	}
}
