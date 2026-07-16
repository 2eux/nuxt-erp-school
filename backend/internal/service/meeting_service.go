package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/opencode/erp-school-backend/internal/domain"
	"github.com/opencode/erp-school-backend/internal/dto"
	"go.uber.org/zap"
)

type MeetingService interface {
	ListMeetings(ctx context.Context, schoolID string, filter dto.PaginationRequest) ([]dto.MeetingResponse, int64, error)
	GetMeeting(ctx context.Context, id string) (*dto.MeetingResponse, error)
	CreateMeeting(ctx context.Context, schoolID, createdBy string, req dto.CreateMeetingRequest) (*dto.MeetingResponse, error)
	UpdateMeetingMinutes(ctx context.Context, id string, minutes string) error

	ListEvents(ctx context.Context, schoolID string, filter dto.PaginationRequest) ([]dto.EventResponse, int64, error)
	CreateEvent(ctx context.Context, schoolID, createdBy string, req dto.CreateEventRequest) (*dto.EventResponse, error)

	ListTasks(ctx context.Context, schoolID string, filter dto.PaginationRequest) ([]dto.TaskResponse, int64, error)
	CreateTask(ctx context.Context, schoolID, createdBy string, req dto.CreateTaskRequest) (*dto.TaskResponse, error)
	UpdateTask(ctx context.Context, id string, req dto.UpdateTaskRequest) (*dto.TaskResponse, error)
}

type meetingService struct {
	db     *sqlx.DB
	logger *zap.Logger
}

func NewMeetingService(db *sqlx.DB, logger *zap.Logger) MeetingService {
	return &meetingService{db: db, logger: logger}
}

func (s *meetingService) ListMeetings(ctx context.Context, schoolID string, filter dto.PaginationRequest) ([]dto.MeetingResponse, int64, error) {
	filter.Defaults()
	var items []domain.Meeting
	var total int64

	s.db.GetContext(ctx, &total, `SELECT COUNT(*) FROM meetings WHERE school_id=$1`, schoolID)
	query := `SELECT * FROM meetings WHERE school_id=$1 ORDER BY date DESC LIMIT $2 OFFSET $3`
	if err := s.db.SelectContext(ctx, &items, query, schoolID, filter.PageSize, filter.Offset()); err != nil {
		return nil, 0, domain.NewInternalError("failed to list meetings", err)
	}

	result := make([]dto.MeetingResponse, len(items))
	for i, m := range items {
		result[i] = dto.MeetingResponse{
			ID:        m.ID,
			Title:     m.Title,
			Agenda:    m.Agenda,
			Date:      m.Date,
			StartTime: m.StartTime,
			EndTime:   m.EndTime,
			Location:  m.Location,
			Minutes:   m.Minutes,
			Status:    m.Status,
			CreatedBy: m.CreatedBy,
			CreatedAt: m.CreatedAt,
		}
	}
	return result, total, nil
}

func (s *meetingService) GetMeeting(ctx context.Context, id string) (*dto.MeetingResponse, error) {
	var m domain.Meeting
	if err := s.db.GetContext(ctx, &m, `SELECT * FROM meetings WHERE id=$1`, id); err != nil {
		return nil, domain.NewNotFoundError("meeting", id)
	}

	var attendees []struct {
		domain.MeetingAttendee
		FullName string `db:"full_name"`
	}

	query := `SELECT ma.*, u.full_name FROM meeting_attendees ma JOIN users u ON ma.user_id = u.id WHERE ma.meeting_id=$1`
	s.db.SelectContext(ctx, &attendees, query, id)

	attResponses := make([]dto.MeetingAttendeeResponse, len(attendees))
	for i, a := range attendees {
		attResponses[i] = dto.MeetingAttendeeResponse{
			ID:       a.ID,
			UserID:   a.UserID,
			FullName: a.FullName,
			Status:   a.Status,
		}
	}

	return &dto.MeetingResponse{
		ID:        m.ID,
		Title:     m.Title,
		Agenda:    m.Agenda,
		Date:      m.Date,
		StartTime: m.StartTime,
		EndTime:   m.EndTime,
		Location:  m.Location,
		Minutes:   m.Minutes,
		Status:    m.Status,
		CreatedBy: m.CreatedBy,
		Attendees: attResponses,
		CreatedAt: m.CreatedAt,
	}, nil
}

func (s *meetingService) CreateMeeting(ctx context.Context, schoolID, createdBy string, req dto.CreateMeetingRequest) (*dto.MeetingResponse, error) {
	m := &domain.Meeting{
		ID:        uuid.New().String(),
		SchoolID:  schoolID,
		Title:     req.Title,
		Agenda:    req.Agenda,
		Date:      req.Date,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Location:  req.Location,
		Status:    "scheduled",
		CreatedBy: createdBy,
	}

	tx, err := s.db.Beginx()
	if err != nil {
		return nil, domain.NewInternalError("failed to begin tx", err)
	}
	defer tx.Rollback()

	mQuery := `INSERT INTO meetings (id, school_id, title, agenda, date, start_time, end_time, location, status, created_by) VALUES (:id, :school_id, :title, :agenda, :date, :start_time, :end_time, :location, :status, :created_by)`
	if _, err := tx.NamedExecContext(ctx, mQuery, m); err != nil {
		return nil, domain.NewInternalError("failed to create meeting", err)
	}

	for _, userID := range req.Attendees {
		att := &domain.MeetingAttendee{
			ID:        uuid.New().String(),
			MeetingID: m.ID,
			UserID:    userID,
			Status:    "pending",
		}
		attQuery := `INSERT INTO meeting_attendees (id, meeting_id, user_id, status) VALUES (:id, :meeting_id, :user_id, :status)`
		tx.NamedExecContext(ctx, attQuery, att)
	}

	if err := tx.Commit(); err != nil {
		return nil, domain.NewInternalError("failed to commit meeting", err)
	}

	return s.GetMeeting(ctx, m.ID)
}

func (s *meetingService) UpdateMeetingMinutes(ctx context.Context, id string, minutes string) error {
	query := `UPDATE meetings SET minutes=$1, updated_at=NOW() WHERE id=$2`
	if _, err := s.db.ExecContext(ctx, query, minutes, id); err != nil {
		return domain.NewInternalError("failed to update minutes", err)
	}
	return nil
}

func (s *meetingService) ListEvents(ctx context.Context, schoolID string, filter dto.PaginationRequest) ([]dto.EventResponse, int64, error) {
	filter.Defaults()

	var items []domain.Event
	var total int64

	s.db.GetContext(ctx, &total, `SELECT COUNT(*) FROM events WHERE school_id=$1`, schoolID)
	query := `SELECT * FROM events WHERE school_id=$1 ORDER BY date DESC LIMIT $2 OFFSET $3`
	if err := s.db.SelectContext(ctx, &items, query, schoolID, filter.PageSize, filter.Offset()); err != nil {
		return nil, 0, domain.NewInternalError("failed to list events", err)
	}

	result := make([]dto.EventResponse, len(items))
	for i, e := range items {
		result[i] = dto.EventResponse{
			ID:          e.ID,
			Title:       e.Title,
			Description: e.Description,
			EventType:   e.EventType,
			Date:        e.Date,
			StartTime:   e.StartTime,
			EndTime:     e.EndTime,
			Location:    e.Location,
			CreatedBy:   e.CreatedBy,
		}
	}
	return result, total, nil
}

func (s *meetingService) CreateEvent(ctx context.Context, schoolID, createdBy string, req dto.CreateEventRequest) (*dto.EventResponse, error) {
	e := &domain.Event{
		ID:          uuid.New().String(),
		SchoolID:    schoolID,
		Title:       req.Title,
		Description: req.Description,
		EventType:   req.EventType,
		Date:        req.Date,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		Location:    req.Location,
		CreatedBy:   createdBy,
	}

	query := `INSERT INTO events (id, school_id, title, description, event_type, date, start_time, end_time, location, created_by) VALUES (:id, :school_id, :title, :description, :event_type, :date, :start_time, :end_time, :location, :created_by)`
	if _, err := s.db.NamedExecContext(ctx, query, e); err != nil {
		return nil, domain.NewInternalError("failed to create event", err)
	}

	return &dto.EventResponse{
		ID:          e.ID,
		Title:       e.Title,
		Description: e.Description,
		EventType:   e.EventType,
		Date:        e.Date,
		StartTime:   e.StartTime,
		EndTime:     e.EndTime,
		Location:    e.Location,
		CreatedBy:   e.CreatedBy,
	}, nil
}

func (s *meetingService) ListTasks(ctx context.Context, schoolID string, filter dto.PaginationRequest) ([]dto.TaskResponse, int64, error) {
	filter.Defaults()

	var items []domain.Task
	var total int64

	s.db.GetContext(ctx, &total, `SELECT COUNT(*) FROM tasks WHERE school_id=$1`, schoolID)
	query := `SELECT * FROM tasks WHERE school_id=$1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	if err := s.db.SelectContext(ctx, &items, query, schoolID, filter.PageSize, filter.Offset()); err != nil {
		return nil, 0, domain.NewInternalError("failed to list tasks", err)
	}

	result := make([]dto.TaskResponse, len(items))
	for i, t := range items {
		result[i] = dto.TaskResponse{
			ID:          t.ID,
			Title:       t.Title,
			Description: t.Description,
			AssignedTo:  t.AssignedTo,
			CreatedBy:   t.CreatedBy,
			Status:      t.Status,
			Priority:    t.Priority,
			DueDate:     t.DueDate,
			CompletedAt: t.CompletedAt,
			CreatedAt:   t.CreatedAt,
		}
	}
	return result, total, nil
}

func (s *meetingService) CreateTask(ctx context.Context, schoolID, createdBy string, req dto.CreateTaskRequest) (*dto.TaskResponse, error) {
	t := &domain.Task{
		ID:          uuid.New().String(),
		SchoolID:    schoolID,
		Title:       req.Title,
		Description: req.Description,
		AssignedTo:  req.AssignedTo,
		CreatedBy:   createdBy,
		Priority:    req.Priority,
		DueDate:     req.DueDate,
		Status:      "pending",
	}

	query := `INSERT INTO tasks (id, school_id, title, description, assigned_to, created_by, priority, due_date, status) VALUES (:id, :school_id, :title, :description, :assigned_to, :created_by, :priority, :due_date, :status)`
	if _, err := s.db.NamedExecContext(ctx, query, t); err != nil {
		return nil, domain.NewInternalError("failed to create task", err)
	}

	return &dto.TaskResponse{
		ID:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		AssignedTo:  t.AssignedTo,
		CreatedBy:   t.CreatedBy,
		Status:      t.Status,
		Priority:    t.Priority,
		DueDate:     t.DueDate,
		CreatedAt:   t.CreatedAt,
	}, nil
}

func (s *meetingService) UpdateTask(ctx context.Context, id string, req dto.UpdateTaskRequest) (*dto.TaskResponse, error) {
	var t domain.Task
	if err := s.db.GetContext(ctx, &t, `SELECT * FROM tasks WHERE id=$1`, id); err != nil {
		return nil, domain.NewNotFoundError("task", id)
	}

	now := time.Now()
	if req.Status == "completed" {
		t.CompletedAt = &now
		t.Status = "completed"
	} else if req.Status != "" {
		t.Status = req.Status
	}
	if req.Priority != "" {
		t.Priority = req.Priority
	}

	query := `UPDATE tasks SET status=$1, priority=$2, completed_at=$3, updated_at=NOW() WHERE id=$4`
	if _, err := s.db.ExecContext(ctx, query, t.Status, t.Priority, t.CompletedAt, id); err != nil {
		return nil, domain.NewInternalError("failed to update task", err)
	}

	return &dto.TaskResponse{
		ID:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		AssignedTo:  t.AssignedTo,
		CreatedBy:   t.CreatedBy,
		Status:      t.Status,
		Priority:    t.Priority,
		DueDate:     t.DueDate,
		CompletedAt: t.CompletedAt,
	}, nil
}
