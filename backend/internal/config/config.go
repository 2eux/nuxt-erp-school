package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig
	DB       DatabaseConfig
	Redis    RedisConfig
	MinIO    MinIOConfig
	JWT      JWTConfig
	AI       AIConfig
	SMTP     SMTPConfig
	SMS      SMSConfig
	Push     PushConfig
	WhatsApp WhatsAppConfig
}

type AppConfig struct {
	Name  string
	Env   string
	Port  int
	Debug bool
}

type DatabaseConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	Name            string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
	PoolSize int
}

type MinIOConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
	Region    string
}

type JWTConfig struct {
	AccessSecret  string
	RefreshSecret string
	AccessTTL     time.Duration
	RefreshTTL    time.Duration
	Issuer        string
}

type AIConfig struct {
	Provider        string
	OpenAIKey       string
	OpenAIEndpoint  string
	GeminiKey       string
	GeminiEndpoint  string
	ClaudeKey       string
	ClaudeEndpoint  string
	OllamaEndpoint  string
	QdrantHost      string
	QdrantPort      int
	QdrantAPIKey    string
	DefaultModel    string
	EmbeddingModel  string
	MaxTokens       int
	Temperature     float64
	RequestTimeout  time.Duration
}

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

type SMSConfig struct {
	Provider  string
	APIKey    string
	APISecret string
	SenderID  string
}

type PushConfig struct {
	FirebaseKey string
}

type WhatsAppConfig struct {
	APIKey    string
	APISecret string
	PhoneID   string
	WABaseURL string
}

func Load() (*Config, error) {
	v := viper.New()

	v.SetConfigName(".env")
	v.SetConfigType("env")
	v.AddConfigPath(".")
	v.AddConfigPath("./..")
	v.AddConfigPath("../../")
	v.AutomaticEnv()

	setDefaults(v)

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	cfg := &Config{
		App: AppConfig{
			Name:  v.GetString("APP_NAME"),
			Env:   v.GetString("APP_ENV"),
			Port:  v.GetInt("APP_PORT"),
			Debug: v.GetBool("APP_DEBUG"),
		},
		DB: DatabaseConfig{
			Host:            v.GetString("DB_HOST"),
			Port:            v.GetInt("DB_PORT"),
			User:            v.GetString("DB_USER"),
			Password:        v.GetString("DB_PASSWORD"),
			Name:            v.GetString("DB_NAME"),
			SSLMode:         v.GetString("DB_SSLMODE"),
			MaxOpenConns:    v.GetInt("DB_MAX_OPEN_CONNS"),
			MaxIdleConns:    v.GetInt("DB_MAX_IDLE_CONNS"),
			ConnMaxLifetime: v.GetDuration("DB_CONN_MAX_LIFETIME"),
		},
		Redis: RedisConfig{
			Host:     v.GetString("REDIS_HOST"),
			Port:     v.GetInt("REDIS_PORT"),
			Password: v.GetString("REDIS_PASSWORD"),
			DB:       v.GetInt("REDIS_DB"),
			PoolSize: v.GetInt("REDIS_POOL_SIZE"),
		},
		MinIO: MinIOConfig{
			Endpoint:  v.GetString("MINIO_ENDPOINT"),
			AccessKey: v.GetString("MINIO_ACCESS_KEY"),
			SecretKey: v.GetString("MINIO_SECRET_KEY"),
			Bucket:    v.GetString("MINIO_BUCKET"),
			UseSSL:    v.GetBool("MINIO_USE_SSL"),
			Region:    v.GetString("MINIO_REGION"),
		},
		JWT: JWTConfig{
			AccessSecret:  v.GetString("JWT_ACCESS_SECRET"),
			RefreshSecret: v.GetString("JWT_REFRESH_SECRET"),
			AccessTTL:     v.GetDuration("JWT_ACCESS_TTL"),
			RefreshTTL:    v.GetDuration("JWT_REFRESH_TTL"),
			Issuer:        v.GetString("JWT_ISSUER"),
		},
		AI: AIConfig{
			Provider:        v.GetString("AI_PROVIDER"),
			OpenAIKey:       v.GetString("AI_OPENAI_KEY"),
			OpenAIEndpoint:  v.GetString("AI_OPENAI_ENDPOINT"),
			GeminiKey:       v.GetString("AI_GEMINI_KEY"),
			GeminiEndpoint:  v.GetString("AI_GEMINI_ENDPOINT"),
			ClaudeKey:       v.GetString("AI_CLAUDE_KEY"),
			ClaudeEndpoint:  v.GetString("AI_CLAUDE_ENDPOINT"),
			OllamaEndpoint:  v.GetString("AI_OLLAMA_ENDPOINT"),
			QdrantHost:      v.GetString("AI_QDRANT_HOST"),
			QdrantPort:      v.GetInt("AI_QDRANT_PORT"),
			QdrantAPIKey:    v.GetString("AI_QDRANT_API_KEY"),
			DefaultModel:    v.GetString("AI_DEFAULT_MODEL"),
			EmbeddingModel:  v.GetString("AI_EMBEDDING_MODEL"),
			MaxTokens:       v.GetInt("AI_MAX_TOKENS"),
			Temperature:     v.GetFloat64("AI_TEMPERATURE"),
			RequestTimeout:  v.GetDuration("AI_REQUEST_TIMEOUT"),
		},
		SMTP: SMTPConfig{
			Host:     v.GetString("SMTP_HOST"),
			Port:     v.GetInt("SMTP_PORT"),
			Username: v.GetString("SMTP_USERNAME"),
			Password: v.GetString("SMTP_PASSWORD"),
			From:     v.GetString("SMTP_FROM"),
		},
		SMS: SMSConfig{
			Provider:  v.GetString("SMS_PROVIDER"),
			APIKey:    v.GetString("SMS_API_KEY"),
			APISecret: v.GetString("SMS_API_SECRET"),
			SenderID:  v.GetString("SMS_SENDER_ID"),
		},
		Push: PushConfig{
			FirebaseKey: v.GetString("PUSH_FIREBASE_KEY"),
		},
		WhatsApp: WhatsAppConfig{
			APIKey:    v.GetString("WA_API_KEY"),
			APISecret: v.GetString("WA_API_SECRET"),
			PhoneID:   v.GetString("WA_PHONE_ID"),
			WABaseURL: v.GetString("WA_BASE_URL"),
		},
	}

	return cfg, nil
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("APP_NAME", "erp-school")
	v.SetDefault("APP_ENV", "development")
	v.SetDefault("APP_PORT", 8080)
	v.SetDefault("APP_DEBUG", true)

	v.SetDefault("DB_HOST", "localhost")
	v.SetDefault("DB_PORT", 5432)
	v.SetDefault("DB_USER", "postgres")
	v.SetDefault("DB_PASSWORD", "postgres")
	v.SetDefault("DB_NAME", "erp_school")
	v.SetDefault("DB_SSLMODE", "disable")
	v.SetDefault("DB_MAX_OPEN_CONNS", 25)
	v.SetDefault("DB_MAX_IDLE_CONNS", 10)
	v.SetDefault("DB_CONN_MAX_LIFETIME", "5m")

	v.SetDefault("REDIS_HOST", "localhost")
	v.SetDefault("REDIS_PORT", 6379)
	v.SetDefault("REDIS_PASSWORD", "")
	v.SetDefault("REDIS_DB", 0)
	v.SetDefault("REDIS_POOL_SIZE", 10)

	v.SetDefault("MINIO_ENDPOINT", "localhost:9000")
	v.SetDefault("MINIO_ACCESS_KEY", "minioadmin")
	v.SetDefault("MINIO_SECRET_KEY", "minioadmin")
	v.SetDefault("MINIO_BUCKET", "erp-school")
	v.SetDefault("MINIO_USE_SSL", false)
	v.SetDefault("MINIO_REGION", "us-east-1")

	v.SetDefault("JWT_ACCESS_SECRET", "change-me-access-secret-key")
	v.SetDefault("JWT_REFRESH_SECRET", "change-me-refresh-secret-key")
	v.SetDefault("JWT_ACCESS_TTL", "15m")
	v.SetDefault("JWT_REFRESH_TTL", "168h")
	v.SetDefault("JWT_ISSUER", "erp-school")

	v.SetDefault("AI_PROVIDER", "openai")
	v.SetDefault("AI_OPENAI_ENDPOINT", "https://api.openai.com/v1")
	v.SetDefault("AI_GEMINI_ENDPOINT", "https://generativelanguage.googleapis.com/v1beta")
	v.SetDefault("AI_CLAUDE_ENDPOINT", "https://api.anthropic.com/v1")
	v.SetDefault("AI_OLLAMA_ENDPOINT", "http://localhost:11434")
	v.SetDefault("AI_QDRANT_HOST", "localhost")
	v.SetDefault("AI_QDRANT_PORT", 6334)
	v.SetDefault("AI_DEFAULT_MODEL", "gpt-4o")
	v.SetDefault("AI_EMBEDDING_MODEL", "text-embedding-3-small")
	v.SetDefault("AI_MAX_TOKENS", 4096)
	v.SetDefault("AI_TEMPERATURE", 0.7)
	v.SetDefault("AI_REQUEST_TIMEOUT", "60s")

	v.SetDefault("SMTP_HOST", "smtp.gmail.com")
	v.SetDefault("SMTP_PORT", 587)
	v.SetDefault("SMTP_FROM", "noreply@school.id")

	v.SetDefault("SMS_PROVIDER", "twilio")

	v.SetDefault("WA_BASE_URL", "https://graph.facebook.com/v21.0")
}
