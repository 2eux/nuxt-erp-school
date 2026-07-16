package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type AuditMiddleware struct {
	db     *sqlx.DB
	logger *zap.Logger
}

func NewAuditMiddleware(db *sqlx.DB, logger *zap.Logger) *AuditMiddleware {
	return &AuditMiddleware{db: db, logger: logger}
}

func (a *AuditMiddleware) Audit() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "GET" || c.Request.Method == "OPTIONS" || c.Request.Method == "HEAD" {
			c.Next()
			return
		}

		start := time.Now()

		var body []byte
		if c.Request.Body != nil {
			body, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		}

		c.Next()

		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/api/v1/auth/") ||
			strings.HasPrefix(path, "/api/v1/ai/") ||
			strings.HasPrefix(path, "/swagger") {
			return
		}

		userID := GetUserID(c)
		schoolID := GetSchoolID(c)
		if userID == "" {
			return
		}

		entity, action := parsePathForAudit(path, c.Request.Method)
		if entity == "" {
			return
		}

		latency := time.Since(start)

		newValues := "{}"
		if len(body) > 0 {
			var compactBody bytes.Buffer
			if err := json.Compact(&compactBody, body); err == nil {
				newValues = compactBody.String()
			}
			if len(newValues) > 2000 {
				newValues = newValues[:2000]
			}
		}

		query := `INSERT INTO audit_logs (id, school_id, user_id, action, entity, entity_id, old_values, new_values, ip_address, user_agent)
				  VALUES ($1, $2, $3, $4, $5, NULL, NULL, $6, $7, $8)`

		_, err := a.db.Exec(query,
			uuid.New().String(),
			schoolID,
			userID,
			action,
			entity,
			newValues,
			c.ClientIP(),
			c.Request.UserAgent(),
		)

		if err != nil {
			a.logger.Error("failed to write audit log", zap.Error(err))
		}

		a.logger.Debug("audit",
			zap.String("user_id", userID),
			zap.String("entity", entity),
			zap.String("action", action),
			zap.Duration("latency", latency),
		)
	}
}

func parsePathForAudit(path, method string) (string, string) {
	path = strings.TrimPrefix(path, "/api/v1/")
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) == 0 {
		return "", ""
	}

	entity := parts[0]

	action := methodToAction(method)
	if len(parts) > 1 {
		if _, err := uuid.Parse(parts[1]); err == nil {
			if method == "PUT" || method == "PATCH" {
				action = "update"
			} else if method == "DELETE" {
				action = "delete"
			}
			return entity, action
		}
		entity = parts[0] + "." + parts[1]
	}

	return entity, action
}

func methodToAction(method string) string {
	switch method {
	case "POST":
		return "create"
	case "PUT", "PATCH":
		return "update"
	case "DELETE":
		return "delete"
	default:
		return "access"
	}
}
