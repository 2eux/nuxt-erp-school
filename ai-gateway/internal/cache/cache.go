package cache

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/opencode/erp-ai-gateway/internal/config"
	"go.uber.org/zap"
)

type Cache struct {
	client *redis.Client
	cfg    config.RedisConfig
	logger *zap.Logger
}

func NewCache(cfg config.RedisConfig, logger *zap.Logger) (*Cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		logger.Warn("redis connection failed, running without cache", zap.Error(err))
		return &Cache{
			client: nil,
			cfg:    cfg,
			logger: logger,
		}, nil
	}

	logger.Info("redis connected successfully", zap.String("addr", fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)))

	return &Cache{
		client: client,
		cfg:    cfg,
		logger: logger,
	}, nil
}

func (c *Cache) IsAvailable() bool {
	return c.client != nil
}

func (c *Cache) Get(ctx context.Context, key string) (string, error) {
	if c.client == nil {
		return "", fmt.Errorf("cache not available")
	}

	result, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return result, nil
}

func (c *Cache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	if c.client == nil {
		return nil
	}
	return c.client.Set(ctx, key, value, ttl).Err()
}

func (c *Cache) Delete(ctx context.Context, key string) error {
	if c.client == nil {
		return nil
	}
	return c.client.Del(ctx, key).Err()
}

func (c *Cache) DeletePattern(ctx context.Context, pattern string) error {
	if c.client == nil {
		return nil
	}

	iter := c.client.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		if err := c.client.Del(ctx, iter.Val()).Err(); err != nil {
			c.logger.Warn("failed to delete cache key", zap.String("key", iter.Val()), zap.Error(err))
		}
	}
	return iter.Err()
}

func (c *Cache) GetJSON(ctx context.Context, key string, dest interface{}) error {
	if c.client == nil {
		return fmt.Errorf("cache not available")
	}

	result, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return fmt.Errorf("key not found: %s", key)
	}
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(result), dest)
}

func (c *Cache) SetJSON(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	if c.client == nil {
		return nil
	}

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, key, string(data), ttl).Err()
}

type SemanticCacheEntry struct {
	OriginalQuery  string        `json:"original_query"`
	Response       string        `json:"response"`
	Embedding      []float32     `json:"embedding"`
	CachedAt       time.Time     `json:"cached_at"`
	AccessCount    int           `json:"access_count"`
}

func (c *Cache) GetSemantic(ctx context.Context, queryEmbedding []float32) (*SemanticCacheEntry, error) {
	if c.client == nil || !c.cfg.SemanticCacheEnabled {
		return nil, fmt.Errorf("semantic cache not available")
	}

	pattern := "semantic:*"
	iter := c.client.Scan(ctx, 0, pattern, 0).Iterator()

	var bestMatch *SemanticCacheEntry
	bestSimilarity := float64(0)

	for iter.Next(ctx) {
		data, err := c.client.Get(ctx, iter.Val()).Result()
		if err != nil {
			continue
		}

		var entry SemanticCacheEntry
		if err := json.Unmarshal([]byte(data), &entry); err != nil {
			continue
		}

		similarity := cosineSimilarity(queryEmbedding, entry.Embedding)
		if similarity > bestSimilarity && similarity >= c.cfg.SemanticCacheThresh {
			bestSimilarity = similarity
			entryCopy := entry
			bestMatch = &entryCopy
		}
	}

	if bestMatch != nil {
		bestMatch.AccessCount++
		return bestMatch, nil
	}

	return nil, fmt.Errorf("no semantic match found")
}

func (c *Cache) SetSemantic(ctx context.Context, query string, response string, embedding []float32, ttl time.Duration) error {
	if c.client == nil || !c.cfg.SemanticCacheEnabled {
		return nil
	}

	entry := SemanticCacheEntry{
		OriginalQuery: query,
		Response:      response,
		Embedding:     embedding,
		CachedAt:      time.Now(),
		AccessCount:   0,
	}

	key := fmt.Sprintf("semantic:%x", sha256.Sum256([]byte(query)))
	return c.SetJSON(ctx, key, entry, ttl)
}

func (c *Cache) CacheChatResponse(ctx context.Context, key string, response interface{}, ttl time.Duration) error {
	chatKey := fmt.Sprintf("chat:%s", key)
	return c.SetJSON(ctx, chatKey, response, ttl)
}

func (c *Cache) GetChatResponse(ctx context.Context, key string) (string, error) {
	chatKey := fmt.Sprintf("chat:%s", key)
	return c.Get(ctx, chatKey)
}

func (c *Cache) CacheEmbedding(ctx context.Context, text string, embedding []float32, ttl time.Duration) error {
	textHash := fmt.Sprintf("%x", sha256.Sum256([]byte(text)))
	embedKey := fmt.Sprintf("embed:%s", textHash)
	return c.SetJSON(ctx, embedKey, embedding, ttl)
}

func (c *Cache) GetCachedEmbedding(ctx context.Context, text string) ([]float32, error) {
	textHash := fmt.Sprintf("%x", sha256.Sum256([]byte(text)))
	embedKey := fmt.Sprintf("embed:%s", textHash)

	var embedding []float32
	if err := c.GetJSON(ctx, embedKey, &embedding); err != nil {
		return nil, err
	}
	return embedding, nil
}

func (c *Cache) Close() error {
	if c.client == nil {
		return nil
	}
	return c.client.Close()
}

func cosineSimilarity(a, b []float32) float64 {
	if len(a) != len(b) || len(a) == 0 {
		return 0
	}

	var dotProduct, normA, normB float64
	for i := range a {
		dotProduct += float64(a[i]) * float64(b[i])
		normA += float64(a[i]) * float64(a[i])
		normB += float64(b[i]) * float64(b[i])
	}

	if normA == 0 || normB == 0 {
		return 0
	}

	return dotProduct / (float64Sqrt(normA) * float64Sqrt(normB))
}

func float64Sqrt(x float64) float64 {
	if x <= 0 {
		return 0
	}
	z := x
	for i := 0; i < 100; i++ {
		z = z - (z*z-x)/(2*z)
	}
	return z
}
