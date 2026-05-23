package onec

import (
	"log/slog"
	"time"
)

// MockClient returns seeded loan data for development (no real 1С required).
type MockClient struct{}

func NewMockClient() *MockClient {
	slog.Warn("sync-svc: using 1С mock client — set ONEC_BASE_URL for production")
	return &MockClient{}
}

func (m *MockClient) GetLoans() ([]Loan, error) {
	now := time.Now()
	return []Loan{
		{
			OneCID:     "1C-LOAN-001",
			BorrowerID: "b1000000-0000-0000-0000-000000000001",
			ProgramID:  "a1000000-0000-0000-0000-000000000001",
			Amount:     5000000,
			Rate:       7.5,
			TermMonths: 12,
			IssuedAt:   now.AddDate(-1, 0, 0),
			ExpiresAt:  now.AddDate(0, 2, 0),
			Status:     "active",
		},
		{
			OneCID:     "1C-LOAN-002",
			BorrowerID: "b1000000-0000-0000-0000-000000000002",
			ProgramID:  "a1000000-0000-0000-0000-000000000001",
			Amount:     8000000,
			Rate:       8.0,
			TermMonths: 18,
			IssuedAt:   now.AddDate(-2, 0, 0),
			ExpiresAt:  now.AddDate(-1, -6, 0),
			Status:     "closed",
		},
		{
			OneCID:     "1C-LOAN-003",
			BorrowerID: "b1000000-0000-0000-0000-000000000003",
			ProgramID:  "a1000000-0000-0000-0000-000000000002",
			Amount:     20000000,
			Rate:       7.0,
			TermMonths: 36,
			IssuedAt:   now.AddDate(0, -6, 0),
			ExpiresAt:  now.AddDate(2, 6, 0),
			Status:     "active",
		},
	}, nil
}

func (m *MockClient) GetSchedule(loanOneCID string) ([]ScheduleItem, error) {
	now := time.Now()
	items := make([]ScheduleItem, 0, 6)
	for i := 1; i <= 6; i++ {
		due := now.AddDate(0, i, 0)
		items = append(items, ScheduleItem{
			LoanOneCID: loanOneCID,
			DueDate:    due,
			Principal:  41666.67,
			Interest:   3125.0,
			Total:      44791.67,
			IsPaid:     i <= 2,
		})
	}
	return items, nil
}

func (m *MockClient) GetDebts(loanOneCID string) ([]DebtItem, error) {
	// Only loan 002 has debts (it's overdue/closed)
	if loanOneCID == "1C-LOAN-002" {
		return []DebtItem{
			{LoanOneCID: loanOneCID, Type: "principal", Amount: 150000, DaysOverdue: 45},
			{LoanOneCID: loanOneCID, Type: "interest", Amount: 12500, DaysOverdue: 45},
		}, nil
	}
	return nil, nil
}
