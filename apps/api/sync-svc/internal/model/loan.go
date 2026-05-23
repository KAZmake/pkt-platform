package model

import (
	"time"

	"github.com/google/uuid"
)

type Loan struct {
	ID         uuid.UUID  `json:"id"`
	OneCID     string     `json:"one_c_id"`
	BorrowerID uuid.UUID  `json:"borrower_id"`
	ProgramID  *uuid.UUID `json:"program_id,omitempty"`
	Amount     float64    `json:"amount"`
	Rate       float64    `json:"rate"`
	TermMonths int        `json:"term_months"`
	IssuedAt   time.Time  `json:"issued_at"`
	ExpiresAt  time.Time  `json:"expires_at"`
	Status     string     `json:"status"`
	SyncedAt   time.Time  `json:"synced_at"`
}

type ScheduleItem struct {
	ID        uuid.UUID  `json:"id"`
	LoanID    uuid.UUID  `json:"loan_id"`
	DueDate   time.Time  `json:"due_date"`
	Principal float64    `json:"principal"`
	Interest  float64    `json:"interest"`
	Total     float64    `json:"total"`
	IsPaid    bool       `json:"is_paid"`
	PaidAt    *time.Time `json:"paid_at,omitempty"`
	SyncedAt  time.Time  `json:"synced_at"`
}

type LoanDebt struct {
	ID          uuid.UUID `json:"id"`
	LoanID      uuid.UUID `json:"loan_id"`
	Type        string    `json:"type"`
	Amount      float64   `json:"amount"`
	DaysOverdue int       `json:"days_overdue"`
	SyncedAt    time.Time `json:"synced_at"`
}
