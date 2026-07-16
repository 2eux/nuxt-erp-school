package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opencode/erp-school-backend/internal/dto"
	"github.com/opencode/erp-school-backend/internal/middleware"
	"github.com/opencode/erp-school-backend/internal/service"
)

type AdmissionHandler struct {
	admissionService service.AdmissionService
}

func NewAdmissionHandler(admissionService service.AdmissionService) *AdmissionHandler {
	return &AdmissionHandler{admissionService: admissionService}
}

// ListApplicants godoc
// @Summary List admission applicants
// @Tags Admissions
// @Produce json
// @Param page query int false "Page"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/admissions/applicants [get]
func (h *AdmissionHandler) ListApplicants(c *gin.Context) {
	schoolID := middleware.GetSchoolID(c)
	filter := getPagination(c)
	items, total, err := h.admissionService.ListApplicants(c.Request.Context(), schoolID, filter)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "success", dto.NewPaginatedResponse(items, total, filter.Page, filter.PageSize)))
}

// CreateApplicant godoc
// @Summary Create admission applicant
// @Tags Admissions
// @Accept json
// @Produce json
// @Param request body dto.CreateAdmissionRequest true "Applicant data"
// @Success 201 {object} dto.APIResponse
// @Router /api/v1/admissions/applicants [post]
func (h *AdmissionHandler) CreateApplicant(c *gin.Context) {
	schoolID := middleware.GetSchoolID(c)
	var req dto.CreateAdmissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "invalid request", err.Error()))
		return
	}
	resp, err := h.admissionService.CreateApplicant(c.Request.Context(), schoolID, req)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusCreated, dto.NewAPIResponse(http.StatusCreated, "applicant created", resp))
}

// AcceptApplicant godoc
// @Summary Accept applicant
// @Tags Admissions
// @Produce json
// @Param id path string true "Applicant ID"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/admissions/applicants/{id}/accept [post]
func (h *AdmissionHandler) AcceptApplicant(c *gin.Context) {
	id := c.Param("id")
	if err := h.admissionService.AcceptApplicant(c.Request.Context(), id); err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "applicant accepted", nil))
}

// EnrollApplicant godoc
// @Summary Enroll accepted applicant
// @Tags Admissions
// @Accept json
// @Produce json
// @Param id path string true "Applicant ID"
// @Param request body object true "Enrollment data {class_id}"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/admissions/applicants/{id}/enroll [post]
func (h *AdmissionHandler) EnrollApplicant(c *gin.Context) {
	schoolID := middleware.GetSchoolID(c)
	id := c.Param("id")
	var req struct {
		ClassID string `json:"class_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "invalid request", err.Error()))
		return
	}
	resp, err := h.admissionService.EnrollApplicant(c.Request.Context(), id, schoolID, req.ClassID)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "applicant enrolled", resp))
}
