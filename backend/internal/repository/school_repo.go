package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/opencode/erp-school-backend/internal/domain"
	"github.com/opencode/erp-school-backend/internal/infrastructure/database"
)

type SchoolRepository interface {
	ListSchools(ctx context.Context) ([]domain.School, error)
	GetSchoolByID(ctx context.Context, id string) (*domain.School, error)
	CreateSchool(ctx context.Context, school *domain.School) error
	UpdateSchool(ctx context.Context, school *domain.School) error
	DeleteSchool(ctx context.Context, id string) error

	ListAcademicYears(ctx context.Context, schoolID string, limit, offset int) ([]domain.AcademicYear, int64, error)
	GetAcademicYearByID(ctx context.Context, id string) (*domain.AcademicYear, error)
	GetActiveAcademicYear(ctx context.Context, schoolID string) (*domain.AcademicYear, error)
	CreateAcademicYear(ctx context.Context, ay *domain.AcademicYear) error
	UpdateAcademicYear(ctx context.Context, ay *domain.AcademicYear) error

	ListSemesters(ctx context.Context, academicYearID string) ([]domain.Semester, error)
	GetSemesterByID(ctx context.Context, id string) (*domain.Semester, error)
	CreateSemester(ctx context.Context, s *domain.Semester) error
	UpdateSemester(ctx context.Context, s *domain.Semester) error

	ListGrades(ctx context.Context, schoolID string) ([]domain.Grade, error)
	GetGradeByID(ctx context.Context, id string) (*domain.Grade, error)
	CreateGrade(ctx context.Context, g *domain.Grade) error
	UpdateGrade(ctx context.Context, g *domain.Grade) error

	ListClasses(ctx context.Context, schoolID string, limit, offset int) ([]domain.Class, int64, error)
	GetClassByID(ctx context.Context, id string) (*domain.Class, error)
	CreateClass(ctx context.Context, c *domain.Class) error
	UpdateClass(ctx context.Context, c *domain.Class) error

	ListSubjects(ctx context.Context, schoolID string, limit, offset int) ([]domain.Subject, int64, error)
	GetSubjectByID(ctx context.Context, id string) (*domain.Subject, error)
	CreateSubject(ctx context.Context, s *domain.Subject) error
	UpdateSubject(ctx context.Context, s *domain.Subject) error

	ListCurriculums(ctx context.Context, gradeID, subjectID, semesterID string) ([]domain.Curriculum, error)
	GetCurriculumByID(ctx context.Context, id string) (*domain.Curriculum, error)
	CreateCurriculum(ctx context.Context, c *domain.Curriculum) error
	UpdateCurriculum(ctx context.Context, c *domain.Curriculum) error

	ListRoles(ctx context.Context, schoolID string) ([]domain.Role, error)
	GetRoleByID(ctx context.Context, id string) (*domain.Role, error)
	CreateRole(ctx context.Context, r *domain.Role, permissionIDs []string) error
	UpdateRole(ctx context.Context, r *domain.Role, permissionIDs []string) error

	ListPermissions(ctx context.Context) ([]domain.Permission, error)
}

type schoolRepo struct {
	db *sqlx.DB
}

func NewSchoolRepository(db *sqlx.DB) SchoolRepository {
	return &schoolRepo{db: db}
}

func (r *schoolRepo) ListSchools(ctx context.Context) ([]domain.School, error) {
	var schools []domain.School
	if err := r.db.SelectContext(ctx, &schools, database.ListSchools); err != nil {
		return nil, fmt.Errorf("list schools: %w", err)
	}
	return schools, nil
}

func (r *schoolRepo) GetSchoolByID(ctx context.Context, id string) (*domain.School, error) {
	var s domain.School
	if err := r.db.GetContext(ctx, &s, database.GetSchoolByID, id); err != nil {
		return nil, fmt.Errorf("get school: %w", err)
	}
	return &s, nil
}

func (r *schoolRepo) CreateSchool(ctx context.Context, school *domain.School) error {
	query := `INSERT INTO schools (name, npsn, address, city, province, postal_code, phone, email, website, type, accreditation, established_date)
		VALUES (:name, :npsn, :address, :city, :province, :postal_code, :phone, :email, :website, :type, :accreditation, :established_date)
		RETURNING id, created_at, updated_at`
	rows, err := r.db.NamedQueryContext(ctx, query, school)
	if err != nil {
		return fmt.Errorf("create school: %w", err)
	}
	defer rows.Close()
	if rows.Next() {
		return rows.Scan(&school.ID, &school.CreatedAt, &school.UpdatedAt)
	}
	return nil
}

func (r *schoolRepo) UpdateSchool(ctx context.Context, school *domain.School) error {
	query := `UPDATE schools SET name=:name, npsn=:npsn, address=:address, city=:city, province=:province,
		postal_code=:postal_code, phone=:phone, email=:email, website=:website, type=:type,
		accreditation=:accreditation, is_active=:is_active, updated_at=NOW() WHERE id=:id`
	_, err := r.db.NamedExecContext(ctx, query, school)
	if err != nil {
		return fmt.Errorf("update school: %w", err)
	}
	return nil
}

func (r *schoolRepo) DeleteSchool(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM schools WHERE id=$1`, id)
	if err != nil {
		return fmt.Errorf("delete school: %w", err)
	}
	return nil
}

func (r *schoolRepo) ListAcademicYears(ctx context.Context, schoolID string, limit, offset int) ([]domain.AcademicYear, int64, error) {
	var items []domain.AcademicYear
	var total int64

	if err := r.db.GetContext(ctx, &total, database.CountAcademicYears, schoolID); err != nil {
		return nil, 0, fmt.Errorf("count academic years: %w", err)
	}

	if err := r.db.SelectContext(ctx, &items, database.ListAcademicYears, schoolID, limit, offset); err != nil {
		return nil, 0, fmt.Errorf("list academic years: %w", err)
	}

	return items, total, nil
}

func (r *schoolRepo) GetAcademicYearByID(ctx context.Context, id string) (*domain.AcademicYear, error) {
	var ay domain.AcademicYear
	query := `SELECT * FROM academic_years WHERE id=$1`
	if err := r.db.GetContext(ctx, &ay, query, id); err != nil {
		return nil, fmt.Errorf("get academic year: %w", err)
	}
	return &ay, nil
}

func (r *schoolRepo) GetActiveAcademicYear(ctx context.Context, schoolID string) (*domain.AcademicYear, error) {
	var ay domain.AcademicYear
	if err := r.db.GetContext(ctx, &ay, database.GetActiveAcademicYear, schoolID); err != nil {
		return nil, fmt.Errorf("get active academic year: %w", err)
	}
	return &ay, nil
}

func (r *schoolRepo) CreateAcademicYear(ctx context.Context, ay *domain.AcademicYear) error {
	query := `INSERT INTO academic_years (school_id, name, start_date, end_date)
		VALUES (:school_id, :name, :start_date, :end_date)
		RETURNING id, created_at, updated_at`
	rows, err := r.db.NamedQueryContext(ctx, query, ay)
	if err != nil {
		return fmt.Errorf("create academic year: %w", err)
	}
	defer rows.Close()
	if rows.Next() {
		return rows.Scan(&ay.ID, &ay.CreatedAt, &ay.UpdatedAt)
	}
	return nil
}

func (r *schoolRepo) UpdateAcademicYear(ctx context.Context, ay *domain.AcademicYear) error {
	query := `UPDATE academic_years SET name=:name, start_date=:start_date, end_date=:end_date, is_active=:is_active, updated_at=NOW() WHERE id=:id`
	_, err := r.db.NamedExecContext(ctx, query, ay)
	if err != nil {
		return fmt.Errorf("update academic year: %w", err)
	}
	return nil
}

func (r *schoolRepo) ListSemesters(ctx context.Context, academicYearID string) ([]domain.Semester, error) {
	var semesters []domain.Semester
	query := `SELECT * FROM semesters WHERE academic_year_id=$1 ORDER BY semester_number`
	if err := r.db.SelectContext(ctx, &semesters, query, academicYearID); err != nil {
		return nil, fmt.Errorf("list semesters: %w", err)
	}
	return semesters, nil
}

func (r *schoolRepo) GetSemesterByID(ctx context.Context, id string) (*domain.Semester, error) {
	var s domain.Semester
	query := `SELECT * FROM semesters WHERE id=$1`
	if err := r.db.GetContext(ctx, &s, query, id); err != nil {
		return nil, fmt.Errorf("get semester: %w", err)
	}
	return &s, nil
}

func (r *schoolRepo) CreateSemester(ctx context.Context, s *domain.Semester) error {
	query := `INSERT INTO semesters (academic_year_id, name, semester_number, start_date, end_date, is_active)
		VALUES (:academic_year_id, :name, :semester_number, :start_date, :end_date, :is_active)
		RETURNING id, created_at, updated_at`
	rows, err := r.db.NamedQueryContext(ctx, query, s)
	if err != nil {
		return fmt.Errorf("create semester: %w", err)
	}
	defer rows.Close()
	if rows.Next() {
		return rows.Scan(&s.ID, &s.CreatedAt, &s.UpdatedAt)
	}
	return nil
}

func (r *schoolRepo) UpdateSemester(ctx context.Context, s *domain.Semester) error {
	query := `UPDATE semesters SET name=:name, start_date=:start_date, end_date=:end_date, is_active=:is_active, updated_at=NOW() WHERE id=:id`
	_, err := r.db.NamedExecContext(ctx, query, s)
	if err != nil {
		return fmt.Errorf("update semester: %w", err)
	}
	return nil
}

func (r *schoolRepo) ListGrades(ctx context.Context, schoolID string) ([]domain.Grade, error) {
	var grades []domain.Grade
	query := `SELECT * FROM grades WHERE school_id=$1 ORDER BY level`
	if err := r.db.SelectContext(ctx, &grades, query, schoolID); err != nil {
		return nil, fmt.Errorf("list grades: %w", err)
	}
	return grades, nil
}

func (r *schoolRepo) GetGradeByID(ctx context.Context, id string) (*domain.Grade, error) {
	var g domain.Grade
	query := `SELECT * FROM grades WHERE id=$1`
	if err := r.db.GetContext(ctx, &g, query, id); err != nil {
		return nil, fmt.Errorf("get grade: %w", err)
	}
	return &g, nil
}

func (r *schoolRepo) CreateGrade(ctx context.Context, g *domain.Grade) error {
	query := `INSERT INTO grades (school_id, name, level) VALUES (:school_id, :name, :level) RETURNING id, created_at, updated_at`
	rows, err := r.db.NamedQueryContext(ctx, query, g)
	if err != nil {
		return fmt.Errorf("create grade: %w", err)
	}
	defer rows.Close()
	if rows.Next() {
		return rows.Scan(&g.ID, &g.CreatedAt, &g.UpdatedAt)
	}
	return nil
}

func (r *schoolRepo) UpdateGrade(ctx context.Context, g *domain.Grade) error {
	query := `UPDATE grades SET name=:name, level=:level, updated_at=NOW() WHERE id=:id`
	_, err := r.db.NamedExecContext(ctx, query, g)
	if err != nil {
		return fmt.Errorf("update grade: %w", err)
	}
	return nil
}

func (r *schoolRepo) ListClasses(ctx context.Context, schoolID string, limit, offset int) ([]domain.Class, int64, error) {
	var items []domain.Class
	var total int64

	if err := r.db.GetContext(ctx, &total, database.CountClasses, schoolID); err != nil {
		return nil, 0, fmt.Errorf("count classes: %w", err)
	}

	if err := r.db.SelectContext(ctx, &items, database.ListClasses, schoolID, limit, offset); err != nil {
		return nil, 0, fmt.Errorf("list classes: %w", err)
	}

	return items, total, nil
}

func (r *schoolRepo) GetClassByID(ctx context.Context, id string) (*domain.Class, error) {
	var c domain.Class
	if err := r.db.GetContext(ctx, &c, database.GetClassByID, id); err != nil {
		return nil, fmt.Errorf("get class: %w", err)
	}
	return &c, nil
}

func (r *schoolRepo) CreateClass(ctx context.Context, c *domain.Class) error {
	query := `INSERT INTO classes (school_id, grade_id, name, capacity, homeroom_teacher_id, academic_year_id)
		VALUES (:school_id, :grade_id, :name, :capacity, :homeroom_teacher_id, :academic_year_id)
		RETURNING id, created_at, updated_at`
	rows, err := r.db.NamedQueryContext(ctx, query, c)
	if err != nil {
		return fmt.Errorf("create class: %w", err)
	}
	defer rows.Close()
	if rows.Next() {
		return rows.Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
	}
	return nil
}

func (r *schoolRepo) UpdateClass(ctx context.Context, c *domain.Class) error {
	query := `UPDATE classes SET name=:name, capacity=:capacity, homeroom_teacher_id=:homeroom_teacher_id, updated_at=NOW() WHERE id=:id`
	_, err := r.db.NamedExecContext(ctx, query, c)
	if err != nil {
		return fmt.Errorf("update class: %w", err)
	}
	return nil
}

func (r *schoolRepo) ListSubjects(ctx context.Context, schoolID string, limit, offset int) ([]domain.Subject, int64, error) {
	var items []domain.Subject
	var total int64

	query := `SELECT * FROM subjects WHERE school_id=$1 ORDER BY code LIMIT $2 OFFSET $3`
	if err := r.db.SelectContext(ctx, &items, query, schoolID, limit, offset); err != nil {
		return nil, 0, fmt.Errorf("list subjects: %w", err)
	}

	if err := r.db.GetContext(ctx, &total, `SELECT COUNT(*) FROM subjects WHERE school_id=$1`, schoolID); err != nil {
		return nil, 0, fmt.Errorf("count subjects: %w", err)
	}

	return items, total, nil
}

func (r *schoolRepo) GetSubjectByID(ctx context.Context, id string) (*domain.Subject, error) {
	var s domain.Subject
	query := `SELECT * FROM subjects WHERE id=$1`
	if err := r.db.GetContext(ctx, &s, query, id); err != nil {
		return nil, fmt.Errorf("get subject: %w", err)
	}
	return &s, nil
}

func (r *schoolRepo) CreateSubject(ctx context.Context, s *domain.Subject) error {
	query := `INSERT INTO subjects (school_id, code, name, category, description, kkm)
		VALUES (:school_id, :code, :name, :category, :description, :kkm)
		RETURNING id, created_at, updated_at`
	rows, err := r.db.NamedQueryContext(ctx, query, s)
	if err != nil {
		return fmt.Errorf("create subject: %w", err)
	}
	defer rows.Close()
	if rows.Next() {
		return rows.Scan(&s.ID, &s.CreatedAt, &s.UpdatedAt)
	}
	return nil
}

func (r *schoolRepo) UpdateSubject(ctx context.Context, s *domain.Subject) error {
	query := `UPDATE subjects SET code=:code, name=:name, category=:category, description=:description, kkm=:kkm, is_active=:is_active, updated_at=NOW() WHERE id=:id`
	_, err := r.db.NamedExecContext(ctx, query, s)
	if err != nil {
		return fmt.Errorf("update subject: %w", err)
	}
	return nil
}

func (r *schoolRepo) ListCurriculums(ctx context.Context, gradeID, subjectID, semesterID string) ([]domain.Curriculum, error) {
	var items []domain.Curriculum
	query := `SELECT * FROM curriculums WHERE grade_id=$1 AND subject_id=$2 AND semester_id=$3`
	if err := r.db.SelectContext(ctx, &items, query, gradeID, subjectID, semesterID); err != nil {
		return nil, fmt.Errorf("list curriculums: %w", err)
	}
	return items, nil
}

func (r *schoolRepo) GetCurriculumByID(ctx context.Context, id string) (*domain.Curriculum, error) {
	var c domain.Curriculum
	if err := r.db.GetContext(ctx, &c, `SELECT * FROM curriculums WHERE id=$1`, id); err != nil {
		return nil, fmt.Errorf("get curriculum: %w", err)
	}
	return &c, nil
}

func (r *schoolRepo) CreateCurriculum(ctx context.Context, c *domain.Curriculum) error {
	query := `INSERT INTO curriculums (school_id, grade_id, subject_id, semester_id, content)
		VALUES (:school_id, :grade_id, :subject_id, :semester_id, :content)
		RETURNING id, created_at, updated_at`
	rows, err := r.db.NamedQueryContext(ctx, query, c)
	if err != nil {
		return fmt.Errorf("create curriculum: %w", err)
	}
	defer rows.Close()
	if rows.Next() {
		return rows.Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
	}
	return nil
}

func (r *schoolRepo) UpdateCurriculum(ctx context.Context, c *domain.Curriculum) error {
	query := `UPDATE curriculums SET content=:content, updated_at=NOW() WHERE id=:id`
	_, err := r.db.NamedExecContext(ctx, query, c)
	if err != nil {
		return fmt.Errorf("update curriculum: %w", err)
	}
	return nil
}

func (r *schoolRepo) ListRoles(ctx context.Context, schoolID string) ([]domain.Role, error) {
	var roles []domain.Role
	query := `SELECT * FROM roles WHERE school_id=$1 ORDER BY name`
	if err := r.db.SelectContext(ctx, &roles, query, schoolID); err != nil {
		return nil, fmt.Errorf("list roles: %w", err)
	}
	return roles, nil
}

func (r *schoolRepo) GetRoleByID(ctx context.Context, id string) (*domain.Role, error) {
	var role domain.Role
	query := `SELECT * FROM roles WHERE id=$1`
	if err := r.db.GetContext(ctx, &role, query, id); err != nil {
		return nil, fmt.Errorf("get role: %w", err)
	}
	return &role, nil
}

func (r *schoolRepo) CreateRole(ctx context.Context, role *domain.Role, permissionIDs []string) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	query := `INSERT INTO roles (school_id, name, slug, description) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`
	row := tx.QueryRowxContext(ctx, query, role.SchoolID, role.Name, role.Slug, role.Description)
	if err := row.Scan(&role.ID, &role.CreatedAt, &role.UpdatedAt); err != nil {
		return fmt.Errorf("create role: %w", err)
	}

	for _, permID := range permissionIDs {
		if _, err := tx.ExecContext(ctx, `INSERT INTO role_permissions (role_id, permission_id) VALUES ($1, $2)`, role.ID, permID); err != nil {
			return fmt.Errorf("assign permission: %w", err)
		}
	}

	return tx.Commit()
}

func (r *schoolRepo) UpdateRole(ctx context.Context, role *domain.Role, permissionIDs []string) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	query := `UPDATE roles SET name=$1, description=$2, updated_at=NOW() WHERE id=$3`
	if _, err := tx.ExecContext(ctx, query, role.Name, role.Description, role.ID); err != nil {
		return fmt.Errorf("update role: %w", err)
	}

	if _, err := tx.ExecContext(ctx, `DELETE FROM role_permissions WHERE role_id=$1`, role.ID); err != nil {
		return fmt.Errorf("clear permissions: %w", err)
	}

	for _, permID := range permissionIDs {
		if _, err := tx.ExecContext(ctx, `INSERT INTO role_permissions (role_id, permission_id) VALUES ($1, $2)`, role.ID, permID); err != nil {
			return fmt.Errorf("assign permission: %w", err)
		}
	}

	return tx.Commit()
}

func (r *schoolRepo) ListPermissions(ctx context.Context) ([]domain.Permission, error) {
	var permissions []domain.Permission
	query := `SELECT * FROM permissions ORDER BY module, name`
	if err := r.db.SelectContext(ctx, &permissions, query); err != nil {
		return nil, fmt.Errorf("list permissions: %w", err)
	}
	return permissions, nil
}
