package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Env      string
	Port     string
	JWT      JWTConfig
	APIKey   string
	AdminKey string

	Redis RedisConfig

	Providers map[string]ProviderConfig
	Enabled   []string

	DefaultProvider string
	FallbackChain   []string

	Qdrant QdrantConfig

	RAG RAGConfig

	RateLimit RateLimitConfig

	LogLevel  string
	LogFormat string
}

type JWTConfig struct {
	Secret string
}

type RedisConfig struct {
	Host                  string
	Port                  string
	Password              string
	DB                    int
	CacheTTL              time.Duration
	SemanticCacheEnabled  bool
	SemanticCacheThresh   float64
}

type ProviderConfig struct {
	Name           string
	APIKey         string
	BaseURL        string
	ChatModel      string
	FastModel      string
	EmbeddingModel string
	BestModel      string
	CheapModel     string
	RateLimitRPM   int
	CostPer1KIn    float64
	CostPer1KOut   float64
	Extra          map[string]string
}

type QdrantConfig struct {
	Host       string
	Port       string
	APIKey     string
	Collection string
}

type RAGConfig struct {
	ChunkSize           int
	ChunkOverlap        int
	TopK                int
	SimilarityThreshold float64
}

type RateLimitConfig struct {
	PerUser       int
	PerTenant     int
	WindowSeconds int
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		Env:  getEnv("AI_GATEWAY_ENV", "development"),
		Port: getEnv("AI_GATEWAY_PORT", "8081"),
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "erp-ai-gateway-secret"),
		},
		APIKey:   getEnv("API_KEY", ""),
		AdminKey: getEnv("ADMIN_API_KEY", ""),

		Redis: RedisConfig{
			Host:                 getEnv("REDIS_HOST", "localhost"),
			Port:                 getEnv("REDIS_PORT", "6379"),
			Password:             getEnv("REDIS_PASSWORD", ""),
			DB:                   getEnvInt("REDIS_DB", 1),
			CacheTTL:             time.Duration(getEnvInt("CACHE_TTL_SECONDS", 3600)) * time.Second,
			SemanticCacheEnabled: getEnvBool("SEMANTIC_CACHE_ENABLED", true),
			SemanticCacheThresh:  getEnvFloat("SEMANTIC_CACHE_THRESHOLD", 0.92),
		},

		DefaultProvider: getEnv("DEFAULT_PROVIDER", "openai"),
		FallbackChain:   strings.Split(getEnv("FALLBACK_CHAIN", "openai,gemini,claude,ollama"), ","),

		Qdrant: QdrantConfig{
			Host:       getEnv("QDRANT_HOST", "localhost"),
			Port:       getEnv("QDRANT_PORT", "6334"),
			APIKey:     getEnv("QDRANT_API_KEY", ""),
			Collection: getEnv("QDRANT_COLLECTION", "erp_documents"),
		},

		RAG: RAGConfig{
			ChunkSize:           getEnvInt("RAG_CHUNK_SIZE", 1500),
			ChunkOverlap:        getEnvInt("RAG_CHUNK_OVERLAP", 200),
			TopK:                getEnvInt("RAG_TOP_K", 5),
			SimilarityThreshold: getEnvFloat("RAG_SIMILARITY_THRESHOLD", 0.7),
		},

		RateLimit: RateLimitConfig{
			PerUser:       getEnvInt("RATE_LIMIT_PER_USER", 100),
			PerTenant:     getEnvInt("RATE_LIMIT_PER_TENANT", 1000),
			WindowSeconds: getEnvInt("RATE_LIMIT_WINDOW_SECONDS", 60),
		},

		LogLevel:  getEnv("LOG_LEVEL", "debug"),
		LogFormat: getEnv("LOG_FORMAT", "json"),
	}

	cfg.Enabled = strings.Split(getEnv("ENABLED_PROVIDERS", "openai,gemini,claude,ollama"), ",")

	cfg.Providers = map[string]ProviderConfig{
		"openai": {
			Name:           "openai",
			APIKey:         getEnv("OPENAI_API_KEY", ""),
			BaseURL:        getEnv("OPENAI_BASE_URL", "https://api.openai.com/v1"),
			ChatModel:      getEnv("OPENAI_MODEL_CHAT", "gpt-4o"),
			FastModel:      getEnv("OPENAI_MODEL_CHEAP", "gpt-4o-mini"),
			EmbeddingModel: getEnv("OPENAI_MODEL_EMBEDDING", "text-embedding-3-large"),
			CheapModel:     getEnv("OPENAI_MODEL_CHEAP", "gpt-4o-mini"),
			RateLimitRPM:   getEnvInt("OPENAI_RATE_LIMIT", 1000),
			CostPer1KIn:    getEnvFloat("OPENAI_COST_PER_1K_INPUT", 0.005),
			CostPer1KOut:   getEnvFloat("OPENAI_COST_PER_1K_OUTPUT", 0.015),
		},
		"gemini": {
			Name:           "gemini",
			APIKey:         getEnv("GEMINI_API_KEY", ""),
			ChatModel:      getEnv("GEMINI_MODEL_CHAT", "gemini-1.5-pro"),
			FastModel:      getEnv("GEMINI_MODEL_FAST", "gemini-1.5-flash"),
			EmbeddingModel: getEnv("GEMINI_MODEL_EMBEDDING", "text-embedding-004"),
			RateLimitRPM:   getEnvInt("GEMINI_RATE_LIMIT", 1500),
			CostPer1KIn:    getEnvFloat("GEMINI_COST_PER_1K_INPUT", 0.0),
			CostPer1KOut:   getEnvFloat("GEMINI_COST_PER_1K_OUTPUT", 0.0),
		},
		"claude": {
			Name:           "claude",
			APIKey:         getEnv("ANTHROPIC_API_KEY", ""),
			BaseURL:        getEnv("ANTHROPIC_BASE_URL", "https://api.anthropic.com"),
			ChatModel:      getEnv("ANTHROPIC_MODEL_CHAT", "claude-3-5-sonnet-20241022"),
			BestModel:      getEnv("ANTHROPIC_MODEL_BEST", "claude-3-opus-20240229"),
			RateLimitRPM:   getEnvInt("ANTHROPIC_RATE_LIMIT", 500),
			CostPer1KIn:    getEnvFloat("ANTHROPIC_COST_PER_1K_INPUT", 0.003),
			CostPer1KOut:   getEnvFloat("ANTHROPIC_COST_PER_1K_OUTPUT", 0.015),
		},
		"ollama": {
			Name:           "ollama",
			BaseURL:        getEnv("OLLAMA_BASE_URL", "http://localhost:11434"),
			ChatModel:      getEnv("OLLAMA_MODEL_CHAT", "llama3.1:8b"),
			EmbeddingModel: getEnv("OLLAMA_MODEL_CHAT", "nomic-embed-text"),
			RateLimitRPM:   999999,
			CostPer1KIn:    0.0,
			CostPer1KOut:   0.0,
			Extra: map[string]string{
				"islamic_model": getEnv("OLLAMA_MODEL_ISLAMIC", "hijaz-ai-arabic-llm"),
				"quran_model":   getEnv("OLLAMA_MODEL_QURAN", "quran-tafsir-model"),
			},
		},
	}

	return cfg, nil
}

func (c *Config) GetProvider(name string) (ProviderConfig, error) {
	provider, ok := c.Providers[name]
	if !ok {
		return ProviderConfig{}, fmt.Errorf("provider %s not configured", name)
	}
	return provider, nil
}

func (c *Config) IsProviderEnabled(name string) bool {
	for _, p := range c.Enabled {
		if strings.EqualFold(strings.TrimSpace(p), name) {
			return true
		}
	}
	return false
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if val, ok := os.LookupEnv(key); ok {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return fallback
}

func getEnvFloat(key string, fallback float64) float64 {
	if val, ok := os.LookupEnv(key); ok {
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			return f
		}
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	if val, ok := os.LookupEnv(key); ok {
		val = strings.ToLower(strings.TrimSpace(val))
		if val == "true" || val == "1" || val == "yes" {
			return true
		}
		if val == "false" || val == "0" || val == "no" {
			return false
		}
	}
	return fallback
}
