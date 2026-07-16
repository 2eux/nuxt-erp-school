package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/opencode/erp-school-backend/internal/domain"
	"github.com/opencode/erp-school-backend/internal/dto"
	"go.uber.org/zap"
)

type MedicalService interface {
	ListRecords(ctx context.Context, studentID string) ([]dto.MedicalRecordResponse, error)
	CreateRecord(ctx context.Context, createdBy string, req dto.CreateMedicalRecordRequest) (*dto.MedicalRecordResponse, error)
}

type medicalService struct {
	db     *sqlx.DB
	logger *zap.Logger
}

func NewMedicalService(db *sqlx.DB, logger *zap.Logger) MedicalService {
	return &medicalService{db: db, logger: logger}
}

func (s *medicalService) ListRecords(ctx context.Context, studentID string) ([]dto.MedicalRecordResponse, error) {
	var items []domain.MedicalRecord
	query := `SELECT * FROM medical_records WHERE student_id=$1 ORDER BY date DESC`
	if err := s.db.SelectContext(ctx, &items, query, studentID); err != nil {
		return nil, domain.NewInternalError("failed to list medical records", err)
	}

	result := make([]dto.MedicalRecordResponse, len(items))
	for i, m := range items {
		result[i] = dto.MedicalRecordResponse{
			ID:         m.ID,
			StudentID:  m.StudentID,
			Date:       m.Date,
			Diagnosis:  m.Diagnosis,
			Treatment:  m.Treatment,
			Medication: m.Medication,
			DocType:    m.DocType,
			Notes:      m.Notes,
			CreatedBy:  m.CreatedBy,
		}
	}
	return result, nil
}

func (s *medicalService) CreateRecord(ctx context.Context, createdBy string, req dto.CreateMedicalRecordRequest) (*dto.MedicalRecordResponse, error) {
	mr := &domain.MedicalRecord{
		ID:         uuid.New().String(),
		StudentID:  req.StudentID,
		Date:       req.Date,
		Diagnosis:  req.Diagnosis,
		Treatment:  req.Treatment,
		Medication: req.Medication,
		DocType:    req.DocType,
		Notes:      req.Notes,
		CreatedBy:  createdBy,
	}

	query := `INSERT INTO medical_records (id, student_id, date, diagnosis, treatment, medication, doc_type, notes, created_by) VALUES (:id, :student_id, :date, :diagnosis, :treatment, :medication, :doc_type, :notes, :created_by)`
	if _, err := s.db.NamedExecContext(ctx, query, mr); err != nil {
		return nil, domain.NewInternalError("failed to create medical record", err)
	}

	return &dto.MedicalRecordResponse{
		ID:         mr.ID,
		StudentID:  mr.StudentID,
		Date:       mr.Date,
		Diagnosis:  mr.Diagnosis,
		Treatment:  mr.Treatment,
		Medication: mr.Medication,
		DocType:    mr.DocType,
		Notes:      mr.Notes,
		CreatedBy:  mr.CreatedBy,
	}, nil
}
