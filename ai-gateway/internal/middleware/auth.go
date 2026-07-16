package middleware

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type AuthMiddleware struct {
	logger       *zap.Logger
	jwtSecret    string
	apiKey       string
	adminKey     string
	rateLimiter  *RateLimiter
}

type RateLimiter struct {
	mu       sync.Mutex
	buckets  map[string]*tokenBucket
	maxPerSec int
}

type tokenBucket struct {
	tokens    float64
	lastCheck time.Time
}

func NewRateLimiter(maxPerSec int) *RateLimiter {
	rl := &RateLimiter{
		buckets:   make(map[string]*tokenBucket),
		maxPerSec: maxPerSec,
	}
	go rl.cleanup()
	return rl
}

func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	bucket, ok := rl.buckets[key]
	if !ok {
		bucket = &tokenBucket{tokens: float64(rl.maxPerSec), lastCheck: time.Now()}
		rl.buckets[key] = bucket
	}

	now := time.Now()
	elapsed := now.Sub(bucket.lastCheck).Seconds()
	bucket.tokens += elapsed * float64(rl.maxPerSec)
	if bucket.tokens > float64(rl.maxPerSec) {
		bucket.tokens = float64(rl.maxPerSec)
	}
	bucket.lastCheck = now

	if bucket.tokens >= 1 {
		bucket.tokens--
		return true
	}
	return false
}

func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for key, bucket := range rl.buckets {
			if now.Sub(bucket.lastCheck) > 10*time.Minute {
				delete(rl.buckets, key)
			}
		}
		rl.mu.Unlock()
	}
}

type Claims struct {
	UserID   string `json:"user_id"`
	TenantID string `json:"tenant_id"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func NewAuthMiddleware(logger *zap.Logger, jwtSecret, apiKey, adminKey string, rateLimitRPS int) *AuthMiddleware {
	return &AuthMiddleware{
		logger:      logger,
		jwtSecret:   jwtSecret,
		apiKey:      apiKey,
		adminKey:    adminKey,
		rateLimiter: NewRateLimiter(rateLimitRPS),
	}
}

func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			apiKey := c.GetHeader("X-API-Key")
			if apiKey == "" {
				apiKey = c.Query("api_key")
			}
			if apiKey != "" {
				m.handleAPIKey(c, apiKey)
				return
			}
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization required"})
			c.Abort()
			return
		}

		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			m.handleJWT(c, tokenString)
			return
		}

		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization format"})
		c.Abort()
	}
}

func (m *AuthMiddleware) handleAPIKey(c *gin.Context, key string) {
	if m.adminKey != "" && key == m.adminKey {
		c.Set("auth_type", "admin_key")
		c.Set("role", "admin")
		c.Next()
		return
	}

	if m.apiKey != "" && key == m.apiKey {
		c.Set("auth_type", "api_key")
		c.Set("role", "service")
		c.Next()
		return
	}

	c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid API key"})
	c.Abort()
}

func (m *AuthMiddleware) handleJWT(c *gin.Context, tokenString string) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(m.jwtSecret), nil
	})

	if err != nil || !token.Valid {
		m.logger.Warn("invalid JWT token", zap.Error(err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
		c.Abort()
		return
	}

	c.Set("auth_type", "jwt")
	c.Set("user_id", claims.UserID)
	c.Set("tenant_id", claims.TenantID)
	c.Set("role", claims.Role)
	c.Next()
}

func (m *AuthMiddleware) RateLimiterMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.ClientIP()

		if userID, exists := c.Get("user_id"); exists {
			key = "user:" + userID.(string)
		}
		if tenantID, exists := c.Get("tenant_id"); exists {
			key = "tenant:" + tenantID.(string)
		}

		if !m.rateLimiter.Allow(key) {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded. Please try again later."})
			c.Abort()
			return
		}
		c.Next()
	}
}

func (m *AuthMiddleware) RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "role not found"})
			c.Abort()
			return
		}

		if userRole.(string) != role && userRole.(string) != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			apiKey := c.GetHeader("X-API-Key")
			if apiKey == "" {
				apiKey = c.Query("api_key")
			}
			if apiKey == "" {
				c.Next()
				return
			}
			m.handleAPIKey(c, apiKey)
			return
		}

		if strings.HasPrefix(authHeader, "Bearer ") {
			m.handleJWT(c, strings.TrimPrefix(authHeader, "Bearer "))
			return
		}
		c.Next()
	}
}
