package model

import (
	"time"

	"github.com/google/uuid"
)

// FSM statuses
const (
	StatusReceived            = "received"
	StatusPrimaryScoring      = "primary_scoring"
	StatusSecurityCheck       = "security_check"
	StatusCollateralExpertise = "collateral_expertise"
	StatusLegalCheck          = "legal_check"
	StatusCreditAnalysis      = "credit_analysis"
	StatusCreditCommittee     = "credit_committee"
	StatusApproved            = "approved"
	StatusRejected            = "rejected"
	StatusRevision            = "revision"
	StatusDocumentation       = "documentation"
	StatusIssued              = "issued"
)

// AllowedTransitions defines valid FSM transitions per current status.
var AllowedTransitions = map[string][]string{
	StatusReceived:            {StatusPrimaryScoring, StatusRejected},
	StatusPrimaryScoring:      {StatusSecurityCheck, StatusRejected, StatusRevision},
	StatusSecurityCheck:       {StatusCollateralExpertise, StatusRejected, StatusRevision},
	StatusCollateralExpertise: {StatusLegalCheck, StatusRejected, StatusRevision},
	StatusLegalCheck:          {StatusCreditAnalysis, StatusRejected, StatusRevision},
	StatusCreditAnalysis:      {StatusCreditCommittee, StatusRejected, StatusRevision},
	StatusCreditCommittee:     {StatusApproved, StatusRejected, StatusRevision},
	StatusApproved:            {StatusDocumentation},
	StatusRevision:            {StatusPrimaryScoring, StatusRejected},
	StatusDocumentation:       {StatusIssued},
	StatusRejected:            {},
	StatusIssued:              {},
}

// IsValidTransition returns true if moving from → to is allowed.
func IsValidTransition(from, to string) bool {
	allowed, ok := AllowedTransitions[from]
	if !ok {
		return false
	}
	for _, s := range allowed {
		if s == to {
			return true
		}
	}
	return false
}

type Application struct {
	ID          uuid.UUID  `json:"id"`
	BorrowerID  uuid.UUID  `json:"borrower_id"`
	ProgramID   uuid.UUID  `json:"program_id"`
	AssigneeID  *uuid.UUID `json:"assignee_id,omitempty"`
	Status      string     `json:"status"`
	Amount      float64    `json:"amount"`
	TermMonths  int        `json:"term_months"`
	PaymentType string     `json:"payment_type"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type ApplicationHistory struct {
	ID            uuid.UUID `json:"id"`
	ApplicationID uuid.UUID `json:"application_id"`
	FromStatus    *string   `json:"from_status,omitempty"`
	ToStatus      string    `json:"to_status"`
	ActorID       uuid.UUID `json:"actor_id"`
	Comment       *string   `json:"comment,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}
