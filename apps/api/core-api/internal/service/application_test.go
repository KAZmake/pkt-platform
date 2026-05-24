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

type mockApplicationRepo struct {
	createFn         func(context.Context, repository.CreateApplicationInput) (*model.Application, error)
	getByIDFn        func(context.Context, uuid.UUID) (*model.Application, error)
	getHistoryFn     func(context.Context, uuid.UUID) ([]*model.ApplicationHistory, error)
	listByBorrowerFn func(context.Context, uuid.UUID) ([]*model.Application, error)
	listAllFn        func(context.Context, string) ([]*model.Application, error)
	updateStatusFn   func(context.Context, uuid.UUID, uuid.UUID, string, *string) (*model.Application, error)
}

func (m *mockApplicationRepo) Create(ctx context.Context, inp repository.CreateApplicationInput) (*model.Application, error) {
	if m.createFn != nil {
		return m.createFn(ctx, inp)
	}
	return nil, nil
}
func (m *mockApplicationRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.Application, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, nil
}
func (m *mockApplicationRepo) GetHistory(ctx context.Context, id uuid.UUID) ([]*model.ApplicationHistory, error) {
	if m.getHistoryFn != nil {
		return m.getHistoryFn(ctx, id)
	}
	return nil, nil
}
func (m *mockApplicationRepo) ListByBorrower(ctx context.Context, borrowerID uuid.UUID) ([]*model.Application, error) {
	if m.listByBorrowerFn != nil {
		return m.listByBorrowerFn(ctx, borrowerID)
	}
	return nil, nil
}
func (m *mockApplicationRepo) ListAll(ctx context.Context, status string) ([]*model.Application, error) {
	if m.listAllFn != nil {
		return m.listAllFn(ctx, status)
	}
	return nil, nil
}
func (m *mockApplicationRepo) UpdateStatus(ctx context.Context, appID, actorID uuid.UUID, toStatus string, comment *string) (*model.Application, error) {
	if m.updateStatusFn != nil {
		return m.updateStatusFn(ctx, appID, actorID, toStatus, comment)
	}
	return nil, nil
}

func stubApplication(status string) *model.Application {
	return &model.Application{
		ID:          uuid.New(),
		BorrowerID:  uuid.New(),
		ProgramID:   uuid.New(),
		Status:      status,
		Amount:      500_000,
		TermMonths:  12,
		PaymentType: "annuity",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// ── Create validation ─────────────────────────────────────────────────────────

func TestApplicationService_Create_Validation(t *testing.T) {
	ctx := context.Background()
	borrowerID := uuid.New()
	programID := uuid.New()

	cases := []struct {
		name    string
		inp     CreateApplicationInput
		wantErr string
	}{
		{
			"zero amount",
			CreateApplicationInput{ProgramID: programID, Amount: 0, TermMonths: 12, PaymentType: "annuity"},
			"amount must be positive",
		},
		{
			"negative amount",
			CreateApplicationInput{ProgramID: programID, Amount: -1, TermMonths: 12, PaymentType: "annuity"},
			"amount must be positive",
		},
		{
			"zero term",
			CreateApplicationInput{ProgramID: programID, Amount: 500_000, TermMonths: 0, PaymentType: "annuity"},
			"term_months must be positive",
		},
		{
			"invalid payment type",
			CreateApplicationInput{ProgramID: programID, Amount: 500_000, TermMonths: 12, PaymentType: "bullet"},
			"payment_type must be annuity or differentiated",
		},
		{
			"empty payment type",
			CreateApplicationInput{ProgramID: programID, Amount: 500_000, TermMonths: 12, PaymentType: ""},
			"payment_type must be annuity or differentiated",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			svc := NewApplicationService(&mockApplicationRepo{}, nil)
			_, err := svc.Create(ctx, borrowerID, tc.inp)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if err.Error() != tc.wantErr {
				t.Errorf("want %q, got %q", tc.wantErr, err.Error())
			}
		})
	}
}

func TestApplicationService_Create_Annuity_Success(t *testing.T) {
	ctx := context.Background()
	borrowerID := uuid.New()
	want := stubApplication(model.StatusReceived)

	mock := &mockApplicationRepo{
		createFn: func(_ context.Context, _ repository.CreateApplicationInput) (*model.Application, error) {
			return want, nil
		},
	}
	svc := NewApplicationService(mock, nil)
	got, err := svc.Create(ctx, borrowerID, CreateApplicationInput{
		ProgramID: want.ProgramID, Amount: 500_000, TermMonths: 12, PaymentType: "annuity",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil || got.ID != want.ID {
		t.Errorf("want ID %v, got %v", want.ID, got)
	}
}

func TestApplicationService_Create_Differentiated_Success(t *testing.T) {
	ctx := context.Background()
	borrowerID := uuid.New()
	want := stubApplication(model.StatusReceived)

	mock := &mockApplicationRepo{
		createFn: func(_ context.Context, _ repository.CreateApplicationInput) (*model.Application, error) {
			return want, nil
		},
	}
	svc := NewApplicationService(mock, nil)
	got, err := svc.Create(ctx, borrowerID, CreateApplicationInput{
		ProgramID: want.ProgramID, Amount: 500_000, TermMonths: 24, PaymentType: "differentiated",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil {
		t.Fatal("expected application, got nil")
	}
}

func TestApplicationService_Create_RepoError(t *testing.T) {
	ctx := context.Background()
	dbErr := errors.New("db error")
	mock := &mockApplicationRepo{
		createFn: func(_ context.Context, _ repository.CreateApplicationInput) (*model.Application, error) {
			return nil, dbErr
		},
	}
	svc := NewApplicationService(mock, nil)
	_, err := svc.Create(ctx, uuid.New(), CreateApplicationInput{
		ProgramID: uuid.New(), Amount: 500_000, TermMonths: 12, PaymentType: "annuity",
	})
	if !errors.Is(err, dbErr) {
		t.Errorf("want %v, got %v", dbErr, err)
	}
}

// ── ChangeStatus ──────────────────────────────────────────────────────────────

func TestApplicationService_ChangeStatus_NotFound(t *testing.T) {
	ctx := context.Background()
	mock := &mockApplicationRepo{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (*model.Application, error) { return nil, nil },
	}
	svc := NewApplicationService(mock, nil)
	_, err := svc.ChangeStatus(ctx, uuid.New(), uuid.New(), ChangeStatusInput{ToStatus: model.StatusPrimaryScoring})
	if err == nil {
		t.Fatal("expected error")
	}
	if err.Error() != "application not found" {
		t.Errorf("want 'application not found', got %q", err.Error())
	}
}

func TestApplicationService_ChangeStatus_GetByIDError(t *testing.T) {
	ctx := context.Background()
	dbErr := errors.New("db error")
	mock := &mockApplicationRepo{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (*model.Application, error) { return nil, dbErr },
	}
	svc := NewApplicationService(mock, nil)
	_, err := svc.ChangeStatus(ctx, uuid.New(), uuid.New(), ChangeStatusInput{ToStatus: model.StatusPrimaryScoring})
	if !errors.Is(err, dbErr) {
		t.Errorf("want %v, got %v", dbErr, err)
	}
}

func TestApplicationService_ChangeStatus_InvalidTransition(t *testing.T) {
	ctx := context.Background()
	app := stubApplication(model.StatusReceived)
	mock := &mockApplicationRepo{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (*model.Application, error) { return app, nil },
	}
	svc := NewApplicationService(mock, nil)
	// received → issued is not allowed
	_, err := svc.ChangeStatus(ctx, app.ID, uuid.New(), ChangeStatusInput{ToStatus: model.StatusIssued})
	if err == nil {
		t.Fatal("expected error for invalid transition")
	}
}

func TestApplicationService_ChangeStatus_InvalidTransition_Terminal(t *testing.T) {
	ctx := context.Background()
	app := stubApplication(model.StatusRejected)
	mock := &mockApplicationRepo{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (*model.Application, error) { return app, nil },
	}
	svc := NewApplicationService(mock, nil)
	_, err := svc.ChangeStatus(ctx, app.ID, uuid.New(), ChangeStatusInput{ToStatus: model.StatusApproved})
	if err == nil {
		t.Fatal("expected error for transition from terminal status")
	}
}

func TestApplicationService_ChangeStatus_ValidTransition(t *testing.T) {
	ctx := context.Background()
	appID := uuid.New()
	actorID := uuid.New()

	app := stubApplication(model.StatusReceived)
	app.ID = appID

	updated := stubApplication(model.StatusPrimaryScoring)
	updated.ID = appID

	var capturedToStatus string
	mock := &mockApplicationRepo{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (*model.Application, error) { return app, nil },
		updateStatusFn: func(_ context.Context, _, _ uuid.UUID, toStatus string, _ *string) (*model.Application, error) {
			capturedToStatus = toStatus
			return updated, nil
		},
	}
	svc := NewApplicationService(mock, nil)
	got, err := svc.ChangeStatus(ctx, appID, actorID, ChangeStatusInput{ToStatus: model.StatusPrimaryScoring})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil || got.Status != model.StatusPrimaryScoring {
		t.Errorf("want status %q, got %v", model.StatusPrimaryScoring, got)
	}
	if capturedToStatus != model.StatusPrimaryScoring {
		t.Errorf("UpdateStatus called with %q, want %q", capturedToStatus, model.StatusPrimaryScoring)
	}
}

func TestApplicationService_ChangeStatus_WithComment(t *testing.T) {
	ctx := context.Background()
	app := stubApplication(model.StatusPrimaryScoring)
	updated := stubApplication(model.StatusRevision)

	comment := "Нужны дополнительные документы"
	var capturedComment *string

	mock := &mockApplicationRepo{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (*model.Application, error) { return app, nil },
		updateStatusFn: func(_ context.Context, _, _ uuid.UUID, _ string, c *string) (*model.Application, error) {
			capturedComment = c
			return updated, nil
		},
	}
	svc := NewApplicationService(mock, nil)
	svc.ChangeStatus(ctx, app.ID, uuid.New(), ChangeStatusInput{ //nolint:errcheck
		ToStatus: model.StatusRevision,
		Comment:  &comment,
	})
	if capturedComment == nil || *capturedComment != comment {
		t.Errorf("want comment %q, got %v", comment, capturedComment)
	}
}

// ── GetByID ───────────────────────────────────────────────────────────────────

func TestApplicationService_GetByID_Found(t *testing.T) {
	ctx := context.Background()
	app := stubApplication(model.StatusReceived)
	mock := &mockApplicationRepo{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (*model.Application, error) { return app, nil },
	}
	svc := NewApplicationService(mock, nil)
	got, err := svc.GetByID(ctx, app.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got == nil || got.ID != app.ID {
		t.Errorf("want ID %v, got %v", app.ID, got)
	}
}

func TestApplicationService_GetByID_NotFound(t *testing.T) {
	ctx := context.Background()
	mock := &mockApplicationRepo{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (*model.Application, error) { return nil, nil },
	}
	svc := NewApplicationService(mock, nil)
	got, err := svc.GetByID(ctx, uuid.New())
	if err != nil || got != nil {
		t.Errorf("expected nil,nil; got %v,%v", got, err)
	}
}

// ── GetWithHistory ────────────────────────────────────────────────────────────

func TestApplicationService_GetWithHistory_NotFound(t *testing.T) {
	ctx := context.Background()
	mock := &mockApplicationRepo{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (*model.Application, error) { return nil, nil },
	}
	svc := NewApplicationService(mock, nil)
	app, hist, err := svc.GetWithHistory(ctx, uuid.New())
	if err != nil {
		t.Fatal(err)
	}
	if app != nil || hist != nil {
		t.Error("expected nil app and history for not-found ID")
	}
}

func TestApplicationService_GetWithHistory_Success(t *testing.T) {
	ctx := context.Background()
	app := stubApplication(model.StatusReceived)
	history := []*model.ApplicationHistory{
		{ID: uuid.New(), ApplicationID: app.ID, ToStatus: model.StatusReceived, ActorID: uuid.New(), CreatedAt: time.Now()},
	}
	mock := &mockApplicationRepo{
		getByIDFn:    func(_ context.Context, _ uuid.UUID) (*model.Application, error) { return app, nil },
		getHistoryFn: func(_ context.Context, _ uuid.UUID) ([]*model.ApplicationHistory, error) { return history, nil },
	}
	svc := NewApplicationService(mock, nil)
	gotApp, gotHist, err := svc.GetWithHistory(ctx, app.ID)
	if err != nil {
		t.Fatal(err)
	}
	if gotApp == nil || gotApp.ID != app.ID {
		t.Errorf("want app ID %v, got %v", app.ID, gotApp)
	}
	if len(gotHist) != 1 {
		t.Errorf("want 1 history record, got %d", len(gotHist))
	}
}

// ── ListForBorrower / ListAll ─────────────────────────────────────────────────

func TestApplicationService_ListForBorrower(t *testing.T) {
	ctx := context.Background()
	borrowerID := uuid.New()
	mock := &mockApplicationRepo{
		listByBorrowerFn: func(_ context.Context, id uuid.UUID) ([]*model.Application, error) {
			if id != borrowerID {
				return nil, nil
			}
			return []*model.Application{stubApplication(model.StatusReceived)}, nil
		},
	}
	svc := NewApplicationService(mock, nil)
	apps, err := svc.ListForBorrower(ctx, borrowerID)
	if err != nil {
		t.Fatal(err)
	}
	if len(apps) != 1 {
		t.Errorf("want 1 application, got %d", len(apps))
	}
}

func TestApplicationService_ListAll_WithStatus(t *testing.T) {
	ctx := context.Background()
	var capturedStatus string
	mock := &mockApplicationRepo{
		listAllFn: func(_ context.Context, status string) ([]*model.Application, error) {
			capturedStatus = status
			return nil, nil
		},
	}
	svc := NewApplicationService(mock, nil)
	svc.ListAll(ctx, model.StatusReceived) //nolint:errcheck
	if capturedStatus != model.StatusReceived {
		t.Errorf("want status %q, got %q", model.StatusReceived, capturedStatus)
	}
}
