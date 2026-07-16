package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opencode/erp-ai-gateway/internal/cache"
	"github.com/opencode/erp-ai-gateway/internal/config"
	"github.com/opencode/erp-ai-gateway/internal/handlers"
	"github.com/opencode/erp-ai-gateway/internal/mcp"
	"github.com/opencode/erp-ai-gateway/internal/middleware"
	"github.com/opencode/erp-ai-gateway/internal/providers"
	"github.com/opencode/erp-ai-gateway/internal/rag"
	"github.com/opencode/erp-ai-gateway/internal/router"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	logger := initLogger(cfg)
	defer func() { _ = logger.Sync() }()

	redisCache, err := cache.NewCache(cfg.Redis, logger)
	if err != nil {
		logger.Warn("redis initialization warning", zap.Error(err))
	}

	providerMap := initProviders(cfg, logger)

	embeddingService := rag.NewEmbeddingService(logger, redisCache)

	qdrantClient := rag.NewQdrantClient(cfg.Qdrant, logger)

	ragService := rag.NewRAGService(qdrantClient, embeddingService, cfg.RAG, logger)

	providerRouter := router.NewProviderRouter(cfg, providerMap)

	mcpServer := mcp.NewMCPServer(logger, providerRouter, redisCache)

	authMiddleware := middleware.NewAuthMiddleware(
		logger,
		cfg.JWT.Secret,
		cfg.APIKey,
		cfg.AdminKey,
		cfg.RateLimit.PerUser,
	)

	loggingMiddleware := middleware.NewLoggingMiddleware(logger)

	chatHandler := handlers.NewChatHandler(logger, providerRouter)
	mcpHandler := handlers.NewMCPHandler(logger, mcpServer)
	ragHandler := handlers.NewRAGHandler(logger, providerRouter, ragService)

	engine := setupRouter(
		cfg,
		loggingMiddleware,
		authMiddleware,
		chatHandler,
		mcpHandler,
		ragHandler,
		providerRouter,
		mcpServer,
		qdrantClient,
		ragService,
		logger,
	)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      engine,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 300 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Info("AI Gateway starting",
			zap.String("env", cfg.Env),
			zap.String("port", cfg.Port),
		)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("server failed to start", zap.Error(err))
		}
	}()

	<-quit
	logger.Info("shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Fatal("server forced to shutdown", zap.Error(err))
	}

	if redisCache != nil {
		if err := redisCache.Close(); err != nil {
			logger.Warn("redis close error", zap.Error(err))
		}
	}

	qdrantClient.Close()

	logger.Info("server exited gracefully")
}

func initLogger(cfg *config.Config) *zap.Logger {
	var level zapcore.Level
	switch cfg.LogLevel {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	var encoder zapcore.Encoder
	if cfg.LogFormat == "console" {
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	}

	core := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level)
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return logger
}

func initProviders(cfg *config.Config, logger *zap.Logger) map[string]providers.Provider {
	providerMap := make(map[string]providers.Provider)

	if cfg.IsProviderEnabled("openai") {
		pc, err := cfg.GetProvider("openai")
		if err == nil && pc.APIKey != "" {
			providerMap["openai"] = providers.NewOpenAIProvider(pc)
			logger.Info("openai provider initialized")
		} else {
			logger.Warn("openai provider skipped (no API key)")
		}
	}

	if cfg.IsProviderEnabled("gemini") {
		pc, err := cfg.GetProvider("gemini")
		if err == nil && pc.APIKey != "" {
			gemini, initErr := providers.NewGeminiProvider(pc)
			if initErr != nil {
				logger.Warn("gemini provider init failed", zap.Error(initErr))
			} else {
				providerMap["gemini"] = gemini
				logger.Info("gemini provider initialized")
			}
		} else {
			logger.Warn("gemini provider skipped (no API key)")
		}
	}

	if cfg.IsProviderEnabled("claude") {
		pc, err := cfg.GetProvider("claude")
		if err == nil && pc.APIKey != "" {
			providerMap["claude"] = providers.NewClaudeProvider(pc)
			logger.Info("claude provider initialized")
		} else {
			logger.Warn("claude provider skipped (no API key)")
		}
	}

	if cfg.IsProviderEnabled("ollama") {
		pc, err := cfg.GetProvider("ollama")
		if err == nil {
			ollama := providers.NewOllamaProvider(pc)
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			if ollama.IsAvailable(ctx) {
				providerMap["ollama"] = ollama
				logger.Info("ollama provider initialized")
			} else {
				logger.Warn("ollama provider skipped (not reachable)")
			}
			cancel()
		}
	}

	return providerMap
}

func setupRouter(
	cfg *config.Config,
	logMw *middleware.LoggingMiddleware,
	authMw *middleware.AuthMiddleware,
	chatHandler *handlers.ChatHandler,
	mcpHandler *handlers.MCPHandler,
	ragHandler *handlers.RAGHandler,
	providerRouter *router.ProviderRouter,
	mcpServer *mcp.MCPServer,
	qdrantClient *rag.QdrantClient,
	ragService *rag.RAGService,
	logger *zap.Logger,
) *gin.Engine {
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()

	engine.Use(gin.Recovery())
	engine.Use(logMw.RequestLogger())

	engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"service":   "erp-ai-gateway",
			"version":   "1.0.0",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	})

	engine.GET("/ready", func(c *gin.Context) {
		readinessCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		checks := make(map[string]string)

		providers := providerRouter.GetEnabledProviders()
		for _, p := range providers {
			if p.IsAvailable(readinessCtx) {
				checks[p.Name()] = "healthy"
			} else {
				checks[p.Name()] = "unavailable"
			}
		}

		if qdrantClient != nil {
			if err := qdrantClient.HealthCheck(readinessCtx); err != nil {
				checks["qdrant"] = "unhealthy"
			} else {
				checks["qdrant"] = "healthy"
			}
		}

		allHealthy := true
		hasProviders := false
		for _, status := range checks {
			if status == "unhealthy" || status == "unavailable" {
				allHealthy = false
			}
			hasProviders = true
		}

		status := http.StatusOK
		readinessStatus := "ready"
		if !hasProviders {
			readinessStatus = "degraded"
			status = http.StatusServiceUnavailable
		} else if !allHealthy {
			readinessStatus = "degraded"
		}

		c.JSON(status, gin.H{
			"status":  readinessStatus,
			"checks":  checks,
		})
	})

	v1 := engine.Group("/api/v1")
	v1.Use(authMw.Authenticate())
	v1.Use(authMw.RateLimiterMiddleware())
	{
		v1.POST("/chat/completions", chatHandler.ChatCompletions)
		v1.POST("/chat/stream", chatHandler.ChatStream)
		v1.POST("/embeddings", chatHandler.Embeddings)

		mcpGroup := v1.Group("/mcp")
		{
			mcpGroup.POST("/messages", mcpHandler.HandleJSONRPC)
			mcpGroup.GET("/sse", mcpHandler.HandleSSE)
			mcpGroup.GET("/tools", mcpHandler.ListTools)
			mcpGroup.POST("/tools/call", mcpHandler.CallTool)
			mcpGroup.GET("/resources", mcpHandler.ListResources)
			mcpGroup.GET("/resources/read", mcpHandler.ReadResource)
			mcpGroup.GET("/prompts", mcpHandler.ListPrompts)
			mcpGroup.POST("/prompts/get", mcpHandler.GetPrompt)
		}

		ragGroup := v1.Group("/rag")
		{
			ragGroup.POST("/documents", ragHandler.UploadDocument)
			ragGroup.POST("/documents/upload", ragHandler.UploadFile)
			ragGroup.POST("/query", ragHandler.Query)
			ragGroup.GET("/documents", ragHandler.ListDocuments)
			ragGroup.DELETE("/documents/:id", ragHandler.DeleteDocument)
		}
	}

	adminGroup := engine.Group("/api/v1/admin")
	adminGroup.Use(authMw.Authenticate())
	adminGroup.Use(authMw.RequireRole("admin"))
	{
		adminGroup.GET("/stats", func(c *gin.Context) {
			stats := providerRouter.GetStats()
			c.JSON(http.StatusOK, gin.H{
				"provider_stats": stats,
			})
		})

		adminGroup.GET("/providers", func(c *gin.Context) {
			providers := providerRouter.GetEnabledProviders()
			result := make([]gin.H, 0, len(providers))
			ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
			defer cancel()
			for _, p := range providers {
				result = append(result, gin.H{
					"name":      p.Name(),
					"models":    p.Models(),
					"available": p.IsAvailable(ctx),
				})
			}
			c.JSON(http.StatusOK, gin.H{"providers": result})
		})

		adminGroup.GET("/rag/stats", func(c *gin.Context) {
			stats, err := ragService.GetStats(c.Request.Context())
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, stats)
		})

		adminGroup.DELETE("/cache", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "cache cleared"})
		})
	}

	return engine
}


