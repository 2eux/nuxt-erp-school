package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/opencode/erp-school-backend/internal/domain"
	"github.com/opencode/erp-school-backend/internal/dto"
	"go.uber.org/zap"
)

type AdmissionService interface {
	ListApplicants(ctx context.Context, schoolID string, filter dto.PaginationRequest) ([]dto.AdmissionApplicantResponse, int64, error)
	CreateApplicant(ctx context.Context, schoolID string, req dto.CreateAdmissionRequest) (*dto.AdmissionApplicantResponse, error)
	UpdateApplicant(ctx context.Context, id string, req dto.UpdateAdmissionRequest) (*dto.AdmissionApplicantResponse, error)
	AcceptApplicant(ctx context.Context, id string) error
	EnrollApplicant(ctx context.Context, id, schoolID, classID string) (*dto.StudentResponse, error)
}

type admissionService struct {
	db     *sqlx.DB
	logger *zap.Logger
}

func NewAdmissionService(db *sqlx.DB, logger *zap.Logger) AdmissionService {
	return &admissionService{db: db, logger: logger}
}

func (s *admissionService) ListApplicants(ctx context.Context, schoolID string, filter dto.PaginationRequest) ([]dto.AdmissionApplicantResponse, int64, error) {
	filter.Defaults()
	var items []domain.AdmissionApplicant
	var total int64

	s.db.GetContext(ctx, &total, `SELECT COUNT(*) FROM admission_applicants WHERE school_id=$1`, schoolID)
	query := `SELECT * FROM admission_applicants WHERE school_id=$1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	if err := s.db.SelectContext(ctx, &items, query, schoolID, filter.PageSize, filter.Offset()); err != nil {
		return nil, 0, domain.NewInternalError("failed to list applicants", err)
	}

	result := make([]dto.AdmissionApplicantResponse, len(items))
	for i, a := range items {
		result[i] = dto.AdmissionApplicantResponse{
			ID:             a.ID,
			FullName:       a.FullName,
			Gender:         string(a.Gender),
			PlaceOfBirth:   a.PlaceOfBirth,
			DateOfBirth:    a.DateOfBirth,
			PreviousSchool: a.PreviousSchool,
			GradeID:        a.GradeID,
			RegistrationNo: a.RegistrationNo,
			ParentName:     a.ParentName,
			ParentPhone:    a.ParentPhone,
			ParentEmail:    a.ParentEmail,
			Status:         a.Status,
			TestScore:      a.TestScore,
			InterviewScore: a.InterviewScore,
			AcceptedAt:     a.AcceptedAt,
			CreatedAt:      a.CreatedAt,
		}
	}
	return result, total, nil
}

func (s *admissionService) CreateApplicant(ctx context.Context, schoolID string, req dto.CreateAdmissionRequest) (*dto.AdmissionApplicantResponse, error) {
	regNo := fmt.Sprintf("ADM-%s-%d", time.Now().Format("20060102"), time.Now().UnixNano()%10000)

	a := &domain.AdmissionApplicant{
		ID:             uuid.New().String(),
		SchoolID:       schoolID,
		FullName:       req.FullName,
		Gender:         req.Gender,
		PlaceOfBirth:   req.PlaceOfBirth,
		DateOfBirth:    req.DateOfBirth,
		PreviousSchool: req.PreviousSchool,
		GradeID:        req.GradeID,
		RegistrationNo: regNo,
		ParentName:     req.ParentName,
		ParentPhone:    req.ParentPhone,
		ParentEmail:    req.ParentEmail,
		Status:         "pending",
	}

	query := `INSERT INTO admission_applicants (id, school_id, full_name, gender, place_of_birth, date_of_birth, previous_school, grade_id, registration_no, parent_name, parent_phone, parent_email, status) VALUES (:id, :school_id, :full_name, :gender, :place_of_birth, :date_of_birth, :previous_school, :grade_id, :registration_no, :parent_name, :parent_phone, :parent_email, :status)`
	if _, err := s.db.NamedExecContext(ctx, query, a); err != nil {
		return nil, domain.NewInternalError("failed to create applicant", err)
	}

	return &dto.AdmissionApplicantResponse{
		ID:             a.ID,
		FullName:       a.FullName,
		Gender:         string(a.Gender),
		PlaceOfBirth:   a.PlaceOfBirth,
		DateOfBirth:    a.DateOfBirth,
		PreviousSchool: a.PreviousSchool,
		GradeID:        a.GradeID,
		RegistrationNo: a.RegistrationNo,
		ParentName:     a.ParentName,
		ParentPhone:    a.ParentPhone,
		ParentEmail:    a.ParentEmail,
		Status:         a.Status,
		CreatedAt:      a.CreatedAt,
	}, nil
}

func (s *admissionService) UpdateApplicant(ctx context.Context, id string, req dto.UpdateAdmissionRequest) (*dto.AdmissionApplicantResponse, error) {
	var a domain.AdmissionApplicant
	if err := s.db.GetContext(ctx, &a, `SELECT * FROM admission_applicants WHERE id=$1`, id); err != nil {
		return nil, domain.NewNotFoundError("applicant", id)
	}

	if req.Status != "" {
		a.Status = req.Status
	}
	if req.TestScore != nil {
		a.TestScore = req.TestScore
	}
	if req.InterviewScore != nil {
		a.InterviewScore = req.InterviewScore
	}

	query := `UPDATE admission_applicants SET status=$1, test_score=$2, interview_score=$3, updated_at=NOW() WHERE id=$4`
	if _, err := s.db.ExecContext(ctx, query, a.Status, a.TestScore, a.InterviewScore, id); err != nil {
		return nil, domain.NewInternalError("failed to update applicant", err)
	}

	return &dto.AdmissionApplicantResponse{
		ID:             a.ID,
		FullName:       a.FullName,
		Gender:         string(a.Gender),
		RegistrationNo: a.RegistrationNo,
		Status:         a.Status,
		TestScore:      a.TestScore,
		InterviewScore: a.InterviewScore,
	}, nil
}

func (s *admissionService) AcceptApplicant(ctx context.Context, id string) error {
	now := time.Now()
	query := `UPDATE admission_applicants SET status='accepted', accepted_at=$1, updated_at=NOW() WHERE id=$2`
	if _, err := s.db.ExecContext(ctx, query, now, id); err != nil {
		return domain.NewInternalError("failed to accept applicant", err)
	}
	return nil
}

func (s *admissionService) EnrollApplicant(ctx context.Context, id, schoolID, classID string) (*dto.StudentResponse, error) {
	var a domain.AdmissionApplicant
	if err := s.db.GetContext(ctx, &a, `SELECT * FROM admission_applicants WHERE id=$1`, id); err != nil {
		return nil, domain.NewNotFoundError("applicant", id)
	}

	studentID := uuid.New().String()
	userID := uuid.New().String()

	tx, err := s.db.Beginx()
	if err != nil {
		return nil, domain.NewInternalError("failed to begin tx", err)
	}
	defer tx.Rollback()

	tx.ExecContext(ctx, `INSERT INTO users (id, school_id, email, username, password_hash, full_name) VALUES ($1, $2, $3, $4, $5, $6)`, userID, schoolID, a.ParentEmail, a.ParentEmail, "temp", a.FullName)
	tx.ExecContext(ctx, `INSERT INTO students (id, user_id, school_id, class_id, academic_year_id, enrollment_date, status) VALUES ($1, $2, $3, $4, $5, $6, $7)`, studentID, userID, schoolID, classID, "00000000-0000-0000-0000-000000000000", time.Now(), "active")
	tx.ExecContext(ctx, `UPDATE admission_applicants SET status='enrolled', updated_at=NOW() WHERE id=$1`, id)

	if err := tx.Commit(); err != nil {
		return nil, domain.NewInternalError("failed to enroll applicant", err)
	}

	return &dto.StudentResponse{
		ID:       studentID,
		UserID:   userID,
		SchoolID: schoolID,
		FullName: a.FullName,
		ClassID:  classID,
		Status:   "active",
	}, nil
}
