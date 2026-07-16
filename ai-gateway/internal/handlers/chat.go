package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opencode/erp-ai-gateway/internal/providers"
	"github.com/opencode/erp-ai-gateway/internal/router"
	"go.uber.org/zap"
)

type ChatHandler struct {
	logger *zap.Logger
	router *router.ProviderRouter
}

func NewChatHandler(logger *zap.Logger, router *router.ProviderRouter) *ChatHandler {
	return &ChatHandler{logger: logger, router: router}
}

type chatCompletionRequest struct {
	Model       string                 `json:"model"`
	Messages    []providers.Message    `json:"messages"`
	Tools       []providers.ToolDefinition `json:"tools,omitempty"`
	Temperature *float64               `json:"temperature,omitempty"`
	MaxTokens   *int                   `json:"max_tokens,omitempty"`
	TopP        *float64               `json:"top_p,omitempty"`
	Stream      bool                   `json:"stream"`
	Provider    string                 `json:"provider,omitempty"`
	TaskType    string                 `json:"task_type,omitempty"`
}

type completionResponse struct {
	ID        string           `json:"id"`
	Model     string           `json:"model"`
	Provider  string           `json:"provider"`
	Message   providers.Message `json:"message"`
	Usage     providers.Usage   `json:"usage"`
	CreatedAt time.Time        `json:"created_at"`
}

type embeddingRequest struct {
	Model  string   `json:"model"`
	Input  []string `json:"input"`
	Provider string `json:"provider,omitempty"`
}

type embeddingResponse struct {
	Model      string      `json:"model"`
	Embeddings [][]float32 `json:"embeddings"`
	Provider   string      `json:"provider"`
}

func (h *ChatHandler) ChatCompletions(c *gin.Context) {
	var req chatCompletionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid request: %v", err)})
		return
	}

	if len(req.Messages) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "messages are required"})
		return
	}

	var provider providers.Provider
	var err error

	if req.Provider != "" {
		provider, err = h.router.GetProvider(req.Provider)
	} else if req.TaskType != "" {
		provider, err = h.router.Route(c.Request.Context(), "", req.TaskType)
	} else {
		provider, err = h.router.Route(c.Request.Context(), "", "default")
	}

	if err != nil {
		h.logger.Error("failed to route provider", zap.Error(err))
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": fmt.Sprintf("no provider available: %v", err)})
		return
	}

	chatReq := &providers.ChatRequest{
		Model:       req.Model,
		Messages:    req.Messages,
		Tools:       req.Tools,
		Temperature: req.Temperature,
		MaxTokens:   req.MaxTokens,
		TopP:        req.TopP,
		Stream:      false,
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 120*time.Second)
	defer cancel()

	resp, err := provider.Chat(ctx, chatReq)
	if err != nil {
		h.logger.Error("chat completion failed", zap.Error(err), zap.String("provider", provider.Name()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("chat failed: %v", err)})
		return
	}

	c.JSON(http.StatusOK, completionResponse{
		ID:        resp.ID,
		Model:     resp.Model,
		Provider:  provider.Name(),
		Message:   resp.Message,
		Usage:     resp.Usage,
		CreatedAt: time.Now(),
	})
}

func (h *ChatHandler) ChatStream(c *gin.Context) {
	var req chatCompletionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid request: %v", err)})
		return
	}

	if len(req.Messages) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "messages are required"})
		return
	}

	var provider providers.Provider
	var err error

	if req.Provider != "" {
		provider, err = h.router.GetProvider(req.Provider)
	} else if req.TaskType != "" {
		provider, err = h.router.Route(c.Request.Context(), "", req.TaskType)
	} else {
		provider, err = h.router.Route(c.Request.Context(), "", "streaming")
	}

	if err != nil {
		h.logger.Error("failed to route provider for streaming", zap.Error(err))
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": fmt.Sprintf("no provider available: %v", err)})
		return
	}

	chatReq := &providers.ChatRequest{
		Model:       req.Model,
		Messages:    req.Messages,
		Tools:       req.Tools,
		Temperature: req.Temperature,
		MaxTokens:   req.MaxTokens,
		TopP:        req.TopP,
		Stream:      true,
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	streamCh := make(chan *providers.StreamingResponse, 10)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				h.logger.Error("panic in stream goroutine", zap.Any("recover", r))
			}
		}()
		if err := provider.ChatStream(c.Request.Context(), chatReq, streamCh); err != nil {
			h.logger.Error("stream error", zap.Error(err))
		}
	}()

	c.Stream(func(w io.Writer) bool {
		select {
		case msg, ok := <-streamCh:
			if !ok {
				return false
			}

			data, _ := json.Marshal(msg)
			fmt.Fprintf(w, "data: %s\n\n", string(data))

			if msg.Done {
				return false
			}
			return true
		case <-c.Request.Context().Done():
			return false
		}
	})
}

func (h *ChatHandler) Embeddings(c *gin.Context) {
	var req embeddingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid request: %v", err)})
		return
	}

	if len(req.Input) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "input is required"})
		return
	}

	var provider providers.Provider
	var err error

	if req.Provider != "" {
		provider, err = h.router.GetProvider(req.Provider)
	} else {
		provider, err = h.router.Route(c.Request.Context(), "", "embedding")
	}

	if err != nil {
		h.logger.Error("failed to route embedding provider", zap.Error(err))
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": fmt.Sprintf("no embedding provider available: %v", err)})
		return
	}

	embedReq := &providers.EmbeddingRequest{
		Model: req.Model,
		Input: req.Input,
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 60*time.Second)
	defer cancel()

	resp, err := provider.Embeddings(ctx, embedReq)
	if err != nil {
		h.logger.Error("embeddings failed", zap.Error(err), zap.String("provider", provider.Name()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("embedding failed: %v", err)})
		return
	}

	c.JSON(http.StatusOK, embeddingResponse{
		Model:      resp.Model,
		Embeddings: resp.Embeddings,
		Provider:   provider.Name(),
	})
}
