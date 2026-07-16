package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opencode/erp-ai-gateway/internal/providers"
	"github.com/opencode/erp-ai-gateway/internal/rag"
	"github.com/opencode/erp-ai-gateway/internal/router"
	"go.uber.org/zap"
)

type RAGHandler struct {
	logger  *zap.Logger
	router  *router.ProviderRouter
	service *rag.RAGService
}

func NewRAGHandler(logger *zap.Logger, router *router.ProviderRouter, service *rag.RAGService) *RAGHandler {
	return &RAGHandler{
		logger:  logger,
		router:  router,
		service: service,
	}
}

type documentUploadRequest struct {
	Title       string            `json:"title"`
	Content     string            `json:"content"`
	Metadata    map[string]string `json:"metadata,omitempty"`
	ChunkSize   int               `json:"chunk_size,omitempty"`
	ChunkOverlap int              `json:"chunk_overlap,omitempty"`
}

type ragQueryRequest struct {
	Query       string            `json:"query"`
	TopK        int               `json:"top_k,omitempty"`
	Filter      map[string]string `json:"filter,omitempty"`
	HybridMode  bool              `json:"hybrid_mode,omitempty"`
	Model       string            `json:"model,omitempty"`
	Temperature *float64          `json:"temperature,omitempty"`
}

type ragQueryResponse struct {
	Query    string               `json:"query"`
	Answer   string               `json:"answer"`
	Sources  []rag.SearchResult   `json:"sources"`
	Provider string               `json:"provider"`
	Usage    providers.Usage      `json:"usage,omitempty"`
}

func (h *RAGHandler) UploadDocument(c *gin.Context) {
	var req documentUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid request: %v", err)})
		return
	}

	if req.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
		return
	}

	if req.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "content is required"})
		return
	}

	provider, err := h.router.Route(c.Request.Context(), "", "embedding")
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": fmt.Sprintf("no embedding provider: %v", err)})
		return
	}

	doc := &rag.Document{
		Title:    req.Title,
		Content:  req.Content,
		Metadata: req.Metadata,
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 120*time.Second)
	defer cancel()

	docID, err := h.service.IngestDocument(ctx, provider, doc)
	if err != nil {
		h.logger.Error("document ingestion failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("ingestion failed: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"id":     docID,
		"title":  req.Title,
	})
}

func (h *RAGHandler) UploadFile(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("file is required: %v", err)})
		return
	}
	defer file.Close()

	title := c.PostForm("title")
	if title == "" {
		title = header.Filename
	}

	metadataStr := c.PostForm("metadata")
	metadata := make(map[string]string)
	if metadataStr != "" {
		var parsed map[string]string
		if jsonErr := json.Unmarshal([]byte(metadataStr), &parsed); jsonErr == nil {
			metadata = parsed
		}
	}

	buf := make([]byte, header.Size)
	n, err := file.Read(buf)
	if err != nil && err.Error() != "EOF" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to read file: %v", err)})
		return
	}
	buf = buf[:n]

	extractedText, err := rag.ExtractText(buf, header.Filename)
	if err != nil {
		h.logger.Error("text extraction failed", zap.Error(err), zap.String("filename", header.Filename))
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("text extraction failed: %v", err)})
		return
	}

	if extractedText == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no extractable text found in document"})
		return
	}

	provider, err := h.router.Route(c.Request.Context(), "", "embedding")
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": fmt.Sprintf("no embedding provider: %v", err)})
		return
	}

	doc := &rag.Document{
		Title:    title,
		Content:  extractedText,
		Metadata: metadata,
	}
	if metadata == nil {
		doc.Metadata = make(map[string]string)
	}
	doc.Metadata["filename"] = header.Filename
	doc.Metadata["file_type"] = header.Header.Get("Content-Type")

	ctx, cancel := context.WithTimeout(c.Request.Context(), 120*time.Second)
	defer cancel()

	docID, err := h.service.IngestDocument(ctx, provider, doc)
	if err != nil {
		h.logger.Error("document ingestion failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("ingestion failed: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":         "success",
		"id":             docID,
		"title":          title,
		"extracted_text_length": len(extractedText),
	})
}

func (h *RAGHandler) Query(c *gin.Context) {
	var req ragQueryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid request: %v", err)})
		return
	}

	if req.Query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query is required"})
		return
	}

	topK := req.TopK
	if topK <= 0 {
		topK = 5
	}

	embedProvider, err := h.router.Route(c.Request.Context(), "", "embedding")
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": fmt.Sprintf("no embedding provider: %v", err)})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 120*time.Second)
	defer cancel()

	embedResp, err := embedProvider.Embeddings(ctx, &providers.EmbeddingRequest{
		Input: []string{req.Query},
	})
	if err != nil {
		h.logger.Error("query embedding failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("embedding failed: %v", err)})
		return
	}

	if len(embedResp.Embeddings) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no embeddings returned"})
		return
	}

	searchResults, err := h.service.Search(ctx, embedResp.Embeddings[0], topK)
	if err != nil {
		h.logger.Error("vector search failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("search failed: %v", err)})
		return
	}

	llmProvider, err := h.router.Route(c.Request.Context(), "", "default")
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": fmt.Sprintf("no LLM provider: %v", err)})
		return
	}

	systemPrompt := `You are an AI assistant for an Islamic school. Answer questions based on the provided context documents.
If the answer cannot be found in the context, say so clearly and suggest where the user might find the information.
Always cite which document source you used in your response.
Be accurate, professional, and helpful.

Context documents:
`

	for i, result := range searchResults {
		systemPrompt += fmt.Sprintf("\n[Document %d: %s]\n%s\n", i+1, result.Title, result.Content)
	}

	temperature := 0.3
	if req.Temperature != nil {
		temperature = *req.Temperature
	}
	maxTokens := 2048

	llmResp, err := llmProvider.Chat(ctx, &providers.ChatRequest{
		Messages: []providers.Message{
			{Role: providers.RoleSystem, Content: systemPrompt},
			{Role: providers.RoleUser, Content: req.Query},
		},
		Temperature: &temperature,
		MaxTokens:   &maxTokens,
	})
	if err != nil {
		h.logger.Error("LLM query failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("LLM call failed: %v", err)})
		return
	}

	c.JSON(http.StatusOK, ragQueryResponse{
		Query:    req.Query,
		Answer:   llmResp.Message.Content,
		Sources:  searchResults,
		Provider: llmProvider.Name(),
		Usage:    llmResp.Usage,
	})
}

func (h *RAGHandler) ListDocuments(c *gin.Context) {
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	docs, total, err := h.service.ListDocuments(c.Request.Context(), offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to list documents: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"documents": docs,
		"total":     total,
		"offset":    offset,
		"limit":     limit,
	})
}

func (h *RAGHandler) DeleteDocument(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "document id is required"})
		return
	}

	if err := h.service.DeleteDocument(c.Request.Context(), id); err != nil {
		h.logger.Error("document deletion failed", zap.Error(err), zap.String("id", id))
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("deletion failed: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "deleted", "id": id})
}
