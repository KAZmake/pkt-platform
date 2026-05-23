package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/KAZmake/pkt-platform/apps/api/core-api/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TicketRepository struct {
	db *pgxpool.Pool
}

func NewTicketRepository(db *pgxpool.Pool) *TicketRepository {
	return &TicketRepository{db: db}
}

const ticketCols = `id, borrower_id, assignee_id, type, subject, status, created_at, updated_at`

type CreateTicketInput struct {
	BorrowerID uuid.UUID
	Type       string
	Subject    string
}

func (r *TicketRepository) Create(ctx context.Context, inp CreateTicketInput) (*model.Ticket, error) {
	row := r.db.QueryRow(ctx, `
		INSERT INTO tickets (borrower_id, type, subject)
		VALUES ($1, $2, $3)
		RETURNING `+ticketCols,
		inp.BorrowerID, inp.Type, inp.Subject,
	)
	return scanTicket(row)
}

func (r *TicketRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Ticket, error) {
	row := r.db.QueryRow(ctx, `SELECT `+ticketCols+` FROM tickets WHERE id = $1`, id)
	return scanTicket(row)
}

func (r *TicketRepository) ListByBorrower(ctx context.Context, borrowerID uuid.UUID) ([]*model.Ticket, error) {
	return r.queryList(ctx,
		`SELECT `+ticketCols+` FROM tickets WHERE borrower_id = $1 ORDER BY created_at DESC`,
		borrowerID)
}

func (r *TicketRepository) ListAll(ctx context.Context, status string) ([]*model.Ticket, error) {
	if status != "" {
		return r.queryList(ctx,
			`SELECT `+ticketCols+` FROM tickets WHERE status = $1 ORDER BY created_at DESC`,
			status)
	}
	return r.queryList(ctx, `SELECT `+ticketCols+` FROM tickets ORDER BY created_at DESC`)
}

func (r *TicketRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) (*model.Ticket, error) {
	row := r.db.QueryRow(ctx, `
		UPDATE tickets SET status = $2, updated_at = NOW()
		WHERE id = $1
		RETURNING `+ticketCols, id, status)
	return scanTicket(row)
}

func (r *TicketRepository) AddMessage(ctx context.Context, ticketID, authorID uuid.UUID, body string) (*model.TicketMessage, error) {
	m := &model.TicketMessage{}
	err := r.db.QueryRow(ctx, `
		INSERT INTO ticket_messages (ticket_id, author_id, body)
		VALUES ($1, $2, $3)
		RETURNING id, ticket_id, author_id, body, attachment_path, created_at`,
		ticketID, authorID, body,
	).Scan(&m.ID, &m.TicketID, &m.AuthorID, &m.Body, &m.AttachmentPath, &m.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("add message: %w", err)
	}
	return m, nil
}

func (r *TicketRepository) GetMessages(ctx context.Context, ticketID uuid.UUID) ([]*model.TicketMessage, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, ticket_id, author_id, body, attachment_path, created_at
		FROM ticket_messages
		WHERE ticket_id = $1
		ORDER BY created_at ASC`, ticketID)
	if err != nil {
		return nil, fmt.Errorf("get messages: %w", err)
	}
	defer rows.Close()

	var msgs []*model.TicketMessage
	for rows.Next() {
		m := &model.TicketMessage{}
		if err := rows.Scan(&m.ID, &m.TicketID, &m.AuthorID, &m.Body, &m.AttachmentPath, &m.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan message: %w", err)
		}
		msgs = append(msgs, m)
	}
	return msgs, nil
}

func (r *TicketRepository) queryList(ctx context.Context, q string, args ...any) ([]*model.Ticket, error) {
	rows, err := r.db.Query(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("tickets query: %w", err)
	}
	defer rows.Close()

	var tickets []*model.Ticket
	for rows.Next() {
		t, err := scanTicket(rows)
		if err != nil {
			return nil, err
		}
		tickets = append(tickets, t)
	}
	return tickets, nil
}

func scanTicket(s scanner) (*model.Ticket, error) {
	t := &model.Ticket{}
	err := s.Scan(
		&t.ID, &t.BorrowerID, &t.AssigneeID, &t.Type,
		&t.Subject, &t.Status, &t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("scan ticket: %w", err)
	}
	return t, nil
}
