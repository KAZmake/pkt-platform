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

type DocumentRepository struct {
	db *pgxpool.Pool
}

func NewDocumentRepository(db *pgxpool.Pool) *DocumentRepository {
	return &DocumentRepository{db: db}
}

const docCols = `id, owner_id, owner_type, bucket, object_key, name, mime_type, size_bytes, created_at`

type CreateDocumentInput struct {
	OwnerID   uuid.UUID
	OwnerType string
	Bucket    string
	ObjectKey string
	Name      string
	MimeType  *string
	SizeBytes *int64
}

func (r *DocumentRepository) Create(ctx context.Context, inp CreateDocumentInput) (*model.Document, error) {
	row := r.db.QueryRow(ctx, `
		INSERT INTO documents (owner_id, owner_type, bucket, object_key, name, mime_type, size_bytes)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING `+docCols,
		inp.OwnerID, inp.OwnerType, inp.Bucket, inp.ObjectKey,
		inp.Name, inp.MimeType, inp.SizeBytes,
	)
	return scanDocument(row)
}

func (r *DocumentRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Document, error) {
	row := r.db.QueryRow(ctx, `SELECT `+docCols+` FROM documents WHERE id = $1`, id)
	return scanDocument(row)
}

func (r *DocumentRepository) ListByOwner(ctx context.Context, ownerID uuid.UUID, ownerType string) ([]*model.Document, error) {
	rows, err := r.db.Query(ctx,
		`SELECT `+docCols+` FROM documents WHERE owner_id = $1 AND owner_type = $2 ORDER BY created_at DESC`,
		ownerID, ownerType)
	if err != nil {
		return nil, fmt.Errorf("list documents: %w", err)
	}
	defer rows.Close()

	var docs []*model.Document
	for rows.Next() {
		d, err := scanDocument(rows)
		if err != nil {
			return nil, err
		}
		docs = append(docs, d)
	}
	return docs, nil
}

func (r *DocumentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM documents WHERE id = $1`, id)
	return err
}

func scanDocument(s scanner) (*model.Document, error) {
	d := &model.Document{}
	err := s.Scan(
		&d.ID, &d.OwnerID, &d.OwnerType, &d.Bucket, &d.ObjectKey,
		&d.Name, &d.MimeType, &d.SizeBytes, &d.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("scan document: %w", err)
	}
	return d, nil
}
