package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/opencode/erp-school-backend/internal/domain"
	"github.com/opencode/erp-school-backend/internal/infrastructure/database"
)

type AuthRepository interface {
	FindUserByEmail(ctx context.Context, email string) (*domain.User, error)
	FindUserByID(ctx context.Context, id string) (*domain.User, error)
	FindUserByUsername(ctx context.Context, username string) (*domain.User, error)
	GetUserRoles(ctx context.Context, userID string) ([]domain.Role, error)
	GetUserPermissions(ctx context.Context, userID string) ([]string, error)
	CreateUser(ctx context.Context, user *domain.User) error
	UpdateUser(ctx context.Context, user *domain.User) error
	UpdatePassword(ctx context.Context, userID, hashedPassword string) error
	UpdateLastLogin(ctx context.Context, userID string) error
	CreateSession(ctx context.Context, session *domain.UserSession) error
	FindSessionByToken(ctx context.Context, refreshToken string) (*domain.UserSession, error)
	DeleteSession(ctx context.Context, sessionID string) error
	DeleteUserSessions(ctx context.Context, userID string) error
	AssignRoles(ctx context.Context, userID string, roleIDs []string) error
}

type authRepo struct {
	db *sqlx.DB
}

func NewAuthRepository(db *sqlx.DB) AuthRepository {
	return &authRepo{db: db}
}

func (r *authRepo) FindUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	if err := r.db.GetContext(ctx, &user, database.GetUserByEmail, email); err != nil {
		return nil, fmt.Errorf("find user by email: %w", err)
	}
	return &user, nil
}

func (r *authRepo) FindUserByID(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User
	query := `SELECT * FROM users WHERE id=$1 AND deleted_at IS NULL`
	if err := r.db.GetContext(ctx, &user, query, id); err != nil {
		return nil, fmt.Errorf("find user by id: %w", err)
	}
	return &user, nil
}

func (r *authRepo) FindUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	var user domain.User
	query := `SELECT * FROM users WHERE username=$1 AND deleted_at IS NULL`
	if err := r.db.GetContext(ctx, &user, query, username); err != nil {
		return nil, fmt.Errorf("find user by username: %w", err)
	}
	return &user, nil
}

func (r *authRepo) GetUserRoles(ctx context.Context, userID string) ([]domain.Role, error) {
	var roles []domain.Role
	query := `
		SELECT r.* FROM roles r
		JOIN user_roles ur ON r.id = ur.role_id
		WHERE ur.user_id = $1`
	if err := r.db.SelectContext(ctx, &roles, query, userID); err != nil {
		return nil, fmt.Errorf("get user roles: %w", err)
	}
	return roles, nil
}

func (r *authRepo) GetUserPermissions(ctx context.Context, userID string) ([]string, error) {
	var permissions []string
	if err := r.db.SelectContext(ctx, &permissions, database.GetUserPermissions, userID); err != nil {
		return nil, fmt.Errorf("get user permissions: %w", err)
	}
	return permissions, nil
}

func (r *authRepo) CreateUser(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (id, school_id, email, username, password_hash, full_name, phone)
		VALUES (:id, :school_id, :email, :username, :password_hash, :full_name, :phone)
		RETURNING created_at, updated_at`
	rows, err := r.db.NamedQueryContext(ctx, query, user)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	defer rows.Close()
	if rows.Next() {
		return rows.Scan(&user.CreatedAt, &user.UpdatedAt)
	}
	return nil
}

func (r *authRepo) UpdateUser(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users SET email=:email, username=:username, full_name=:full_name,
		phone=:phone, is_active=:is_active, updated_at=NOW()
		WHERE id=:id`
	_, err := r.db.NamedExecContext(ctx, query, user)
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	return nil
}

func (r *authRepo) UpdatePassword(ctx context.Context, userID, hashedPassword string) error {
	query := `UPDATE users SET password_hash=$1, password_changed_at=NOW(), updated_at=NOW() WHERE id=$2`
	_, err := r.db.ExecContext(ctx, query, hashedPassword, userID)
	if err != nil {
		return fmt.Errorf("update password: %w", err)
	}
	return nil
}

func (r *authRepo) UpdateLastLogin(ctx context.Context, userID string) error {
	query := `UPDATE users SET last_login_at=NOW(), updated_at=NOW() WHERE id=$1`
	_, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("update last login: %w", err)
	}
	return nil
}

func (r *authRepo) CreateSession(ctx context.Context, session *domain.UserSession) error {
	query := `
		INSERT INTO user_sessions (id, user_id, refresh_token, user_agent, ip_address, expires_at)
		VALUES (:id, :user_id, :refresh_token, :user_agent, :ip_address, :expires_at)
		RETURNING created_at`
	rows, err := r.db.NamedQueryContext(ctx, query, session)
	if err != nil {
		return fmt.Errorf("create session: %w", err)
	}
	defer rows.Close()
	if rows.Next() {
		return rows.Scan(&session.CreatedAt)
	}
	return nil
}

func (r *authRepo) FindSessionByToken(ctx context.Context, refreshToken string) (*domain.UserSession, error) {
	var session domain.UserSession
	if err := r.db.GetContext(ctx, &session, database.GetSessionByToken, refreshToken); err != nil {
		return nil, fmt.Errorf("find session by token: %w", err)
	}
	return &session, nil
}

func (r *authRepo) DeleteSession(ctx context.Context, sessionID string) error {
	_, err := r.db.ExecContext(ctx, database.DeleteSession, sessionID)
	if err != nil {
		return fmt.Errorf("delete session: %w", err)
	}
	return nil
}

func (r *authRepo) DeleteUserSessions(ctx context.Context, userID string) error {
	_, err := r.db.ExecContext(ctx, database.DeleteUserSessions, userID)
	if err != nil {
		return fmt.Errorf("delete user sessions: %w", err)
	}
	return nil
}

func (r *authRepo) AssignRoles(ctx context.Context, userID string, roleIDs []string) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `DELETE FROM user_roles WHERE user_id=$1`, userID); err != nil {
		return fmt.Errorf("delete existing roles: %w", err)
	}

	for _, roleID := range roleIDs {
		if _, err := tx.ExecContext(ctx, `INSERT INTO user_roles (user_id, role_id) VALUES ($1, $2)`, userID, roleID); err != nil {
			return fmt.Errorf("insert role: %w", err)
		}
	}

	return tx.Commit()
}
