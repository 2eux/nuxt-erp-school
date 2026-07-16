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

type NotificationService interface {
	ListNotifications(ctx context.Context, userID string, filter dto.PaginationRequest) ([]dto.NotificationResponse, int64, error)
	MarkAsRead(ctx context.Context, id, userID string) error
	MarkAllAsRead(ctx context.Context, userID string) error
	CreateNotification(ctx context.Context, schoolID string, req dto.CreateNotificationRequest) error
}

type notificationService struct {
	db     *sqlx.DB
	logger *zap.Logger
}

func NewNotificationService(db *sqlx.DB, logger *zap.Logger) NotificationService {
	return &notificationService{db: db, logger: logger}
}

func (s *notificationService) ListNotifications(ctx context.Context, userID string, filter dto.PaginationRequest) ([]dto.NotificationResponse, int64, error) {
	filter.Defaults()

	var items []domain.Notification
	if err := s.db.SelectContext(ctx, &items, database.ListNotifications, userID, filter.PageSize, filter.Offset()); err != nil {
		return nil, 0, domain.NewInternalError("failed to list notifications", err)
	}

	var total int64
	s.db.GetContext(ctx, &total, database.CountNotifications, userID)

	result := make([]dto.NotificationResponse, len(items))
	for i, n := range items {
		result[i] = dto.NotificationResponse{
			ID:        n.ID,
			Title:     n.Title,
			Message:   n.Message,
			Type:      n.Type,
			RefType:   n.RefType,
			RefID:     n.RefID,
			IsRead:    n.IsRead,
			ReadAt:    n.ReadAt,
			CreatedAt: n.CreatedAt,
		}
	}
	return result, total, nil
}

func (s *notificationService) MarkAsRead(ctx context.Context, id, userID string) error {
	now := time.Now()
	query := `UPDATE notifications SET is_read=true, read_at=$1 WHERE id=$2 AND user_id=$3`
	if _, err := s.db.ExecContext(ctx, query, now, id, userID); err != nil {
		return domain.NewInternalError("failed to mark as read", err)
	}
	return nil
}

func (s *notificationService) MarkAllAsRead(ctx context.Context, userID string) error {
	now := time.Now()
	query := `UPDATE notifications SET is_read=true, read_at=$1 WHERE user_id=$2 AND is_read=false`
	if _, err := s.db.ExecContext(ctx, query, now, userID); err != nil {
		return domain.NewInternalError("failed to mark all as read", err)
	}
	return nil
}

func (s *notificationService) CreateNotification(ctx context.Context, schoolID string, req dto.CreateNotificationRequest) error {
	tx, err := s.db.Beginx()
	if err != nil {
		return domain.NewInternalError("failed to begin tx", err)
	}
	defer tx.Rollback()

	for _, userID := range req.UserIDs {
		n := &domain.Notification{
			ID:      uuid.New().String(),
			UserID:  userID,
			Title:   req.Title,
			Message: req.Message,
			Type:    req.Type,
			RefType: req.RefType,
			RefID:   req.RefID,
		}
		query := `INSERT INTO notifications (id, user_id, title, message, type, ref_type, ref_id) VALUES (:id, :user_id, :title, :message, :type, :ref_type, :ref_id)`
		if _, err := tx.NamedExecContext(ctx, query, n); err != nil {
			return domain.NewInternalError("failed to create notification", err)
		}
	}

	return tx.Commit()
}
