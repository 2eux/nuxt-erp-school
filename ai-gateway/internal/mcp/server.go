package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/opencode/erp-ai-gateway/internal/providers"
	"github.com/opencode/erp-ai-gateway/internal/router"
	"go.uber.org/zap"
)

type MCPServer struct {
	logger    *zap.Logger
	router    *router.ProviderRouter
	tools     map[string]*Tool
	resources map[string]*Resource
	prompts   map[string]*Prompt
	cache     CacheProvider
	mu        sync.RWMutex
}

type CacheProvider interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, ttl time.Duration) error
}

type Tool struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Schema      map[string]any `json:"inputSchema"`
	Handler     ToolHandler    `json:"-"`
	Category    string         `json:"category"`
}

type ToolHandler func(ctx context.Context, args map[string]any) (string, error)

type Resource struct {
	URI         string          `json:"uri"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	MimeType    string          `json:"mimeType"`
	Handler     ResourceHandler `json:"-"`
}

type ResourceHandler func(ctx context.Context, uri string) (string, string, error)

type Prompt struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Arguments   []PromptArg     `json:"arguments"`
	Content     string          `json:"-"`
	Handler     PromptHandler   `json:"-"`
}

type PromptArg struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Required    bool   `json:"required"`
}

type PromptHandler func(ctx context.Context, args map[string]string) (string, error)

type JSONRPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      any             `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type JSONRPCResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      any         `json:"id,omitempty"`
	Result  any         `json:"result,omitempty"`
	Error   *JSONRPCErr `json:"error,omitempty"`
}

type JSONRPCErr struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func NewMCPServer(logger *zap.Logger, router *router.ProviderRouter, cache CacheProvider) *MCPServer {
	s := &MCPServer{
		logger:    logger,
		router:    router,
		tools:     make(map[string]*Tool),
		resources: make(map[string]*Resource),
		prompts:   make(map[string]*Prompt),
		cache:     cache,
	}
	s.RegisterTools()
	s.RegisterResources()
	s.RegisterPrompts()
	return s
}

func (s *MCPServer) GetProvider(ctx context.Context) (providers.Provider, error) {
	return s.router.Route(ctx, "", "default")
}

func (s *MCPServer) GetIslamicProvider(ctx context.Context) (providers.Provider, error) {
	return s.router.Route(ctx, "ollama", "islamic")
}

func (s *MCPServer) RegisterTool(tool *Tool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tools[tool.Name] = tool
}

func (s *MCPServer) GetTool(name string) (*Tool, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.tools[name]
	return t, ok
}

func (s *MCPServer) ListTools() []*Tool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]*Tool, 0, len(s.tools))
	for _, t := range s.tools {
		result = append(result, t)
	}
	return result
}

func (s *MCPServer) RegisterResource(res *Resource) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.resources[res.URI] = res
}

func (s *MCPServer) GetResource(uri string) (*Resource, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.resources[uri]
	return r, ok
}

func (s *MCPServer) ListResources() []*Resource {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]*Resource, 0, len(s.resources))
	for _, r := range s.resources {
		result = append(result, r)
	}
	return result
}

func (s *MCPServer) RegisterPrompt(prompt *Prompt) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.prompts[prompt.Name] = prompt
}

func (s *MCPServer) GetPrompt(name string) (*Prompt, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, ok := s.prompts[name]
	return p, ok
}

func (s *MCPServer) ListPrompts() []*Prompt {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]*Prompt, 0, len(s.prompts))
	for _, p := range s.prompts {
		result = append(result, p)
	}
	return result
}

func (s *MCPServer) HandleJSONRPC(ctx context.Context, raw json.RawMessage) *JSONRPCResponse {
	var req JSONRPCRequest
	if err := json.Unmarshal(raw, &req); err != nil {
		return &JSONRPCResponse{
			JSONRPC: "2.0",
			Error: &JSONRPCErr{
				Code:    -32700,
				Message: "Parse error",
				Data:    err.Error(),
			},
		}
	}

	if req.JSONRPC != "2.0" {
		return &JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &JSONRPCErr{
				Code:    -32600,
				Message: "Invalid Request: jsonrpc must be 2.0",
			},
		}
	}

	switch req.Method {
	case "tools/list":
		return s.handleListTools(req.ID)
	case "tools/call":
		return s.handleToolCall(ctx, req.ID, req.Params)
	case "resources/list":
		return s.handleListResources(req.ID)
	case "resources/read":
		return s.handleReadResource(ctx, req.ID, req.Params)
	case "prompts/list":
		return s.handleListPrompts(req.ID)
	case "prompts/get":
		return s.handleGetPrompt(ctx, req.ID, req.Params)
	case "ping":
		return s.handlePing(req.ID)
	case "initialize":
		return s.handleInitialize(req.ID)
	default:
		return &JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &JSONRPCErr{
				Code:    -32601,
				Message: fmt.Sprintf("Method not found: %s", req.Method),
			},
		}
	}
}

func (s *MCPServer) handleListTools(id any) *JSONRPCResponse {
	tools := s.ListTools()
	type toolInfo struct {
		Name        string         `json:"name"`
		Description string         `json:"description"`
		InputSchema map[string]any `json:"inputSchema"`
	}
	info := make([]toolInfo, len(tools))
	for i, t := range tools {
		info[i] = toolInfo{
			Name:        t.Name,
			Description: t.Description,
			InputSchema: t.Schema,
		}
	}
	return &JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result:  map[string]any{"tools": info},
	}
}

func (s *MCPServer) handleToolCall(ctx context.Context, id any, params json.RawMessage) *JSONRPCResponse {
	var call struct {
		Name      string         `json:"name"`
		Arguments map[string]any `json:"arguments"`
	}
	if err := json.Unmarshal(params, &call); err != nil {
		return &JSONRPCResponse{
			JSONRPC: "2.0", ID: id,
			Error: &JSONRPCErr{Code: -32602, Message: "Invalid params", Data: err.Error()},
		}
	}

	tool, ok := s.GetTool(call.Name)
	if !ok {
		return &JSONRPCResponse{
			JSONRPC: "2.0", ID: id,
			Error: &JSONRPCErr{Code: -32602, Message: fmt.Sprintf("Tool not found: %s", call.Name)},
		}
	}

	s.logger.Info("MCP tool call", zap.String("tool", call.Name), zap.Any("args", call.Arguments))

	if call.Arguments == nil {
		call.Arguments = make(map[string]any)
	}

	result, err := tool.Handler(ctx, call.Arguments)
	if err != nil {
		return &JSONRPCResponse{
			JSONRPC: "2.0", ID: id,
			Error: &JSONRPCErr{Code: -32000, Message: err.Error()},
		}
	}

	return &JSONRPCResponse{
		JSONRPC: "2.0", ID: id,
		Result: map[string]any{
			"content": []map[string]string{
				{"type": "text", "text": result},
			},
		},
	}
}

func (s *MCPServer) handleListResources(id any) *JSONRPCResponse {
	resources := s.ListResources()
	type resInfo struct {
		URI         string `json:"uri"`
		Name        string `json:"name"`
		Description string `json:"description"`
		MimeType    string `json:"mimeType"`
	}
	info := make([]resInfo, len(resources))
	for i, r := range resources {
		info[i] = resInfo{
			URI:         r.URI,
			Name:        r.Name,
			Description: r.Description,
			MimeType:    r.MimeType,
		}
	}
	return &JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result:  map[string]any{"resources": info},
	}
}

func (s *MCPServer) handleReadResource(ctx context.Context, id any, params json.RawMessage) *JSONRPCResponse {
	var read struct {
		URI string `json:"uri"`
	}
	if err := json.Unmarshal(params, &read); err != nil {
		return &JSONRPCResponse{
			JSONRPC: "2.0", ID: id,
			Error: &JSONRPCErr{Code: -32602, Message: "Invalid params", Data: err.Error()},
		}
	}

	res, ok := s.GetResource(read.URI)
	if !ok {
		return &JSONRPCResponse{
			JSONRPC: "2.0", ID: id,
			Error: &JSONRPCErr{Code: -32602, Message: fmt.Sprintf("Resource not found: %s", read.URI)},
		}
	}

	content, mimeType, err := res.Handler(ctx, read.URI)
	if err != nil {
		return &JSONRPCResponse{
			JSONRPC: "2.0", ID: id,
			Error: &JSONRPCErr{Code: -32000, Message: err.Error()},
		}
	}

	return &JSONRPCResponse{
		JSONRPC: "2.0", ID: id,
		Result: map[string]any{
			"contents": []map[string]string{
				{
					"uri":      read.URI,
					"mimeType": mimeType,
					"text":     content,
				},
			},
		},
	}
}

func (s *MCPServer) handleListPrompts(id any) *JSONRPCResponse {
	prompts := s.ListPrompts()
	type promptInfo struct {
		Name        string      `json:"name"`
		Description string      `json:"description"`
		Arguments   []PromptArg `json:"arguments"`
	}
	info := make([]promptInfo, len(prompts))
	for i, p := range prompts {
		info[i] = promptInfo{
			Name:        p.Name,
			Description: p.Description,
			Arguments:   p.Arguments,
		}
	}
	return &JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result:  map[string]any{"prompts": info},
	}
}

func (s *MCPServer) handleGetPrompt(ctx context.Context, id any, params json.RawMessage) *JSONRPCResponse {
	var get struct {
		Name      string            `json:"name"`
		Arguments map[string]string `json:"arguments"`
	}
	if err := json.Unmarshal(params, &get); err != nil {
		return &JSONRPCResponse{
			JSONRPC: "2.0", ID: id,
			Error: &JSONRPCErr{Code: -32602, Message: "Invalid params", Data: err.Error()},
		}
	}

	prompt, ok := s.GetPrompt(get.Name)
	if !ok {
		return &JSONRPCResponse{
			JSONRPC: "2.0", ID: id,
			Error: &JSONRPCErr{Code: -32602, Message: fmt.Sprintf("Prompt not found: %s", get.Name)},
		}
	}

	if get.Arguments == nil {
		get.Arguments = make(map[string]string)
	}

	content, err := prompt.Handler(ctx, get.Arguments)
	if err != nil {
		return &JSONRPCResponse{
			JSONRPC: "2.0", ID: id,
			Error: &JSONRPCErr{Code: -32000, Message: err.Error()},
		}
	}

	return &JSONRPCResponse{
		JSONRPC: "2.0", ID: id,
		Result: map[string]any{
			"description": prompt.Description,
			"messages": []map[string]any{
				{
					"role": "user",
					"content": map[string]string{
						"type": "text",
						"text": content,
					},
				},
			},
		},
	}
}

func (s *MCPServer) handlePing(id any) *JSONRPCResponse {
	return &JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result:  map[string]any{},
	}
}

func (s *MCPServer) handleInitialize(id any) *JSONRPCResponse {
	return &JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result: map[string]any{
			"protocolVersion": "2024-11-05",
			"serverInfo": map[string]string{
				"name":    "erp-ai-gateway",
				"version": "1.0.0",
			},
			"capabilities": map[string]any{
				"tools":     map[string]bool{},
				"resources": map[string]bool{},
				"prompts":   map[string]bool{},
				"logging":   map[string]bool{},
			},
		},
	}
}

func (s *MCPServer) ExecuteTool(ctx context.Context, name string, args map[string]any) (string, error) {
	tool, ok := s.GetTool(name)
	if !ok {
		return "", fmt.Errorf("tool %s not found", name)
	}

	if args == nil {
		args = make(map[string]any)
	}

	return tool.Handler(ctx, args)
}
