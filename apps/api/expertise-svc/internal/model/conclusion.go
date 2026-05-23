package model

import (
	"time"

	"github.com/google/uuid"
)

// Expert conclusion stages
const (
	StageCollateralExpertise = "collateral_expertise"
	StageLegalCheck          = "legal_check"
	StageCreditAnalysis      = "credit_analysis"
)

// Conclusion results
const (
	ResultApproved = "approved"
	ResultRejected = "rejected"
	ResultRevision = "revision"
)

// Vote values
const (
	VoteApproved  = "approved"
	VoteRejected  = "rejected"
	VoteAbstained = "abstained"
)

type ExpertConclusion struct {
	ID             uuid.UUID `json:"id"`
	ApplicationID  uuid.UUID `json:"application_id"`
	ExpertID       uuid.UUID `json:"expert_id"`
	Stage          string    `json:"stage"`
	Risks          *string   `json:"risks,omitempty"` // JSONB stored as string
	ConclusionText *string   `json:"conclusion_text,omitempty"`
	Result         string    `json:"result"`
	FilePath       *string   `json:"file_path,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}

type CommitteeVote struct {
	ID            uuid.UUID `json:"id"`
	ApplicationID uuid.UUID `json:"application_id"`
	ExpertID      uuid.UUID `json:"expert_id"`
	Vote          string    `json:"vote"`
	Comment       *string   `json:"comment,omitempty"`
	SignedAt      time.Time `json:"signed_at"`
}
