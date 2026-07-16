package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opencode/erp-school-backend/internal/domain"
	"github.com/opencode/erp-school-backend/internal/infrastructure/cache"
	"go.uber.org/zap"
)

type RateLimiter struct {
	redis  *cache.RedisClient
	logger *zap.Logger
}

func NewRateLimiter(redis *cache.RedisClient, logger *zap.Logger) *RateLimiter {
	return &RateLimiter{redis: redis, logger: logger}
}

func (rl *RateLimiter) Limit(maxRequests int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := "ratelimit:" + c.ClientIP()

		count, err := rl.redis.Incr(c.Request.Context(), key)
		if err != nil {
			rl.logger.Error("rate limit check failed", zap.Error(err))
			c.Next()
			return
		}

		if count == 1 {
			rl.redis.Expire(c.Request.Context(), key, window)
		}

		c.Header("X-RateLimit-Limit", strconv.Itoa(maxRequests))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(maxRequests-int(count)))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(window).Unix(), 10))

		if int(count) > maxRequests {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"code":    http.StatusTooManyRequests,
				"message": domain.ErrRateLimited.Error(),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
