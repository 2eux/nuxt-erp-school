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

type ClaudeProvider struct {
	cfg        config.ProviderConfig
	httpClient *http.Client
}

func NewClaudeProvider(cfg config.ProviderConfig) *ClaudeProvider {
	return &ClaudeProvider{
		cfg: cfg,
		httpClient: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

func (p *ClaudeProvider) Name() string { return "claude" }

func (p *ClaudeProvider) IsAvailable(ctx context.Context) bool {
	req, _ := http.NewRequestWithContext(ctx, "GET", p.cfg.BaseURL+"/v1/models", nil)
	req.Header.Set("x-api-key", p.cfg.APIKey)
	req.Header.Set("anthropic-version", "2023-06-01")
	resp, err := p.httpClient.Do(req)
	if err != nil {
		return false
	}
	resp.Body.Close()
	return resp.StatusCode == 200
}

func (p *ClaudeProvider) Models() []string {
	return []string{"claude-3-5-sonnet-20241022", "claude-3-opus-20240229", "claude-3-haiku-20240307"}
}

type claudeRequest struct {
	Model       string          `json:"model"`
	MaxTokens   int             `json:"max_tokens"`
	Messages    []claudeMessage `json:"messages"`
	System      string          `json:"system,omitempty"`
	Temperature *float64        `json:"temperature,omitempty"`
	TopP        *float64        `json:"top_p,omitempty"`
	Stream      bool            `json:"stream,omitempty"`
	StopSequences []string      `json:"stop_sequences,omitempty"`
}

type claudeMessage struct {
	Role    string        `json:"role"`
	Content []claudeBlock `json:"content"`
}

type claudeBlock struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
}

type claudeResponse struct {
	ID      string `json:"id"`
	Model   string `json:"model"`
	Type    string `json:"type"`
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text,omitempty"`
	}  `json:"content"`
	Usage struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
	StopReason string `json:"stop_reason"`
}

type claudeStreamEvent struct {
	Type  string `json:"type"`
	Index int    `json:"index"`
	Delta *struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"delta"`
	ContentBlock *struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content_block,omitempty"`
	Usage *struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"usage,omitempty"`
}

func (p *ClaudeProvider) Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	model := p.selectModel(req.Model)

	maxTokens := 4096
	if req.MaxTokens != nil {
		maxTokens = *req.MaxTokens
	}

	var systemPrompt string
	messages := make([]claudeMessage, 0)
	for _, msg := range req.Messages {
		if msg.Role == RoleSystem {
			systemPrompt += msg.Content + "\n"
			continue
		}

		role := "user"
		if msg.Role == RoleAssistant {
			role = "assistant"
		}

		messages = append(messages, claudeMessage{
			Role: role,
			Content: []claudeBlock{{
				Type: "text",
				Text: msg.Content,
			}},
		})
	}

	claudeReq := claudeRequest{
		Model:     model,
		MaxTokens: maxTokens,
		Messages:  messages,
		Stream:    false,
	}

	if systemPrompt != "" {
		claudeReq.System = strings.TrimSpace(systemPrompt)
	}
	if req.Temperature != nil {
		claudeReq.Temperature = req.Temperature
	}
	if req.TopP != nil {
		claudeReq.TopP = req.TopP
	}
	if len(req.Stop) > 0 {
		claudeReq.StopSequences = req.Stop
	}

	body, _ := json.Marshal(claudeReq)

	httpReq, _ := http.NewRequestWithContext(ctx, "POST", p.cfg.BaseURL+"/v1/messages", bytes.NewReader(body))
	httpReq.Header.Set("x-api-key", p.cfg.APIKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("claude request error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("claude API error %d: %s", resp.StatusCode, string(respBody))
	}

	var result claudeResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("claude decode error: %w", err)
	}

	var texts []string
	for _, c := range result.Content {
		if c.Text != "" {
			texts = append(texts, c.Text)
		}
	}

	usage := p.computeUsage(result.Usage.InputTokens, result.Usage.OutputTokens)

	return &ChatResponse{
		ID:       result.ID,
		Model:    result.Model,
		Provider: p.Name(),
		Message: Message{
			Role:    RoleAssistant,
			Content: strings.Join(texts, "\n"),
		},
		Usage:     usage,
		CreatedAt: time.Now(),
	}, nil
}

func (p *ClaudeProvider) ChatStream(ctx context.Context, req *ChatRequest, streamCh chan<- *StreamingResponse) error {
	defer close(streamCh)

	model := p.selectModel(req.Model)

	maxTokens := 4096
	if req.MaxTokens != nil {
		maxTokens = *req.MaxTokens
	}

	var systemPrompt string
	messages := make([]claudeMessage, 0)
	for _, msg := range req.Messages {
		if msg.Role == RoleSystem {
			systemPrompt += msg.Content + "\n"
			continue
		}
		role := "user"
		if msg.Role == RoleAssistant {
			role = "assistant"
		}
		messages = append(messages, claudeMessage{
			Role: role,
			Content: []claudeBlock{{
				Type: "text",
				Text: msg.Content,
			}},
		})
	}

	claudeReq := claudeRequest{
		Model:     model,
		MaxTokens: maxTokens,
		Messages:  messages,
		Stream:    true,
	}

	if systemPrompt != "" {
		claudeReq.System = strings.TrimSpace(systemPrompt)
	}
	if req.Temperature != nil {
		claudeReq.Temperature = req.Temperature
	}

	body, _ := json.Marshal(claudeReq)

	httpReq, _ := http.NewRequestWithContext(ctx, "POST", p.cfg.BaseURL+"/v1/messages", bytes.NewReader(body))
	httpReq.Header.Set("x-api-key", p.cfg.APIKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("claude stream request error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("claude stream error %d: %s", resp.StatusCode, string(respBody))
	}

	scanner := bufio.NewScanner(resp.Body)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	var fullText strings.Builder

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if !strings.HasPrefix(line, "data: ") {
			continue
		}

		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" {
			break
		}

		var event claudeStreamEvent
		if err := json.Unmarshal([]byte(data), &event); err != nil {
			continue
		}

		switch event.Type {
		case "content_block_delta":
			if event.Delta != nil && event.Delta.Text != "" {
				fullText.WriteString(event.Delta.Text)
				streamCh <- &StreamingResponse{
					ID:    fmt.Sprintf("claude-%d", time.Now().UnixNano()),
					Model: model,
					Delta: Message{Role: RoleAssistant, Content: event.Delta.Text},
				}
			}
		case "content_block_start":
			if event.ContentBlock != nil && event.ContentBlock.Text != "" {
				fullText.WriteString(event.ContentBlock.Text)
				streamCh <- &StreamingResponse{
					ID:    fmt.Sprintf("claude-%d", time.Now().UnixNano()),
					Model: model,
					Delta: Message{Role: RoleAssistant, Content: event.ContentBlock.Text},
				}
			}
		case "message_delta":
			if event.Usage != nil {
				streamCh <- &StreamingResponse{
					Done: true,
					Usage: &Usage{
						InputTokens:     event.Usage.InputTokens,
						OutputTokens:    event.Usage.OutputTokens,
						CompletionTokens: fullText.Len(),
						TotalTokens:      event.Usage.InputTokens + fullText.Len(),
					},
				}
			}
		case "message_stop":
			streamCh <- &StreamingResponse{
				Done: true,
				Usage: &Usage{
					CompletionTokens: fullText.Len(),
					TotalTokens:      fullText.Len(),
				},
			}
		}
	}

	streamCh <- &StreamingResponse{
		Done: true,
		Usage: &Usage{
			CompletionTokens: fullText.Len(),
			TotalTokens:      fullText.Len(),
		},
	}

	return nil
}

func (p *ClaudeProvider) Embeddings(ctx context.Context, req *EmbeddingRequest) (*EmbeddingResponse, error) {
	return nil, fmt.Errorf("claude does not support embeddings, use openai or gemini")
}

func (p *ClaudeProvider) selectModel(requested string) string {
	if requested != "" {
		return requested
	}
	return p.cfg.ChatModel
}

func (p *ClaudeProvider) computeUsage(inputTokens, outputTokens int) Usage {
	cost := (float64(inputTokens)/1000)*p.cfg.CostPer1KIn +
		(float64(outputTokens)/1000)*p.cfg.CostPer1KOut
	return Usage{
		PromptTokens:     inputTokens,
		CompletionTokens: outputTokens,
		TotalTokens:      inputTokens + outputTokens,
		Cost:             cost,
	}
}
