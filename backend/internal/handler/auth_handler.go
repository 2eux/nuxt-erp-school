package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opencode/erp-school-backend/internal/dto"
	"github.com/opencode/erp-school-backend/internal/middleware"
	"github.com/opencode/erp-school-backend/internal/service"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Login godoc
// @Summary Login user
// @Description Authenticate user with email and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login credentials"
// @Success 200 {object} dto.APIResponse{data=dto.LoginResponse}
// @Failure 401 {object} dto.APIResponse
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "invalid request", err.Error()))
		return
	}

	userAgent := c.GetHeader("User-Agent")
	ipAddress := c.ClientIP()

	resp, err := h.authService.Login(c.Request.Context(), req, userAgent, ipAddress)
	if err != nil {
		domainErr := parseError(err)
		c.JSON(domainErr.Code, dto.NewErrorResponse(domainErr.Code, domainErr.Message, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "login successful", resp))
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Get new access token using refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} dto.APIResponse{data=dto.TokenResponse}
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "invalid request", err.Error()))
		return
	}

	userAgent := c.GetHeader("User-Agent")
	ipAddress := c.ClientIP()

	resp, err := h.authService.RefreshToken(c.Request.Context(), req, userAgent, ipAddress)
	if err != nil {
		domainErr := parseError(err)
		c.JSON(domainErr.Code, dto.NewErrorResponse(domainErr.Code, domainErr.Message, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "token refreshed", resp))
}

// Logout godoc
// @Summary Logout user
// @Description Revoke refresh token and end session
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "invalid request", err.Error()))
		return
	}

	if err := h.authService.Logout(c.Request.Context(), req.RefreshToken); err != nil {
		domainErr := parseError(err)
		c.JSON(domainErr.Code, dto.NewErrorResponse(domainErr.Code, domainErr.Message, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "logout successful", nil))
}

// ForgotPassword godoc
// @Summary Forgot password
// @Description Send password reset token to email
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.ForgotPasswordRequest true "Email address"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/auth/forgot-password [post]
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req dto.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "invalid request", err.Error()))
		return
	}

	if err := h.authService.ForgotPassword(c.Request.Context(), req); err != nil {
		domainErr := parseError(err)
		c.JSON(domainErr.Code, dto.NewErrorResponse(domainErr.Code, domainErr.Message, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "if email exists, reset instructions sent", nil))
}

// ResetPassword godoc
// @Summary Reset password
// @Description Reset password using token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.ResetPasswordRequest true "Reset password data"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/auth/reset-password [post]
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req dto.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "invalid request", err.Error()))
		return
	}

	if err := h.authService.ResetPassword(c.Request.Context(), req); err != nil {
		domainErr := parseError(err)
		c.JSON(domainErr.Code, dto.NewErrorResponse(domainErr.Code, domainErr.Message, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "password reset successful", nil))
}

// GetMe godoc
// @Summary Get current user profile
// @Description Get authenticated user's profile and permissions
// @Tags Auth
// @Produce json
// @Security Bearer
// @Success 200 {object} dto.APIResponse{data=dto.UserDetail}
// @Router /api/v1/auth/me [get]
func (h *AuthHandler) GetMe(c *gin.Context) {
	userID := middleware.GetUserID(c)
	resp, err := h.authService.GetMe(c.Request.Context(), userID)
	if err != nil {
		domainErr := parseError(err)
		c.JSON(domainErr.Code, dto.NewErrorResponse(domainErr.Code, domainErr.Message, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "user profile", resp))
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update authenticated user's profile
// @Tags Auth
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body dto.UpdateUserRequest true "Profile data"
// @Success 200 {object} dto.APIResponse{data=dto.UserBrief}
// @Router /api/v1/auth/profile [put]
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "invalid request", err.Error()))
		return
	}

	userID := middleware.GetUserID(c)
	resp, err := h.authService.UpdateProfile(c.Request.Context(), userID, req)
	if err != nil {
		domainErr := parseError(err)
		c.JSON(domainErr.Code, dto.NewErrorResponse(domainErr.Code, domainErr.Message, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "profile updated", resp))
}
