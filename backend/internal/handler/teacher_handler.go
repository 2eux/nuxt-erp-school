package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opencode/erp-school-backend/internal/dto"
	"github.com/opencode/erp-school-backend/internal/middleware"
	"github.com/opencode/erp-school-backend/internal/service"
)

type TeacherHandler struct {
	teacherService service.TeacherService
}

func NewTeacherHandler(teacherService service.TeacherService) *TeacherHandler {
	return &TeacherHandler{teacherService: teacherService}
}

// ListTeachers godoc
// @Summary List teachers
// @Tags Teachers
// @Produce json
// @Param page query int false "Page"
// @Param page_size query int false "Page size"
// @Param search query string false "Search"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/teachers [get]
func (h *TeacherHandler) ListTeachers(c *gin.Context) {
	schoolID := middleware.GetSchoolID(c)
	filter := getPagination(c)
	items, total, err := h.teacherService.ListTeachers(c.Request.Context(), schoolID, filter)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "success", dto.NewPaginatedResponse(items, total, filter.Page, filter.PageSize)))
}

// CreateTeacher godoc
// @Summary Create teacher
// @Tags Teachers
// @Accept json
// @Produce json
// @Param request body dto.CreateTeacherRequest true "Teacher data"
// @Success 201 {object} dto.APIResponse
// @Router /api/v1/teachers [post]
func (h *TeacherHandler) CreateTeacher(c *gin.Context) {
	schoolID := middleware.GetSchoolID(c)
	var req dto.CreateTeacherRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "invalid request", err.Error()))
		return
	}
	resp, err := h.teacherService.CreateTeacher(c.Request.Context(), schoolID, req)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusCreated, dto.NewAPIResponse(http.StatusCreated, "teacher created", resp))
}

// GetTeacher godoc
// @Summary Get teacher
// @Tags Teachers
// @Produce json
// @Param id path string true "Teacher ID"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/teachers/{id} [get]
func (h *TeacherHandler) GetTeacher(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.teacherService.GetTeacher(c.Request.Context(), id)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "success", resp))
}

// UpdateTeacher godoc
// @Summary Update teacher
// @Tags Teachers
// @Accept json
// @Produce json
// @Param id path string true "Teacher ID"
// @Param request body dto.CreateTeacherRequest true "Teacher data"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/teachers/{id} [put]
func (h *TeacherHandler) UpdateTeacher(c *gin.Context) {
	id := c.Param("id")
	var req dto.CreateTeacherRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "invalid request", err.Error()))
		return
	}
	resp, err := h.teacherService.UpdateTeacher(c.Request.Context(), id, req)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "teacher updated", resp))
}

// AssignSubjects godoc
// @Summary Assign subjects to teacher
// @Tags Teachers
// @Accept json
// @Produce json
// @Param id path string true "Teacher ID"
// @Param request body object true "Subject assignment data"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/teachers/{id}/subjects [post]
func (h *TeacherHandler) AssignSubjects(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		SubjectIDs []string `json:"subject_ids" binding:"required"`
		ClassIDs   []string `json:"class_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "invalid request", err.Error()))
		return
	}
	if err := h.teacherService.AssignSubjects(c.Request.Context(), id, req.SubjectIDs, req.ClassIDs); err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "subjects assigned", nil))
}

// GetSchedule godoc
// @Summary Get teacher schedule
// @Tags Teachers
// @Produce json
// @Param id path string true "Teacher ID"
// @Param semester_id query string false "Semester ID"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/teachers/{id}/schedule [get]
func (h *TeacherHandler) GetSchedule(c *gin.Context) {
	id := c.Param("id")
	semesterID := c.Query("semester_id")
	items, err := h.teacherService.GetSchedule(c.Request.Context(), id, semesterID)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "success", items))
}
