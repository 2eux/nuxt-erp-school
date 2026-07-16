package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/opencode/erp-school-backend/internal/domain"
	"github.com/opencode/erp-school-backend/internal/dto"
	"github.com/opencode/erp-school-backend/internal/infrastructure/database"
	"go.uber.org/zap"
)

type AcademicService interface {
	ListSchedules(ctx context.Context, semesterID string) ([]dto.ScheduleResponse, error)
	CreateSchedule(ctx context.Context, req dto.CreateScheduleRequest) (*dto.ScheduleResponse, error)
	GetSchedule(ctx context.Context, id string) (*dto.ScheduleResponse, error)

	ListAttendances(ctx context.Context, studentID string, startDate, endDate time.Time) ([]dto.AttendanceResponse, error)
	CreateAttendance(ctx context.Context, createdBy string, req dto.CreateAttendanceRequest) (*dto.AttendanceResponse, error)
	BulkCreateAttendance(ctx context.Context, createdBy, scheduleID string, req dto.BatchAttendanceRequest) error

	ListExams(ctx context.Context, classID, semesterID string) ([]dto.ExamResponse, error)
	CreateExam(ctx context.Context, createdBy string, req dto.CreateExamRequest) (*dto.ExamResponse, error)
	GetExamResults(ctx context.Context, examID string) ([]dto.ExamResultResponse, error)
	CreateExamResult(ctx context.Context, req dto.CreateExamResultRequest) (*dto.ExamResultResponse, error)

	ListGradebooks(ctx context.Context, classID, subjectID, semesterID string) ([]dto.GradebookResponse, error)
	CreateGradebook(ctx context.Context, teacherID string, req dto.CreateGradebookRequest) (*dto.GradebookResponse, error)

	ListReportCards(ctx context.Context, classID, semesterID string) ([]dto.ReportCardResponse, error)
	GenerateReportCard(ctx context.Context, userID string, req dto.CreateReportCardRequest) (*dto.ReportCardResponse, error)

	ListAssignments(ctx context.Context, classID, subjectID string) ([]dto.AssignmentResponse, error)
	CreateAssignment(ctx context.Context, teacherID string, req dto.CreateAssignmentRequest) (*dto.AssignmentResponse, error)
	ListSubmissions(ctx context.Context, assignmentID string) ([]dto.AssignmentSubmissionResponse, error)
	SubmitAssignment(ctx context.Context, studentID, assignmentID string, req dto.SubmitAssignmentRequest) (*dto.AssignmentSubmissionResponse, error)

	ListLessonPlans(ctx context.Context, teacherID string, filter dto.PaginationRequest) ([]dto.LessonPlanResponse, int64, error)
	CreateLessonPlan(ctx context.Context, teacherID string, req dto.CreateLessonPlanRequest) (*dto.LessonPlanResponse, error)

	ListTeachingJournals(ctx context.Context, teacherID string) ([]dto.TeachingJournalResponse, error)
	CreateTeachingJournal(ctx context.Context, teacherID string, req dto.CreateTeachingJournalRequest) (*dto.TeachingJournalResponse, error)
}

type academicService struct {
	db     *sqlx.DB
	logger *zap.Logger
}

func NewAcademicService(db *sqlx.DB, logger *zap.Logger) AcademicService {
	return &academicService{db: db, logger: logger}
}

func (s *academicService) ListSchedules(ctx context.Context, semesterID string) ([]dto.ScheduleResponse, error) {
	var items []struct {
		domain.Schedule
		ClassName   string `db:"class_name"`
		SubjectName string `db:"subject_name"`
		TeacherName string `db:"teacher_name"`
	}

	if err := s.db.SelectContext(ctx, &items, database.ListSchedules, semesterID); err != nil {
		return nil, domain.NewInternalError("failed to list schedules", err)
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

func (s *academicService) CreateSchedule(ctx context.Context, req dto.CreateScheduleRequest) (*dto.ScheduleResponse, error) {
	sc := &domain.Schedule{
		ID:         uuid.New().String(),
		ClassID:    req.ClassID,
		SubjectID:  req.SubjectID,
		TeacherID:  req.TeacherID,
		Day:        req.Day,
		StartTime:  req.StartTime,
		EndTime:    req.EndTime,
		Room:       req.Room,
		SemesterID: req.SemesterID,
	}

	query := `INSERT INTO schedules (id, class_id, subject_id, teacher_id, day, start_time, end_time, room, semester_id) VALUES (:id, :class_id, :subject_id, :teacher_id, :day, :start_time, :end_time, :room, :semester_id)`
	if _, err := s.db.NamedExecContext(ctx, query, sc); err != nil {
		return nil, domain.NewInternalError("failed to create schedule", err)
	}

	return &dto.ScheduleResponse{
		ID:         sc.ID,
		ClassID:    sc.ClassID,
		SubjectID:  sc.SubjectID,
		TeacherID:  sc.TeacherID,
		Day:        sc.Day,
		StartTime:  sc.StartTime,
		EndTime:    sc.EndTime,
		Room:       sc.Room,
		SemesterID: sc.SemesterID,
	}, nil
}

func (s *academicService) GetSchedule(ctx context.Context, id string) (*dto.ScheduleResponse, error) {
	var sc struct {
		domain.Schedule
		ClassName   string `db:"class_name"`
		SubjectName string `db:"subject_name"`
		TeacherName string `db:"teacher_name"`
	}

	query := `SELECT sc.*, c.name as class_name, sub.name as subject_name, u.full_name as teacher_name FROM schedules sc JOIN classes c ON sc.class_id = c.id JOIN subjects sub ON sc.subject_id = sub.id JOIN teachers t ON sc.teacher_id = t.id JOIN users u ON t.user_id = u.id WHERE sc.id=$1`
	if err := s.db.GetContext(ctx, &sc, query, id); err != nil {
		return nil, domain.NewNotFoundError("schedule", id)
	}

	return &dto.ScheduleResponse{
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
	}, nil
}

func (s *academicService) ListAttendances(ctx context.Context, studentID string, startDate, endDate time.Time) ([]dto.AttendanceResponse, error) {
	var items []struct {
		domain.Attendance
		StudentName string `db:"student_name"`
		SubjectName string `db:"subject_name"`
	}

	query := database.ListStudents + ` AND s.id=$1` + ` ORDER BY a.date DESC`
	_ = query

	var atts []domain.Attendance
	attQuery := `SELECT a.* FROM attendances a WHERE a.student_id=$1 AND a.date >= $2 AND a.date <= $3 ORDER BY a.date DESC`
	if err := s.db.SelectContext(ctx, &atts, attQuery, studentID, startDate, endDate); err != nil {
		return nil, domain.NewInternalError("failed to list attendances", err)
	}

	result := make([]dto.AttendanceResponse, len(atts))
	for i, a := range atts {
		result[i] = dto.AttendanceResponse{
			ID:         a.ID,
			StudentID:  a.StudentID,
			ScheduleID: a.ScheduleID,
			Date:       a.Date,
			Status:     a.Status,
			Notes:      a.Notes,
			CreatedBy:  a.CreatedBy,
		}
	}
	_ = items
	return result, nil
}

func (s *academicService) CreateAttendance(ctx context.Context, createdBy string, req dto.CreateAttendanceRequest) (*dto.AttendanceResponse, error) {
	att := &domain.Attendance{
		ID:         uuid.New().String(),
		StudentID:  req.StudentID,
		ScheduleID: req.ScheduleID,
		Date:       req.Date,
		Status:     req.Status,
		Notes:      req.Notes,
		CreatedBy:  createdBy,
	}

	query := `INSERT INTO attendances (id, student_id, schedule_id, date, status, notes, created_by) VALUES (:id, :student_id, :schedule_id, :date, :status, :notes, :created_by)`
	if _, err := s.db.NamedExecContext(ctx, query, att); err != nil {
		return nil, domain.NewInternalError("failed to create attendance", err)
	}

	return &dto.AttendanceResponse{
		ID:         att.ID,
		StudentID:  att.StudentID,
		ScheduleID: att.ScheduleID,
		Date:       att.Date,
		Status:     att.Status,
		Notes:      att.Notes,
		CreatedBy:  att.CreatedBy,
	}, nil
}

func (s *academicService) BulkCreateAttendance(ctx context.Context, createdBy, scheduleID string, req dto.BatchAttendanceRequest) error {
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return domain.NewInvalidInputError("invalid date format")
	}

	tx, err := s.db.Beginx()
	if err != nil {
		return domain.NewInternalError("failed to begin tx", err)
	}
	defer tx.Rollback()

	for _, a := range req.Attendances {
		id := uuid.New().String()
		query := `INSERT INTO attendances (id, student_id, schedule_id, date, status, notes, created_by) VALUES ($1, $2, $3, $4, $5, $6, $7)`
		tx.ExecContext(ctx, query, id, a.StudentID, scheduleID, date, a.Status, a.Notes, createdBy)
	}

	return tx.Commit()
}

func (s *academicService) ListExams(ctx context.Context, classID, semesterID string) ([]dto.ExamResponse, error) {
	var items []struct {
		domain.Exam
		SubjectName string `db:"subject_name"`
		ClassName   string `db:"class_name"`
	}

	if err := s.db.SelectContext(ctx, &items, database.ListExams, classID, semesterID); err != nil {
		return nil, domain.NewInternalError("failed to list exams", err)
	}

	result := make([]dto.ExamResponse, len(items))
	for i, e := range items {
		result[i] = dto.ExamResponse{
			ID:          e.ID,
			SubjectID:   e.SubjectID,
			SubjectName: e.SubjectName,
			ClassID:     e.ClassID,
			ClassName:   e.ClassName,
			SemesterID:  e.SemesterID,
			Name:        e.Name,
			ExamType:    e.ExamType,
			Date:        e.Date,
			Duration:    e.Duration,
			TotalScore:  e.TotalScore,
		}
	}
	return result, nil
}

func (s *academicService) CreateExam(ctx context.Context, createdBy string, req dto.CreateExamRequest) (*dto.ExamResponse, error) {
	exam := &domain.Exam{
		ID:         uuid.New().String(),
		SubjectID:  req.SubjectID,
		ClassID:    req.ClassID,
		SemesterID: req.SemesterID,
		Name:       req.Name,
		ExamType:   req.ExamType,
		Date:       req.Date,
		Duration:   req.Duration,
		TotalScore: req.TotalScore,
		CreatedBy:  createdBy,
	}

	query := `INSERT INTO exams (id, subject_id, class_id, semester_id, name, exam_type, date, duration, total_score, created_by) VALUES (:id, :subject_id, :class_id, :semester_id, :name, :exam_type, :date, :duration, :total_score, :created_by)`
	if _, err := s.db.NamedExecContext(ctx, query, exam); err != nil {
		return nil, domain.NewInternalError("failed to create exam", err)
	}

	return &dto.ExamResponse{
		ID:         exam.ID,
		SubjectID:  exam.SubjectID,
		ClassID:    exam.ClassID,
		SemesterID: exam.SemesterID,
		Name:       exam.Name,
		ExamType:   exam.ExamType,
		Date:       exam.Date,
		Duration:   exam.Duration,
		TotalScore: exam.TotalScore,
	}, nil
}

func (s *academicService) GetExamResults(ctx context.Context, examID string) ([]dto.ExamResultResponse, error) {
	var items []struct {
		domain.ExamResult
		StudentName string `db:"student_name"`
	}

	query := `SELECT er.*, u.full_name as student_name FROM exam_results er JOIN students st ON er.student_id = st.id JOIN users u ON st.user_id = u.id WHERE er.exam_id=$1`
	if err := s.db.SelectContext(ctx, &items, query, examID); err != nil {
		return nil, domain.NewInternalError("failed to get exam results", err)
	}

	result := make([]dto.ExamResultResponse, len(items))
	for i, r := range items {
		result[i] = dto.ExamResultResponse{
			ID:          r.ID,
			ExamID:      r.ExamID,
			StudentID:   r.StudentID,
			StudentName: r.StudentName,
			Score:       r.Score,
			Grade:       r.Grade,
			Notes:       r.Notes,
		}
	}
	return result, nil
}

func (s *academicService) CreateExamResult(ctx context.Context, req dto.CreateExamResultRequest) (*dto.ExamResultResponse, error) {
	result := &domain.ExamResult{
		ID:        uuid.New().String(),
		ExamID:    req.ExamID,
		StudentID: req.StudentID,
		Score:     req.Score,
		Grade:     req.Grade,
		Notes:     req.Notes,
	}

	query := `INSERT INTO exam_results (id, exam_id, student_id, score, grade, notes) VALUES (:id, :exam_id, :student_id, :score, :grade, :notes)`
	if _, err := s.db.NamedExecContext(ctx, query, result); err != nil {
		return nil, domain.NewInternalError("failed to create exam result", err)
	}

	return &dto.ExamResultResponse{
		ID:        result.ID,
		ExamID:    result.ExamID,
		StudentID: result.StudentID,
		Score:     result.Score,
		Grade:     result.Grade,
		Notes:     result.Notes,
	}, nil
}

func (s *academicService) ListGradebooks(ctx context.Context, classID, subjectID, semesterID string) ([]dto.GradebookResponse, error) {
	var items []struct {
		domain.Gradebook
		SubjectName string `db:"subject_name"`
		StudentName string `db:"student_name"`
	}

	if err := s.db.SelectContext(ctx, &items, database.ListGradebooks, classID, subjectID, semesterID); err != nil {
		return nil, domain.NewInternalError("failed to list gradebooks", err)
	}

	result := make([]dto.GradebookResponse, len(items))
	for i, gb := range items {
		result[i] = dto.GradebookResponse{
			ID:            gb.ID,
			ClassID:       gb.ClassID,
			SubjectID:     gb.SubjectID,
			SubjectName:   gb.SubjectName,
			StudentID:     gb.StudentID,
			StudentName:   gb.StudentName,
			SemesterID:    gb.SemesterID,
			DailyScore:    gb.DailyScore,
			MidScore:      gb.MidScore,
			FinalScore:    gb.FinalScore,
			PracticeScore: gb.PracticeScore,
			Attitude:      gb.Attitude,
			Notes:         gb.Notes,
		}
	}
	return result, nil
}

func (s *academicService) CreateGradebook(ctx context.Context, teacherID string, req dto.CreateGradebookRequest) (*dto.GradebookResponse, error) {
	gb := &domain.Gradebook{
		ID:            uuid.New().String(),
		ClassID:       req.ClassID,
		SubjectID:     req.SubjectID,
		TeacherID:     teacherID,
		SemesterID:    req.SemesterID,
		StudentID:     req.StudentID,
		DailyScore:    req.DailyScore,
		MidScore:      req.MidScore,
		FinalScore:    req.FinalScore,
		PracticeScore: req.PracticeScore,
		Attitude:      req.Attitude,
		Notes:         req.Notes,
	}

	query := `INSERT INTO gradebooks (id, class_id, subject_id, teacher_id, semester_id, student_id, daily_score, mid_score, final_score, practice_score, attitude, notes) VALUES (:id, :class_id, :subject_id, :teacher_id, :semester_id, :student_id, :daily_score, :mid_score, :final_score, :practice_score, :attitude, :notes)`
	if _, err := s.db.NamedExecContext(ctx, query, gb); err != nil {
		return nil, domain.NewInternalError("failed to create gradebook", err)
	}

	return &dto.GradebookResponse{
		ID:            gb.ID,
		ClassID:       gb.ClassID,
		SubjectID:     gb.SubjectID,
		StudentID:     gb.StudentID,
		SemesterID:    gb.SemesterID,
		DailyScore:    gb.DailyScore,
		MidScore:      gb.MidScore,
		FinalScore:    gb.FinalScore,
		PracticeScore: gb.PracticeScore,
		Attitude:      gb.Attitude,
		Notes:         gb.Notes,
	}, nil
}

func (s *academicService) ListReportCards(ctx context.Context, classID, semesterID string) ([]dto.ReportCardResponse, error) {
	var items []domain.ReportCard
	query := `SELECT * FROM report_cards WHERE class_id=$1 AND semester_id=$2`
	if err := s.db.SelectContext(ctx, &items, query, classID, semesterID); err != nil {
		return nil, domain.NewInternalError("failed to list report cards", err)
	}

	result := make([]dto.ReportCardResponse, len(items))
	for i, rc := range items {
		result[i] = dto.ReportCardResponse{
			ID:              rc.ID,
			StudentID:       rc.StudentID,
			SemesterID:      rc.SemesterID,
			ClassName:        "",
			AverageScore:    rc.AverageScore,
			Rank:            rc.Rank,
			AbsentCount:     rc.AbsentCount,
			SickCount:       rc.SickCount,
			PermitCount:     rc.PermitCount,
			HomeroomComment: rc.HomeroomComment,
			ParentSignature: rc.ParentSignature,
			PublishedAt:     rc.PublishedAt,
		}
	}
	return result, nil
}

func (s *academicService) GenerateReportCard(ctx context.Context, userID string, req dto.CreateReportCardRequest) (*dto.ReportCardResponse, error) {
	rc := &domain.ReportCard{
		ID:              uuid.New().String(),
		StudentID:       req.StudentID,
		SemesterID:      req.SemesterID,
		ClassID:         req.ClassID,
		AbsentCount:     req.AbsentCount,
		SickCount:       req.SickCount,
		PermitCount:     req.PermitCount,
		HomeroomComment: req.HomeroomComment,
		ApprovedBy:      userID,
	}

	query := `INSERT INTO report_cards (id, student_id, semester_id, class_id, absent_count, sick_count, permit_count, homeroom_comment, approved_by) VALUES (:id, :student_id, :semester_id, :class_id, :absent_count, :sick_count, :permit_count, :homeroom_comment, :approved_by)`
	if _, err := s.db.NamedExecContext(ctx, query, rc); err != nil {
		return nil, domain.NewInternalError("failed to generate report card", err)
	}

	return &dto.ReportCardResponse{
		ID:              rc.ID,
		StudentID:       rc.StudentID,
		SemesterID:      rc.SemesterID,
		AbsentCount:     rc.AbsentCount,
		SickCount:       rc.SickCount,
		PermitCount:     rc.PermitCount,
		HomeroomComment: rc.HomeroomComment,
	}, nil
}

func (s *academicService) ListAssignments(ctx context.Context, classID, subjectID string) ([]dto.AssignmentResponse, error) {
	var items []domain.Assignment
	query := `SELECT * FROM assignments WHERE class_id=$1 AND subject_id=$2 ORDER BY created_at DESC`
	if err := s.db.SelectContext(ctx, &items, query, classID, subjectID); err != nil {
		return nil, domain.NewInternalError("failed to list assignments", err)
	}

	result := make([]dto.AssignmentResponse, len(items))
	for i, a := range items {
		result[i] = dto.AssignmentResponse{
			ID:          a.ID,
			SubjectID:   a.SubjectID,
			ClassID:     a.ClassID,
			Title:       a.Title,
			Description: a.Description,
			DueDate:     a.DueDate,
			MaxScore:    a.MaxScore,
			CreatedAt:   a.CreatedAt,
		}
	}
	return result, nil
}

func (s *academicService) CreateAssignment(ctx context.Context, teacherID string, req dto.CreateAssignmentRequest) (*dto.AssignmentResponse, error) {
	a := &domain.Assignment{
		ID:          uuid.New().String(),
		SubjectID:   req.SubjectID,
		ClassID:     req.ClassID,
		TeacherID:   teacherID,
		Title:       req.Title,
		Description: req.Description,
		DueDate:     req.DueDate,
		MaxScore:    req.MaxScore,
	}

	query := `INSERT INTO assignments (id, subject_id, class_id, teacher_id, title, description, due_date, max_score) VALUES (:id, :subject_id, :class_id, :teacher_id, :title, :description, :due_date, :max_score)`
	if _, err := s.db.NamedExecContext(ctx, query, a); err != nil {
		return nil, domain.NewInternalError("failed to create assignment", err)
	}

	return &dto.AssignmentResponse{
		ID:          a.ID,
		SubjectID:   a.SubjectID,
		ClassID:     a.ClassID,
		Title:       a.Title,
		Description: a.Description,
		DueDate:     a.DueDate,
		MaxScore:    a.MaxScore,
		CreatedAt:   a.CreatedAt,
	}, nil
}

func (s *academicService) ListSubmissions(ctx context.Context, assignmentID string) ([]dto.AssignmentSubmissionResponse, error) {
	var items []struct {
		domain.AssignmentSubmission
		StudentName string `db:"student_name"`
	}

	query := `SELECT ass.*, u.full_name as student_name FROM assignment_submissions ass JOIN students st ON ass.student_id = st.id JOIN users u ON st.user_id = u.id WHERE ass.assignment_id=$1 ORDER BY ass.submitted_at DESC`
	if err := s.db.SelectContext(ctx, &items, query, assignmentID); err != nil {
		return nil, domain.NewInternalError("failed to list submissions", err)
	}

	result := make([]dto.AssignmentSubmissionResponse, len(items))
	for i, sub := range items {
		result[i] = dto.AssignmentSubmissionResponse{
			ID:           sub.ID,
			AssignmentID: sub.AssignmentID,
			StudentID:    sub.StudentID,
			StudentName:  sub.StudentName,
			Content:      sub.Content,
			FileURL:      sub.FileURL,
			Score:        sub.Score,
			Feedback:     sub.Feedback,
			SubmittedAt:  sub.SubmittedAt,
			GradedAt:     sub.GradedAt,
		}
	}
	return result, nil
}

func (s *academicService) SubmitAssignment(ctx context.Context, studentID, assignmentID string, req dto.SubmitAssignmentRequest) (*dto.AssignmentSubmissionResponse, error) {
	sub := &domain.AssignmentSubmission{
		ID:           uuid.New().String(),
		AssignmentID: assignmentID,
		StudentID:    studentID,
		Content:      req.Content,
		SubmittedAt:  time.Now(),
	}

	query := `INSERT INTO assignment_submissions (id, assignment_id, student_id, content, submitted_at) VALUES (:id, :assignment_id, :student_id, :content, :submitted_at)`
	if _, err := s.db.NamedExecContext(ctx, query, sub); err != nil {
		return nil, domain.NewInternalError("failed to submit assignment", err)
	}

	return &dto.AssignmentSubmissionResponse{
		ID:           sub.ID,
		AssignmentID: sub.AssignmentID,
		StudentID:    sub.StudentID,
		Content:      sub.Content,
		SubmittedAt:  sub.SubmittedAt,
	}, nil
}

func (s *academicService) ListLessonPlans(ctx context.Context, teacherID string, filter dto.PaginationRequest) ([]dto.LessonPlanResponse, int64, error) {
	filter.Defaults()

	var items []domain.LessonPlan
	var total int64

	s.db.GetContext(ctx, &total, `SELECT COUNT(*) FROM lesson_plans WHERE teacher_id=$1`, teacherID)
	query := `SELECT * FROM lesson_plans WHERE teacher_id=$1 ORDER BY date DESC LIMIT $2 OFFSET $3`
	if err := s.db.SelectContext(ctx, &items, query, teacherID, filter.PageSize, filter.Offset()); err != nil {
		return nil, 0, domain.NewInternalError("failed to list lesson plans", err)
	}

	result := make([]dto.LessonPlanResponse, len(items))
	for i, lp := range items {
		result[i] = dto.LessonPlanResponse{
			ID:         lp.ID,
			SubjectID:  lp.SubjectID,
			ClassID:    lp.ClassID,
			Title:      lp.Title,
			Date:       lp.Date,
			Objectives: lp.Objectives,
			Materials:  lp.Materials,
			Activities: lp.Activities,
			Assessment: lp.Assessment,
			Reflection: lp.Reflection,
			Status:     lp.Status,
			ApprovedBy: lp.ApprovedBy,
			ApprovedAt: lp.ApprovedAt,
			CreatedAt:  lp.CreatedAt,
		}
	}
	return result, total, nil
}

func (s *academicService) CreateLessonPlan(ctx context.Context, teacherID string, req dto.CreateLessonPlanRequest) (*dto.LessonPlanResponse, error) {
	lp := &domain.LessonPlan{
		ID:         uuid.New().String(),
		SubjectID:  req.SubjectID,
		ClassID:    req.ClassID,
		TeacherID:  teacherID,
		Title:      req.Title,
		Date:       req.Date,
		Objectives: req.Objectives,
		Materials:  req.Materials,
		Activities: req.Activities,
		Assessment: req.Assessment,
		Reflection: req.Reflection,
		Status:     "draft",
	}

	query := `INSERT INTO lesson_plans (id, subject_id, class_id, teacher_id, title, date, objectives, materials, activities, assessment, reflection, status) VALUES (:id, :subject_id, :class_id, :teacher_id, :title, :date, :objectives, :materials, :activities, :assessment, :reflection, :status)`
	if _, err := s.db.NamedExecContext(ctx, query, lp); err != nil {
		return nil, domain.NewInternalError("failed to create lesson plan", err)
	}

	return &dto.LessonPlanResponse{
		ID:         lp.ID,
		SubjectID:  lp.SubjectID,
		ClassID:    lp.ClassID,
		Title:      lp.Title,
		Date:       lp.Date,
		Objectives: lp.Objectives,
		Materials:  lp.Materials,
		Activities: lp.Activities,
		Assessment: lp.Assessment,
		Reflection: lp.Reflection,
		Status:     lp.Status,
		CreatedAt:  lp.CreatedAt,
	}, nil
}

func (s *academicService) ListTeachingJournals(ctx context.Context, teacherID string) ([]dto.TeachingJournalResponse, error) {
	var items []domain.TeachingJournal
	query := `SELECT * FROM teaching_journals WHERE teacher_id=$1 ORDER BY date DESC`
	if err := s.db.SelectContext(ctx, &items, query, teacherID); err != nil {
		return nil, domain.NewInternalError("failed to list teaching journals", err)
	}

	result := make([]dto.TeachingJournalResponse, len(items))
	for i, tj := range items {
		result[i] = dto.TeachingJournalResponse{
			ID:          tj.ID,
			TeacherID:   tj.TeacherID,
			ScheduleID:  tj.ScheduleID,
			Date:        tj.Date,
			Material:    tj.Material,
			Method:      tj.Method,
			AttendCount: tj.AttendCount,
			Notes:       tj.Notes,
		}
	}
	return result, nil
}

func (s *academicService) CreateTeachingJournal(ctx context.Context, teacherID string, req dto.CreateTeachingJournalRequest) (*dto.TeachingJournalResponse, error) {
	tj := &domain.TeachingJournal{
		ID:         uuid.New().String(),
		TeacherID:  teacherID,
		ScheduleID: req.ScheduleID,
		Date:       req.Date,
		Material:   req.Material,
		Method:     req.Method,
		AttendCount: req.AttendCount,
		Notes:      req.Notes,
	}

	query := `INSERT INTO teaching_journals (id, teacher_id, schedule_id, date, material, method, attend_count, notes) VALUES (:id, :teacher_id, :schedule_id, :date, :material, :method, :attend_count, :notes)`
	if _, err := s.db.NamedExecContext(ctx, query, tj); err != nil {
		return nil, domain.NewInternalError("failed to create teaching journal", err)
	}

	return &dto.TeachingJournalResponse{
		ID:          tj.ID,
		TeacherID:   tj.TeacherID,
		ScheduleID:  tj.ScheduleID,
		Date:        tj.Date,
		Material:    tj.Material,
		Method:      tj.Method,
		AttendCount: tj.AttendCount,
		Notes:       tj.Notes,
	}, nil
}
