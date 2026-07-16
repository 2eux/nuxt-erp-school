package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type LoggingMiddleware struct {
	logger *zap.Logger
}

type requestLog struct {
	Timestamp    string        `json:"timestamp"`
	Method       string        `json:"method"`
	Path         string        `json:"path"`
	Status       int           `json:"status"`
	Latency      string        `json:"latency"`
	ClientIP     string        `json:"client_ip"`
	UserID       string        `json:"user_id,omitempty"`
	TenantID     string        `json:"tenant_id,omitempty"`
	RequestBody  any           `json:"request_body,omitempty"`
	ResponseSize int           `json:"response_size"`
	UserAgent    string        `json:"user_agent"`
	TokensUsed   int           `json:"tokens_used,omitempty"`
	Cost         float64       `json:"cost,omitempty"`
	Provider     string        `json:"provider,omitempty"`
	Model        string        `json:"model,omitempty"`
	Error        string        `json:"error,omitempty"`
}

func NewLoggingMiddleware(logger *zap.Logger) *LoggingMiddleware {
	return &LoggingMiddleware{logger: logger}
}

func (m *LoggingMiddleware) RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		var requestBody map[string]any
		if c.Request.Body != nil {
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err == nil {
				c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
				var bodyMap map[string]any
				if json.Valid(bodyBytes) {
					_ = json.Unmarshal(bodyBytes, &bodyMap)
					sanitizedBody := sanitizeBody(bodyMap)
					requestBody = sanitizedBody
				}
			}
		}

		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		c.Next()

		latency := time.Since(start)

		userID, _ := c.Get("user_id")
		tenantID, _ := c.Get("tenant_id")

		log := requestLog{
			Timestamp:    start.Format(time.RFC3339),
			Method:       c.Request.Method,
			Path:         c.Request.URL.Path,
			Status:       c.Writer.Status(),
			Latency:      latency.String(),
			ClientIP:     c.ClientIP(),
			UserID:       getString(userID),
			TenantID:     getString(tenantID),
			RequestBody:  requestBody,
			ResponseSize: blw.body.Len(),
			UserAgent:    c.Request.UserAgent(),
		}

		var responseMap map[string]any
		if blw.body.Len() > 0 {
			respBytes := blw.body.Bytes()
			if json.Valid(respBytes) {
				var respData map[string]any
				if err := json.Unmarshal(respBytes, &respData); err == nil {
					responseMap = respData
				}
			}
		}

		if responseMap != nil {
			if tok, ok := responseMap["tokens_used"]; ok {
				log.TokensUsed = int(getFloat(tok))
			}
			if usage, ok := responseMap["usage"]; ok {
				if usageMap, ok := usage.(map[string]any); ok {
					if tot, ok := usageMap["total_tokens"]; ok {
						log.TokensUsed = int(getFloat(tot))
					}
					if cst, ok := usageMap["cost"]; ok {
						log.Cost = getFloat(cst)
					}
				}
			}
			if provider, ok := responseMap["provider"]; ok {
				log.Provider = getString(provider)
			}
			if model, ok := responseMap["model"]; ok {
				log.Model = getString(model)
			}
			if errResp, ok := responseMap["error"]; ok {
				log.Error = getString(errResp)
			}
		}

		if c.Writer.Status() >= 400 {
			log.Error = string(blw.body.Bytes())
		}

		logData, _ := json.Marshal(log)
		if c.Writer.Status() >= 500 {
			m.logger.Error("request", zap.String("log", string(logData)))
		} else if c.Writer.Status() >= 400 {
			m.logger.Warn("request", zap.String("log", string(logData)))
		} else {
			m.logger.Debug("request", zap.String("log", string(logData)))
		}
	}
}

func sanitizeBody(body map[string]any) map[string]any {
	sanitized := make(map[string]any)
	for k, v := range body {
		if isSensitiveKey(k) {
			sanitized[k] = "[REDACTED]"
		} else if subMap, ok := v.(map[string]any); ok {
			sanitized[k] = sanitizeBody(subMap)
		} else {
			sanitized[k] = v
		}
	}
	return sanitized
}

func isSensitiveKey(key string) bool {
	sensitive := []string{"password", "secret", "api_key", "token", "authorization", "credit_card", "ssn"}
	for _, s := range sensitive {
		if key == s {
			return true
		}
	}
	return false
}

func getString(v interface{}) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func getFloat(v interface{}) float64 {
	if v == nil {
		return 0
	}
	switch val := v.(type) {
	case float64:
		return val
	case float32:
		return float64(val)
	case int:
		return float64(val)
	case int64:
		return float64(val)
	case json.Number:
		f, _ := val.Float64()
		return f
	}
	return 0
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
