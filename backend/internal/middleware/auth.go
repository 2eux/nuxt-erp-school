package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/opencode/erp-school-backend/internal/config"
	"github.com/opencode/erp-school-backend/internal/domain"
	"go.uber.org/zap"
)

type Claims struct {
	UserID   string   `json:"user_id"`
	SchoolID string   `json:"school_id"`
	Email    string   `json:"email"`
	Roles    []string `json:"roles"`
	jwt.RegisteredClaims
}

type AuthMiddleware struct {
	cfg    config.JWTConfig
	logger *zap.Logger
}

func NewAuthMiddleware(cfg config.JWTConfig, logger *zap.Logger) *AuthMiddleware {
	return &AuthMiddleware{cfg: cfg, logger: logger}
}

func (a *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"code":    http.StatusUnauthorized,
				"message": "missing authorization header",
				"error":   domain.ErrUnauthorized.Error(),
			})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"code":    http.StatusUnauthorized,
				"message": "invalid authorization format",
				"error":   domain.ErrUnauthorized.Error(),
			})
			c.Abort()
			return
		}

		claims, err := a.validateToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"code":    http.StatusUnauthorized,
				"message": "invalid or expired token",
				"error":   err.Error(),
			})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("school_id", claims.SchoolID)
		c.Set("email", claims.Email)
		c.Set("roles", claims.Roles)
		c.Next()
	}
}

func (a *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
			c.Next()
			return
		}

		claims, err := a.validateToken(parts[1])
		if err != nil {
			c.Next()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("school_id", claims.SchoolID)
		c.Set("email", claims.Email)
		c.Set("roles", claims.Roles)
		c.Next()
	}
}

func (a *AuthMiddleware) RequirePermission(perm string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roles, exists := c.Get("roles")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"code":    http.StatusForbidden,
				"message": "not authorized",
				"error":   domain.ErrForbidden.Error(),
			})
			c.Abort()
			return
		}

		roleList, ok := roles.([]string)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"code":    http.StatusForbidden,
				"message": "invalid role data",
				"error":   domain.ErrForbidden.Error(),
			})
			c.Abort()
			return
		}

		for _, role := range roleList {
			if role == "superadmin" || role == "admin" {
				c.Next()
				return
			}
			if role == perm {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"code":    http.StatusForbidden,
			"message": "insufficient permissions",
			"error":   domain.ErrForbidden.Error(),
		})
		c.Abort()
	}
}

func (a *AuthMiddleware) RequireRoles(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRoles, exists := c.Get("roles")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{})
			c.Abort()
			return
		}

		roleList, ok := userRoles.([]string)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{})
			c.Abort()
			return
		}

		roleSet := make(map[string]bool)
		for _, r := range roleList {
			roleSet[r] = true
		}

		for _, required := range roles {
			if roleSet[required] || roleSet["superadmin"] {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"code":    http.StatusForbidden,
			"message": "insufficient role permissions",
		})
		c.Abort()
	}
}

func (a *AuthMiddleware) SchoolContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		schoolID, exists := c.Get("school_id")
		if exists && schoolID != "" {
			c.Next()
			return
		}

		schoolID = c.GetHeader("X-School-ID")
		if schoolID != "" {
			c.Set("school_id", schoolID)
		}
		c.Next()
	}
}

func (a *AuthMiddleware) validateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, domain.ErrTokenInvalid
		}
		return []byte(a.cfg.AccessSecret), nil
	})
	if err != nil {
		return nil, domain.ErrTokenExpired
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, domain.ErrTokenInvalid
	}

	return claims, nil
}

func (a *AuthMiddleware) GenerateAccessToken(userID, schoolID, email string, roles []string) (string, time.Time, error) {
	expiresAt := time.Now().Add(a.cfg.AccessTTL)
	claims := &Claims{
		UserID:   userID,
		SchoolID: schoolID,
		Email:    email,
		Roles:    roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    a.cfg.Issuer,
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(a.cfg.AccessSecret))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}

func (a *AuthMiddleware) GenerateRefreshToken(userID string) (string, time.Time, error) {
	expiresAt := time.Now().Add(a.cfg.RefreshTTL)
	claims := jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    a.cfg.Issuer,
		ID:        userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(a.cfg.RefreshSecret))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}

func (a *AuthMiddleware) ValidateRefreshToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, domain.ErrTokenInvalid
		}
		return []byte(a.cfg.RefreshSecret), nil
	})
	if err != nil {
		return "", domain.ErrTokenExpired
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", domain.ErrTokenInvalid
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		return "", domain.ErrTokenInvalid
	}

	return sub, nil
}

func GetUserID(c *gin.Context) string {
	id, _ := c.Get("user_id")
	if id == nil {
		return ""
	}
	return id.(string)
}

func GetSchoolID(c *gin.Context) string {
	id, _ := c.Get("school_id")
	if id == nil {
		return ""
	}
	return id.(string)
}

func GetEmail(c *gin.Context) string {
	email, _ := c.Get("email")
	if email == nil {
		return ""
	}
	return email.(string)
}

func GetRoles(c *gin.Context) []string {
	roles, _ := c.Get("roles")
	if roles == nil {
		return nil
	}
	return roles.([]string)
}
