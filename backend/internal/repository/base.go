package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/opencode/erp-school-backend/internal/dto"
)

type Repository[T any] interface {
	FindAll(ctx context.Context, filter dto.PaginationRequest, schoolID string) ([]T, int64, error)
	FindByID(ctx context.Context, id string) (*T, error)
	Create(ctx context.Context, entity *T) error
	Update(ctx context.Context, id string, updates map[string]interface{}) error
	Delete(ctx context.Context, id string) error
}

type BaseRepository[T any] struct {
	db    *sqlx.DB
	table string
}

func NewBaseRepository[T any](db *sqlx.DB, table string) *BaseRepository[T] {
	return &BaseRepository[T]{db: db, table: table}
}

func (r *BaseRepository[T]) DB() *sqlx.DB {
	return r.db
}

func (r *BaseRepository[T]) Table() string {
	return r.table
}

func (r *BaseRepository[T]) FindAll(ctx context.Context, filter dto.PaginationRequest, schoolID string) ([]T, int64, error) {
	filter.Defaults()
	var entities []T
	var total int64

	where := "WHERE 1=1"

	var args []interface{}
	argIdx := 1

	where += fmt.Sprintf(" AND school_id=$%d", argIdx)
	args = append(args, schoolID)
	argIdx++

	if filter.Search != "" {
		where += fmt.Sprintf(" AND (name ILIKE $%d OR title ILIKE $%d)", argIdx, argIdx+1)
		args = append(args, "%"+filter.Search+"%", "%"+filter.Search+"%")
		argIdx += 2
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s %s", r.table, where)
	if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return nil, 0, fmt.Errorf("count query failed: %w", err)
	}

	sortDir := "ASC"
	if strings.ToUpper(filter.SortDir) == "DESC" {
		sortDir = "DESC"
	}
	sortBy := r.sanitizeSortBy(filter.SortBy)

	query := fmt.Sprintf("SELECT * FROM %s %s ORDER BY %s %s LIMIT $%d OFFSET $%d",
		r.table, where, sortBy, sortDir, argIdx, argIdx+1)
	args = append(args, filter.PageSize, filter.Offset())

	if err := r.db.SelectContext(ctx, &entities, query, args...); err != nil {
		return nil, 0, fmt.Errorf("select query failed: %w", err)
	}

	return entities, total, nil
}

func (r *BaseRepository[T]) FindByID(ctx context.Context, id string) (*T, error) {
	var entity T
	query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1", r.table)
	if err := r.db.GetContext(ctx, &entity, query, id); err != nil {
		return nil, fmt.Errorf("find by id failed: %w", err)
	}
	return &entity, nil
}

func (r *BaseRepository[T]) Create(ctx context.Context, entity *T) error {
	query := fmt.Sprintf("INSERT INTO %s VALUES (:id, :school_id) RETURNING *", r.table)
	rows, err := r.db.NamedQueryContext(ctx, query, entity)
	if err != nil {
		return fmt.Errorf("create failed: %w", err)
	}
	defer rows.Close()
	if rows.Next() {
		return rows.StructScan(entity)
	}
	return nil
}

func (r *BaseRepository[T]) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return nil
	}

	setClauses := make([]string, 0, len(updates))
	args := make([]interface{}, 0, len(updates)+1)
	argIdx := 1

	for col, val := range updates {
		setClauses = append(setClauses, fmt.Sprintf("%s=$%d", r.sanitizeColumn(col), argIdx))
		args = append(args, val)
		argIdx++
	}

	args = append(args, id)

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id=$%d",
		r.table, strings.Join(setClauses, ", "), argIdx)

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("update failed: %w", err)
	}
	return nil
}

func (r *BaseRepository[T]) Delete(ctx context.Context, id string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", r.table)
	if _, err := r.db.ExecContext(ctx, query, id); err != nil {
		return fmt.Errorf("delete failed: %w", err)
	}
	return nil
}

func (r *BaseRepository[T]) SoftDelete(ctx context.Context, id string) error {
	query := fmt.Sprintf("UPDATE %s SET deleted_at=NOW() WHERE id=$1", r.table)
	if _, err := r.db.ExecContext(ctx, query, id); err != nil {
		return fmt.Errorf("soft delete failed: %w", err)
	}
	return nil
}

func (r *BaseRepository[T]) sanitizeSortBy(field string) string {
	allowed := map[string]bool{
		"id": true, "name": true, "title": true, "code": true,
		"created_at": true, "updated_at": true, "date": true,
		"start_date": true, "end_date": true, "status": true,
		"level": true, "due_date": true, "full_name": true,
		"email": true, "username": true, "amount": true,
	}
	if allowed[field] {
		return field
	}
	return "created_at"
}

func (r *BaseRepository[T]) sanitizeColumn(col string) string {
	allowed := map[string]bool{
		"id": true, "school_id": true, "name": true, "title": true, "code": true,
		"email": true, "username": true, "full_name": true, "phone": true,
		"description": true, "content": true, "status": true, "type": true,
		"is_active": true, "address": true, "city": true, "province": true,
		"postal_code": true, "website": true, "accreditation": true,
		"npsn": true, "nis": true, "nisn": true, "nik": true, "nip": true,
		"nupk": true, "position": true, "department": true, "base_salary": true,
		"bank_account": true, "bank_name": true, "capacity": true,
		"homeroom_teacher_id": true, "academic_year_id": true, "class_id": true,
		"subject_id": true, "teacher_id": true, "grade_id": true,
		"category": true, "kkm": true, "amount": true, "frequency": true,
		"start_date": true, "end_date": true, "due_date": true, "date": true,
		"day": true, "start_time": true, "end_time": true, "room": true,
		"semester_id": true, "semester_number": true, "level": true,
		"notes": true, "quality": true, "relation": true, "is_primary": true,
		"occupation": true, "institution": true, "income": true,
		"enrollment_date": true, "join_date": true, "leave_type": true,
		"reason": true, "approved_by": true, "approved_at": true,
		"password_changed_at": true, "avatar_url": true, "logo_url": true,
		"is_system": true, "role": true, "message": true, "target_role": true,
		"priority": true, "location": true, "module": true, "key": true,
		"value": true, "biography": true, "bio": true,
	}
	if allowed[col] {
		return col
	}
	return col
}
