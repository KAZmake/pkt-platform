package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/KAZmake/pkt-platform/apps/api/core-api/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProgramRepository struct {
	db *pgxpool.Pool
}

func NewProgramRepository(db *pgxpool.Pool) *ProgramRepository {
	return &ProgramRepository{db: db}
}

const programCols = `id, name, name_kz, name_en, rate,
	min_amount, max_amount, min_term_months, max_term_months,
	activity_types, is_active, created_at, updated_at`

func (r *ProgramRepository) List(ctx context.Context, activeOnly bool) ([]*model.LoanProgram, error) {
	q := `SELECT ` + programCols + ` FROM loan_programs`
	if activeOnly {
		q += ` WHERE is_active = TRUE`
	}
	q += ` ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("programs list: %w", err)
	}
	defer rows.Close()

	var programs []*model.LoanProgram
	for rows.Next() {
		p, err := scanProgram(rows)
		if err != nil {
			return nil, err
		}
		programs = append(programs, p)
	}
	return programs, nil
}

func (r *ProgramRepository) GetByID(ctx context.Context, id string) (*model.LoanProgram, error) {
	row := r.db.QueryRow(ctx,
		`SELECT `+programCols+` FROM loan_programs WHERE id = $1`, id)
	return scanProgram(row)
}

type CreateProgramInput struct {
	Name          string   `json:"name"`
	NameKZ        *string  `json:"name_kz"`
	NameEN        *string  `json:"name_en"`
	Rate          float64  `json:"rate"`
	MinAmount     float64  `json:"min_amount"`
	MaxAmount     float64  `json:"max_amount"`
	MinTermMonths int      `json:"min_term_months"`
	MaxTermMonths int      `json:"max_term_months"`
	ActivityTypes []string `json:"activity_types"`
}

func (r *ProgramRepository) Create(ctx context.Context, inp CreateProgramInput) (*model.LoanProgram, error) {
	if inp.ActivityTypes == nil {
		inp.ActivityTypes = []string{}
	}
	row := r.db.QueryRow(ctx, `
		INSERT INTO loan_programs
		  (name, name_kz, name_en, rate, min_amount, max_amount,
		   min_term_months, max_term_months, activity_types)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		RETURNING `+programCols,
		inp.Name, inp.NameKZ, inp.NameEN, inp.Rate,
		inp.MinAmount, inp.MaxAmount,
		inp.MinTermMonths, inp.MaxTermMonths,
		inp.ActivityTypes,
	)
	return scanProgram(row)
}

type UpdateProgramInput struct {
	Name          *string  `json:"name"`
	NameKZ        *string  `json:"name_kz"`
	NameEN        *string  `json:"name_en"`
	Rate          *float64 `json:"rate"`
	MinAmount     *float64 `json:"min_amount"`
	MaxAmount     *float64 `json:"max_amount"`
	MinTermMonths *int     `json:"min_term_months"`
	MaxTermMonths *int     `json:"max_term_months"`
	ActivityTypes []string `json:"activity_types"`
	IsActive      *bool    `json:"is_active"`
}

func (r *ProgramRepository) Update(ctx context.Context, id string, inp UpdateProgramInput) (*model.LoanProgram, error) {
	row := r.db.QueryRow(ctx, `
		UPDATE loan_programs SET
		  name           = COALESCE($2,  name),
		  name_kz        = COALESCE($3,  name_kz),
		  name_en        = COALESCE($4,  name_en),
		  rate           = COALESCE($5,  rate),
		  min_amount     = COALESCE($6,  min_amount),
		  max_amount     = COALESCE($7,  max_amount),
		  min_term_months= COALESCE($8,  min_term_months),
		  max_term_months= COALESCE($9,  max_term_months),
		  activity_types = COALESCE($10, activity_types),
		  is_active      = COALESCE($11, is_active),
		  updated_at     = NOW()
		WHERE id = $1
		RETURNING `+programCols,
		id,
		inp.Name, inp.NameKZ, inp.NameEN, inp.Rate,
		inp.MinAmount, inp.MaxAmount,
		inp.MinTermMonths, inp.MaxTermMonths,
		inp.ActivityTypes, inp.IsActive,
	)
	return scanProgram(row)
}

func (r *ProgramRepository) SetActive(ctx context.Context, id string, active bool) error {
	_, err := r.db.Exec(ctx,
		`UPDATE loan_programs SET is_active = $2, updated_at = NOW() WHERE id = $1`,
		id, active)
	return err
}

// ── helpers ──────────────────────────────────────────────────────────────────

func scanProgram(s scanner) (*model.LoanProgram, error) {
	p := &model.LoanProgram{}
	err := s.Scan(
		&p.ID, &p.Name, &p.NameKZ, &p.NameEN, &p.Rate,
		&p.MinAmount, &p.MaxAmount, &p.MinTermMonths, &p.MaxTermMonths,
		&p.ActivityTypes, &p.IsActive, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("scan program: %w", err)
	}
	return p, nil
}
