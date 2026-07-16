package providers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/opencode/erp-ai-gateway/internal/config"
	openai "github.com/sashabaranov/go-openai"
)

type OpenAIProvider struct {
	cfg    config.ProviderConfig
	client *openai.Client
}

func NewOpenAIProvider(cfg config.ProviderConfig) *OpenAIProvider {
	clientCfg := openai.DefaultConfig(cfg.APIKey)
	if cfg.BaseURL != "" {
		clientCfg.BaseURL = cfg.BaseURL
	}
	return &OpenAIProvider{
		cfg:    cfg,
		client: openai.NewClientWithConfig(clientCfg),
	}
}

func (p *OpenAIProvider) Name() string { return "openai" }

func (p *OpenAIProvider) IsAvailable(ctx context.Context) bool {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	_, err := p.client.ListModels(ctx)
	return err == nil
}

func (p *OpenAIProvider) Models() []string {
	return []string{"gpt-4o", "gpt-4o-mini", "gpt-4-turbo", "gpt-4", "gpt-3.5-turbo", "text-embedding-3-large", "text-embedding-3-small"}
}

func (p *OpenAIProvider) Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	model := p.selectModel(req.Model)

	msgs := p.convertMessages(req.Messages)
	tools := p.convertTools(req.Tools)

	openaiReq := openai.ChatCompletionRequest{
		Model:    model,
		Messages: msgs,
		Tools:    tools,
	}

	if req.Temperature != nil {
		openaiReq.Temperature = float32(*req.Temperature)
	}
	if req.MaxTokens != nil {
		openaiReq.MaxTokens = *req.MaxTokens
	}
	if req.TopP != nil {
		openaiReq.TopP = float32(*req.TopP)
	}
	if len(req.Stop) > 0 {
		openaiReq.Stop = req.Stop
	}

	if req.ToolChoice != "" {
		openaiReq.ToolChoice = req.ToolChoice
	}

	ctx, cancel := context.WithTimeout(ctx, 120*time.Second)
	defer cancel()

	resp, err := p.client.CreateChatCompletion(ctx, openaiReq)
	if err != nil {
		return nil, fmt.Errorf("openai chat error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("openai returned no choices")
	}

	choice := resp.Choices[0]
	msg := Message{
		Role:    RoleAssistant,
		Content: choice.Message.Content,
	}

	for _, tc := range choice.Message.ToolCalls {
		msg.ToolCalls = append(msg.ToolCalls, ToolCall{
			ID:   tc.ID,
			Type: string(tc.Type),
			Function: FunctionCall{
				Name:      tc.Function.Name,
				Arguments: tc.Function.Arguments,
			},
		})
	}

	usage := p.computeUsage(resp.Usage, false)

	return &ChatResponse{
		ID:        resp.ID,
		Model:     resp.Model,
		Provider:  p.Name(),
		Message:   msg,
		Usage:     usage,
		CreatedAt: time.Now(),
	}, nil
}

func (p *OpenAIProvider) ChatStream(ctx context.Context, req *ChatRequest, streamCh chan<- *StreamingResponse) error {
	defer close(streamCh)

	model := p.selectModel(req.Model)
	msgs := p.convertMessages(req.Messages)
	tools := p.convertTools(req.Tools)

	openaiReq := openai.ChatCompletionRequest{
		Model:    model,
		Messages: msgs,
		Tools:    tools,
		Stream:   true,
	}

	if req.Temperature != nil {
		openaiReq.Temperature = float32(*req.Temperature)
	}
	if req.MaxTokens != nil {
		openaiReq.MaxTokens = *req.MaxTokens
	}

	ctx, cancel := context.WithTimeout(ctx, 120*time.Second)
	defer cancel()

	stream, err := p.client.CreateChatCompletionStream(ctx, openaiReq)
	if err != nil {
		return fmt.Errorf("openai stream error: %w", err)
	}
	defer stream.Close()

	var fullContent strings.Builder
	var fullToolCalls []ToolCall
	promptTokens := 0

	for {
		response, err := stream.Recv()
		if err != nil {
			if strings.Contains(err.Error(), "EOF") {
				streamCh <- &StreamingResponse{Done: true, Usage: &Usage{
					PromptTokens: promptTokens, CompletionTokens: fullContent.Len(), TotalTokens: promptTokens + fullContent.Len(),
				}}
				return nil
			}
			streamCh <- &StreamingResponse{Error: err.Error()}
			return fmt.Errorf("stream recv error: %w", err)
		}

		if len(response.Choices) == 0 {
			continue
		}

		delta := response.Choices[0].Delta
		if delta.Content != "" {
			fullContent.WriteString(delta.Content)
			streamCh <- &StreamingResponse{
				ID:    response.ID,
				Model: response.Model,
				Delta: Message{Role: RoleAssistant, Content: delta.Content},
			}
		}

		for _, tc := range delta.ToolCalls {
			fullToolCalls = append(fullToolCalls, ToolCall{
				ID:   tc.ID,
				Type: string(tc.Type),
				Function: FunctionCall{
					Name:      tc.Function.Name,
					Arguments: tc.Function.Arguments,
				},
			})

			streamCh <- &StreamingResponse{
				ID:    response.ID,
				Model: response.Model,
				Delta: Message{
					Role: RoleAssistant,
					ToolCalls: []ToolCall{{
						ID:   tc.ID,
						Type: string(tc.Type),
						Function: FunctionCall{
							Name:      tc.Function.Name,
							Arguments: tc.Function.Arguments,
						},
					}},
				},
			}
		}
	}
}

func (p *OpenAIProvider) Embeddings(ctx context.Context, req *EmbeddingRequest) (*EmbeddingResponse, error) {
	model := p.cfg.EmbeddingModel
	if model == "" {
		model = "text-embedding-3-large"
	}

	var input []string
	if len(req.Input) > 0 {
		input = req.Input
	} else {
		input = req.Inputs
	}

	if len(input) == 0 {
		return nil, fmt.Errorf("no input provided for embeddings")
	}

	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	embedReq := openai.EmbeddingRequest{
		Model: openai.EmbeddingModel(model),
		Input: input,
	}

	resp, err := p.client.CreateEmbeddings(ctx, embedReq)
	if err != nil {
		return nil, fmt.Errorf("openai embeddings error: %w", err)
	}

	embeddings := make([][]float32, len(resp.Data))
	for i, d := range resp.Data {
		embeddings[i] = d.Embedding
	}

	return &EmbeddingResponse{
		Model:      string(resp.Model),
		Embeddings: embeddings,
		Usage: Usage{
			PromptTokens: resp.Usage.PromptTokens,
			TotalTokens:  resp.Usage.TotalTokens,
		},
	}, nil
}

func (p *OpenAIProvider) selectModel(requested string) string {
	if requested != "" {
		return requested
	}
	return p.cfg.ChatModel
}

func (p *OpenAIProvider) convertMessages(messages []Message) []openai.ChatCompletionMessage {
	result := make([]openai.ChatCompletionMessage, len(messages))
	for i, msg := range messages {
		m := openai.ChatCompletionMessage{
			Role:    string(msg.Role),
			Content: msg.Content,
			Name:    msg.Name,
		}

		for _, tc := range msg.ToolCalls {
			m.ToolCalls = append(m.ToolCalls, openai.ToolCall{
				ID:   tc.ID,
				Type: openai.ToolType(tc.Type),
				Function: openai.FunctionCall{
					Name:      tc.Function.Name,
					Arguments: tc.Function.Arguments,
				},
			})
		}

		if msg.ToolCallID != "" {
			m.ToolCallID = msg.ToolCallID
		}

		result[i] = m
	}
	return result
}

func (p *OpenAIProvider) convertTools(tools []ToolDefinition) []openai.Tool {
	result := make([]openai.Tool, len(tools))
	for i, t := range tools {
		result[i] = openai.Tool{
			Type: openai.ToolType(t.Type),
			Function: &openai.FunctionDefinition{
				Name:        t.Function.Name,
				Description: t.Function.Description,
				Parameters:  t.Function.Parameters,
			},
		}
	}
	return result
}

func (p *OpenAIProvider) computeUsage(usage openai.Usage, isStream bool) Usage {
	u := Usage{
		PromptTokens:     usage.PromptTokens,
		CompletionTokens: usage.CompletionTokens,
		TotalTokens:      usage.TotalTokens,
	}
	u.Cost = (float64(u.PromptTokens)/1000)*p.cfg.CostPer1KIn +
		(float64(u.CompletionTokens)/1000)*p.cfg.CostPer1KOut
	return u
}
