package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opencode/erp-school-backend/internal/dto"
	"github.com/opencode/erp-school-backend/internal/middleware"
	"github.com/opencode/erp-school-backend/internal/service"
)

type AnalyticsHandler struct {
	analyticsService service.AnalyticsService
}

func NewAnalyticsHandler(analyticsService service.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{analyticsService: analyticsService}
}

// GetDashboard godoc
// @Summary Get dashboard stats
// @Tags Analytics
// @Produce json
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/analytics/dashboard [get]
func (h *AnalyticsHandler) GetDashboard(c *gin.Context) {
	schoolID := middleware.GetSchoolID(c)
	resp, err := h.analyticsService.GetDashboard(c.Request.Context(), schoolID)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "success", resp))
}

// GetAcademicAnalytics godoc
// @Summary Get academic analytics
// @Tags Analytics
// @Produce json
// @Param class_id query string false "Class ID"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/analytics/academic [get]
func (h *AnalyticsHandler) GetAcademicAnalytics(c *gin.Context) {
	schoolID := middleware.GetSchoolID(c)
	classID := c.Query("class_id")
	resp, err := h.analyticsService.GetAcademicAnalytics(c.Request.Context(), schoolID, classID)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "success", resp))
}

// GetFinanceAnalytics godoc
// @Summary Get finance analytics
// @Tags Analytics
// @Produce json
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/analytics/finance [get]
func (h *AnalyticsHandler) GetFinanceAnalytics(c *gin.Context) {
	schoolID := middleware.GetSchoolID(c)
	resp, err := h.analyticsService.GetFinanceAnalytics(c.Request.Context(), schoolID)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "success", resp))
}

// GetTahfidzAnalytics godoc
// @Summary Get tahfidz analytics
// @Tags Analytics
// @Produce json
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/analytics/tahfidz [get]
func (h *AnalyticsHandler) GetTahfidzAnalytics(c *gin.Context) {
	schoolID := middleware.GetSchoolID(c)
	resp, err := h.analyticsService.GetTahfidzAnalytics(c.Request.Context(), schoolID)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "success", resp))
}

// GetAdmissionAnalytics godoc
// @Summary Get admission analytics
// @Tags Analytics
// @Produce json
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/analytics/admissions [get]
func (h *AnalyticsHandler) GetAdmissionAnalytics(c *gin.Context) {
	schoolID := middleware.GetSchoolID(c)
	resp, err := h.analyticsService.GetAdmissionAnalytics(c.Request.Context(), schoolID)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "success", resp))
}

// GetAttendanceAnalytics godoc
// @Summary Get attendance analytics
// @Tags Analytics
// @Produce json
// @Param class_id query string false "Class ID"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/analytics/attendance [get]
func (h *AnalyticsHandler) GetAttendanceAnalytics(c *gin.Context) {
	schoolID := middleware.GetSchoolID(c)
	classID := c.Query("class_id")
	resp, err := h.analyticsService.GetAttendanceAnalytics(c.Request.Context(), schoolID, classID)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "success", resp))
}
