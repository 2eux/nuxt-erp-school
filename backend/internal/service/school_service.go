package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/opencode/erp-school-backend/internal/domain"
	"github.com/opencode/erp-school-backend/internal/dto"
	"github.com/opencode/erp-school-backend/internal/repository"
	"go.uber.org/zap"
)

type SchoolService interface {
	ListSchools(ctx context.Context, filter dto.PaginationRequest) ([]dto.SchoolResponse, int64, error)
	GetSchool(ctx context.Context, id string) (*dto.SchoolResponse, error)
	CreateSchool(ctx context.Context, req dto.CreateSchoolRequest) (*dto.SchoolResponse, error)
	UpdateSchool(ctx context.Context, id string, req dto.UpdateSchoolRequest) (*dto.SchoolResponse, error)
	DeleteSchool(ctx context.Context, id string) error

	ListAcademicYears(ctx context.Context, schoolID string, filter dto.PaginationRequest) ([]dto.AcademicYearResponse, int64, error)
	GetAcademicYear(ctx context.Context, id string) (*dto.AcademicYearResponse, error)
	CreateAcademicYear(ctx context.Context, schoolID string, req dto.CreateAcademicYearRequest) (*dto.AcademicYearResponse, error)
	UpdateAcademicYear(ctx context.Context, id string, req dto.UpdateAcademicYearRequest) (*dto.AcademicYearResponse, error)
	SetActiveAcademicYear(ctx context.Context, schoolID, ayID string) error

	ListSemesters(ctx context.Context, academicYearID string) ([]dto.SemesterResponse, error)
	GetSemester(ctx context.Context, id string) (*dto.SemesterResponse, error)
	CreateSemester(ctx context.Context, req dto.CreateSemesterRequest) (*dto.SemesterResponse, error)
	UpdateSemester(ctx context.Context, id string, req dto.CreateSemesterRequest) (*dto.SemesterResponse, error)

	ListGrades(ctx context.Context, schoolID string) ([]dto.GradeResponse, error)
	GetGrade(ctx context.Context, id string) (*dto.GradeResponse, error)
	CreateGrade(ctx context.Context, schoolID string, req dto.CreateGradeRequest) (*dto.GradeResponse, error)
	UpdateGrade(ctx context.Context, id string, req dto.CreateGradeRequest) (*dto.GradeResponse, error)

	ListClasses(ctx context.Context, schoolID string, filter dto.PaginationRequest) ([]dto.ClassResponse, int64, error)
	GetClass(ctx context.Context, id string) (*dto.ClassResponse, error)
	CreateClass(ctx context.Context, schoolID string, req dto.CreateClassRequest) (*dto.ClassResponse, error)
	UpdateClass(ctx context.Context, id string, req dto.UpdateClassRequest) (*dto.ClassResponse, error)

	ListSubjects(ctx context.Context, schoolID string, filter dto.PaginationRequest) ([]dto.SubjectResponse, int64, error)
	GetSubject(ctx context.Context, id string) (*dto.SubjectResponse, error)
	CreateSubject(ctx context.Context, schoolID string, req dto.CreateSubjectRequest) (*dto.SubjectResponse, error)
	UpdateSubject(ctx context.Context, id string, req dto.UpdateSubjectRequest) (*dto.SubjectResponse, error)

	ListCurriculums(ctx context.Context, gradeID, subjectID, semesterID string) ([]dto.CurriculumResponse, error)
	GetCurriculum(ctx context.Context, id string) (*dto.CurriculumResponse, error)
	CreateCurriculum(ctx context.Context, schoolID string, req dto.CreateCurriculumRequest) (*dto.CurriculumResponse, error)
	UpdateCurriculum(ctx context.Context, id string, req dto.CreateCurriculumRequest) (*dto.CurriculumResponse, error)
}

type schoolService struct {
	schoolRepo repository.SchoolRepository
	logger     *zap.Logger
}

func NewSchoolService(schoolRepo repository.SchoolRepository, logger *zap.Logger) SchoolService {
	return &schoolService{schoolRepo: schoolRepo, logger: logger}
}

func (s *schoolService) ListSchools(ctx context.Context, filter dto.PaginationRequest) ([]dto.SchoolResponse, int64, error) {
	schools, err := s.schoolRepo.ListSchools(ctx)
	if err != nil {
		return nil, 0, domain.NewInternalError("failed to list schools", err)
	}

	result := make([]dto.SchoolResponse, len(schools))
	for i, sc := range schools {
		result[i] = toSchoolResponse(sc)
	}
	return result, int64(len(result)), nil
}

func (s *schoolService) GetSchool(ctx context.Context, id string) (*dto.SchoolResponse, error) {
	school, err := s.schoolRepo.GetSchoolByID(ctx, id)
	if err != nil {
		return nil, domain.NewNotFoundError("school", id)
	}
	r := toSchoolResponse(*school)
	return &r, nil
}

func (s *schoolService) CreateSchool(ctx context.Context, req dto.CreateSchoolRequest) (*dto.SchoolResponse, error) {
	school := &domain.School{
		Name:            req.Name,
		NPSN:            req.NPSN,
		Address:         req.Address,
		City:            req.City,
		Province:        req.Province,
		PostalCode:      req.PostalCode,
		Phone:           req.Phone,
		Email:           req.Email,
		Website:         req.Website,
		Type:            req.Type,
		Accreditation:   req.Accreditation,
		EstablishedDate: req.EstablishedDate,
		IsActive:        true,
	}

	if err := s.schoolRepo.CreateSchool(ctx, school); err != nil {
		return nil, domain.NewInternalError("failed to create school", err)
	}

	r := toSchoolResponse(*school)
	return &r, nil
}

func (s *schoolService) UpdateSchool(ctx context.Context, id string, req dto.UpdateSchoolRequest) (*dto.SchoolResponse, error) {
	existing, err := s.schoolRepo.GetSchoolByID(ctx, id)
	if err != nil {
		return nil, domain.NewNotFoundError("school", id)
	}

	if req.Name != "" {
		existing.Name = req.Name
	}
	if req.NPSN != "" {
		existing.NPSN = req.NPSN
	}
	if req.Address != "" {
		existing.Address = req.Address
	}
	if req.City != "" {
		existing.City = req.City
	}
	if req.Province != "" {
		existing.Province = req.Province
	}
	if req.PostalCode != "" {
		existing.PostalCode = req.PostalCode
	}
	if req.Phone != "" {
		existing.Phone = req.Phone
	}
	if req.Email != "" {
		existing.Email = req.Email
	}
	if req.Website != "" {
		existing.Website = req.Website
	}
	if req.Type != "" {
		existing.Type = req.Type
	}
	if req.Accreditation != "" {
		existing.Accreditation = req.Accreditation
	}
	if req.IsActive != nil {
		existing.IsActive = *req.IsActive
	}

	if err := s.schoolRepo.UpdateSchool(ctx, existing); err != nil {
		return nil, domain.NewInternalError("failed to update school", err)
	}

	r := toSchoolResponse(*existing)
	return &r, nil
}

func (s *schoolService) DeleteSchool(ctx context.Context, id string) error {
	if err := s.schoolRepo.DeleteSchool(ctx, id); err != nil {
		return domain.NewInternalError("failed to delete school", err)
	}
	return nil
}

func (s *schoolService) ListAcademicYears(ctx context.Context, schoolID string, filter dto.PaginationRequest) ([]dto.AcademicYearResponse, int64, error) {
	filter.Defaults()
	items, total, err := s.schoolRepo.ListAcademicYears(ctx, schoolID, filter.PageSize, filter.Offset())
	if err != nil {
		return nil, 0, domain.NewInternalError("failed to list academic years", err)
	}

	result := make([]dto.AcademicYearResponse, len(items))
	for i, ay := range items {
		result[i] = dto.AcademicYearResponse{
			ID:        ay.ID,
			SchoolID:  ay.SchoolID,
			Name:      ay.Name,
			StartDate: ay.StartDate,
			EndDate:   ay.EndDate,
			IsActive:  ay.IsActive,
			CreatedAt: ay.CreatedAt,
			UpdatedAt: ay.UpdatedAt,
		}
	}
	return result, total, nil
}

func (s *schoolService) GetAcademicYear(ctx context.Context, id string) (*dto.AcademicYearResponse, error) {
	ay, err := s.schoolRepo.GetAcademicYearByID(ctx, id)
	if err != nil {
		return nil, domain.NewNotFoundError("academic year", id)
	}
	return &dto.AcademicYearResponse{
		ID:        ay.ID,
		SchoolID:  ay.SchoolID,
		Name:      ay.Name,
		StartDate: ay.StartDate,
		EndDate:   ay.EndDate,
		IsActive:  ay.IsActive,
		CreatedAt: ay.CreatedAt,
		UpdatedAt: ay.UpdatedAt,
	}, nil
}

func (s *schoolService) CreateAcademicYear(ctx context.Context, schoolID string, req dto.CreateAcademicYearRequest) (*dto.AcademicYearResponse, error) {
	ay := &domain.AcademicYear{
		ID:        uuid.New().String(),
		SchoolID:  schoolID,
		Name:      req.Name,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		IsActive:  false,
	}

	if err := s.schoolRepo.CreateAcademicYear(ctx, ay); err != nil {
		return nil, domain.NewInternalError("failed to create academic year", err)
	}

	return &dto.AcademicYearResponse{
		ID:        ay.ID,
		SchoolID:  ay.SchoolID,
		Name:      ay.Name,
		StartDate: ay.StartDate,
		EndDate:   ay.EndDate,
		IsActive:  ay.IsActive,
		CreatedAt: ay.CreatedAt,
		UpdatedAt: ay.UpdatedAt,
	}, nil
}

func (s *schoolService) UpdateAcademicYear(ctx context.Context, id string, req dto.UpdateAcademicYearRequest) (*dto.AcademicYearResponse, error) {
	ay, err := s.schoolRepo.GetAcademicYearByID(ctx, id)
	if err != nil {
		return nil, domain.NewNotFoundError("academic year", id)
	}

	if req.Name != "" {
		ay.Name = req.Name
	}
	if req.StartDate != nil {
		ay.StartDate = *req.StartDate
	}
	if req.EndDate != nil {
		ay.EndDate = *req.EndDate
	}
	if req.IsActive != nil {
		ay.IsActive = *req.IsActive
	}

	if err := s.schoolRepo.UpdateAcademicYear(ctx, ay); err != nil {
		return nil, domain.NewInternalError("failed to update academic year", err)
	}

	return &dto.AcademicYearResponse{
		ID:        ay.ID,
		SchoolID:  ay.SchoolID,
		Name:      ay.Name,
		StartDate: ay.StartDate,
		EndDate:   ay.EndDate,
		IsActive:  ay.IsActive,
		CreatedAt: ay.CreatedAt,
		UpdatedAt: ay.UpdatedAt,
	}, nil
}

func (s *schoolService) SetActiveAcademicYear(ctx context.Context, schoolID, ayID string) error {
	active, err := s.schoolRepo.GetActiveAcademicYear(ctx, schoolID)
	if err == nil && active != nil {
		active.IsActive = false
		s.schoolRepo.UpdateAcademicYear(ctx, active)
	}

	ay, err := s.schoolRepo.GetAcademicYearByID(ctx, ayID)
	if err != nil {
		return domain.NewNotFoundError("academic year", ayID)
	}

	ay.IsActive = true
	if err := s.schoolRepo.UpdateAcademicYear(ctx, ay); err != nil {
		return domain.NewInternalError("failed to set active academic year", err)
	}

	return nil
}

func (s *schoolService) ListSemesters(ctx context.Context, academicYearID string) ([]dto.SemesterResponse, error) {
	items, err := s.schoolRepo.ListSemesters(ctx, academicYearID)
	if err != nil {
		return nil, domain.NewInternalError("failed to list semesters", err)
	}

	result := make([]dto.SemesterResponse, len(items))
	for i, sem := range items {
		result[i] = dto.SemesterResponse{
			ID:             sem.ID,
			AcademicYearID: sem.AcademicYearID,
			Name:           sem.Name,
			SemesterNumber: sem.SemesterNumber,
			StartDate:      sem.StartDate,
			EndDate:        sem.EndDate,
			IsActive:       sem.IsActive,
		}
	}
	return result, nil
}

func (s *schoolService) GetSemester(ctx context.Context, id string) (*dto.SemesterResponse, error) {
	sem, err := s.schoolRepo.GetSemesterByID(ctx, id)
	if err != nil {
		return nil, domain.NewNotFoundError("semester", id)
	}
	return &dto.SemesterResponse{
		ID:             sem.ID,
		AcademicYearID: sem.AcademicYearID,
		Name:           sem.Name,
		SemesterNumber: sem.SemesterNumber,
		StartDate:      sem.StartDate,
		EndDate:        sem.EndDate,
		IsActive:       sem.IsActive,
	}, nil
}

func (s *schoolService) CreateSemester(ctx context.Context, req dto.CreateSemesterRequest) (*dto.SemesterResponse, error) {
	sem := &domain.Semester{
		ID:             uuid.New().String(),
		AcademicYearID: req.AcademicYearID,
		Name:           req.Name,
		SemesterNumber: req.SemesterNumber,
		StartDate:      req.StartDate,
		EndDate:        req.EndDate,
		IsActive:       true,
	}

	if err := s.schoolRepo.CreateSemester(ctx, sem); err != nil {
		return nil, domain.NewInternalError("failed to create semester", err)
	}

	return &dto.SemesterResponse{
		ID:             sem.ID,
		AcademicYearID: sem.AcademicYearID,
		Name:           sem.Name,
		SemesterNumber: sem.SemesterNumber,
		StartDate:      sem.StartDate,
		EndDate:        sem.EndDate,
		IsActive:       sem.IsActive,
	}, nil
}

func (s *schoolService) UpdateSemester(ctx context.Context, id string, req dto.CreateSemesterRequest) (*dto.SemesterResponse, error) {
	sem, err := s.schoolRepo.GetSemesterByID(ctx, id)
	if err != nil {
		return nil, domain.NewNotFoundError("semester", id)
	}

	sem.Name = req.Name
	sem.StartDate = req.StartDate
	sem.EndDate = req.EndDate

	if err := s.schoolRepo.UpdateSemester(ctx, sem); err != nil {
		return nil, domain.NewInternalError("failed to update semester", err)
	}

	return &dto.SemesterResponse{
		ID:             sem.ID,
		AcademicYearID: sem.AcademicYearID,
		Name:           sem.Name,
		SemesterNumber: sem.SemesterNumber,
		StartDate:      sem.StartDate,
		EndDate:        sem.EndDate,
		IsActive:       sem.IsActive,
	}, nil
}

func (s *schoolService) ListGrades(ctx context.Context, schoolID string) ([]dto.GradeResponse, error) {
	items, err := s.schoolRepo.ListGrades(ctx, schoolID)
	if err != nil {
		return nil, domain.NewInternalError("failed to list grades", err)
	}

	result := make([]dto.GradeResponse, len(items))
	for i, g := range items {
		result[i] = dto.GradeResponse{
			ID:        g.ID,
			SchoolID:  g.SchoolID,
			Name:      g.Name,
			Level:     g.Level,
			CreatedAt: g.CreatedAt,
			UpdatedAt: g.UpdatedAt,
		}
	}
	return result, nil
}

func (s *schoolService) GetGrade(ctx context.Context, id string) (*dto.GradeResponse, error) {
	g, err := s.schoolRepo.GetGradeByID(ctx, id)
	if err != nil {
		return nil, domain.NewNotFoundError("grade", id)
	}
	return &dto.GradeResponse{
		ID:        g.ID,
		SchoolID:  g.SchoolID,
		Name:      g.Name,
		Level:     g.Level,
		CreatedAt: g.CreatedAt,
		UpdatedAt: g.UpdatedAt,
	}, nil
}

func (s *schoolService) CreateGrade(ctx context.Context, schoolID string, req dto.CreateGradeRequest) (*dto.GradeResponse, error) {
	g := &domain.Grade{
		ID:       uuid.New().String(),
		SchoolID: schoolID,
		Name:     req.Name,
		Level:    req.Level,
	}

	if err := s.schoolRepo.CreateGrade(ctx, g); err != nil {
		return nil, domain.NewInternalError("failed to create grade", err)
	}

	return &dto.GradeResponse{
		ID:        g.ID,
		SchoolID:  g.SchoolID,
		Name:      g.Name,
		Level:     g.Level,
		CreatedAt: g.CreatedAt,
		UpdatedAt: g.UpdatedAt,
	}, nil
}

func (s *schoolService) UpdateGrade(ctx context.Context, id string, req dto.CreateGradeRequest) (*dto.GradeResponse, error) {
	g, err := s.schoolRepo.GetGradeByID(ctx, id)
	if err != nil {
		return nil, domain.NewNotFoundError("grade", id)
	}

	g.Name = req.Name
	g.Level = req.Level

	if err := s.schoolRepo.UpdateGrade(ctx, g); err != nil {
		return nil, domain.NewInternalError("failed to update grade", err)
	}

	return &dto.GradeResponse{
		ID:        g.ID,
		SchoolID:  g.SchoolID,
		Name:      g.Name,
		Level:     g.Level,
		CreatedAt: g.CreatedAt,
		UpdatedAt: g.UpdatedAt,
	}, nil
}

func (s *schoolService) ListClasses(ctx context.Context, schoolID string, filter dto.PaginationRequest) ([]dto.ClassResponse, int64, error) {
	filter.Defaults()
	items, total, err := s.schoolRepo.ListClasses(ctx, schoolID, filter.PageSize, filter.Offset())
	if err != nil {
		return nil, 0, domain.NewInternalError("failed to list classes", err)
	}

	result := make([]dto.ClassResponse, len(items))
	for i, c := range items {
		result[i] = dto.ClassResponse{
			ID:                c.ID,
			SchoolID:          c.SchoolID,
			GradeID:           c.GradeID,
			Name:              c.Name,
			Capacity:          c.Capacity,
			HomeroomTeacherID: c.HomeroomTeacherID,
			AcademicYearID:    c.AcademicYearID,
			CreatedAt:         c.CreatedAt,
			UpdatedAt:         c.UpdatedAt,
		}
	}
	return result, total, nil
}

func (s *schoolService) GetClass(ctx context.Context, id string) (*dto.ClassResponse, error) {
	c, err := s.schoolRepo.GetClassByID(ctx, id)
	if err != nil {
		return nil, domain.NewNotFoundError("class", id)
	}
	return &dto.ClassResponse{
		ID:                c.ID,
		SchoolID:          c.SchoolID,
		GradeID:           c.GradeID,
		Name:              c.Name,
		Capacity:          c.Capacity,
		HomeroomTeacherID: c.HomeroomTeacherID,
		AcademicYearID:    c.AcademicYearID,
		CreatedAt:         c.CreatedAt,
		UpdatedAt:         c.UpdatedAt,
	}, nil
}

func (s *schoolService) CreateClass(ctx context.Context, schoolID string, req dto.CreateClassRequest) (*dto.ClassResponse, error) {
	c := &domain.Class{
		ID:                uuid.New().String(),
		SchoolID:          schoolID,
		GradeID:           req.GradeID,
		Name:              req.Name,
		Capacity:          req.Capacity,
		HomeroomTeacherID: req.HomeroomTeacherID,
		AcademicYearID:    req.AcademicYearID,
	}

	if err := s.schoolRepo.CreateClass(ctx, c); err != nil {
		return nil, domain.NewInternalError("failed to create class", err)
	}

	return &dto.ClassResponse{
		ID:                c.ID,
		SchoolID:          c.SchoolID,
		GradeID:           c.GradeID,
		Name:              c.Name,
		Capacity:          c.Capacity,
		HomeroomTeacherID: c.HomeroomTeacherID,
		AcademicYearID:    c.AcademicYearID,
		CreatedAt:         c.CreatedAt,
		UpdatedAt:         c.UpdatedAt,
	}, nil
}

func (s *schoolService) UpdateClass(ctx context.Context, id string, req dto.UpdateClassRequest) (*dto.ClassResponse, error) {
	c, err := s.schoolRepo.GetClassByID(ctx, id)
	if err != nil {
		return nil, domain.NewNotFoundError("class", id)
	}

	if req.Name != "" {
		c.Name = req.Name
	}
	if req.Capacity != nil {
		c.Capacity = *req.Capacity
	}
	if req.HomeroomTeacherID != nil {
		c.HomeroomTeacherID = req.HomeroomTeacherID
	}

	if err := s.schoolRepo.UpdateClass(ctx, c); err != nil {
		return nil, domain.NewInternalError("failed to update class", err)
	}

	return &dto.ClassResponse{
		ID:                c.ID,
		SchoolID:          c.SchoolID,
		GradeID:           c.GradeID,
		Name:              c.Name,
		Capacity:          c.Capacity,
		HomeroomTeacherID: c.HomeroomTeacherID,
		AcademicYearID:    c.AcademicYearID,
		CreatedAt:         c.CreatedAt,
		UpdatedAt:         c.UpdatedAt,
	}, nil
}

func (s *schoolService) ListSubjects(ctx context.Context, schoolID string, filter dto.PaginationRequest) ([]dto.SubjectResponse, int64, error) {
	filter.Defaults()
	items, total, err := s.schoolRepo.ListSubjects(ctx, schoolID, filter.PageSize, filter.Offset())
	if err != nil {
		return nil, 0, domain.NewInternalError("failed to list subjects", err)
	}

	result := make([]dto.SubjectResponse, len(items))
	for i, sub := range items {
		result[i] = dto.SubjectResponse{
			ID:          sub.ID,
			SchoolID:    sub.SchoolID,
			Code:        sub.Code,
			Name:        sub.Name,
			Category:    sub.Category,
			Description: sub.Description,
			KKM:         sub.KKM,
			IsActive:    sub.IsActive,
			CreatedAt:   sub.CreatedAt,
			UpdatedAt:   sub.UpdatedAt,
		}
	}
	return result, total, nil
}

func (s *schoolService) GetSubject(ctx context.Context, id string) (*dto.SubjectResponse, error) {
	sub, err := s.schoolRepo.GetSubjectByID(ctx, id)
	if err != nil {
		return nil, domain.NewNotFoundError("subject", id)
	}
	return &dto.SubjectResponse{
		ID:          sub.ID,
		SchoolID:    sub.SchoolID,
		Code:        sub.Code,
		Name:        sub.Name,
		Category:    sub.Category,
		Description: sub.Description,
		KKM:         sub.KKM,
		IsActive:    sub.IsActive,
		CreatedAt:   sub.CreatedAt,
		UpdatedAt:   sub.UpdatedAt,
	}, nil
}

func (s *schoolService) CreateSubject(ctx context.Context, schoolID string, req dto.CreateSubjectRequest) (*dto.SubjectResponse, error) {
	sub := &domain.Subject{
		ID:          uuid.New().String(),
		SchoolID:    schoolID,
		Code:        req.Code,
		Name:        req.Name,
		Category:    req.Category,
		Description: req.Description,
		KKM:         req.KKM,
		IsActive:    true,
	}

	if err := s.schoolRepo.CreateSubject(ctx, sub); err != nil {
		return nil, domain.NewInternalError("failed to create subject", err)
	}

	return &dto.SubjectResponse{
		ID:          sub.ID,
		SchoolID:    sub.SchoolID,
		Code:        sub.Code,
		Name:        sub.Name,
		Category:    sub.Category,
		Description: sub.Description,
		KKM:         sub.KKM,
		IsActive:    sub.IsActive,
		CreatedAt:   sub.CreatedAt,
		UpdatedAt:   sub.UpdatedAt,
	}, nil
}

func (s *schoolService) UpdateSubject(ctx context.Context, id string, req dto.UpdateSubjectRequest) (*dto.SubjectResponse, error) {
	sub, err := s.schoolRepo.GetSubjectByID(ctx, id)
	if err != nil {
		return nil, domain.NewNotFoundError("subject", id)
	}

	if req.Code != "" {
		sub.Code = req.Code
	}
	if req.Name != "" {
		sub.Name = req.Name
	}
	if req.Category != "" {
		sub.Category = req.Category
	}
	if req.Description != "" {
		sub.Description = req.Description
	}
	if req.KKM != nil {
		sub.KKM = *req.KKM
	}
	if req.IsActive != nil {
		sub.IsActive = *req.IsActive
	}

	if err := s.schoolRepo.UpdateSubject(ctx, sub); err != nil {
		return nil, domain.NewInternalError("failed to update subject", err)
	}

	return &dto.SubjectResponse{
		ID:          sub.ID,
		SchoolID:    sub.SchoolID,
		Code:        sub.Code,
		Name:        sub.Name,
		Category:    sub.Category,
		Description: sub.Description,
		KKM:         sub.KKM,
		IsActive:    sub.IsActive,
		CreatedAt:   sub.CreatedAt,
		UpdatedAt:   sub.UpdatedAt,
	}, nil
}

func (s *schoolService) ListCurriculums(ctx context.Context, gradeID, subjectID, semesterID string) ([]dto.CurriculumResponse, error) {
	items, err := s.schoolRepo.ListCurriculums(ctx, gradeID, subjectID, semesterID)
	if err != nil {
		return nil, domain.NewInternalError("failed to list curriculums", err)
	}

	result := make([]dto.CurriculumResponse, len(items))
	for i, c := range items {
		result[i] = dto.CurriculumResponse{
			ID:         c.ID,
			SchoolID:   c.SchoolID,
			GradeID:    c.GradeID,
			SubjectID:  c.SubjectID,
			SemesterID: c.SemesterID,
			Content:    c.Content,
			CreatedAt:  c.CreatedAt,
			UpdatedAt:  c.UpdatedAt,
		}
	}
	return result, nil
}

func (s *schoolService) GetCurriculum(ctx context.Context, id string) (*dto.CurriculumResponse, error) {
	c, err := s.schoolRepo.GetCurriculumByID(ctx, id)
	if err != nil {
		return nil, domain.NewNotFoundError("curriculum", id)
	}
	return &dto.CurriculumResponse{
		ID:         c.ID,
		SchoolID:   c.SchoolID,
		GradeID:    c.GradeID,
		SubjectID:  c.SubjectID,
		SemesterID: c.SemesterID,
		Content:    c.Content,
		CreatedAt:  c.CreatedAt,
		UpdatedAt:  c.UpdatedAt,
	}, nil
}

func (s *schoolService) CreateCurriculum(ctx context.Context, schoolID string, req dto.CreateCurriculumRequest) (*dto.CurriculumResponse, error) {
	c := &domain.Curriculum{
		ID:         uuid.New().String(),
		SchoolID:   schoolID,
		GradeID:    req.GradeID,
		SubjectID:  req.SubjectID,
		SemesterID: req.SemesterID,
		Content:    req.Content,
	}

	if err := s.schoolRepo.CreateCurriculum(ctx, c); err != nil {
		return nil, domain.NewInternalError("failed to create curriculum", err)
	}

	return &dto.CurriculumResponse{
		ID:         c.ID,
		SchoolID:   c.SchoolID,
		GradeID:    c.GradeID,
		SubjectID:  c.SubjectID,
		SemesterID: c.SemesterID,
		Content:    c.Content,
		CreatedAt:  c.CreatedAt,
		UpdatedAt:  c.UpdatedAt,
	}, nil
}

func (s *schoolService) UpdateCurriculum(ctx context.Context, id string, req dto.CreateCurriculumRequest) (*dto.CurriculumResponse, error) {
	c, err := s.schoolRepo.GetCurriculumByID(ctx, id)
	if err != nil {
		return nil, domain.NewNotFoundError("curriculum", id)
	}

	c.GradeID = req.GradeID
	c.SubjectID = req.SubjectID
	c.SemesterID = req.SemesterID
	c.Content = req.Content

	if err := s.schoolRepo.UpdateCurriculum(ctx, c); err != nil {
		return nil, domain.NewInternalError("failed to update curriculum", err)
	}

	return &dto.CurriculumResponse{
		ID:         c.ID,
		SchoolID:   c.SchoolID,
		GradeID:    c.GradeID,
		SubjectID:  c.SubjectID,
		SemesterID: c.SemesterID,
		Content:    c.Content,
		CreatedAt:  c.CreatedAt,
		UpdatedAt:  c.UpdatedAt,
	}, nil
}

func toSchoolResponse(s domain.School) dto.SchoolResponse {
	return dto.SchoolResponse{
		ID:             s.ID,
		Name:           s.Name,
		NPSN:           s.NPSN,
		Address:        s.Address,
		City:           s.City,
		Province:       s.Province,
		PostalCode:     s.PostalCode,
		Phone:          s.Phone,
		Email:          s.Email,
		Website:        s.Website,
		LogoURL:        s.LogoURL,
		Type:           s.Type,
		Accreditation:  s.Accreditation,
		IsActive:       s.IsActive,
		CreatedAt:      s.CreatedAt,
		UpdatedAt:      s.UpdatedAt,
	}
}
