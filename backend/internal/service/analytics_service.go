package service

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/opencode/erp-school-backend/internal/domain"
	"github.com/opencode/erp-school-backend/internal/dto"
	"github.com/opencode/erp-school-backend/internal/infrastructure/database"
	"go.uber.org/zap"
)

type AnalyticsService interface {
	GetDashboard(ctx context.Context, schoolID string) (*dto.DashboardStats, error)
	GetAcademicAnalytics(ctx context.Context, schoolID, classID string) (map[string]interface{}, error)
	GetFinanceAnalytics(ctx context.Context, schoolID string) (*dto.FinancialSummary, error)
	GetTahfidzAnalytics(ctx context.Context, schoolID string) (map[string]interface{}, error)
	GetAdmissionAnalytics(ctx context.Context, schoolID string) (map[string]interface{}, error)
	GetAttendanceAnalytics(ctx context.Context, schoolID, classID string) (*dto.StudentAttendanceSummary, error)
}

type analyticsService struct {
	db     *sqlx.DB
	logger *zap.Logger
}

func NewAnalyticsService(db *sqlx.DB, logger *zap.Logger) AnalyticsService {
	return &analyticsService{db: db, logger: logger}
}

func (s *analyticsService) GetDashboard(ctx context.Context, schoolID string) (*dto.DashboardStats, error) {
	var stats dto.DashboardStats
	query := database.DashboardStatsQuery
	row := s.db.QueryRowxContext(ctx, query, schoolID)
	if err := row.Scan(&stats.TotalStudents, &stats.TotalTeachers, &stats.TotalEmployees, &stats.TotalClasses); err != nil {
		return nil, domain.NewInternalError("failed to get dashboard stats", err)
	}
	return &stats, nil
}

func (s *analyticsService) GetAcademicAnalytics(ctx context.Context, schoolID, classID string) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	var totalStudents int64
	query := `SELECT COUNT(*) FROM students WHERE school_id=$1 AND status='active'`
	args := []interface{}{schoolID}
	if classID != "" {
		query += ` AND class_id=$2`
		args = append(args, classID)
	}
	s.db.GetContext(ctx, &totalStudents, query, args...)
	result["total_students"] = totalStudents

	var avgScore float64
	s.db.GetContext(ctx, &avgScore, `SELECT COALESCE(AVG(average_score), 0) FROM report_cards WHERE class_id IN (SELECT id FROM classes WHERE school_id=$1)`, schoolID)
	result["average_score"] = avgScore

	return result, nil
}

func (s *analyticsService) GetFinanceAnalytics(ctx context.Context, schoolID string) (*dto.FinancialSummary, error) {
	var summary dto.FinancialSummary

	var revenue, expense, outstanding float64
	row := s.db.QueryRowxContext(ctx, database.FinancialSummaryQuery, schoolID)
	if err := row.Scan(&revenue, &expense, &outstanding); err != nil {
		return nil, domain.NewInternalError("failed to get financial summary", err)
	}

	summary.TotalRevenue = revenue
	summary.TotalExpense = expense
	summary.Outstanding = outstanding
	summary.Balance = revenue - expense

	if summary.TotalRevenue > 0 {
		summary.CollectionRate = (summary.TotalRevenue - summary.Outstanding) / summary.TotalRevenue * 100
	}

	return &summary, nil
}

func (s *analyticsService) GetTahfidzAnalytics(ctx context.Context, schoolID string) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	var totalJuz int
	s.db.GetContext(ctx, &totalJuz, `SELECT COALESCE(SUM(juz_count), 0) FROM (SELECT COUNT(DISTINCT juz) as juz_count FROM tahfidz_progress tp JOIN students st ON tp.student_id = st.id WHERE st.school_id=$1 AND tp.status='memorized' GROUP BY tp.student_id) sub`, schoolID)
	result["total_juz_memorized"] = totalJuz

	var totalStudents int
	s.db.GetContext(ctx, &totalStudents, `SELECT COUNT(DISTINCT student_id) FROM tahfidz_progress`)
	result["students_in_tahfidz"] = totalStudents

	return result, nil
}

func (s *analyticsService) GetAdmissionAnalytics(ctx context.Context, schoolID string) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	var total, accepted, rejected, pending int64
	s.db.GetContext(ctx, &total, `SELECT COUNT(*) FROM admission_applicants WHERE school_id=$1`, schoolID)
	s.db.GetContext(ctx, &accepted, `SELECT COUNT(*) FROM admission_applicants WHERE school_id=$1 AND status='accepted'`, schoolID)
	s.db.GetContext(ctx, &rejected, `SELECT COUNT(*) FROM admission_applicants WHERE school_id=$1 AND status='rejected'`, schoolID)
	s.db.GetContext(ctx, &pending, `SELECT COUNT(*) FROM admission_applicants WHERE school_id=$1 AND status='pending'`, schoolID)

	result["total_applicants"] = total
	result["accepted"] = accepted
	result["rejected"] = rejected
	result["pending"] = pending

	if total > 0 {
		result["acceptance_rate"] = float64(accepted) / float64(total) * 100
	}

	return result, nil
}

func (s *analyticsService) GetAttendanceAnalytics(ctx context.Context, schoolID, classID string) (*dto.StudentAttendanceSummary, error) {
	var summary dto.StudentAttendanceSummary

	args := []interface{}{schoolID}
	query := `SELECT COALESCE(SUM(CASE WHEN a.status='present' THEN 1 ELSE 0 END), 0),
		COALESCE(SUM(CASE WHEN a.status='absent' THEN 1 ELSE 0 END), 0),
		COALESCE(SUM(CASE WHEN a.status='sick' THEN 1 ELSE 0 END), 0),
		COALESCE(SUM(CASE WHEN a.status='permit' THEN 1 ELSE 0 END), 0),
		COALESCE(SUM(CASE WHEN a.status='late' THEN 1 ELSE 0 END), 0)
		FROM attendances a
		JOIN students s ON a.student_id = s.id
		WHERE s.school_id=$1`

	if classID != "" {
		query += ` AND s.class_id=$2`
		args = append(args, classID)
	}

	if err := s.db.QueryRowxContext(ctx, query, args...).Scan(&summary.Present, &summary.Absent, &summary.Sick, &summary.Permit, &summary.Late); err != nil {
		return nil, domain.NewInternalError("failed to get attendance summary", err)
	}

	return &summary, nil
}
