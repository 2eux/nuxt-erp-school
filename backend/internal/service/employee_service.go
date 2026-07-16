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
	"golang.org/x/crypto/bcrypt"
)

type EmployeeService interface {
	ListEmployees(ctx context.Context, schoolID string, filter dto.PaginationRequest) ([]dto.EmployeeResponse, int64, error)
	GetEmployee(ctx context.Context, id string) (*dto.EmployeeResponse, error)
	CreateEmployee(ctx context.Context, schoolID string, req dto.CreateEmployeeRequest) (*dto.EmployeeResponse, error)
	UpdateEmployee(ctx context.Context, id string, req dto.CreateEmployeeRequest) (*dto.EmployeeResponse, error)
	ListAttendances(ctx context.Context, employeeID string, startDate, endDate time.Time) ([]domain.EmployeeAttendance, error)
	CreateAttendance(ctx context.Context, employeeID string, req dto.CreateEmployeeAttendanceRequest) (*domain.EmployeeAttendance, error)
	ListLeaveRequests(ctx context.Context, schoolID string, filter dto.PaginationRequest) ([]dto.LeaveRequestResponse, int64, error)
	SubmitLeaveRequest(ctx context.Context, employeeID, schoolID string, req dto.CreateLeaveRequest) (*dto.LeaveRequestResponse, error)
	ApproveLeave(ctx context.Context, id, approverID string) error
	RejectLeave(ctx context.Context, id, approverID string) error
}

type employeeService struct {
	db     *sqlx.DB
	logger *zap.Logger
}

func NewEmployeeService(db *sqlx.DB, logger *zap.Logger) EmployeeService {
	return &employeeService{db: db, logger: logger}
}

func (s *employeeService) ListEmployees(ctx context.Context, schoolID string, filter dto.PaginationRequest) ([]dto.EmployeeResponse, int64, error) {
	filter.Defaults()

	var total int64
	if err := s.db.GetContext(ctx, &total, database.CountEmployees, schoolID); err != nil {
		return nil, 0, domain.NewInternalError("failed to count employees", err)
	}

	var items []struct {
		domain.Employee
		FullName string `db:"full_name"`
		Email    string `db:"email"`
		Phone    string `db:"phone"`
	}

	if err := s.db.SelectContext(ctx, &items, database.ListEmployees, schoolID, filter.PageSize, filter.Offset()); err != nil {
		return nil, 0, domain.NewInternalError("failed to list employees", err)
	}

	result := make([]dto.EmployeeResponse, len(items))
	for i, e := range items {
		result[i] = dto.EmployeeResponse{
			ID:         e.ID,
			UserID:     e.UserID,
			SchoolID:   e.SchoolID,
			Email:      e.Email,
			FullName:   e.FullName,
			NIP:        e.NIP,
			NIK:        e.NIK,
			Position:   e.Position,
			Department: e.Department,
			JoinDate:   e.JoinDate,
			Status:     e.Status,
			BaseSalary: e.BaseSalary,
			CreatedAt:  e.CreatedAt,
			UpdatedAt:  e.UpdatedAt,
		}
	}
	return result, total, nil
}

func (s *employeeService) GetEmployee(ctx context.Context, id string) (*dto.EmployeeResponse, error) {
	var e domain.Employee
	if err := s.db.GetContext(ctx, &e, `SELECT * FROM employees WHERE id=$1`, id); err != nil {
		return nil, domain.NewNotFoundError("employee", id)
	}

	var u domain.User
	s.db.GetContext(ctx, &u, `SELECT * FROM users WHERE id=$1`, e.UserID)

	return &dto.EmployeeResponse{
		ID:         e.ID,
		UserID:     e.UserID,
		SchoolID:   e.SchoolID,
		Email:      u.Email,
		FullName:   u.FullName,
		NIP:        e.NIP,
		NIK:        e.NIK,
		Position:   e.Position,
		Department: e.Department,
		JoinDate:   e.JoinDate,
		Status:     e.Status,
		BaseSalary: e.BaseSalary,
		CreatedAt:  e.CreatedAt,
		UpdatedAt:  e.UpdatedAt,
	}, nil
}

func (s *employeeService) CreateEmployee(ctx context.Context, schoolID string, req dto.CreateEmployeeRequest) (*dto.EmployeeResponse, error) {
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
		return nil, domain.NewInternalError("failed to create user for employee", err)
	}

	empID := uuid.New().String()
	empQuery := `INSERT INTO employees (id, user_id, school_id, nip, nik, position, department, join_date, base_salary, bank_account, bank_name) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	if _, err := s.db.ExecContext(ctx, empQuery, empID, userID, schoolID, req.NIP, req.NIK, req.Position, req.Department, req.JoinDate, req.BaseSalary, req.BankAccount, req.BankName); err != nil {
		return nil, domain.NewInternalError("failed to create employee", err)
	}

	return s.GetEmployee(ctx, empID)
}

func (s *employeeService) UpdateEmployee(ctx context.Context, id string, req dto.CreateEmployeeRequest) (*dto.EmployeeResponse, error) {
	var e domain.Employee
	if err := s.db.GetContext(ctx, &e, `SELECT * FROM employees WHERE id=$1`, id); err != nil {
		return nil, domain.NewNotFoundError("employee", id)
	}

	query := `UPDATE employees SET nip=$1, nik=$2, position=$3, department=$4, base_salary=$5, bank_account=$6, bank_name=$7, updated_at=NOW() WHERE id=$8`
	if _, err := s.db.ExecContext(ctx, query, req.NIP, req.NIK, req.Position, req.Department, req.BaseSalary, req.BankAccount, req.BankName, id); err != nil {
		return nil, domain.NewInternalError("failed to update employee", err)
	}

	return s.GetEmployee(ctx, id)
}

func (s *employeeService) ListAttendances(ctx context.Context, employeeID string, startDate, endDate time.Time) ([]domain.EmployeeAttendance, error) {
	var items []domain.EmployeeAttendance
	query := `SELECT * FROM employee_attendances WHERE employee_id=$1 AND date >= $2 AND date <= $3 ORDER BY date DESC`
	if err := s.db.SelectContext(ctx, &items, query, employeeID, startDate, endDate); err != nil {
		return nil, domain.NewInternalError("failed to list attendances", err)
	}
	return items, nil
}

func (s *employeeService) CreateAttendance(ctx context.Context, employeeID string, req dto.CreateEmployeeAttendanceRequest) (*domain.EmployeeAttendance, error) {
	attendance := &domain.EmployeeAttendance{
		ID:         uuid.New().String(),
		EmployeeID: employeeID,
		Date:       req.Date,
		CheckIn:    req.Date,
		Status:     req.Status,
		Notes:      req.Notes,
	}

	query := `INSERT INTO employee_attendances (id, employee_id, date, check_in, status, notes) VALUES (:id, :employee_id, :date, :check_in, :status, :notes)`
	if _, err := s.db.NamedExecContext(ctx, query, attendance); err != nil {
		return nil, domain.NewInternalError("failed to create attendance", err)
	}
	return attendance, nil
}

func (s *employeeService) ListLeaveRequests(ctx context.Context, schoolID string, filter dto.PaginationRequest) ([]dto.LeaveRequestResponse, int64, error) {
	filter.Defaults()

	type leaveRow struct {
		domain.LeaveRequest
		FullName string `db:"full_name"`
	}

	var rows []leaveRow
	if err := s.db.SelectContext(ctx, &rows, database.ListLeaveRequests, schoolID, filter.PageSize, filter.Offset()); err != nil {
		return nil, 0, domain.NewInternalError("failed to list leave requests", err)
	}

	var total int64
	s.db.GetContext(ctx, &total, `SELECT COUNT(*) FROM leave_requests lr JOIN employees e ON lr.employee_id = e.id WHERE e.school_id=$1`, schoolID)

	result := make([]dto.LeaveRequestResponse, len(rows))
	for i, r := range rows {
		result[i] = dto.LeaveRequestResponse{
			ID:         r.ID,
			EmployeeID: r.EmployeeID,
			FullName:   r.FullName,
			LeaveType:  r.LeaveType,
			StartDate:  r.StartDate,
			EndDate:    r.EndDate,
			Reason:     r.Reason,
			Status:     r.Status,
			ApprovedBy: r.ApprovedBy,
			ApprovedAt: r.ApprovedAt,
			CreatedAt:  r.CreatedAt,
		}
	}
	return result, total, nil
}

func (s *employeeService) SubmitLeaveRequest(ctx context.Context, employeeID, schoolID string, req dto.CreateLeaveRequest) (*dto.LeaveRequestResponse, error) {
	lr := &domain.LeaveRequest{
		ID:         uuid.New().String(),
		EmployeeID: employeeID,
		LeaveType:  req.LeaveType,
		StartDate:  req.StartDate,
		EndDate:    req.EndDate,
		Reason:     req.Reason,
		Status:     "pending",
	}

	query := `INSERT INTO leave_requests (id, employee_id, leave_type, start_date, end_date, reason, status) VALUES (:id, :employee_id, :leave_type, :start_date, :end_date, :reason, :status)`
	if _, err := s.db.NamedExecContext(ctx, query, lr); err != nil {
		return nil, domain.NewInternalError("failed to submit leave request", err)
	}

	return &dto.LeaveRequestResponse{
		ID:         lr.ID,
		EmployeeID: lr.EmployeeID,
		LeaveType:  lr.LeaveType,
		StartDate:  lr.StartDate,
		EndDate:    lr.EndDate,
		Reason:     lr.Reason,
		Status:     lr.Status,
		CreatedAt:  lr.CreatedAt,
	}, nil
}

func (s *employeeService) ApproveLeave(ctx context.Context, id, approverID string) error {
	now := time.Now()
	query := `UPDATE leave_requests SET status='approved', approved_by=$1, approved_at=$2, updated_at=NOW() WHERE id=$3`
	if _, err := s.db.ExecContext(ctx, query, approverID, now, id); err != nil {
		return domain.NewInternalError("failed to approve leave", err)
	}
	return nil
}

func (s *employeeService) RejectLeave(ctx context.Context, id, approverID string) error {
	now := time.Now()
	query := `UPDATE leave_requests SET status='rejected', approved_by=$1, approved_at=$2, updated_at=NOW() WHERE id=$3`
	if _, err := s.db.ExecContext(ctx, query, approverID, now, id); err != nil {
		return domain.NewInternalError("failed to reject leave", err)
	}
	return nil
}
