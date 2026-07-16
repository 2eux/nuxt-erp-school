package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opencode/erp-school-backend/internal/dto"
	"github.com/opencode/erp-school-backend/internal/middleware"
	"github.com/opencode/erp-school-backend/internal/service"
)

type SchoolHandler struct {
	schoolService service.SchoolService
}

func NewSchoolHandler(schoolService service.SchoolService) *SchoolHandler {
	return &SchoolHandler{schoolService: schoolService}
}

// ListSchools godoc
// @Summary List schools
// @Tags Schools
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param search query string false "Search term"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/schools [get]
func (h *SchoolHandler) ListSchools(c *gin.Context) {
	filter := getPagination(c)
	items, total, err := h.schoolService.ListSchools(c.Request.Context(), filter)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "success", dto.NewPaginatedResponse(items, total, filter.Page, filter.PageSize)))
}

// CreateSchool godoc
// @Summary Create school
// @Tags Schools
// @Accept json
// @Produce json
// @Param request body dto.CreateSchoolRequest true "School data"
// @Success 201 {object} dto.APIResponse
// @Router /api/v1/schools [post]
func (h *SchoolHandler) CreateSchool(c *gin.Context) {
	var req dto.CreateSchoolRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "invalid request", err.Error()))
		return
	}
	resp, err := h.schoolService.CreateSchool(c.Request.Context(), req)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusCreated, dto.NewAPIResponse(http.StatusCreated, "school created", resp))
}

// GetSchool godoc
// @Summary Get school
// @Tags Schools
// @Produce json
// @Param id path string true "School ID"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/schools/{id} [get]
func (h *SchoolHandler) GetSchool(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.schoolService.GetSchool(c.Request.Context(), id)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "success", resp))
}

// UpdateSchool godoc
// @Summary Update school
// @Tags Schools
// @Accept json
// @Produce json
// @Param id path string true "School ID"
// @Param request body dto.UpdateSchoolRequest true "School data"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/schools/{id} [put]
func (h *SchoolHandler) UpdateSchool(c *gin.Context) {
	id := c.Param("id")
	var req dto.UpdateSchoolRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "invalid request", err.Error()))
		return
	}
	resp, err := h.schoolService.UpdateSchool(c.Request.Context(), id, req)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "school updated", resp))
}

// DeleteSchool godoc
// @Summary Delete school
// @Tags Schools
// @Produce json
// @Param id path string true "School ID"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/schools/{id} [delete]
func (h *SchoolHandler) DeleteSchool(c *gin.Context) {
	id := c.Param("id")
	if err := h.schoolService.DeleteSchool(c.Request.Context(), id); err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "school deleted", nil))
}

// ListAcademicYears godoc
// @Summary List academic years
// @Tags Schools
// @Produce json
// @Param id path string true "School ID"
// @Param page query int false "Page"
// @Param page_size query int false "Page size"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/schools/{id}/academic-years [get]
func (h *SchoolHandler) ListAcademicYears(c *gin.Context) {
	schoolID := c.Param("id")
	filter := getPagination(c)
	items, total, err := h.schoolService.ListAcademicYears(c.Request.Context(), schoolID, filter)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "success", dto.NewPaginatedResponse(items, total, filter.Page, filter.PageSize)))
}

// CreateAcademicYear godoc
// @Summary Create academic year
// @Tags Schools
// @Accept json
// @Produce json
// @Param id path string true "School ID"
// @Param request body dto.CreateAcademicYearRequest true "Academic year data"
// @Success 201 {object} dto.APIResponse
// @Router /api/v1/schools/{id}/academic-years [post]
func (h *SchoolHandler) CreateAcademicYear(c *gin.Context) {
	schoolID := c.Param("id")
	var req dto.CreateAcademicYearRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "invalid request", err.Error()))
		return
	}
	resp, err := h.schoolService.CreateAcademicYear(c.Request.Context(), schoolID, req)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusCreated, dto.NewAPIResponse(http.StatusCreated, "academic year created", resp))
}

// UpdateAcademicYear godoc
// @Summary Update academic year
// @Tags Schools
// @Accept json
// @Produce json
// @Param id path string true "School ID"
// @Param ayid path string true "Academic Year ID"
// @Param request body dto.UpdateAcademicYearRequest true "Academic year data"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/schools/{id}/academic-years/{ayid} [put]
func (h *SchoolHandler) UpdateAcademicYear(c *gin.Context) {
	ayID := c.Param("ayid")
	var req dto.UpdateAcademicYearRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "invalid request", err.Error()))
		return
	}
	resp, err := h.schoolService.UpdateAcademicYear(c.Request.Context(), ayID, req)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "academic year updated", resp))
}

// ListSemesters godoc
// @Summary List semesters
// @Tags Schools
// @Produce json
// @Param id path string true "Academic Year ID"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/academic-years/{id}/semesters [get]
func (h *SchoolHandler) ListSemesters(c *gin.Context) {
	ayID := c.Param("id")
	items, err := h.schoolService.ListSemesters(c.Request.Context(), ayID)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "success", items))
}

// CreateSemester godoc
// @Summary Create semester
// @Tags Schools
// @Accept json
// @Produce json
// @Param id path string true "Academic Year ID"
// @Param request body dto.CreateSemesterRequest true "Semester data"
// @Success 201 {object} dto.APIResponse
// @Router /api/v1/academic-years/{id}/semesters [post]
func (h *SchoolHandler) CreateSemester(c *gin.Context) {
	var req dto.CreateSemesterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "invalid request", err.Error()))
		return
	}
	req.AcademicYearID = c.Param("id")
	resp, err := h.schoolService.CreateSemester(c.Request.Context(), req)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusCreated, dto.NewAPIResponse(http.StatusCreated, "semester created", resp))
}

// ListGrades godoc
// @Summary List grades
// @Tags Schools
// @Produce json
// @Param id path string true "School ID"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/schools/{id}/grades [get]
func (h *SchoolHandler) ListGrades(c *gin.Context) {
	schoolID := c.Param("id")
	items, err := h.schoolService.ListGrades(c.Request.Context(), schoolID)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "success", items))
}

// CreateGrade godoc
// @Summary Create grade
// @Tags Schools
// @Accept json
// @Produce json
// @Param id path string true "School ID"
// @Param request body dto.CreateGradeRequest true "Grade data"
// @Success 201 {object} dto.APIResponse
// @Router /api/v1/schools/{id}/grades [post]
func (h *SchoolHandler) CreateGrade(c *gin.Context) {
	schoolID := c.Param("id")
	var req dto.CreateGradeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "invalid request", err.Error()))
		return
	}
	resp, err := h.schoolService.CreateGrade(c.Request.Context(), schoolID, req)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusCreated, dto.NewAPIResponse(http.StatusCreated, "grade created", resp))
}

// ListClasses godoc
// @Summary List classes
// @Tags Schools
// @Produce json
// @Param id path string true "School ID"
// @Param page query int false "Page"
// @Param page_size query int false "Page size"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/schools/{id}/classes [get]
func (h *SchoolHandler) ListClasses(c *gin.Context) {
	schoolID := c.Param("id")
	filter := getPagination(c)
	items, total, err := h.schoolService.ListClasses(c.Request.Context(), schoolID, filter)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "success", dto.NewPaginatedResponse(items, total, filter.Page, filter.PageSize)))
}

// CreateClass godoc
// @Summary Create class
// @Tags Schools
// @Accept json
// @Produce json
// @Param id path string true "School ID"
// @Param request body dto.CreateClassRequest true "Class data"
// @Success 201 {object} dto.APIResponse
// @Router /api/v1/schools/{id}/classes [post]
func (h *SchoolHandler) CreateClass(c *gin.Context) {
	schoolID := c.Param("id")
	var req dto.CreateClassRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "invalid request", err.Error()))
		return
	}
	resp, err := h.schoolService.CreateClass(c.Request.Context(), schoolID, req)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusCreated, dto.NewAPIResponse(http.StatusCreated, "class created", resp))
}

// ListSubjects godoc
// @Summary List subjects
// @Tags Schools
// @Produce json
// @Param id path string true "School ID"
// @Param page query int false "Page"
// @Param page_size query int false "Page size"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/schools/{id}/subjects [get]
func (h *SchoolHandler) ListSubjects(c *gin.Context) {
	schoolID := c.Param("id")
	filter := getPagination(c)
	items, total, err := h.schoolService.ListSubjects(c.Request.Context(), schoolID, filter)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "success", dto.NewPaginatedResponse(items, total, filter.Page, filter.PageSize)))
}

// CreateSubject godoc
// @Summary Create subject
// @Tags Schools
// @Accept json
// @Produce json
// @Param id path string true "School ID"
// @Param request body dto.CreateSubjectRequest true "Subject data"
// @Success 201 {object} dto.APIResponse
// @Router /api/v1/schools/{id}/subjects [post]
func (h *SchoolHandler) CreateSubject(c *gin.Context) {
	schoolID := c.Param("id")
	var req dto.CreateSubjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "invalid request", err.Error()))
		return
	}
	resp, err := h.schoolService.CreateSubject(c.Request.Context(), schoolID, req)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusCreated, dto.NewAPIResponse(http.StatusCreated, "subject created", resp))
}
