package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/opencode/erp-school-backend/internal/domain"
	"github.com/opencode/erp-school-backend/internal/dto"
	"github.com/opencode/erp-school-backend/internal/infrastructure/database"
	"go.uber.org/zap"
)

type AIService interface {
	CreateConversation(ctx context.Context, userID, schoolID, title, model string) (*dto.AIConversationResponse, error)
	GetConversations(ctx context.Context, userID string, filter dto.PaginationRequest) ([]dto.AIConversationResponse, int64, error)
	SendMessage(ctx context.Context, conversationID, userID string, req dto.AIChatRequest) (*dto.AIChatResponse, error)
	GetMessages(ctx context.Context, conversationID string) ([]dto.AIMessageResponse, error)
	Generate(ctx context.Context, userID, schoolID string, req dto.AIGenerateRequest) (*dto.AIGeneratedResponse, error)
	UploadKnowledge(ctx context.Context, schoolID, createdBy, title, docType, module, content string) (*domain.KnowledgeDocument, error)
}

type aiService struct {
	db     *sqlx.DB
	logger *zap.Logger
}

func NewAIService(db *sqlx.DB, logger *zap.Logger) AIService {
	return &aiService{db: db, logger: logger}
}

func (s *aiService) CreateConversation(ctx context.Context, userID, schoolID, title, model string) (*dto.AIConversationResponse, error) {
	if model == "" {
		model = "gpt-4o"
	}

	conv := &domain.AIConversation{
		ID:       uuid.New().String(),
		UserID:   userID,
		SchoolID: schoolID,
		Title:    title,
		Model:    model,
	}

	query := `INSERT INTO ai_conversations (id, user_id, school_id, title, model) VALUES (:id, :user_id, :school_id, :title, :model)`
	if _, err := s.db.NamedExecContext(ctx, query, conv); err != nil {
		return nil, domain.NewInternalError("failed to create conversation", err)
	}

	return &dto.AIConversationResponse{
		ID:        conv.ID,
		Title:     conv.Title,
		Model:     conv.Model,
		CreatedAt: conv.CreatedAt,
		UpdatedAt: conv.UpdatedAt,
	}, nil
}

func (s *aiService) GetConversations(ctx context.Context, userID string, filter dto.PaginationRequest) ([]dto.AIConversationResponse, int64, error) {
	filter.Defaults()

	var items []struct {
		domain.AIConversation
		MessageCount int `db:"message_count"`
	}

	var total int64
	s.db.GetContext(ctx, &total, `SELECT COUNT(*) FROM ai_conversations WHERE user_id=$1`, userID)

	query := database.ListAIConversations
	if err := s.db.SelectContext(ctx, &items, query, userID, filter.PageSize, filter.Offset()); err != nil {
		return nil, 0, domain.NewInternalError("failed to list conversations", err)
	}

	result := make([]dto.AIConversationResponse, len(items))
	for i, c := range items {
		result[i] = dto.AIConversationResponse{
			ID:           c.ID,
			Title:        c.Title,
			Model:        c.Model,
			MessageCount: c.MessageCount,
			CreatedAt:    c.CreatedAt,
			UpdatedAt:    c.UpdatedAt,
		}
	}
	return result, total, nil
}

func (s *aiService) SendMessage(ctx context.Context, conversationID, userID string, req dto.AIChatRequest) (*dto.AIChatResponse, error) {
	if conversationID == "" {
		conv, err := s.CreateConversation(ctx, userID, "", "New Chat", req.Model)
		if err != nil {
			return nil, err
		}
		conversationID = conv.ID
	}

	userMsg := &domain.AIMessage{
		ID:             uuid.New().String(),
		ConversationID: conversationID,
		Role:           "user",
		Content:        req.Message,
	}

	msgQuery := `INSERT INTO ai_messages (id, conversation_id, role, content) VALUES (:id, :conversation_id, :role, :content)`
	s.db.NamedExecContext(ctx, msgQuery, userMsg)

	response := "I understand your message: " + req.Message

	assistantMsg := &domain.AIMessage{
		ID:             uuid.New().String(),
		ConversationID: conversationID,
		Role:           "assistant",
		Content:        response,
		TokenCount:     len(response) / 4,
	}
	s.db.NamedExecContext(ctx, msgQuery, assistantMsg)

	return &dto.AIChatResponse{
		ConversationID: conversationID,
		Message:        response,
		Role:           "assistant",
		TokenCount:     assistantMsg.TokenCount,
	}, nil
}

func (s *aiService) GetMessages(ctx context.Context, conversationID string) ([]dto.AIMessageResponse, error) {
	var items []domain.AIMessage
	if err := s.db.SelectContext(ctx, &items, database.ListAIMessages, conversationID); err != nil {
		return nil, domain.NewInternalError("failed to list messages", err)
	}

	result := make([]dto.AIMessageResponse, len(items))
	for i, m := range items {
		result[i] = dto.AIMessageResponse{
			ID:         m.ID,
			Role:       m.Role,
			Content:    m.Content,
			TokenCount: m.TokenCount,
			CreatedAt:  m.CreatedAt,
		}
	}
	return result, nil
}

func (s *aiService) Generate(ctx context.Context, userID, schoolID string, req dto.AIGenerateRequest) (*dto.AIGeneratedResponse, error) {
	conv, err := s.CreateConversation(ctx, userID, schoolID, req.Type+" generation", req.Model)
	if err != nil {
		return nil, err
	}

	var content interface{}
	content = map[string]interface{}{
		"title":    "Generated " + req.Type,
		"type":     req.Type,
		"template": "This is a placeholder generated content for " + req.Type,
	}

	return &dto.AIGeneratedResponse{
		Type:    req.Type,
		Content: content,
	}, err
}

func (s *aiService) UploadKnowledge(ctx context.Context, schoolID, createdBy, title, docType, module, content string) (*domain.KnowledgeDocument, error) {
	doc := &domain.KnowledgeDocument{
		ID:              uuid.New().String(),
		SchoolID:        schoolID,
		Title:           title,
		Content:         content,
		DocType:         docType,
		Module:          module,
		EmbeddingStatus: "pending",
		CreatedBy:       createdBy,
	}

	query := `INSERT INTO knowledge_documents (id, school_id, title, content, doc_type, module, embedding_status, created_by) VALUES (:id, :school_id, :title, :content, :doc_type, :module, :embedding_status, :created_by)`
	if _, err := s.db.NamedExecContext(ctx, query, doc); err != nil {
		return nil, domain.NewInternalError("failed to upload knowledge", err)
	}

	return doc, nil
}
