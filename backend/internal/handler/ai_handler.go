package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opencode/erp-school-backend/internal/dto"
	"github.com/opencode/erp-school-backend/internal/middleware"
	"github.com/opencode/erp-school-backend/internal/service"
)

type AIHandler struct {
	aiService service.AIService
}

func NewAIHandler(aiService service.AIService) *AIHandler {
	return &AIHandler{aiService: aiService}
}

// Chat godoc
// @Summary AI Chat
// @Tags AI
// @Accept json
// @Produce json
// @Param request body dto.AIChatRequest true "Chat message"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/ai/chat [post]
func (h *AIHandler) Chat(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req dto.AIChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "invalid request", err.Error()))
		return
	}
	resp, err := h.aiService.SendMessage(c.Request.Context(), req.ConversationID, userID, req)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "success", resp))
}

// Generate godoc
// @Summary AI Generate
// @Tags AI
// @Accept json
// @Produce json
// @Param request body dto.AIGenerateRequest true "Generate request"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/ai/generate [post]
func (h *AIHandler) Generate(c *gin.Context) {
	userID := middleware.GetUserID(c)
	schoolID := middleware.GetSchoolID(c)
	var req dto.AIGenerateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "invalid request", err.Error()))
		return
	}
	resp, err := h.aiService.Generate(c.Request.Context(), userID, schoolID, req)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "success", resp))
}

// ListConversations godoc
// @Summary List AI conversations
// @Tags AI
// @Produce json
// @Param page query int false "Page"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/ai/conversations [get]
func (h *AIHandler) ListConversations(c *gin.Context) {
	userID := middleware.GetUserID(c)
	filter := getPagination(c)
	items, total, err := h.aiService.GetConversations(c.Request.Context(), userID, filter)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "success", dto.NewPaginatedResponse(items, total, filter.Page, filter.PageSize)))
}

// GetMessages godoc
// @Summary Get conversation messages
// @Tags AI
// @Produce json
// @Param id path string true "Conversation ID"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/ai/conversations/{id}/messages [get]
func (h *AIHandler) GetMessages(c *gin.Context) {
	convID := c.Param("id")
	items, err := h.aiService.GetMessages(c.Request.Context(), convID)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "success", items))
}

// SendMessage godoc
// @Summary Send message to conversation
// @Tags AI
// @Accept json
// @Produce json
// @Param id path string true "Conversation ID"
// @Param request body dto.AIChatRequest true "Message data"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/ai/conversations/{id}/messages [post]
func (h *AIHandler) SendMessage(c *gin.Context) {
	userID := middleware.GetUserID(c)
	convID := c.Param("id")
	var req dto.AIChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "invalid request", err.Error()))
		return
	}
	resp, err := h.aiService.SendMessage(c.Request.Context(), convID, userID, req)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "success", resp))
}

// UploadKnowledge godoc
// @Summary Upload knowledge document
// @Tags AI
// @Accept json
// @Produce json
// @Param request body object true "Knowledge doc {title, doc_type, module, content}"
// @Success 201 {object} dto.APIResponse
// @Router /api/v1/ai/knowledge/upload [post]
func (h *AIHandler) UploadKnowledge(c *gin.Context) {
	schoolID := middleware.GetSchoolID(c)
	createdBy := middleware.GetUserID(c)
	var req struct {
		Title   string `json:"title" binding:"required"`
		DocType string `json:"doc_type" binding:"required"`
		Module  string `json:"module" binding:"required"`
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "invalid request", err.Error()))
		return
	}
	resp, err := h.aiService.UploadKnowledge(c.Request.Context(), schoolID, createdBy, req.Title, req.DocType, req.Module, req.Content)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusCreated, dto.NewAPIResponse(http.StatusCreated, "knowledge uploaded", resp))
}
