package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/opencode/erp-school-backend/internal/domain"
	"github.com/opencode/erp-school-backend/internal/dto"
	"github.com/opencode/erp-school-backend/internal/repository"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type StudentService interface {
	ListStudents(ctx context.Context, schoolID string, filter dto.PaginationRequest, classID, status string) ([]dto.StudentResponse, int64, error)
	GetStudent(ctx context.Context, id string) (*dto.StudentDetail, error)
	CreateStudent(ctx context.Context, schoolID string, req dto.CreateStudentRequest) (*dto.StudentDetail, error)
	UpdateStudent(ctx context.Context, id string, req dto.UpdateStudentRequest) (*dto.StudentResponse, error)
	DeleteStudent(ctx context.Context, id string) error
	PromoteStudent(ctx context.Context, id string, newClassID, newAcademicYearID string) (*dto.StudentResponse, error)
	TransferStudent(ctx context.Context, id string, newSchoolID, newClassID string) (*dto.StudentResponse, error)
	ListParents(ctx context.Context, studentID string) ([]dto.StudentParentResponse, error)
	LinkParent(ctx context.Context, studentID, schoolID string, req dto.CreateStudentParentRequest) (*dto.StudentParentResponse, error)
	ListDocuments(ctx context.Context, studentID string) ([]domain.StudentDocument, error)
	UploadDocument(ctx context.Context, studentID, name, docType, fileURL string) (*domain.StudentDocument, error)
}

type studentService struct {
	studentRepo repository.StudentRepository
	logger      *zap.Logger
}

func NewStudentService(studentRepo repository.StudentRepository, logger *zap.Logger) StudentService {
	return &studentService{studentRepo: studentRepo, logger: logger}
}

func (s *studentService) ListStudents(ctx context.Context, schoolID string, filter dto.PaginationRequest, classID, status string) ([]dto.StudentResponse, int64, error) {
	filter.Defaults()
	filters := map[string]interface{}{
		"class_id": classID,
		"status":   status,
		"search":   filter.Search,
	}
	items, total, err := s.studentRepo.ListStudents(ctx, schoolID, filter.PageSize, filter.Offset(), filters)
	if err != nil {
		return nil, 0, domain.NewInternalError("failed to list students", err)
	}

	result := make([]dto.StudentResponse, len(items))
	for i, st := range items {
		result[i] = dto.StudentResponse{
			ID:             st.ID,
			UserID:         st.UserID,
			SchoolID:       st.SchoolID,
			NIS:            st.NIS,
			NISN:           st.NISN,
			NIK:            st.NIK,
			ClassID:        st.ClassID,
			AcademicYearID: st.AcademicYearID,
			EnrollmentDate: st.EnrollmentDate,
			Status:         st.Status,
			CreatedAt:      st.CreatedAt,
			UpdatedAt:      st.UpdatedAt,
		}
	}
	return result, total, nil
}

func (s *studentService) GetStudent(ctx context.Context, id string) (*dto.StudentDetail, error) {
	st, parents, err := s.studentRepo.GetStudentDetail(ctx, id)
	if err != nil {
		return nil, domain.NewNotFoundError("student", id)
	}

	parentResponses := make([]dto.StudentParentResponse, len(parents))
	for i, p := range parents {
		parentResponses[i] = dto.StudentParentResponse{
			ID:          p.ID,
			UserID:      p.UserID,
			Relation:    p.Relation,
			IsPrimary:   p.IsPrimary,
			Occupation:  p.Occupation,
			Institution: p.Institution,
			CreatedAt:   p.CreatedAt,
		}
	}

	return &dto.StudentDetail{
		StudentResponse: dto.StudentResponse{
			ID:             st.ID,
			UserID:         st.UserID,
			SchoolID:       st.SchoolID,
			NIS:            st.NIS,
			NISN:           st.NISN,
			NIK:            st.NIK,
			ClassID:        st.ClassID,
			AcademicYearID: st.AcademicYearID,
			EnrollmentDate: st.EnrollmentDate,
			Status:         st.Status,
			CreatedAt:      st.CreatedAt,
			UpdatedAt:      st.UpdatedAt,
		},
		Parents: parentResponses,
	}, nil
}

func (s *studentService) CreateStudent(ctx context.Context, schoolID string, req dto.CreateStudentRequest) (*dto.StudentDetail, error) {
	userID := req.UserID
	if userID == "" {
		userID = uuid.New().String()
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, domain.NewInternalError("failed to hash password", err)
	}

	user := &domain.User{
		ID:           userID,
		SchoolID:     schoolID,
		Email:        req.Email,
		Username:     req.Email,
		PasswordHash: string(hash),
		FullName:     req.FullName,
		IsActive:     true,
	}

	profile := &domain.UserProfile{
		ID:           uuid.New().String(),
		UserID:       userID,
		Gender:       req.Gender,
		PlaceOfBirth: req.PlaceOfBirth,
		DateOfBirth:  req.DateOfBirth,
		Address:      req.Address,
	}

	student := &domain.Student{
		ID:             uuid.New().String(),
		UserID:         userID,
		SchoolID:       schoolID,
		NIS:            req.NIS,
		NISN:           req.NISN,
		NIK:            req.NIK,
		ClassID:        req.ClassID,
		AcademicYearID: req.AcademicYearID,
		EnrollmentDate: req.EnrollmentDate,
		Status:         "active",
	}

	if err := s.studentRepo.CreateStudent(ctx, student, user, profile); err != nil {
		return nil, domain.NewInternalError("failed to create student", err)
	}

	return s.GetStudent(ctx, student.ID)
}

func (s *studentService) UpdateStudent(ctx context.Context, id string, req dto.UpdateStudentRequest) (*dto.StudentResponse, error) {
	st, err := s.studentRepo.GetStudentByID(ctx, id)
	if err != nil {
		return nil, domain.NewNotFoundError("student", id)
	}

	if req.NIS != "" {
		st.NIS = req.NIS
	}
	if req.NISN != "" {
		st.NISN = req.NISN
	}
	if req.NIK != "" {
		st.NIK = req.NIK
	}
	if req.ClassID != "" {
		st.ClassID = req.ClassID
	}
	if req.Status != "" {
		st.Status = req.Status
	}

	if err := s.studentRepo.UpdateStudent(ctx, st); err != nil {
		return nil, domain.NewInternalError("failed to update student", err)
	}

	return &dto.StudentResponse{
		ID:             st.ID,
		UserID:         st.UserID,
		SchoolID:       st.SchoolID,
		NIS:            st.NIS,
		NISN:           st.NISN,
		NIK:            st.NIK,
		ClassID:        st.ClassID,
		AcademicYearID: st.AcademicYearID,
		EnrollmentDate: st.EnrollmentDate,
		Status:         st.Status,
		CreatedAt:      st.CreatedAt,
		UpdatedAt:      st.UpdatedAt,
	}, nil
}

func (s *studentService) DeleteStudent(ctx context.Context, id string) error {
	st, err := s.studentRepo.GetStudentByID(ctx, id)
	if err != nil {
		return domain.NewNotFoundError("student", id)
	}
	st.Status = "inactive"
	return s.studentRepo.UpdateStudent(ctx, st)
}

func (s *studentService) PromoteStudent(ctx context.Context, id string, newClassID, newAcademicYearID string) (*dto.StudentResponse, error) {
	st, err := s.studentRepo.GetStudentByID(ctx, id)
	if err != nil {
		return nil, domain.NewNotFoundError("student", id)
	}

	st.ClassID = newClassID
	st.AcademicYearID = newAcademicYearID

	if err := s.studentRepo.UpdateStudent(ctx, st); err != nil {
		return nil, domain.NewInternalError("failed to promote student", err)
	}

	return &dto.StudentResponse{
		ID:             st.ID,
		UserID:         st.UserID,
		SchoolID:       st.SchoolID,
		NIS:            st.NIS,
		NISN:           st.NISN,
		NIK:            st.NIK,
		ClassID:        st.ClassID,
		AcademicYearID: st.AcademicYearID,
		EnrollmentDate: st.EnrollmentDate,
		Status:         st.Status,
	}, nil
}

func (s *studentService) TransferStudent(ctx context.Context, id string, newSchoolID, newClassID string) (*dto.StudentResponse, error) {
	st, err := s.studentRepo.GetStudentByID(ctx, id)
	if err != nil {
		return nil, domain.NewNotFoundError("student", id)
	}

	st.SchoolID = newSchoolID
	st.ClassID = newClassID
	st.Status = "transferred"

	if err := s.studentRepo.UpdateStudent(ctx, st); err != nil {
		return nil, domain.NewInternalError("failed to transfer student", err)
	}

	return &dto.StudentResponse{
		ID:       st.ID,
		SchoolID: st.SchoolID,
		ClassID:  st.ClassID,
		Status:   st.Status,
	}, nil
}

func (s *studentService) ListParents(ctx context.Context, studentID string) ([]dto.StudentParentResponse, error) {
	parents, err := s.studentRepo.GetParents(ctx, studentID)
	if err != nil {
		return nil, domain.NewInternalError("failed to list parents", err)
	}

	result := make([]dto.StudentParentResponse, len(parents))
	for i, p := range parents {
		result[i] = dto.StudentParentResponse{
			ID:          p.ID,
			UserID:      p.UserID,
			Relation:    p.Relation,
			IsPrimary:   p.IsPrimary,
			Occupation:  p.Occupation,
			Institution: p.Institution,
			CreatedAt:   p.CreatedAt,
		}
	}
	return result, nil
}

func (s *studentService) LinkParent(ctx context.Context, studentID, schoolID string, req dto.CreateStudentParentRequest) (*dto.StudentParentResponse, error) {
	userID := req.UserID
	if userID == "" {
		userID = uuid.New().String()
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	user := &domain.User{
		ID:           userID,
		SchoolID:     schoolID,
		Email:        req.Email,
		Username:     req.Email,
		PasswordHash: string(hash),
		FullName:     req.FullName,
		Phone:        req.Phone,
		IsActive:     true,
	}

	parent := &domain.StudentParent{
		ID:          uuid.New().String(),
		StudentID:   studentID,
		UserID:      userID,
		Relation:    req.Relation,
		IsPrimary:   req.IsPrimary,
		Occupation:  req.Occupation,
		Institution: req.Institution,
		Income:      req.Income,
	}

	if err := s.studentRepo.LinkParent(ctx, parent, user); err != nil {
		return nil, domain.NewInternalError("failed to link parent", err)
	}

	return &dto.StudentParentResponse{
		ID:          parent.ID,
		UserID:      parent.UserID,
		Relation:    parent.Relation,
		IsPrimary:   parent.IsPrimary,
		Occupation:  parent.Occupation,
		Institution: parent.Institution,
		CreatedAt:   time.Now(),
	}, nil
}

func (s *studentService) ListDocuments(ctx context.Context, studentID string) ([]domain.StudentDocument, error) {
	return s.studentRepo.ListDocuments(ctx, studentID)
}

func (s *studentService) UploadDocument(ctx context.Context, studentID, name, docType, fileURL string) (*domain.StudentDocument, error) {
	doc := &domain.StudentDocument{
		ID:        uuid.New().String(),
		StudentID: studentID,
		Name:      name,
		DocType:   docType,
		FileURL:   fileURL,
		Status:    "pending",
	}

	if err := s.studentRepo.CreateDocument(ctx, doc); err != nil {
		return nil, domain.NewInternalError("failed to create document", err)
	}

	return doc, nil
}
