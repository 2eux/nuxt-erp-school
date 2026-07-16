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
	"github.com/opencode/erp-school-backend/internal/config"
	"github.com/opencode/erp-school-backend/internal/handler"
	"github.com/opencode/erp-school-backend/internal/infrastructure/cache"
	"github.com/opencode/erp-school-backend/internal/infrastructure/database"
	"github.com/opencode/erp-school-backend/internal/infrastructure/storage"
	"github.com/opencode/erp-school-backend/internal/middleware"
	"github.com/opencode/erp-school-backend/internal/repository"
	"github.com/opencode/erp-school-backend/internal/service"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("failed to load config", zap.Error(err))
	}

	db, err := database.NewPostgresDB(cfg.DB, logger)
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	if err := database.RunMigrations(db, logger); err != nil {
		logger.Fatal("failed to run migrations", zap.Error(err))
	}

	redisClient, err := cache.NewRedisClient(cfg.Redis, logger)
	if err != nil {
		logger.Warn("failed to connect to redis", zap.Error(err))
	}
	if redisClient != nil {
		defer redisClient.Close()
	}

	minioClient, err := storage.NewMinioClient(cfg.MinIO, logger)
	if err != nil {
		logger.Warn("failed to connect to minio", zap.Error(err))
	}
	_ = minioClient

	// Repositories
	authRepo := repository.NewAuthRepository(db)
	schoolRepo := repository.NewSchoolRepository(db)
	studentRepo := repository.NewStudentRepository(db)

	// Services
	authService := service.NewAuthService(authRepo, schoolRepo, redisClient, cfg.JWT, logger)
	schoolService := service.NewSchoolService(schoolRepo, logger)
	userService := service.NewUserService(authRepo, schoolRepo, logger)
	studentService := service.NewStudentService(studentRepo, logger)
	teacherService := service.NewTeacherService(db, logger)
	employeeService := service.NewEmployeeService(db, logger)
	academicService := service.NewAcademicService(db, logger)
	islamicService := service.NewIslamicService(db, logger)
	financeService := service.NewFinanceService(db, logger)
	notificationService := service.NewNotificationService(db, logger)
	documentService := service.NewDocumentService(db, logger)
	admissionService := service.NewAdmissionService(db, logger)
	analyticsService := service.NewAnalyticsService(db, logger)
	aiService := service.NewAIService(db, logger)

	// Middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg.JWT, logger)

	// Handlers
	authHandler := handler.NewAuthHandler(authService)
	schoolHandler := handler.NewSchoolHandler(schoolService)
	userHandler := handler.NewUserHandler(userService)
	studentHandler := handler.NewStudentHandler(studentService)
	teacherHandler := handler.NewTeacherHandler(teacherService)
	employeeHandler := handler.NewEmployeeHandler(employeeService)
	academicHandler := handler.NewAcademicHandler(academicService)
	islamicHandler := handler.NewIslamicHandler(islamicService)
	financeHandler := handler.NewFinanceHandler(financeService)
	notificationHandler := handler.NewNotificationHandler(notificationService)
	documentHandler := handler.NewDocumentHandler(documentService)
	admissionHandler := handler.NewAdmissionHandler(admissionService)
	analyticsHandler := handler.NewAnalyticsHandler(analyticsService)
	aiHandler := handler.NewAIHandler(aiService)

	// Gin Engine
	engine := gin.Default()

	engine.Use(middleware.NewCORSMiddleware())
	engine.Use(middleware.NewLoggerMiddleware(logger))

	// Health check
	engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	})

	// Routes
	router := handler.NewRouter(
		authHandler, schoolHandler, userHandler, studentHandler,
		teacherHandler, employeeHandler, academicHandler, islamicHandler,
		financeHandler, notificationHandler, documentHandler, admissionHandler,
		analyticsHandler, aiHandler, authMiddleware,
	)
	router.Setup(engine)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.App.Port),
		Handler: engine,
	}

	go func() {
		logger.Info("server starting", zap.Int("port", cfg.App.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("server failed", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("server forced to shutdown", zap.Error(err))
	}

	logger.Info("server exited gracefully")
}
