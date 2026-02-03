// Package service — бизнес-логика Document Service.
package service

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/vokzal-tech/document-service/internal/config"
	"github.com/vokzal-tech/document-service/internal/models"
	"github.com/vokzal-tech/document-service/internal/pdf"
	"github.com/vokzal-tech/document-service/internal/repository"
	"go.uber.org/zap"
)

// DocumentService — интерфейс сервиса документов.
type DocumentService interface {
	GenerateTicket(ctx context.Context, data *pdf.TicketData) (*models.GeneratedDocument, error)
	GeneratePD2(ctx context.Context, data *pdf.PD2Data) (*models.GeneratedDocument, error)
	GetDocument(ctx context.Context, id string) (*models.GeneratedDocument, error)
	ListDocuments(ctx context.Context, limit int) ([]*models.GeneratedDocument, error)
}

type documentService struct {
	repo      repository.DocumentRepository
	generator *pdf.Generator
	minio     *minio.Client
	cfg       *config.MinIOConfig
	logger    *zap.Logger
}

// NewDocumentService создаёт новый DocumentService.
func NewDocumentService(
	repo repository.DocumentRepository,
	generator *pdf.Generator,
	cfg *config.MinIOConfig,
	logger *zap.Logger,
) (DocumentService, error) {
	minioClient, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
	}

	return &documentService{
		repo:      repo,
		generator: generator,
		minio:     minioClient,
		cfg:       cfg,
		logger:    logger,
	}, nil
}

func (s *documentService) uploadToMinIO(ctx context.Context, fileName string, data []byte) (string, error) {
	reader := bytes.NewReader(data)
	_, err := s.minio.PutObject(ctx, s.cfg.Bucket, fileName, reader, int64(len(data)), minio.PutObjectOptions{
		ContentType: "application/pdf",
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload to MinIO: %w", err)
	}

	url := fmt.Sprintf("http://%s/%s/%s", s.cfg.Endpoint, s.cfg.Bucket, fileName)
	return url, nil
}

func (s *documentService) GenerateTicket(ctx context.Context, data *pdf.TicketData) (*models.GeneratedDocument, error) {
	pdfData, err := s.generator.GenerateTicket(data)
	if err != nil {
		return nil, fmt.Errorf("failed to generate ticket PDF: %w", err)
	}

	fileName := fmt.Sprintf("tickets/%s_%s.pdf", data.TicketID, time.Now().Format("20060102150405"))
	fileURL, err := s.uploadToMinIO(ctx, fileName, pdfData)
	if err != nil {
		return nil, err
	}

	doc := &models.GeneratedDocument{
		DocumentType: "ticket",
		EntityID:     &data.TicketID,
		FileURL:      fileURL,
		FileName:     fileName,
		Status:       "generated",
	}

	if err := s.repo.CreateDocument(ctx, doc); err != nil {
		return nil, fmt.Errorf("failed to save document record: %w", err)
	}

	s.logger.Info("Ticket PDF generated", zap.String("ticket_id", data.TicketID), zap.String("url", fileURL))
	return doc, nil
}

func (s *documentService) GeneratePD2(ctx context.Context, data *pdf.PD2Data) (*models.GeneratedDocument, error) {
	pdfData, err := s.generator.GeneratePD2(data)
	if err != nil {
		return nil, fmt.Errorf("failed to generate PD-2 PDF: %w", err)
	}

	fileName := fmt.Sprintf("pd2/%s_%s_%s.pdf", data.Series, data.Number, time.Now().Format("20060102150405"))
	fileURL, err := s.uploadToMinIO(ctx, fileName, pdfData)
	if err != nil {
		return nil, err
	}

	doc := &models.GeneratedDocument{
		DocumentType: "pd2",
		FileURL:      fileURL,
		FileName:     fileName,
		Status:       "generated",
	}

	if err := s.repo.CreateDocument(ctx, doc); err != nil {
		return nil, fmt.Errorf("failed to save document record: %w", err)
	}

	s.logger.Info("PD-2 PDF generated", zap.String("number", data.Number), zap.String("url", fileURL))
	return doc, nil
}

func (s *documentService) GetDocument(ctx context.Context, id string) (*models.GeneratedDocument, error) {
	return s.repo.FindDocumentByID(ctx, id)
}

func (s *documentService) ListDocuments(ctx context.Context, limit int) ([]*models.GeneratedDocument, error) {
	return s.repo.ListDocuments(ctx, limit)
}
