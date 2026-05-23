package model

import (
	"time"

	"github.com/google/uuid"
)

type Ticket struct {
	ID         uuid.UUID  `json:"id"`
	BorrowerID uuid.UUID  `json:"borrower_id"`
	AssigneeID *uuid.UUID `json:"assignee_id,omitempty"`
	Type       string     `json:"type"` // early_repayment | restructuring | prolongation | other
	Subject    string     `json:"subject"`
	Status     string     `json:"status"` // open | in_progress | resolved | closed
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

type TicketMessage struct {
	ID             uuid.UUID `json:"id"`
	TicketID       uuid.UUID `json:"ticket_id"`
	AuthorID       uuid.UUID `json:"author_id"`
	Body           string    `json:"body"`
	AttachmentPath *string   `json:"attachment_path,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}
