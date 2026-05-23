package model

import (
	"time"

	"github.com/google/uuid"
)

type Document struct {
	ID        uuid.UUID `json:"id"`
	OwnerID   uuid.UUID `json:"owner_id"`
	OwnerType string    `json:"owner_type"` // borrower | application | collateral
	Bucket    string    `json:"bucket"`
	ObjectKey string    `json:"object_key"`
	Name      string    `json:"name"`
	MimeType  *string   `json:"mime_type,omitempty"`
	SizeBytes *int64    `json:"size_bytes,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
