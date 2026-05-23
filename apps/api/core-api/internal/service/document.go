package service

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/KAZmake/pkt-platform/apps/api/core-api/internal/model"
	"github.com/KAZmake/pkt-platform/apps/api/core-api/internal/repository"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

const (
	uploadExpiry   = 15 * time.Minute
	downloadExpiry = 1 * time.Hour
	docsBucket     = "documents"
)

type DocumentService struct {
	repo  *repository.DocumentRepository
	minio *minio.Client
}

func NewDocumentService(repo *repository.DocumentRepository, mc *minio.Client) *DocumentService {
	return &DocumentService{repo: repo, minio: mc}
}

type PresignedUploadResult struct {
	DocumentID uuid.UUID `json:"document_id"`
	UploadURL  string    `json:"upload_url"`
	ExpiresIn  int       `json:"expires_in_seconds"`
}

type PresignedDownloadResult struct {
	DownloadURL string `json:"download_url"`
	ExpiresIn   int    `json:"expires_in_seconds"`
}

// InitiateUpload pre-registers a document in the DB and returns a presigned PUT URL.
// The client uploads the file directly to MinIO using the returned URL.
func (s *DocumentService) InitiateUpload(ctx context.Context, ownerID uuid.UUID, ownerType, name string, mimeType *string) (*PresignedUploadResult, error) {
	objectKey := fmt.Sprintf("%s/%s/%s/%s", ownerType, ownerID, uuid.New(), name)

	// Pre-register in DB (size unknown until upload completes)
	doc, err := s.repo.Create(ctx, repository.CreateDocumentInput{
		OwnerID:   ownerID,
		OwnerType: ownerType,
		Bucket:    docsBucket,
		ObjectKey: objectKey,
		Name:      name,
		MimeType:  mimeType,
	})
	if err != nil {
		return nil, fmt.Errorf("register document: %w", err)
	}

	uploadURL, err := s.minio.PresignedPutObject(ctx, docsBucket, objectKey, uploadExpiry)
	if err != nil {
		// Roll back DB record if we can't generate URL
		_ = s.repo.Delete(ctx, doc.ID)
		return nil, fmt.Errorf("presign upload: %w", err)
	}

	return &PresignedUploadResult{
		DocumentID: doc.ID,
		UploadURL:  uploadURL.String(),
		ExpiresIn:  int(uploadExpiry.Seconds()),
	}, nil
}

// GetDownloadURL generates a short-lived presigned GET URL for an existing document.
func (s *DocumentService) GetDownloadURL(ctx context.Context, docID uuid.UUID) (*PresignedDownloadResult, error) {
	doc, err := s.repo.GetByID(ctx, docID)
	if err != nil {
		return nil, err
	}
	if doc == nil {
		return nil, nil
	}

	downloadURL, err := s.minio.PresignedGetObject(ctx, doc.Bucket, doc.ObjectKey, downloadExpiry, url.Values{})
	if err != nil {
		return nil, fmt.Errorf("presign download: %w", err)
	}

	return &PresignedDownloadResult{
		DownloadURL: downloadURL.String(),
		ExpiresIn:   int(downloadExpiry.Seconds()),
	}, nil
}

// ListByOwner returns all documents for a given owner.
func (s *DocumentService) ListByOwner(ctx context.Context, ownerID uuid.UUID, ownerType string) ([]*model.Document, error) {
	return s.repo.ListByOwner(ctx, ownerID, ownerType)
}

// Delete removes the document from MinIO and the DB.
func (s *DocumentService) Delete(ctx context.Context, docID uuid.UUID) error {
	doc, err := s.repo.GetByID(ctx, docID)
	if err != nil {
		return err
	}
	if doc == nil {
		return nil
	}

	if err := s.minio.RemoveObject(ctx, doc.Bucket, doc.ObjectKey, minio.RemoveObjectOptions{}); err != nil {
		return fmt.Errorf("remove from minio: %w", err)
	}
	return s.repo.Delete(ctx, docID)
}
