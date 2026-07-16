package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/opencode/erp-school-backend/internal/domain"
	"github.com/opencode/erp-school-backend/internal/dto"
	"github.com/opencode/erp-school-backend/internal/infrastructure/database"
	"go.uber.org/zap"
)

type DocumentService interface {
	ListDocuments(ctx context.Context, schoolID string, filter dto.PaginationRequest) ([]dto.DocumentResponse, int64, error)
	GetDocument(ctx context.Context, id string) (*dto.DocumentResponse, error)
	UploadDocument(ctx context.Context, schoolID, createdBy, title, docType, fileURL string, fileSize int64, mimeType string) (*dto.DocumentResponse, error)
	DeleteDocument(ctx context.Context, id string) error

	ListLetters(ctx context.Context, schoolID string, filter dto.PaginationRequest) ([]domain.Letter, int64, error)
	GetLetter(ctx context.Context, id string) (*domain.Letter, error)
	CreateLetter(ctx context.Context, schoolID, createdBy string, req dto.CreateLetterRequest) (*domain.Letter, error)
}

type documentService struct {
	db     *sqlx.DB
	logger *zap.Logger
}

func NewDocumentService(db *sqlx.DB, logger *zap.Logger) DocumentService {
	return &documentService{db: db, logger: logger}
}

func (s *documentService) ListDocuments(ctx context.Context, schoolID string, filter dto.PaginationRequest) ([]dto.DocumentResponse, int64, error) {
	filter.Defaults()

	var items []domain.Document
	if err := s.db.SelectContext(ctx, &items, database.ListDocuments, schoolID, filter.PageSize, filter.Offset()); err != nil {
		return nil, 0, domain.NewInternalError("failed to list documents", err)
	}

	var total int64
	s.db.GetContext(ctx, &total, `SELECT COUNT(*) FROM documents WHERE school_id=$1`, schoolID)

	result := make([]dto.DocumentResponse, len(items))
	for i, d := range items {
		result[i] = dto.DocumentResponse{
			ID:        d.ID,
			SchoolID:  d.SchoolID,
			Title:     d.Title,
			DocType:   d.DocType,
			FileURL:   d.FileURL,
			FileSize:  d.FileSize,
			MimeType:  d.MimeType,
			CreatedBy: d.CreatedBy,
			CreatedAt: d.CreatedAt,
		}
	}
	return result, total, nil
}

func (s *documentService) GetDocument(ctx context.Context, id string) (*dto.DocumentResponse, error) {
	var doc domain.Document
	if err := s.db.GetContext(ctx, &doc, `SELECT * FROM documents WHERE id=$1`, id); err != nil {
		return nil, domain.NewNotFoundError("document", id)
	}
	return &dto.DocumentResponse{
		ID:        doc.ID,
		SchoolID:  doc.SchoolID,
		Title:     doc.Title,
		DocType:   doc.DocType,
		FileURL:   doc.FileURL,
		FileSize:  doc.FileSize,
		MimeType:  doc.MimeType,
		CreatedBy: doc.CreatedBy,
		CreatedAt: doc.CreatedAt,
	}, nil
}

func (s *documentService) UploadDocument(ctx context.Context, schoolID, createdBy, title, docType, fileURL string, fileSize int64, mimeType string) (*dto.DocumentResponse, error) {
	doc := &domain.Document{
		ID:        uuid.New().String(),
		SchoolID:  schoolID,
		Title:     title,
		DocType:   docType,
		FileURL:   fileURL,
		FileSize:  fileSize,
		MimeType:  mimeType,
		CreatedBy: createdBy,
	}

	query := `INSERT INTO documents (id, school_id, title, doc_type, file_url, file_size, mime_type, created_by) VALUES (:id, :school_id, :title, :doc_type, :file_url, :file_size, :mime_type, :created_by)`
	if _, err := s.db.NamedExecContext(ctx, query, doc); err != nil {
		return nil, domain.NewInternalError("failed to upload document", err)
	}

	return &dto.DocumentResponse{
		ID:        doc.ID,
		SchoolID:  doc.SchoolID,
		Title:     doc.Title,
		DocType:   doc.DocType,
		FileURL:   doc.FileURL,
		FileSize:  doc.FileSize,
		MimeType:  doc.MimeType,
		CreatedBy: doc.CreatedBy,
		CreatedAt: doc.CreatedAt,
	}, nil
}

func (s *documentService) DeleteDocument(ctx context.Context, id string) error {
	if _, err := s.db.ExecContext(ctx, `DELETE FROM documents WHERE id=$1`, id); err != nil {
		return domain.NewInternalError("failed to delete document", err)
	}
	return nil
}

func (s *documentService) ListLetters(ctx context.Context, schoolID string, filter dto.PaginationRequest) ([]domain.Letter, int64, error) {
	filter.Defaults()
	var items []domain.Letter
	var total int64

	s.db.GetContext(ctx, &total, `SELECT COUNT(*) FROM letters WHERE school_id=$1`, schoolID)
	query := `SELECT * FROM letters WHERE school_id=$1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	if err := s.db.SelectContext(ctx, &items, query, schoolID, filter.PageSize, filter.Offset()); err != nil {
		return nil, 0, domain.NewInternalError("failed to list letters", err)
	}
	return items, total, nil
}

func (s *documentService) GetLetter(ctx context.Context, id string) (*domain.Letter, error) {
	var letter domain.Letter
	if err := s.db.GetContext(ctx, &letter, `SELECT * FROM letters WHERE id=$1`, id); err != nil {
		return nil, domain.NewNotFoundError("letter", id)
	}
	return &letter, nil
}

func (s *documentService) CreateLetter(ctx context.Context, schoolID, createdBy string, req dto.CreateLetterRequest) (*domain.Letter, error) {
	letter := &domain.Letter{
		ID:         uuid.New().String(),
		SchoolID:   schoolID,
		LetterNo:   uuid.New().String()[:8],
		Title:      req.Title,
		Content:    req.Content,
		LetterType: req.LetterType,
		To:         req.To,
		Status:     "draft",
		CreatedBy:  createdBy,
	}

	query := `INSERT INTO letters (id, school_id, letter_no, title, content, letter_type, "to", status, created_by) VALUES (:id, :school_id, :letter_no, :title, :content, :letter_type, :to, :status, :created_by)`
	if _, err := s.db.NamedExecContext(ctx, query, letter); err != nil {
		return nil, domain.NewInternalError("failed to create letter", err)
	}
	return letter, nil
}
