package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/KAZmake/pkt-platform/apps/api/sync-svc/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LoanRepository struct {
	db *pgxpool.Pool
}

func NewLoanRepository(db *pgxpool.Pool) *LoanRepository {
	return &LoanRepository{db: db}
}

// UpsertLoan inserts or updates a loan record from 1С (2.3.1).
func (r *LoanRepository) UpsertLoan(ctx context.Context, l model.Loan) (*model.Loan, error) {
	var programID *string
	if l.ProgramID != nil {
		s := l.ProgramID.String()
		programID = &s
	}

	out := &model.Loan{}
	err := r.db.QueryRow(ctx, `
		INSERT INTO loans (one_c_id, borrower_id, program_id, amount, rate, term_months,
		                   issued_at, expires_at, status, synced_at)
		VALUES ($1, $2, $3::uuid, $4, $5, $6, $7, $8, $9, NOW())
		ON CONFLICT (one_c_id) DO UPDATE SET
			amount      = EXCLUDED.amount,
			rate        = EXCLUDED.rate,
			status      = EXCLUDED.status,
			synced_at   = NOW()
		RETURNING id, one_c_id, borrower_id, program_id, amount, rate, term_months,
		          issued_at, expires_at, status, synced_at`,
		l.OneCID, l.BorrowerID, programID, l.Amount, l.Rate, l.TermMonths,
		l.IssuedAt, l.ExpiresAt, l.Status,
	).Scan(&out.ID, &out.OneCID, &out.BorrowerID, &out.ProgramID, &out.Amount, &out.Rate,
		&out.TermMonths, &out.IssuedAt, &out.ExpiresAt, &out.Status, &out.SyncedAt)
	if err != nil {
		return nil, fmt.Errorf("upsert loan: %w", err)
	}
	return out, nil
}

// ReplaceSchedule deletes existing schedule rows and inserts fresh ones (2.3.1).
func (r *LoanRepository) ReplaceSchedule(ctx context.Context, loanID uuid.UUID, items []model.ScheduleItem) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	if _, err := tx.Exec(ctx, `DELETE FROM payment_schedule WHERE loan_id = $1`, loanID); err != nil {
		return fmt.Errorf("delete schedule: %w", err)
	}

	for _, item := range items {
		_, err := tx.Exec(ctx, `
			INSERT INTO payment_schedule (loan_id, due_date, principal, interest, total, is_paid, paid_at, synced_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())`,
			loanID, item.DueDate, item.Principal, item.Interest, item.Total, item.IsPaid, item.PaidAt)
		if err != nil {
			return fmt.Errorf("insert schedule row: %w", err)
		}
	}
	return tx.Commit(ctx)
}

// ReplaceDebts replaces debt records for a loan (2.3.1).
func (r *LoanRepository) ReplaceDebts(ctx context.Context, loanID uuid.UUID, debts []model.LoanDebt) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	if _, err := tx.Exec(ctx, `DELETE FROM loan_debts WHERE loan_id = $1`, loanID); err != nil {
		return fmt.Errorf("delete debts: %w", err)
	}

	for _, d := range debts {
		_, err := tx.Exec(ctx, `
			INSERT INTO loan_debts (loan_id, type, amount, days_overdue, synced_at)
			VALUES ($1, $2, $3, $4, NOW())`,
			loanID, d.Type, d.Amount, d.DaysOverdue)
		if err != nil {
			return fmt.Errorf("insert debt: %w", err)
		}
	}
	return tx.Commit(ctx)
}

// GetByBorrower returns all loans for a borrower (2.3.2).
func (r *LoanRepository) GetByBorrower(ctx context.Context, borrowerID uuid.UUID) ([]*model.Loan, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, one_c_id, borrower_id, program_id, amount, rate, term_months,
		       issued_at, expires_at, status, synced_at
		FROM loans WHERE borrower_id = $1 ORDER BY issued_at DESC`, borrowerID)
	if err != nil {
		return nil, fmt.Errorf("get loans by borrower: %w", err)
	}
	defer rows.Close()
	return scanLoans(rows)
}

// GetAll returns all loans (2.3.2).
func (r *LoanRepository) GetAll(ctx context.Context, status string) ([]*model.Loan, error) {
	q := `SELECT id, one_c_id, borrower_id, program_id, amount, rate, term_months,
	             issued_at, expires_at, status, synced_at FROM loans`
	if status != "" {
		q += ` WHERE status = $1 ORDER BY issued_at DESC`
		rows, err := r.db.Query(ctx, q, status)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		return scanLoans(rows)
	}
	q += ` ORDER BY issued_at DESC`
	rows, err := r.db.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanLoans(rows)
}

// GetByID returns a single loan (2.3.2).
func (r *LoanRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Loan, error) {
	out := &model.Loan{}
	err := r.db.QueryRow(ctx, `
		SELECT id, one_c_id, borrower_id, program_id, amount, rate, term_months,
		       issued_at, expires_at, status, synced_at
		FROM loans WHERE id = $1`, id,
	).Scan(&out.ID, &out.OneCID, &out.BorrowerID, &out.ProgramID, &out.Amount, &out.Rate,
		&out.TermMonths, &out.IssuedAt, &out.ExpiresAt, &out.Status, &out.SyncedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get loan: %w", err)
	}
	return out, nil
}

// GetSchedule returns the payment schedule for a loan (2.3.2).
func (r *LoanRepository) GetSchedule(ctx context.Context, loanID uuid.UUID) ([]*model.ScheduleItem, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, loan_id, due_date, principal, interest, total, is_paid, paid_at, synced_at
		FROM payment_schedule WHERE loan_id = $1 ORDER BY due_date ASC`, loanID)
	if err != nil {
		return nil, fmt.Errorf("get schedule: %w", err)
	}
	defer rows.Close()

	var items []*model.ScheduleItem
	for rows.Next() {
		s := &model.ScheduleItem{}
		if err := rows.Scan(&s.ID, &s.LoanID, &s.DueDate, &s.Principal, &s.Interest,
			&s.Total, &s.IsPaid, &s.PaidAt, &s.SyncedAt); err != nil {
			return nil, fmt.Errorf("scan schedule: %w", err)
		}
		items = append(items, s)
	}
	return items, nil
}

// GetDebts returns debt records for a loan (2.3.2).
func (r *LoanRepository) GetDebts(ctx context.Context, loanID uuid.UUID) ([]*model.LoanDebt, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, loan_id, type, amount, days_overdue, synced_at
		FROM loan_debts WHERE loan_id = $1 ORDER BY days_overdue DESC`, loanID)
	if err != nil {
		return nil, fmt.Errorf("get debts: %w", err)
	}
	defer rows.Close()

	var debts []*model.LoanDebt
	for rows.Next() {
		d := &model.LoanDebt{}
		if err := rows.Scan(&d.ID, &d.LoanID, &d.Type, &d.Amount, &d.DaysOverdue, &d.SyncedAt); err != nil {
			return nil, fmt.Errorf("scan debt: %w", err)
		}
		debts = append(debts, d)
	}
	return debts, nil
}

func scanLoans(rows pgx.Rows) ([]*model.Loan, error) {
	var loans []*model.Loan
	for rows.Next() {
		l := &model.Loan{}
		if err := rows.Scan(&l.ID, &l.OneCID, &l.BorrowerID, &l.ProgramID, &l.Amount, &l.Rate,
			&l.TermMonths, &l.IssuedAt, &l.ExpiresAt, &l.Status, &l.SyncedAt); err != nil {
			return nil, fmt.Errorf("scan loan: %w", err)
		}
		loans = append(loans, l)
	}
	return loans, nil
}
