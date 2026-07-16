package rag

import (
	"context"
	"fmt"
	"time"

	"github.com/opencode/erp-ai-gateway/internal/providers"
	"go.uber.org/zap"
)

type EmbeddingService struct {
	logger *zap.Logger
	cache  EmbeddingCache
}

type EmbeddingCache interface {
	GetCachedEmbedding(ctx context.Context, text string) ([]float32, error)
	CacheEmbedding(ctx context.Context, text string, embedding []float32, ttl time.Duration) error
}

func NewEmbeddingService(logger *zap.Logger, cache EmbeddingCache) *EmbeddingService {
	return &EmbeddingService{
		logger: logger,
		cache:  cache,
	}
}

func (s *EmbeddingService) GenerateEmbedding(ctx context.Context, provider providers.Provider, text string) ([]float32, error) {
	if s.cache != nil {
		if cached, err := s.cache.GetCachedEmbedding(ctx, text); err == nil && len(cached) > 0 {
			return cached, nil
		}
	}

	resp, err := provider.Embeddings(ctx, &providers.EmbeddingRequest{
		Input: []string{text},
	})
	if err != nil {
		return nil, fmt.Errorf("embedding generation failed: %w", err)
	}

	if len(resp.Embeddings) == 0 {
		return nil, fmt.Errorf("no embeddings returned")
	}

	embedding := resp.Embeddings[0]

	if s.cache != nil {
		_ = s.cache.CacheEmbedding(ctx, text, embedding, 24*time.Hour)
	}

	return embedding, nil
}

func (s *EmbeddingService) GenerateBatchEmbeddings(ctx context.Context, provider providers.Provider, texts []string) ([][]float32, error) {
	if len(texts) == 0 {
		return nil, fmt.Errorf("no texts provided")
	}

	results := make([][]float32, len(texts))
	uncachedTexts := make([]string, 0)
	uncachedIndices := make([]int, 0)

	if s.cache != nil {
		for i, text := range texts {
			if cached, err := s.cache.GetCachedEmbedding(ctx, text); err == nil && len(cached) > 0 {
				results[i] = cached
			} else {
				uncachedTexts = append(uncachedTexts, text)
				uncachedIndices = append(uncachedIndices, i)
			}
		}
	} else {
		uncachedTexts = texts
		for i := range texts {
			uncachedIndices = append(uncachedIndices, i)
		}
	}

	if len(uncachedTexts) == 0 {
		return results, nil
	}

	resp, err := provider.Embeddings(ctx, &providers.EmbeddingRequest{
		Input: uncachedTexts,
	})
	if err != nil {
		return nil, fmt.Errorf("batch embedding failed: %w", err)
	}

	for i, idx := range uncachedIndices {
		if i < len(resp.Embeddings) {
			results[idx] = resp.Embeddings[i]
		}
	}

	if s.cache != nil {
		for i, idx := range uncachedIndices {
			if i < len(resp.Embeddings) {
				_ = s.cache.CacheEmbedding(ctx, texts[idx], resp.Embeddings[i], 24*time.Hour)
			}
		}
	}

	return results, nil
}
