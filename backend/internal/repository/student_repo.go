package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/opencode/erp-school-backend/internal/domain"
	"github.com/opencode/erp-school-backend/internal/infrastructure/database"
)

type StudentRepository interface {
	ListStudents(ctx context.Context, schoolID string, limit, offset int, filters map[string]interface{}) ([]domain.Student, int64, error)
	GetStudentByID(ctx context.Context, id string) (*domain.Student, error)
	GetStudentByUserID(ctx context.Context, userID string) (*domain.Student, error)
	GetStudentDetail(ctx context.Context, id string) (*domain.Student, []domain.StudentParent, error)
	CreateStudent(ctx context.Context, student *domain.Student, user *domain.User, profile *domain.UserProfile) error
	UpdateStudent(ctx context.Context, student *domain.Student) error
	LinkParent(ctx context.Context, parent *domain.StudentParent, user *domain.User) error
	UnlinkParent(ctx context.Context, parentID string) error
	GetParents(ctx context.Context, studentID string) ([]domain.StudentParent, error)

	CreateDocument(ctx context.Context, doc *domain.StudentDocument) error
	ListDocuments(ctx context.Context, studentID string) ([]domain.StudentDocument, error)
	VerifyDocument(ctx context.Context, docID, verifiedBy string) error
}

type studentRepo struct {
	db *sqlx.DB
}

func NewStudentRepository(db *sqlx.DB) StudentRepository {
	return &studentRepo{db: db}
}

func (r *studentRepo) ListStudents(ctx context.Context, schoolID string, limit, offset int, filters map[string]interface{}) ([]domain.Student, int64, error) {
	var items []domain.Student
	var total int64

	baseQuery := `FROM students s JOIN users u ON s.user_id = u.id LEFT JOIN user_profiles up ON u.id = up.user_id JOIN classes c ON s.class_id = c.id WHERE s.school_id=$1 AND u.deleted_at IS NULL`
	args := []interface{}{schoolID}
	argIdx := 2

	if classID, ok := filters["class_id"].(string); ok && classID != "" {
		baseQuery += fmt.Sprintf(" AND s.class_id=$%d", argIdx)
		args = append(args, classID)
		argIdx++
	}

	if status, ok := filters["status"].(string); ok && status != "" {
		baseQuery += fmt.Sprintf(" AND s.status=$%d", argIdx)
		args = append(args, status)
		argIdx++
	}

	if search, ok := filters["search"].(string); ok && search != "" {
		baseQuery += fmt.Sprintf(" AND (u.full_name ILIKE $%d OR s.nis ILIKE $%d OR s.nisn ILIKE $%d)", argIdx, argIdx+1, argIdx+2)
		args = append(args, "%"+search+"%", "%"+search+"%", "%"+search+"%")
		argIdx += 3
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) %s", baseQuery)
	if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return nil, 0, fmt.Errorf("count students: %w", err)
	}

	selectQuery := fmt.Sprintf(`SELECT s.*, u.full_name, u.email, up.gender, up.place_of_birth, up.date_of_birth, up.address, u.phone, c.name as class_name %s ORDER BY u.full_name LIMIT $%d OFFSET $%d`, baseQuery, argIdx, argIdx+1)
	args = append(args, limit, offset)

	if err := r.db.SelectContext(ctx, &items, selectQuery, args...); err != nil {
		return nil, 0, fmt.Errorf("list students: %w", err)
	}

	return items, total, nil
}

func (r *studentRepo) GetStudentByID(ctx context.Context, id string) (*domain.Student, error) {
	var s domain.Student
	query := `SELECT * FROM students WHERE id=$1`
	if err := r.db.GetContext(ctx, &s, query, id); err != nil {
		return nil, fmt.Errorf("get student: %w", err)
	}
	return &s, nil
}

func (r *studentRepo) GetStudentByUserID(ctx context.Context, userID string) (*domain.Student, error) {
	var s domain.Student
	query := `SELECT * FROM students WHERE user_id=$1`
	if err := r.db.GetContext(ctx, &s, query, userID); err != nil {
		return nil, fmt.Errorf("get student by user: %w", err)
	}
	return &s, nil
}

func (r *studentRepo) GetStudentDetail(ctx context.Context, id string) (*domain.Student, []domain.StudentParent, error) {
	var s domain.Student
	query := `SELECT * FROM students WHERE id=$1`
	if err := r.db.GetContext(ctx, &s, query, id); err != nil {
		return nil, nil, fmt.Errorf("get student: %w", err)
	}

	var parents []domain.StudentParent
	if err := r.db.SelectContext(ctx, &parents, database.ListStudentParents, id); err != nil {
		return nil, nil, fmt.Errorf("list parents: %w", err)
	}

	return &s, parents, nil
}

func (r *studentRepo) CreateStudent(ctx context.Context, student *domain.Student, user *domain.User, profile *domain.UserProfile) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	userQuery := `INSERT INTO users (id, school_id, email, username, password_hash, full_name, phone) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	if _, err := tx.ExecContext(ctx, userQuery, user.ID, user.SchoolID, user.Email, user.Username, user.PasswordHash, user.FullName, user.Phone); err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	studentQuery := `INSERT INTO students (id, user_id, school_id, nis, nisn, nik, class_id, academic_year_id, enrollment_date, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	if _, err := tx.ExecContext(ctx, studentQuery, student.ID, user.ID, student.SchoolID, student.NIS, student.NISN, student.NIK, student.ClassID, student.AcademicYearID, student.EnrollmentDate, student.Status); err != nil {
		return fmt.Errorf("create student: %w", err)
	}

	if profile != nil {
		profileQuery := `INSERT INTO user_profiles (id, user_id, place_of_birth, date_of_birth, gender, address) VALUES ($1, $2, $3, $4, $5, $6)`
		if _, err := tx.ExecContext(ctx, profileQuery, profile.ID, user.ID, profile.PlaceOfBirth, profile.DateOfBirth, profile.Gender, profile.Address); err != nil {
			return fmt.Errorf("create profile: %w", err)
		}
	}

	return tx.Commit()
}

func (r *studentRepo) UpdateStudent(ctx context.Context, student *domain.Student) error {
	query := `UPDATE students SET nis=:nis, nisn=:nisn, nik=:nik, class_id=:class_id, status=:status, updated_at=NOW() WHERE id=:id`
	_, err := r.db.NamedExecContext(ctx, query, student)
	if err != nil {
		return fmt.Errorf("update student: %w", err)
	}
	return nil
}

func (r *studentRepo) LinkParent(ctx context.Context, parent *domain.StudentParent, user *domain.User) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	userQuery := `INSERT INTO users (id, school_id, email, username, password_hash, full_name, phone) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	if _, err := tx.ExecContext(ctx, userQuery, user.ID, user.SchoolID, user.Email, user.Username, user.PasswordHash, user.FullName, user.Phone); err != nil {
		return fmt.Errorf("create parent user: %w", err)
	}

	parentQuery := `INSERT INTO student_parents (id, student_id, user_id, relation, is_primary, occupation, institution, income) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	if _, err := tx.ExecContext(ctx, parentQuery, parent.ID, parent.StudentID, user.ID, parent.Relation, parent.IsPrimary, parent.Occupation, parent.Institution, parent.Income); err != nil {
		return fmt.Errorf("link parent: %w", err)
	}

	return tx.Commit()
}

func (r *studentRepo) UnlinkParent(ctx context.Context, parentID string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM student_parents WHERE id=$1`, parentID)
	if err != nil {
		return fmt.Errorf("unlink parent: %w", err)
	}
	return nil
}

func (r *studentRepo) GetParents(ctx context.Context, studentID string) ([]domain.StudentParent, error) {
	var parents []domain.StudentParent
	if err := r.db.SelectContext(ctx, &parents, database.ListStudentParents, studentID); err != nil {
		return nil, fmt.Errorf("list parents: %w", err)
	}
	return parents, nil
}

func (r *studentRepo) CreateDocument(ctx context.Context, doc *domain.StudentDocument) error {
	query := `INSERT INTO student_documents (student_id, name, doc_type, file_url, status, notes) VALUES (:student_id, :name, :doc_type, :file_url, :status, :notes) RETURNING id, created_at, updated_at`
	rows, err := r.db.NamedQueryContext(ctx, query, doc)
	if err != nil {
		return fmt.Errorf("create document: %w", err)
	}
	defer rows.Close()
	if rows.Next() {
		return rows.Scan(&doc.ID, &doc.CreatedAt, &doc.UpdatedAt)
	}
	return nil
}

func (r *studentRepo) ListDocuments(ctx context.Context, studentID string) ([]domain.StudentDocument, error) {
	var docs []domain.StudentDocument
	query := `SELECT * FROM student_documents WHERE student_id=$1 ORDER BY created_at DESC`
	if err := r.db.SelectContext(ctx, &docs, query, studentID); err != nil {
		return nil, fmt.Errorf("list documents: %w", err)
	}
	return docs, nil
}

func (r *studentRepo) VerifyDocument(ctx context.Context, docID, verifiedBy string) error {
	query := `UPDATE student_documents SET status='verified', verified_at=NOW(), verified_by=$1, updated_at=NOW() WHERE id=$2`
	_, err := r.db.ExecContext(ctx, query, verifiedBy, docID)
	if err != nil {
		return fmt.Errorf("verify document: %w", err)
	}
	return nil
}
