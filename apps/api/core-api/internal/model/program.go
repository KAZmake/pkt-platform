package model

import "time"

type LoanProgram struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	NameKZ        *string   `json:"name_kz,omitempty"`
	NameEN        *string   `json:"name_en,omitempty"`
	Rate          float64   `json:"rate"`
	MinAmount     float64   `json:"min_amount"`
	MaxAmount     float64   `json:"max_amount"`
	MinTermMonths int       `json:"min_term_months"`
	MaxTermMonths int       `json:"max_term_months"`
	ActivityTypes []string  `json:"activity_types"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
