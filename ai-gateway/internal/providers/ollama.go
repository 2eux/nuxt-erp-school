package providers

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/opencode/erp-ai-gateway/internal/config"
)

type OllamaProvider struct {
	cfg        config.ProviderConfig
	httpClient *http.Client
}

func NewOllamaProvider(cfg config.ProviderConfig) *OllamaProvider {
	return &OllamaProvider{
		cfg: cfg,
		httpClient: &http.Client{
			Timeout: 300 * time.Second,
		},
	}
}

func (p *OllamaProvider) Name() string { return "ollama" }

func (p *OllamaProvider) IsAvailable(ctx context.Context) bool {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	req, _ := http.NewRequestWithContext(ctx, "GET", p.cfg.BaseURL+"/api/tags", nil)
	resp, err := p.httpClient.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == 200
}

func (p *OllamaProvider) Models() []string {
	models := []string{"llama3.1:8b", "llama3.1:70b", "mistral:7b", "mistral:latest", "codellama:latest", "nomic-embed-text"}
	if islamicModel, ok := p.cfg.Extra["islamic_model"]; ok && islamicModel != "" {
		models = append(models, islamicModel)
	}
	if quranModel, ok := p.cfg.Extra["quran_model"]; ok && quranModel != "" {
		models = append(models, quranModel)
	}
	return models
}

func (p *OllamaProvider) getIslamicModel() string {
	if model, ok := p.cfg.Extra["islamic_model"]; ok && model != "" {
		return model
	}
	return p.cfg.ChatModel
}

func (p *OllamaProvider) getQuranModel() string {
	if model, ok := p.cfg.Extra["quran_model"]; ok && model != "" {
		return model
	}
	return p.cfg.ChatModel
}

type ollamaChatRequest struct {
	Model    string          `json:"model"`
	Messages []ollamaMessage `json:"messages"`
	Stream   bool            `json:"stream"`
	Options  ollamaOptions   `json:"options,omitempty"`
}

type ollamaMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ollamaOptions struct {
	Temperature float64 `json:"temperature,omitempty"`
	NumPredict  int     `json:"num_predict,omitempty"`
	TopP        float64 `json:"top_p,omitempty"`
	Seed        int     `json:"seed,omitempty"`
}

type ollamaChatResponse struct {
	Model     string         `json:"model"`
	CreatedAt time.Time      `json:"created_at"`
	Message   ollamaMessage  `json:"message"`
	Done      bool           `json:"done"`
	TotalDuration int64      `json:"total_duration"`
	EvalCount     int        `json:"eval_count"`
	PromptEvalCount int      `json:"prompt_eval_count"`
}

type ollamaEmbedRequest struct {
	Model  string `json:"model"`
	Input  any    `json:"input"`
}

type ollamaEmbedResponse struct {
	Model      string      `json:"model"`
	Embeddings [][]float32 `json:"embeddings"`
}

func (p *OllamaProvider) Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	model := p.selectModel(req.Model)

	messages := make([]ollamaMessage, 0, len(req.Messages))
	for _, msg := range req.Messages {
		role := "user"
		switch msg.Role {
		case RoleSystem:
			role = "system"
		case RoleAssistant:
			role = "assistant"
		case RoleUser:
			role = "user"
		}
		messages = append(messages, ollamaMessage{
			Role:    role,
			Content: msg.Content,
		})
	}

	ollReq := ollamaChatRequest{
		Model:    model,
		Messages: messages,
		Stream:   false,
		Options: ollamaOptions{
			NumPredict: 4096,
		},
	}

	if req.Temperature != nil {
		ollReq.Options.Temperature = *req.Temperature
	} else {
		ollReq.Options.Temperature = 0.7
	}
	if req.MaxTokens != nil {
		ollReq.Options.NumPredict = *req.MaxTokens
	}
	if req.TopP != nil {
		ollReq.Options.TopP = *req.TopP
	}

	body, _ := json.Marshal(ollReq)

	ctx, cancel := context.WithTimeout(ctx, 300*time.Second)
	defer cancel()

	httpReq, _ := http.NewRequestWithContext(ctx, "POST", p.cfg.BaseURL+"/api/chat", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("ollama request error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ollama API error %d: %s", resp.StatusCode, string(respBody))
	}

	var result ollamaChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("ollama decode error: %w", err)
	}

	return &ChatResponse{
		ID:       fmt.Sprintf("ollama-%d", time.Now().UnixNano()),
		Model:    result.Model,
		Provider: p.Name(),
		Message: Message{
			Role:    RoleAssistant,
			Content: result.Message.Content,
		},
		Usage: Usage{
			PromptTokens:     result.PromptEvalCount,
			CompletionTokens: result.EvalCount,
			TotalTokens:      result.PromptEvalCount + result.EvalCount,
			Cost:             0,
		},
		CreatedAt: time.Now(),
	}, nil
}

func (p *OllamaProvider) ChatStream(ctx context.Context, req *ChatRequest, streamCh chan<- *StreamingResponse) error {
	defer close(streamCh)

	model := p.selectModel(req.Model)

	messages := make([]ollamaMessage, 0, len(req.Messages))
	for _, msg := range req.Messages {
		role := "user"
		switch msg.Role {
		case RoleSystem:
			role = "system"
		case RoleAssistant:
			role = "assistant"
		case RoleUser:
			role = "user"
		}
		messages = append(messages, ollamaMessage{
			Role:    role,
			Content: msg.Content,
		})
	}

	ollReq := ollamaChatRequest{
		Model:    model,
		Messages: messages,
		Stream:   true,
		Options: ollamaOptions{
			NumPredict: 4096,
		},
	}

	if req.Temperature != nil {
		ollReq.Options.Temperature = *req.Temperature
	} else {
		ollReq.Options.Temperature = 0.7
	}
	if req.MaxTokens != nil {
		ollReq.Options.NumPredict = *req.MaxTokens
	}

	body, _ := json.Marshal(ollReq)

	ctx, cancel := context.WithTimeout(ctx, 300*time.Second)
	defer cancel()

	httpReq, _ := http.NewRequestWithContext(ctx, "POST", p.cfg.BaseURL+"/api/chat", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("ollama stream request error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("ollama stream error %d: %s", resp.StatusCode, string(respBody))
	}

	scanner := bufio.NewScanner(resp.Body)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		var chunk ollamaChatResponse
		if err := json.Unmarshal([]byte(line), &chunk); err != nil {
			continue
		}

		if chunk.Message.Content != "" {
			streamCh <- &StreamingResponse{
				ID:    fmt.Sprintf("ollama-%d", time.Now().UnixNano()),
				Model: chunk.Model,
				Delta: Message{Role: RoleAssistant, Content: chunk.Message.Content},
			}
		}

		if chunk.Done {
			streamCh <- &StreamingResponse{
				Done:  true,
				Usage: &Usage{CompletionTokens: chunk.EvalCount, TotalTokens: chunk.EvalCount},
			}
			return nil
		}
	}

	streamCh <- &StreamingResponse{
		Done: true,
		Usage: &Usage{},
	}

	return nil
}

func (p *OllamaProvider) Embeddings(ctx context.Context, req *EmbeddingRequest) (*EmbeddingResponse, error) {
	model := p.cfg.EmbeddingModel
	if model == "" {
		model = "nomic-embed-text"
	}

	var input any
	if len(req.Input) > 0 {
		if len(req.Input) == 1 {
			input = req.Input[0]
		} else {
			input = req.Input
		}
	} else if len(req.Inputs) > 0 {
		if len(req.Inputs) == 1 {
			input = req.Inputs[0]
		} else {
			input = req.Inputs
		}
	}

	if input == nil {
		return nil, fmt.Errorf("no input for embeddings")
	}

	embedReq := ollamaEmbedRequest{
		Model: model,
		Input: input,
	}

	body, _ := json.Marshal(embedReq)

	ctx, cancel := context.WithTimeout(ctx, 120*time.Second)
	defer cancel()

	httpReq, _ := http.NewRequestWithContext(ctx, "POST", p.cfg.BaseURL+"/api/embed", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("ollama embed error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ollama embed error %d: %s", resp.StatusCode, string(respBody))
	}

	var result ollamaEmbedResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("ollama embed decode error: %w", err)
	}

	return &EmbeddingResponse{
		Model:      result.Model,
		Embeddings: result.Embeddings,
		Usage:      Usage{},
	}, nil
}

func (p *OllamaProvider) selectModel(requested string) string {
	if requested != "" {
		return requested
	}
	if p.cfg.ChatModel != "" {
		return p.cfg.ChatModel
	}
	return "llama3.1:8b"
}

func (p *OllamaProvider) SelectIslamicModel() string {
	return p.getIslamicModel()
}

func (p *OllamaProvider) SelectQuranModel() string {
	return p.getQuranModel()
}

func (p *OllamaProvider) ChatWithIslamicModel(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	m := p.getIslamicModel()
	modifiedReq := *req
	modifiedReq.Model = m
	return p.Chat(ctx, &modifiedReq)
}

func (p *OllamaProvider) ChatWithQuranModel(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	m := p.getQuranModel()
	modifiedReq := *req
	modifiedReq.Model = m
	return p.Chat(ctx, &modifiedReq)
}
