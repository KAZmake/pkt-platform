// Package onec defines the 1С HTTP-service client interface.
package onec

import "time"

// Loan is the raw loan record as received from 1С.
type Loan struct {
	OneCID     string    `json:"id"`
	BorrowerID string    `json:"borrower_id"` // maps to our borrowers.id
	ProgramID  string    `json:"program_id"`
	Amount     float64   `json:"amount"`
	Rate       float64   `json:"rate"`
	TermMonths int       `json:"term_months"`
	IssuedAt   time.Time `json:"issued_at"`
	ExpiresAt  time.Time `json:"expires_at"`
	Status     string    `json:"status"` // active | overdue | closed
}

// ScheduleItem is one row of the payment schedule from 1С.
type ScheduleItem struct {
	LoanOneCID string     `json:"loan_id"`
	DueDate    time.Time  `json:"due_date"`
	Principal  float64    `json:"principal"`
	Interest   float64    `json:"interest"`
	Total      float64    `json:"total"`
	IsPaid     bool       `json:"is_paid"`
	PaidAt     *time.Time `json:"paid_at,omitempty"`
}

// DebtItem is one overdue debt record from 1С.
type DebtItem struct {
	LoanOneCID  string  `json:"loan_id"`
	Type        string  `json:"type"` // principal | interest | penalty
	Amount      float64 `json:"amount"`
	DaysOverdue int     `json:"days_overdue"`
}

// Client is the interface that both the real and mock 1С clients implement.
type Client interface {
	GetLoans() ([]Loan, error)
	GetSchedule(loanOneCID string) ([]ScheduleItem, error)
	GetDebts(loanOneCID string) ([]DebtItem, error)
}
