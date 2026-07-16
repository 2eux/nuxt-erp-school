package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opencode/erp-school-backend/internal/dto"
	"github.com/opencode/erp-school-backend/internal/middleware"
	"github.com/opencode/erp-school-backend/internal/service"
)

type StudentHandler struct {
	studentService service.StudentService
}

func NewStudentHandler(studentService service.StudentService) *StudentHandler {
	return &StudentHandler{studentService: studentService}
}

// ListStudents godoc
// @Summary List students
// @Tags Students
// @Produce json
// @Param page query int false "Page"
// @Param page_size query int false "Page size"
// @Param class_id query string false "Class ID"
// @Param status query string false "Status"
// @Param search query string false "Search"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/students [get]
func (h *StudentHandler) ListStudents(c *gin.Context) {
	schoolID := middleware.GetSchoolID(c)
	filter := getPagination(c)
	classID := c.Query("class_id")
	status := c.Query("status")

	items, total, err := h.studentService.ListStudents(c.Request.Context(), schoolID, filter, classID, status)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "success", dto.NewPaginatedResponse(items, total, filter.Page, filter.PageSize)))
}

// CreateStudent godoc
// @Summary Create student
// @Tags Students
// @Accept json
// @Produce json
// @Param request body dto.CreateStudentRequest true "Student data"
// @Success 201 {object} dto.APIResponse
// @Router /api/v1/students [post]
func (h *StudentHandler) CreateStudent(c *gin.Context) {
	schoolID := middleware.GetSchoolID(c)
	var req dto.CreateStudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "invalid request", err.Error()))
		return
	}
	resp, err := h.studentService.CreateStudent(c.Request.Context(), schoolID, req)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusCreated, dto.NewAPIResponse(http.StatusCreated, "student created", resp))
}

// GetStudent godoc
// @Summary Get student
// @Tags Students
// @Produce json
// @Param id path string true "Student ID"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/students/{id} [get]
func (h *StudentHandler) GetStudent(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.studentService.GetStudent(c.Request.Context(), id)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "success", resp))
}

// UpdateStudent godoc
// @Summary Update student
// @Tags Students
// @Accept json
// @Produce json
// @Param id path string true "Student ID"
// @Param request body dto.UpdateStudentRequest true "Student data"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/students/{id} [put]
func (h *StudentHandler) UpdateStudent(c *gin.Context) {
	id := c.Param("id")
	var req dto.UpdateStudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "invalid request", err.Error()))
		return
	}
	resp, err := h.studentService.UpdateStudent(c.Request.Context(), id, req)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "student updated", resp))
}

// DeleteStudent godoc
// @Summary Delete student
// @Tags Students
// @Produce json
// @Param id path string true "Student ID"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/students/{id} [delete]
func (h *StudentHandler) DeleteStudent(c *gin.Context) {
	id := c.Param("id")
	if err := h.studentService.DeleteStudent(c.Request.Context(), id); err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "student deleted", nil))
}

// PromoteStudent godoc
// @Summary Promote student
// @Tags Students
// @Accept json
// @Produce json
// @Param id path string true "Student ID"
// @Param request body object true "Promote data {new_class_id, new_academic_year_id}"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/students/{id}/promote [post]
func (h *StudentHandler) PromoteStudent(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		NewClassID        string `json:"new_class_id" binding:"required"`
		NewAcademicYearID string `json:"new_academic_year_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "invalid request", err.Error()))
		return
	}
	resp, err := h.studentService.PromoteStudent(c.Request.Context(), id, req.NewClassID, req.NewAcademicYearID)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "student promoted", resp))
}

// TransferStudent godoc
// @Summary Transfer student
// @Tags Students
// @Accept json
// @Produce json
// @Param id path string true "Student ID"
// @Param request body object true "Transfer data"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/students/{id}/transfer [post]
func (h *StudentHandler) TransferStudent(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		NewSchoolID string `json:"new_school_id" binding:"required"`
		NewClassID  string `json:"new_class_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "invalid request", err.Error()))
		return
	}
	resp, err := h.studentService.TransferStudent(c.Request.Context(), id, req.NewSchoolID, req.NewClassID)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "student transferred", resp))
}

// ListParents godoc
// @Summary List student parents
// @Tags Students
// @Produce json
// @Param id path string true "Student ID"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/students/{id}/parents [get]
func (h *StudentHandler) ListParents(c *gin.Context) {
	id := c.Param("id")
	items, err := h.studentService.ListParents(c.Request.Context(), id)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "success", items))
}

// LinkParent godoc
// @Summary Link parent to student
// @Tags Students
// @Accept json
// @Produce json
// @Param id path string true "Student ID"
// @Param request body dto.CreateStudentParentRequest true "Parent data"
// @Success 201 {object} dto.APIResponse
// @Router /api/v1/students/{id}/parents [post]
func (h *StudentHandler) LinkParent(c *gin.Context) {
	studentID := c.Param("id")
	schoolID := middleware.GetSchoolID(c)
	var req dto.CreateStudentParentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "invalid request", err.Error()))
		return
	}
	resp, err := h.studentService.LinkParent(c.Request.Context(), studentID, schoolID, req)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusCreated, dto.NewAPIResponse(http.StatusCreated, "parent linked", resp))
}

// ListDocuments godoc
// @Summary List student documents
// @Tags Students
// @Produce json
// @Param id path string true "Student ID"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/students/{id}/documents [get]
func (h *StudentHandler) ListDocuments(c *gin.Context) {
	id := c.Param("id")
	items, err := h.studentService.ListDocuments(c.Request.Context(), id)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "success", items))
}

// UploadDocument godoc
// @Summary Upload student document
// @Tags Students
// @Accept json
// @Produce json
// @Param id path string true "Student ID"
// @Param request body object true "Document data"
// @Success 201 {object} dto.APIResponse
// @Router /api/v1/students/{id}/documents [post]
func (h *StudentHandler) UploadDocument(c *gin.Context) {
	studentID := c.Param("id")
	var req struct {
		Name    string `json:"name" binding:"required"`
		DocType string `json:"doc_type" binding:"required"`
		FileURL string `json:"file_url" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "invalid request", err.Error()))
		return
	}
	resp, err := h.studentService.UploadDocument(c.Request.Context(), studentID, req.Name, req.DocType, req.FileURL)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusCreated, dto.NewAPIResponse(http.StatusCreated, "document uploaded", resp))
}
