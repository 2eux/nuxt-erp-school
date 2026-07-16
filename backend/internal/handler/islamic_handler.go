package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opencode/erp-school-backend/internal/dto"
	"github.com/opencode/erp-school-backend/internal/middleware"
	"github.com/opencode/erp-school-backend/internal/service"
)

type IslamicHandler struct {
	islamicService service.IslamicService
}

func NewIslamicHandler(islamicService service.IslamicService) *IslamicHandler {
	return &IslamicHandler{islamicService: islamicService}
}

// ListTahfidzProgress godoc
// @Summary List tahfidz progress
// @Tags Islamic
// @Produce json
// @Param student_id query string true "Student ID"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/tahfidz/progress [get]
func (h *IslamicHandler) ListTahfidzProgress(c *gin.Context) {
	studentID := c.Query("student_id")
	items, err := h.islamicService.ListTahfidzProgress(c.Request.Context(), studentID)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "success", items))
}

// CreateTahfidzProgress godoc
// @Summary Create tahfidz progress
// @Tags Islamic
// @Accept json
// @Produce json
// @Param request body dto.CreateTahfidzProgressRequest true "Progress data"
// @Success 201 {object} dto.APIResponse
// @Router /api/v1/tahfidz/progress [post]
func (h *IslamicHandler) CreateTahfidzProgress(c *gin.Context) {
	teacherID := middleware.GetUserID(c)
	var req dto.CreateTahfidzProgressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "invalid request", err.Error()))
		return
	}
	resp, err := h.islamicService.CreateTahfidzProgress(c.Request.Context(), teacherID, req)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusCreated, dto.NewAPIResponse(http.StatusCreated, "progress recorded", resp))
}

// ListMutabaah godoc
// @Summary List mutabaah
// @Tags Islamic
// @Produce json
// @Param student_id query string true "Student ID"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/mutabaah [get]
func (h *IslamicHandler) ListMutabaah(c *gin.Context) {
	studentID := c.Query("student_id")
	startDate := time.Now().AddDate(0, -1, 0)
	endDate := time.Now()
	items, err := h.islamicService.ListMutabaah(c.Request.Context(), studentID, startDate, endDate)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "success", items))
}

// CreateMutabaah godoc
// @Summary Create mutabaah
// @Tags Islamic
// @Accept json
// @Produce json
// @Param request body dto.CreateMutabaahRequest true "Mutabaah data"
// @Success 201 {object} dto.APIResponse
// @Router /api/v1/mutabaah [post]
func (h *IslamicHandler) CreateMutabaah(c *gin.Context) {
	var req dto.CreateMutabaahRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "invalid request", err.Error()))
		return
	}
	resp, err := h.islamicService.CreateMutabaah(c.Request.Context(), req)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusCreated, dto.NewAPIResponse(http.StatusCreated, "mutabaah created", resp))
}

// ListPrayerAttendance godoc
// @Summary List prayer attendance
// @Tags Islamic
// @Produce json
// @Param student_id query string true "Student ID"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/prayer-attendance [get]
func (h *IslamicHandler) ListPrayerAttendance(c *gin.Context) {
	studentID := c.Query("student_id")
	startDate := time.Now().AddDate(0, -1, 0)
	endDate := time.Now()
	items, err := h.islamicService.ListPrayerAttendance(c.Request.Context(), studentID, startDate, endDate)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "success", items))
}

// CreatePrayerAttendance godoc
// @Summary Create prayer attendance
// @Tags Islamic
// @Accept json
// @Produce json
// @Param request body dto.CreatePrayerAttendanceRequest true "Attendance data"
// @Success 201 {object} dto.APIResponse
// @Router /api/v1/prayer-attendance [post]
func (h *IslamicHandler) CreatePrayerAttendance(c *gin.Context) {
	var req dto.CreatePrayerAttendanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "invalid request", err.Error()))
		return
	}
	resp, err := h.islamicService.CreatePrayerAttendance(c.Request.Context(), req)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusCreated, dto.NewAPIResponse(http.StatusCreated, "prayer attendance created", resp))
}

// ListHalaqahGroups godoc
// @Summary List halaqah groups
// @Tags Islamic
// @Produce json
// @Param page query int false "Page"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/halaqah/groups [get]
func (h *IslamicHandler) ListHalaqahGroups(c *gin.Context) {
	schoolID := middleware.GetSchoolID(c)
	filter := getPagination(c)
	items, total, err := h.islamicService.ListHalaqahGroups(c.Request.Context(), schoolID, filter)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "success", dto.NewPaginatedResponse(items, total, filter.Page, filter.PageSize)))
}

// CreateHalaqahGroup godoc
// @Summary Create halaqah group
// @Tags Islamic
// @Accept json
// @Produce json
// @Param request body dto.CreateHalaqahGroupRequest true "Halaqah data"
// @Success 201 {object} dto.APIResponse
// @Router /api/v1/halaqah/groups [post]
func (h *IslamicHandler) CreateHalaqahGroup(c *gin.Context) {
	schoolID := middleware.GetSchoolID(c)
	var req dto.CreateHalaqahGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "invalid request", err.Error()))
		return
	}
	resp, err := h.islamicService.CreateHalaqahGroup(c.Request.Context(), schoolID, req)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusCreated, dto.NewAPIResponse(http.StatusCreated, "halaqah group created", resp))
}

// GetHalaqahGroup godoc
// @Summary Get halaqah group
// @Tags Islamic
// @Produce json
// @Param id path string true "Halaqah Group ID"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/halaqah/groups/{id} [get]
func (h *IslamicHandler) GetHalaqahGroup(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.islamicService.GetHalaqahGroup(c.Request.Context(), id)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "success", resp))
}

// AddHalaqahMember godoc
// @Summary Add halaqah member
// @Tags Islamic
// @Accept json
// @Produce json
// @Param id path string true "Halaqah Group ID"
// @Param request body dto.AddHalaqahMemberRequest true "Member data"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/halaqah/groups/{id}/members [post]
func (h *IslamicHandler) AddHalaqahMember(c *gin.Context) {
	halaqahID := c.Param("id")
	var req dto.AddHalaqahMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "invalid request", err.Error()))
		return
	}
	if err := h.islamicService.AddHalaqahMember(c.Request.Context(), halaqahID, req); err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "member added", nil))
}
