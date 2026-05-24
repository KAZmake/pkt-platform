package service

import (
	"context"
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/KAZmake/pkt-platform/apps/api/core-api/internal/model"
	"github.com/google/uuid"
)

// ── mocks ─────────────────────────────────────────────────────────────────────

type mockScheduleAppRepo struct {
	getByIDFn func(context.Context, uuid.UUID) (*model.Application, error)
}

func (m *mockScheduleAppRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.Application, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, nil
}

type mockScheduleProgramRepo struct {
	getByIDFn func(context.Context, string) (*model.LoanProgram, error)
}

func (m *mockScheduleProgramRepo) GetByID(ctx context.Context, id string) (*model.LoanProgram, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, nil
}

func stubScheduleApp(programID uuid.UUID, payType string, amount float64, months int) *model.Application {
	return &model.Application{
		ID:          uuid.New(),
		BorrowerID:  uuid.New(),
		ProgramID:   programID,
		Status:      model.StatusReceived,
		Amount:      amount,
		TermMonths:  months,
		PaymentType: payType,
		CreatedAt:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt:   time.Now(),
	}
}

func stubScheduleProgram(rate float64) *model.LoanProgram {
	return &model.LoanProgram{
		ID:            uuid.New().String(),
		Name:          "Test Program",
		Rate:          rate,
		MinAmount:     100_000,
		MaxAmount:     5_000_000,
		MinTermMonths: 6,
		MaxTermMonths: 60,
		IsActive:      true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

// ── round2 ─────────────────────────────────────────────────────────────────────

func TestRound2(t *testing.T) {
	cases := []struct {
		in   float64
		want float64
	}{
		{1.234, 1.23},
		{1.237, 1.24},
		{3.141, 3.14},
		{2.716, 2.72},
		{100.0, 100.0},
		{0.001, 0.0},
		{0.0, 0.0},
	}

	for _, tc := range cases {
		got := round2(tc.in)
		if math.Abs(got-tc.want) > 1e-10 {
			t.Errorf("round2(%v) = %v, want %v", tc.in, got, tc.want)
		}
	}
}

// ── paymentDate ───────────────────────────────────────────────────────────────

func TestPaymentDate(t *testing.T) {
	start := time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC)

	cases := []struct {
		offset int
		want   string
	}{
		{1, "2025-02-01"},
		{2, "2025-03-01"},
		{12, "2026-01-01"},
		{24, "2027-01-01"},
	}

	for _, tc := range cases {
		got := paymentDate(start, tc.offset)
		if got != tc.want {
			t.Errorf("paymentDate(start, %d) = %q, want %q", tc.offset, got, tc.want)
		}
	}
}

func TestPaymentDate_Format(t *testing.T) {
	start := time.Date(2025, 11, 1, 0, 0, 0, 0, time.UTC)
	got := paymentDate(start, 3)
	expected := fmt.Sprintf("%04d-%02d-01", 2026, 2)
	if got != expected {
		t.Errorf("want %q, got %q", expected, got)
	}
}

// ── Calculate ─────────────────────────────────────────────────────────────────

func TestScheduleService_Calculate_AppNotFound(t *testing.T) {
	ctx := context.Background()
	svc := NewScheduleService(
		&mockScheduleAppRepo{getByIDFn: func(_ context.Context, _ uuid.UUID) (*model.Application, error) { return nil, nil }},
		&mockScheduleProgramRepo{},
	)
	result, err := svc.Calculate(ctx, uuid.New())
	if err != nil || result != nil {
		t.Errorf("want nil,nil for not-found app; got %v,%v", result, err)
	}
}

func TestScheduleService_Calculate_ProgramNotFound(t *testing.T) {
	ctx := context.Background()
	app := stubScheduleApp(uuid.New(), "annuity", 500_000, 12)
	svc := NewScheduleService(
		&mockScheduleAppRepo{getByIDFn: func(_ context.Context, _ uuid.UUID) (*model.Application, error) { return app, nil }},
		&mockScheduleProgramRepo{getByIDFn: func(_ context.Context, _ string) (*model.LoanProgram, error) { return nil, nil }},
	)
	_, err := svc.Calculate(ctx, app.ID)
	if err == nil {
		t.Fatal("expected error when program not found")
	}
}

func TestScheduleService_Calculate_Annuity(t *testing.T) {
	ctx := context.Background()
	prog := stubScheduleProgram(12.0) // 12% annual
	app := stubScheduleApp(uuid.MustParse(prog.ID), "annuity", 1_200_000, 12)

	svc := NewScheduleService(
		&mockScheduleAppRepo{getByIDFn: func(_ context.Context, _ uuid.UUID) (*model.Application, error) { return app, nil }},
		&mockScheduleProgramRepo{getByIDFn: func(_ context.Context, _ string) (*model.LoanProgram, error) { return prog, nil }},
	)
	result, err := svc.Calculate(ctx, app.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("expected result, got nil")
	}
	if result.TotalPayments != 12 {
		t.Errorf("want 12 payments, got %d", result.TotalPayments)
	}
	if len(result.Schedule) != 12 {
		t.Errorf("want 12 rows, got %d", len(result.Schedule))
	}
	if result.PaymentType != "annuity" {
		t.Errorf("want annuity, got %q", result.PaymentType)
	}
	// Final balance should be ~0
	last := result.Schedule[len(result.Schedule)-1]
	if math.Abs(last.Balance) > 1 {
		t.Errorf("final balance should be ~0, got %v", last.Balance)
	}
}

func TestScheduleService_Calculate_Differentiated(t *testing.T) {
	ctx := context.Background()
	prog := stubScheduleProgram(7.5)
	app := stubScheduleApp(uuid.MustParse(prog.ID), "differentiated", 600_000, 6)

	svc := NewScheduleService(
		&mockScheduleAppRepo{getByIDFn: func(_ context.Context, _ uuid.UUID) (*model.Application, error) { return app, nil }},
		&mockScheduleProgramRepo{getByIDFn: func(_ context.Context, _ string) (*model.LoanProgram, error) { return prog, nil }},
	)
	result, err := svc.Calculate(ctx, app.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.TotalPayments != 6 {
		t.Errorf("want 6 payments, got %d", result.TotalPayments)
	}
	if len(result.Schedule) != 6 {
		t.Errorf("want 6 rows, got %d", len(result.Schedule))
	}
	// Principal should decrease each month in differentiated schedule
	for i := 0; i < len(result.Schedule)-1; i++ {
		if result.Schedule[i].Balance < result.Schedule[i+1].Balance {
			t.Errorf("balance should decrease: row %d balance %v >= row %d balance %v",
				i, result.Schedule[i].Balance, i+1, result.Schedule[i+1].Balance)
		}
	}
}

func TestScheduleService_Calculate_ZeroRate_Annuity(t *testing.T) {
	ctx := context.Background()
	prog := stubScheduleProgram(0.0) // 0% rate — annuity = equal principal payments
	app := stubScheduleApp(uuid.MustParse(prog.ID), "annuity", 1_200_000, 12)

	svc := NewScheduleService(
		&mockScheduleAppRepo{getByIDFn: func(_ context.Context, _ uuid.UUID) (*model.Application, error) { return app, nil }},
		&mockScheduleProgramRepo{getByIDFn: func(_ context.Context, _ string) (*model.LoanProgram, error) { return prog, nil }},
	)
	result, err := svc.Calculate(ctx, app.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil || len(result.Schedule) != 12 {
		t.Fatal("expected 12 schedule rows")
	}
	if result.TotalInterest != 0 {
		t.Errorf("expected 0 total interest at 0%% rate, got %v", result.TotalInterest)
	}
}

func TestScheduleService_Calculate_InvalidPaymentType(t *testing.T) {
	ctx := context.Background()
	prog := stubScheduleProgram(7.5)
	app := stubScheduleApp(uuid.MustParse(prog.ID), "bullet", 500_000, 12)

	svc := NewScheduleService(
		&mockScheduleAppRepo{getByIDFn: func(_ context.Context, _ uuid.UUID) (*model.Application, error) { return app, nil }},
		&mockScheduleProgramRepo{getByIDFn: func(_ context.Context, _ string) (*model.LoanProgram, error) { return prog, nil }},
	)
	_, err := svc.Calculate(ctx, app.ID)
	if err == nil {
		t.Fatal("expected error for invalid payment type")
	}
}

func TestScheduleService_Calculate_MonthlyPaymentConsistency_Annuity(t *testing.T) {
	ctx := context.Background()
	prog := stubScheduleProgram(12.0)
	app := stubScheduleApp(uuid.MustParse(prog.ID), "annuity", 1_000_000, 12)

	svc := NewScheduleService(
		&mockScheduleAppRepo{getByIDFn: func(_ context.Context, _ uuid.UUID) (*model.Application, error) { return app, nil }},
		&mockScheduleProgramRepo{getByIDFn: func(_ context.Context, _ string) (*model.LoanProgram, error) { return prog, nil }},
	)
	result, err := svc.Calculate(ctx, app.ID)
	if err != nil {
		t.Fatal(err)
	}
	// All non-last payments should be equal (annuity property)
	first := result.Schedule[0].Payment
	for i := 1; i < len(result.Schedule)-1; i++ {
		if math.Abs(result.Schedule[i].Payment-first) > 0.01 {
			t.Errorf("annuity payment inconsistency at row %d: first=%v, got=%v", i, first, result.Schedule[i].Payment)
		}
	}
}
