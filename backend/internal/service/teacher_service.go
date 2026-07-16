package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/opencode/erp-school-backend/internal/domain"
	"github.com/opencode/erp-school-backend/internal/dto"
	"github.com/opencode/erp-school-backend/internal/infrastructure/database"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type TeacherService interface {
	ListTeachers(ctx context.Context, schoolID string, filter dto.PaginationRequest) ([]dto.TeacherResponse, int64, error)
	GetTeacher(ctx context.Context, id string) (*dto.TeacherResponse, error)
	CreateTeacher(ctx context.Context, schoolID string, req dto.CreateTeacherRequest) (*dto.TeacherResponse, error)
	UpdateTeacher(ctx context.Context, id string, req dto.CreateTeacherRequest) (*dto.TeacherResponse, error)
	AssignSubjects(ctx context.Context, teacherID string, subjectIDs, classIDs []string) error
	GetSchedule(ctx context.Context, teacherID, semesterID string) ([]dto.ScheduleResponse, error)
}

type teacherService struct {
	db     *sqlx.DB
	logger *zap.Logger
}

func NewTeacherService(db *sqlx.DB, logger *zap.Logger) TeacherService {
	return &teacherService{db: db, logger: logger}
}

func (s *teacherService) ListTeachers(ctx context.Context, schoolID string, filter dto.PaginationRequest) ([]dto.TeacherResponse, int64, error) {
	filter.Defaults()

	var total int64
	if err := s.db.GetContext(ctx, &total, database.CountTeachers, schoolID); err != nil {
		return nil, 0, domain.NewInternalError("failed to count teachers", err)
	}

	var items []struct {
		domain.Teacher
		FullName string `db:"full_name"`
		Email    string `db:"email"`
		Phone    string `db:"phone"`
	}

	if err := s.db.SelectContext(ctx, &items, database.ListTeachers, schoolID, filter.PageSize, filter.Offset()); err != nil {
		return nil, 0, domain.NewInternalError("failed to list teachers", err)
	}

	result := make([]dto.TeacherResponse, len(items))
	for i, t := range items {
		result[i] = dto.TeacherResponse{
			ID:             t.ID,
			UserID:         t.UserID,
			SchoolID:       t.SchoolID,
			Email:          t.Email,
			FullName:       t.FullName,
			NIP:            t.NIP,
			NIK:            t.NIK,
			NUPTK:          t.NUPTK,
			Status:         t.Status,
			JoinDate:       t.JoinDate,
			EducationLevel: t.EducationLevel,
			Major:          t.Major,
			Phone:          t.Phone,
			CreatedAt:      t.CreatedAt,
			UpdatedAt:      t.UpdatedAt,
		}
	}
	return result, total, nil
}

func (s *teacherService) GetTeacher(ctx context.Context, id string) (*dto.TeacherResponse, error) {
	var t domain.Teacher
	query := `SELECT * FROM teachers WHERE id=$1`
	if err := s.db.GetContext(ctx, &t, query, id); err != nil {
		return nil, domain.NewNotFoundError("teacher", id)
	}

	var u domain.User
	if err := s.db.GetContext(ctx, &u, `SELECT * FROM users WHERE id=$1`, t.UserID); err != nil {
		return nil, domain.NewInternalError("failed to get teacher user", err)
	}

	var subjectCount int
	s.db.GetContext(ctx, &subjectCount, `SELECT COUNT(*) FROM teacher_subjects WHERE teacher_id=$1`, t.ID)

	return &dto.TeacherResponse{
		ID:             t.ID,
		UserID:         t.UserID,
		SchoolID:       t.SchoolID,
		Email:          u.Email,
		FullName:       u.FullName,
		NIP:            t.NIP,
		NIK:            t.NIK,
		NUPTK:          t.NUPTK,
		Status:         t.Status,
		JoinDate:       t.JoinDate,
		EducationLevel: t.EducationLevel,
		Major:          t.Major,
		Phone:          u.Phone,
		SubjectsCount:  subjectCount,
		CreatedAt:      t.CreatedAt,
		UpdatedAt:      t.UpdatedAt,
	}, nil
}

func (s *teacherService) CreateTeacher(ctx context.Context, schoolID string, req dto.CreateTeacherRequest) (*dto.TeacherResponse, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, domain.NewInternalError("failed to hash password", err)
	}

	userID := req.UserID
	if userID == "" {
		userID = uuid.New().String()
	}

	userQuery := `INSERT INTO users (id, school_id, email, username, password_hash, full_name, phone) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	if _, err := s.db.ExecContext(ctx, userQuery, userID, schoolID, req.Email, req.Email, string(hash), req.FullName, ""); err != nil {
		return nil, domain.NewInternalError("failed to create user for teacher", err)
	}

	teacherID := uuid.New().String()
	teacherQuery := `INSERT INTO teachers (id, user_id, school_id, nip, nik, nupk, status, join_date, education_level, major) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	if _, err := s.db.ExecContext(ctx, teacherQuery, teacherID, userID, schoolID, req.NIP, req.NIK, req.NUPTK, req.Status, req.JoinDate, req.EducationLevel, req.Major); err != nil {
		return nil, domain.NewInternalError("failed to create teacher", err)
	}

	return s.GetTeacher(ctx, teacherID)
}

func (s *teacherService) UpdateTeacher(ctx context.Context, id string, req dto.CreateTeacherRequest) (*dto.TeacherResponse, error) {
	var t domain.Teacher
	if err := s.db.GetContext(ctx, &t, `SELECT * FROM teachers WHERE id=$1`, id); err != nil {
		return nil, domain.NewNotFoundError("teacher", id)
	}

	query := `UPDATE teachers SET nip=$1, nik=$2, nupk=$3, status=$4, education_level=$5, major=$6, updated_at=NOW() WHERE id=$7`
	if _, err := s.db.ExecContext(ctx, query, req.NIP, req.NIK, req.NUPTK, req.Status, req.EducationLevel, req.Major, id); err != nil {
		return nil, domain.NewInternalError("failed to update teacher", err)
	}

	return s.GetTeacher(ctx, id)
}

func (s *teacherService) AssignSubjects(ctx context.Context, teacherID string, subjectIDs, classIDs []string) error {
	tx, err := s.db.Beginx()
	if err != nil {
		return domain.NewInternalError("failed to begin transaction", err)
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `DELETE FROM teacher_subjects WHERE teacher_id=$1`, teacherID); err != nil {
		return domain.NewInternalError("failed to clear subjects", err)
	}

	for i := range subjectIDs {
		tsID := uuid.New().String()
		classID := classIDs[i]
		query := `INSERT INTO teacher_subjects (id, teacher_id, subject_id, class_id) VALUES ($1, $2, $3, $4)`
		if _, err := tx.ExecContext(ctx, query, tsID, teacherID, subjectIDs[i], classID); err != nil {
			return domain.NewInternalError("failed to assign subject", err)
		}
	}

	return tx.Commit()
}

func (s *teacherService) GetSchedule(ctx context.Context, teacherID, semesterID string) ([]dto.ScheduleResponse, error) {
	query := `SELECT sc.*, c.name as class_name, sub.name as subject_name, u.full_name as teacher_name FROM schedules sc JOIN classes c ON sc.class_id = c.id JOIN subjects sub ON sc.subject_id = sub.id JOIN teachers t ON sc.teacher_id = t.id JOIN users u ON t.user_id = u.id WHERE sc.teacher_id=$1 AND sc.semester_id=$2`

	var items []struct {
		domain.Schedule
		ClassName   string `db:"class_name"`
		SubjectName string `db:"subject_name"`
		TeacherName string `db:"teacher_name"`
	}

	if err := s.db.SelectContext(ctx, &items, query, teacherID, semesterID); err != nil {
		return nil, domain.NewInternalError("failed to get schedule", err)
	}

	result := make([]dto.ScheduleResponse, len(items))
	for i, sc := range items {
		result[i] = dto.ScheduleResponse{
			ID:          sc.ID,
			ClassID:     sc.ClassID,
			ClassName:   sc.ClassName,
			SubjectID:   sc.SubjectID,
			SubjectName: sc.SubjectName,
			TeacherID:   sc.TeacherID,
			TeacherName: sc.TeacherName,
			Day:         sc.Day,
			StartTime:   sc.StartTime,
			EndTime:     sc.EndTime,
			Room:        sc.Room,
			SemesterID:  sc.SemesterID,
		}
	}
	return result, nil
}
