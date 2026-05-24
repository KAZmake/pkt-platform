package service

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/KAZmake/pkt-platform/apps/api/core-api/internal/model"
	"github.com/google/uuid"
)

type ScheduleRow struct {
	Month     int     `json:"month"`
	Date      string  `json:"date"` // YYYY-MM-DD, 1st of each month
	Principal float64 `json:"principal"`
	Interest  float64 `json:"interest"`
	Payment   float64 `json:"payment"`
	Balance   float64 `json:"balance"`
}

type ScheduleResult struct {
	ApplicationID uuid.UUID     `json:"application_id"`
	PaymentType   string        `json:"payment_type"`
	TotalPayments int           `json:"total_payments"`
	TotalInterest float64       `json:"total_interest"`
	Schedule      []ScheduleRow `json:"schedule"`
}

type scheduleAppRepo interface {
	GetByID(ctx context.Context, id uuid.UUID) (*model.Application, error)
}

type scheduleProgramRepo interface {
	GetByID(ctx context.Context, id string) (*model.LoanProgram, error)
}

type ScheduleService struct {
	appRepo     scheduleAppRepo
	programRepo scheduleProgramRepo
}

func NewScheduleService(appRepo scheduleAppRepo, programRepo scheduleProgramRepo) *ScheduleService {
	return &ScheduleService{appRepo: appRepo, programRepo: programRepo}
}

// Calculate builds the payment schedule for an application.
// Annual rate is stored in loan_programs.rate (e.g. 7.0 = 7%).
func (s *ScheduleService) Calculate(ctx context.Context, appID uuid.UUID) (*ScheduleResult, error) {
	app, err := s.appRepo.GetByID(ctx, appID)
	if err != nil {
		return nil, err
	}
	if app == nil {
		return nil, nil
	}

	prog, err := s.programRepo.GetByID(ctx, app.ProgramID.String())
	if err != nil || prog == nil {
		return nil, fmt.Errorf("loan program not found")
	}

	monthlyRate := prog.Rate / 100 / 12
	n := app.TermMonths
	principal := app.Amount

	var rows []ScheduleRow
	balance := principal
	totalInterest := 0.0
	startDate := app.CreatedAt

	switch app.PaymentType {
	case "annuity":
		// M = P * r * (1+r)^n / ((1+r)^n - 1)
		var monthlyPayment float64
		if monthlyRate == 0 {
			monthlyPayment = principal / float64(n)
		} else {
			factor := math.Pow(1+monthlyRate, float64(n))
			monthlyPayment = principal * monthlyRate * factor / (factor - 1)
		}
		monthlyPayment = round2(monthlyPayment)

		for i := 1; i <= n; i++ {
			interest := round2(balance * monthlyRate)
			pmt := monthlyPayment
			if i == n {
				// last payment absorbs rounding residual
				pmt = round2(balance + interest)
			}
			pmtPrincipal := round2(pmt - interest)
			balance = round2(balance - pmtPrincipal)
			totalInterest += interest

			rows = append(rows, ScheduleRow{
				Month:     i,
				Date:      paymentDate(startDate, i),
				Principal: pmtPrincipal,
				Interest:  interest,
				Payment:   pmt,
				Balance:   math.Max(balance, 0),
			})
		}

	case "differentiated":
		pmtPrincipal := round2(principal / float64(n))
		for i := 1; i <= n; i++ {
			interest := round2(balance * monthlyRate)
			pmt := round2(pmtPrincipal + interest)
			if i == n {
				pmtPrincipal = balance
				interest = round2(balance * monthlyRate)
				pmt = round2(pmtPrincipal + interest)
			}
			balance = round2(balance - pmtPrincipal)
			totalInterest += interest

			rows = append(rows, ScheduleRow{
				Month:     i,
				Date:      paymentDate(startDate, i),
				Principal: pmtPrincipal,
				Interest:  interest,
				Payment:   pmt,
				Balance:   math.Max(balance, 0),
			})
		}

	default:
		return nil, fmt.Errorf("unknown payment_type: %s", app.PaymentType)
	}

	return &ScheduleResult{
		ApplicationID: appID,
		PaymentType:   app.PaymentType,
		TotalPayments: n,
		TotalInterest: round2(totalInterest),
		Schedule:      rows,
	}, nil
}

func paymentDate(start time.Time, monthOffset int) string {
	d := start.AddDate(0, monthOffset, 0)
	return fmt.Sprintf("%04d-%02d-01", d.Year(), d.Month())
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}
