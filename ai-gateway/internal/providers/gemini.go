package providers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/generative-ai-go/genai"
	"github.com/opencode/erp-ai-gateway/internal/config"
	"google.golang.org/api/option"
)

type GeminiProvider struct {
	cfg    config.ProviderConfig
	client *genai.Client
}

func NewGeminiProvider(cfg config.ProviderConfig) (*GeminiProvider, error) {
	client, err := genai.NewClient(context.Background(), option.WithAPIKey(cfg.APIKey))
	if err != nil {
		return nil, fmt.Errorf("gemini client init failed: %w", err)
	}
	return &GeminiProvider{
		cfg:    cfg,
		client: client,
	}, nil
}

func (p *GeminiProvider) Name() string { return "gemini" }

func (p *GeminiProvider) IsAvailable(ctx context.Context) bool {
	return p.client != nil
}

func (p *GeminiProvider) Models() []string {
	return []string{"gemini-1.5-pro", "gemini-1.5-flash", "gemini-1.0-pro", "text-embedding-004"}
}

func (p *GeminiProvider) Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	model := p.selectModel(req.Model)
	gm := p.client.GenerativeModel(model)
	p.configureModel(gm, req)

	contents := p.convertContents(req.Messages)
	ctx, cancel := context.WithTimeout(ctx, 120*time.Second)
	defer cancel()

	resp, err := gm.GenerateContent(ctx, contents...)
	if err != nil {
		return nil, fmt.Errorf("gemini chat error: %w", err)
	}

	if len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil {
		return nil, fmt.Errorf("gemini returned no candidates")
	}

	content := resp.Candidates[0].Content
	text := p.extractText(content)

	tokenCount := int(gm.SafetySettings[0].Threshold)

	return &ChatResponse{
		ID:       fmt.Sprintf("gemini-%d", time.Now().UnixNano()),
		Model:    model,
		Provider: p.Name(),
		Message: Message{
			Role:    RoleAssistant,
			Content: text,
		},
		Usage: Usage{
			PromptTokens:     tokenCount,
			CompletionTokens: len(strings.Fields(text)),
			TotalTokens:      tokenCount + len(strings.Fields(text)),
			Cost:             0,
		},
		CreatedAt: time.Now(),
	}, nil
}

func (p *GeminiProvider) ChatStream(ctx context.Context, req *ChatRequest, streamCh chan<- *StreamingResponse) error {
	defer close(streamCh)

	model := p.selectModel(req.Model)
	gm := p.client.GenerativeModel(model)
	p.configureModel(gm, req)

	contents := p.convertContents(req.Messages)

	ctx, cancel := context.WithTimeout(ctx, 120*time.Second)
	defer cancel()

	iter := gm.GenerateContentStream(ctx, contents...)
	var fullText strings.Builder

	for {
		resp, err := iter.Next()
		if err != nil {
			if strings.Contains(err.Error(), "iterator done") {
				streamCh <- &StreamingResponse{
					Done: true,
					Usage: &Usage{
						CompletionTokens: fullText.Len(),
						TotalTokens:      fullText.Len(),
					},
				}
				return nil
			}
			streamCh <- &StreamingResponse{Error: err.Error()}
			return fmt.Errorf("gemini stream error: %w", err)
		}

		if len(resp.Candidates) > 0 && resp.Candidates[0].Content != nil {
			text := p.extractText(resp.Candidates[0].Content)
			fullText.WriteString(text)
			streamCh <- &StreamingResponse{
				ID:    fmt.Sprintf("gemini-%d", time.Now().UnixNano()),
				Model: model,
				Delta: Message{Role: RoleAssistant, Content: text},
			}
		}
	}
}

func (p *GeminiProvider) Embeddings(ctx context.Context, req *EmbeddingRequest) (*EmbeddingResponse, error) {
	model := p.cfg.EmbeddingModel
	if model == "" {
		model = "text-embedding-004"
	}

	var input []string
	if len(req.Input) > 0 {
		input = req.Input
	} else {
		input = req.Inputs
	}
	if len(input) == 0 {
		return nil, fmt.Errorf("no input for embeddings")
	}

	em := p.client.EmbeddingModel(model)

	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	embeddings := make([][]float32, 0)

	for _, text := range input {
		result, err := em.EmbedContent(ctx, genai.Text(text))
		if err != nil {
			return nil, fmt.Errorf("gemini embedding error: %w", err)
		}
		embeddings = append(embeddings, result.Embedding.Values)
	}

	return &EmbeddingResponse{
		Model:      model,
		Embeddings: embeddings,
		Usage:      Usage{TotalTokens: len(input)},
	}, nil
}

func (p *GeminiProvider) selectModel(requested string) string {
	if requested != "" {
		return requested
	}
	return p.cfg.ChatModel
}

func (p *GeminiProvider) configureModel(gm *genai.GenerativeModel, req *ChatRequest) {
	if req.Temperature != nil {
		gm.SetTemperature(float32(*req.Temperature))
	}
	if req.MaxTokens != nil {
		gm.SetMaxOutputTokens(int32(*req.MaxTokens))
	}
	if req.TopP != nil {
		gm.SetTopP(float32(*req.TopP))
	}

	gm.SafetySettings = []*genai.SafetySetting{
		{
			Category:  genai.HarmCategoryHateSpeech,
			Threshold: genai.HarmBlockNone,
		},
		{
			Category:  genai.HarmCategoryHarassment,
			Threshold: genai.HarmBlockNone,
		},
	}
}

func (p *GeminiProvider) convertContents(messages []Message) []*genai.Content {
	contents := make([]*genai.Content, 0)

	for _, msg := range messages {
		role := "user"
		if msg.Role == RoleAssistant || msg.Role == "model" {
			role = "model"
		} else if msg.Role == RoleSystem {
			role = "user"
		}

		content := &genai.Content{
			Role: role,
			Parts: []genai.Part{
				genai.Text(msg.Content),
			},
		}
		contents = append(contents, content)
	}

	return contents
}

func (p *GeminiProvider) extractText(content *genai.Content) string {
	var sb strings.Builder
	for _, part := range content.Parts {
		switch v := part.(type) {
		case genai.Text:
			sb.WriteString(string(v))
		case *genai.Text:
			sb.WriteString(string(*v))
		}
	}
	return sb.String()
}
