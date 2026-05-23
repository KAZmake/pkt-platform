package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/KAZmake/pkt-platform/apps/api/expertise-svc/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ConclusionRepository struct {
	db *pgxpool.Pool
}

func NewConclusionRepository(db *pgxpool.Pool) *ConclusionRepository {
	return &ConclusionRepository{db: db}
}

type CreateConclusionInput struct {
	ApplicationID  uuid.UUID
	ExpertID       uuid.UUID
	Stage          string
	Risks          *string
	ConclusionText *string
	Result         string
	FilePath       *string
}

func (r *ConclusionRepository) Create(ctx context.Context, inp CreateConclusionInput) (*model.ExpertConclusion, error) {
	c := &model.ExpertConclusion{}
	err := r.db.QueryRow(ctx, `
		INSERT INTO expert_conclusions
		  (application_id, expert_id, stage, risks, conclusion_text, result, file_path)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, application_id, expert_id, stage, risks, conclusion_text, result, file_path, created_at`,
		inp.ApplicationID, inp.ExpertID, inp.Stage, inp.Risks,
		inp.ConclusionText, inp.Result, inp.FilePath,
	).Scan(&c.ID, &c.ApplicationID, &c.ExpertID, &c.Stage, &c.Risks,
		&c.ConclusionText, &c.Result, &c.FilePath, &c.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("create conclusion: %w", err)
	}
	return c, nil
}

func (r *ConclusionRepository) ListByApplication(ctx context.Context, appID uuid.UUID) ([]*model.ExpertConclusion, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, application_id, expert_id, stage, risks, conclusion_text, result, file_path, created_at
		FROM expert_conclusions
		WHERE application_id = $1
		ORDER BY created_at ASC`, appID)
	if err != nil {
		return nil, fmt.Errorf("list conclusions: %w", err)
	}
	defer rows.Close()

	var list []*model.ExpertConclusion
	for rows.Next() {
		c := &model.ExpertConclusion{}
		if err := rows.Scan(&c.ID, &c.ApplicationID, &c.ExpertID, &c.Stage, &c.Risks,
			&c.ConclusionText, &c.Result, &c.FilePath, &c.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan conclusion: %w", err)
		}
		list = append(list, c)
	}
	return list, nil
}

type AddVoteInput struct {
	ApplicationID uuid.UUID
	ExpertID      uuid.UUID
	Vote          string
	Comment       *string
}

func (r *ConclusionRepository) AddVote(ctx context.Context, inp AddVoteInput) (*model.CommitteeVote, error) {
	v := &model.CommitteeVote{}
	err := r.db.QueryRow(ctx, `
		INSERT INTO committee_votes (application_id, expert_id, vote, comment)
		VALUES ($1, $2, $3, $4)
		RETURNING id, application_id, expert_id, vote, comment, signed_at`,
		inp.ApplicationID, inp.ExpertID, inp.Vote, inp.Comment,
	).Scan(&v.ID, &v.ApplicationID, &v.ExpertID, &v.Vote, &v.Comment, &v.SignedAt)
	if err != nil {
		return nil, fmt.Errorf("add vote: %w", err)
	}
	return v, nil
}

func (r *ConclusionRepository) GetVotes(ctx context.Context, appID uuid.UUID) ([]*model.CommitteeVote, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, application_id, expert_id, vote, comment, signed_at
		FROM committee_votes
		WHERE application_id = $1
		ORDER BY signed_at ASC`, appID)
	if err != nil {
		return nil, fmt.Errorf("get votes: %w", err)
	}
	defer rows.Close()

	var votes []*model.CommitteeVote
	for rows.Next() {
		v := &model.CommitteeVote{}
		if err := rows.Scan(&v.ID, &v.ApplicationID, &v.ExpertID, &v.Vote, &v.Comment, &v.SignedAt); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				break
			}
			return nil, fmt.Errorf("scan vote: %w", err)
		}
		votes = append(votes, v)
	}
	return votes, nil
}
