package service

import (
	"context"
	"fmt"

	"github.com/KAZmake/pkt-platform/apps/api/expertise-svc/internal/model"
	"github.com/KAZmake/pkt-platform/apps/api/expertise-svc/internal/repository"
	"github.com/google/uuid"
)

var validStages = map[string]bool{
	model.StageCollateralExpertise: true,
	model.StageLegalCheck:          true,
	model.StageCreditAnalysis:      true,
}

var validResults = map[string]bool{
	model.ResultApproved: true,
	model.ResultRejected: true,
	model.ResultRevision: true,
}

var validVotes = map[string]bool{
	model.VoteApproved:  true,
	model.VoteRejected:  true,
	model.VoteAbstained: true,
}

type ConclusionService struct {
	repo *repository.ConclusionRepository
}

func NewConclusionService(repo *repository.ConclusionRepository) *ConclusionService {
	return &ConclusionService{repo: repo}
}

type SubmitConclusionInput struct {
	Stage          string  `json:"stage"`
	Risks          *string `json:"risks"`
	ConclusionText *string `json:"conclusion_text"`
	Result         string  `json:"result"`
	FilePath       *string `json:"file_path"`
}

func (s *ConclusionService) Submit(ctx context.Context, appID, expertID uuid.UUID, inp SubmitConclusionInput) (*model.ExpertConclusion, error) {
	if !validStages[inp.Stage] {
		return nil, fmt.Errorf("invalid stage: must be collateral_expertise, legal_check, or credit_analysis")
	}
	if !validResults[inp.Result] {
		return nil, fmt.Errorf("invalid result: must be approved, rejected, or revision")
	}
	return s.repo.Create(ctx, repository.CreateConclusionInput{
		ApplicationID:  appID,
		ExpertID:       expertID,
		Stage:          inp.Stage,
		Risks:          inp.Risks,
		ConclusionText: inp.ConclusionText,
		Result:         inp.Result,
		FilePath:       inp.FilePath,
	})
}

func (s *ConclusionService) ListByApplication(ctx context.Context, appID uuid.UUID) ([]*model.ExpertConclusion, error) {
	return s.repo.ListByApplication(ctx, appID)
}

type AddVoteInput struct {
	Vote    string  `json:"vote"`
	Comment *string `json:"comment"`
}

func (s *ConclusionService) AddVote(ctx context.Context, appID, expertID uuid.UUID, inp AddVoteInput) (*model.CommitteeVote, error) {
	if !validVotes[inp.Vote] {
		return nil, fmt.Errorf("invalid vote: must be approved, rejected, or abstained")
	}
	return s.repo.AddVote(ctx, repository.AddVoteInput{
		ApplicationID: appID,
		ExpertID:      expertID,
		Vote:          inp.Vote,
		Comment:       inp.Comment,
	})
}

func (s *ConclusionService) GetVotes(ctx context.Context, appID uuid.UUID) ([]*model.CommitteeVote, error) {
	return s.repo.GetVotes(ctx, appID)
}
