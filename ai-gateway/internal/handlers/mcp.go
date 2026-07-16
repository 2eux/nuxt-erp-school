package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opencode/erp-ai-gateway/internal/mcp"
	"go.uber.org/zap"
)

type MCPHandler struct {
	logger *zap.Logger
	server *mcp.MCPServer
}

func NewMCPHandler(logger *zap.Logger, server *mcp.MCPServer) *MCPHandler {
	return &MCPHandler{logger: logger, server: server}
}

func (h *MCPHandler) HandleJSONRPC(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read body"})
		return
	}

	h.logger.Debug("MCP JSON-RPC request", zap.String("body", string(body)))

	response := h.server.HandleJSONRPC(c.Request.Context(), body)

	c.JSON(http.StatusOK, response)
}

func (h *MCPHandler) HandleSSE(c *gin.Context) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "streaming not supported"})
		return
	}

	fmt.Fprintf(c.Writer, "event: endpoint\ndata: %s\n\n", "/api/v1/mcp/messages")
	flusher.Flush()

	<-c.Request.Context().Done()
}

func (h *MCPHandler) ListTools(c *gin.Context) {
	tools := h.server.ListTools()

	type toolInfo struct {
		Name        string         `json:"name"`
		Description string         `json:"description"`
		InputSchema map[string]any `json:"inputSchema"`
		Category    string         `json:"category"`
	}

	result := make([]toolInfo, len(tools))
	for i, t := range tools {
		result[i] = toolInfo{
			Name:        t.Name,
			Description: t.Description,
			InputSchema: t.Schema,
			Category:    t.Category,
		}
	}

	c.JSON(http.StatusOK, gin.H{"tools": result, "count": len(result)})
}

func (h *MCPHandler) CallTool(c *gin.Context) {
	var req struct {
		Name      string         `json:"name"`
		Arguments map[string]any `json:"arguments"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid request: %v", err)})
		return
	}

	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tool name is required"})
		return
	}

	result, err := h.server.ExecuteTool(c.Request.Context(), req.Name, req.Arguments)
	if err != nil {
		h.logger.Error("tool execution failed", zap.String("tool", req.Name), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tool":    req.Name,
		"content": result,
	})
}

func (h *MCPHandler) ListResources(c *gin.Context) {
	resources := h.server.ListResources()

	type resInfo struct {
		URI         string `json:"uri"`
		Name        string `json:"name"`
		Description string `json:"description"`
		MimeType    string `json:"mimeType"`
	}

	result := make([]resInfo, len(resources))
	for i, r := range resources {
		result[i] = resInfo{
			URI:         r.URI,
			Name:        r.Name,
			Description: r.Description,
			MimeType:    r.MimeType,
		}
	}

	c.JSON(http.StatusOK, gin.H{"resources": result, "count": len(result)})
}

func (h *MCPHandler) ReadResource(c *gin.Context) {
	uri := c.Query("uri")
	if uri == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "uri query parameter is required"})
		return
	}

	res, ok := h.server.GetResource(uri)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("resource not found: %s", uri)})
		return
	}

	content, mimeType, err := res.Handler(c.Request.Context(), uri)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"uri":      uri,
		"mimeType": mimeType,
		"content":  content,
	})
}

func (h *MCPHandler) ListPrompts(c *gin.Context) {
	prompts := h.server.ListPrompts()

	type promptInfo struct {
		Name        string          `json:"name"`
		Description string          `json:"description"`
		Arguments   []mcp.PromptArg `json:"arguments"`
	}

	result := make([]promptInfo, len(prompts))
	for i, p := range prompts {
		result[i] = promptInfo{
			Name:        p.Name,
			Description: p.Description,
			Arguments:   p.Arguments,
		}
	}

	c.JSON(http.StatusOK, gin.H{"prompts": result, "count": len(result)})
}

func (h *MCPHandler) GetPrompt(c *gin.Context) {
	var req struct {
		Name      string            `json:"name"`
		Arguments map[string]string `json:"arguments"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid request: %v", err)})
		return
	}

	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "prompt name is required"})
		return
	}

	prompt, ok := h.server.GetPrompt(req.Name)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("prompt not found: %s", req.Name)})
		return
	}

	if req.Arguments == nil {
		req.Arguments = make(map[string]string)
	}

	result, err := prompt.Handler(c.Request.Context(), req.Arguments)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"prompt":      req.Name,
		"description": prompt.Description,
		"content":     result,
	})
}
