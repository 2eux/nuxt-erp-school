package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opencode/erp-school-backend/internal/dto"
	"github.com/opencode/erp-school-backend/internal/middleware"
	"github.com/opencode/erp-school-backend/internal/service"
)

type EmployeeHandler struct {
	employeeService service.EmployeeService
}

func NewEmployeeHandler(employeeService service.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{employeeService: employeeService}
}

// ListEmployees godoc
// @Summary List employees
// @Tags Employees
// @Produce json
// @Param page query int false "Page"
// @Param page_size query int false "Page size"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/employees [get]
func (h *EmployeeHandler) ListEmployees(c *gin.Context) {
	schoolID := middleware.GetSchoolID(c)
	filter := getPagination(c)
	items, total, err := h.employeeService.ListEmployees(c.Request.Context(), schoolID, filter)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "success", dto.NewPaginatedResponse(items, total, filter.Page, filter.PageSize)))
}

// CreateEmployee godoc
// @Summary Create employee
// @Tags Employees
// @Accept json
// @Produce json
// @Param request body dto.CreateEmployeeRequest true "Employee data"
// @Success 201 {object} dto.APIResponse
// @Router /api/v1/employees [post]
func (h *EmployeeHandler) CreateEmployee(c *gin.Context) {
	schoolID := middleware.GetSchoolID(c)
	var req dto.CreateEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "invalid request", err.Error()))
		return
	}
	resp, err := h.employeeService.CreateEmployee(c.Request.Context(), schoolID, req)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusCreated, dto.NewAPIResponse(http.StatusCreated, "employee created", resp))
}

// GetEmployee godoc
// @Summary Get employee
// @Tags Employees
// @Produce json
// @Param id path string true "Employee ID"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/employees/{id} [get]
func (h *EmployeeHandler) GetEmployee(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.employeeService.GetEmployee(c.Request.Context(), id)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "success", resp))
}

// UpdateEmployee godoc
// @Summary Update employee
// @Tags Employees
// @Accept json
// @Produce json
// @Param id path string true "Employee ID"
// @Param request body dto.CreateEmployeeRequest true "Employee data"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/employees/{id} [put]
func (h *EmployeeHandler) UpdateEmployee(c *gin.Context) {
	id := c.Param("id")
	var req dto.CreateEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "invalid request", err.Error()))
		return
	}
	resp, err := h.employeeService.UpdateEmployee(c.Request.Context(), id, req)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "employee updated", resp))
}

// ListAttendances godoc
// @Summary List employee attendances
// @Tags Employees
// @Produce json
// @Param id path string true "Employee ID"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/employees/{id}/attendances [get]
func (h *EmployeeHandler) ListAttendances(c *gin.Context) {
	employeeID := c.Param("id")
	startDate := time.Now().AddDate(0, -1, 0)
	endDate := time.Now()
	items, err := h.employeeService.ListAttendances(c.Request.Context(), employeeID, startDate, endDate)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "success", items))
}

// CreateAttendance godoc
// @Summary Create employee attendance
// @Tags Employees
// @Accept json
// @Produce json
// @Param id path string true "Employee ID"
// @Param request body dto.CreateEmployeeAttendanceRequest true "Attendance data"
// @Success 201 {object} dto.APIResponse
// @Router /api/v1/employees/{id}/attendances [post]
func (h *EmployeeHandler) CreateAttendance(c *gin.Context) {
	employeeID := c.Param("id")
	var req dto.CreateEmployeeAttendanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "invalid request", err.Error()))
		return
	}
	resp, err := h.employeeService.CreateAttendance(c.Request.Context(), employeeID, req)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusCreated, dto.NewAPIResponse(http.StatusCreated, "attendance created", resp))
}

// ListLeaveRequests godoc
// @Summary List leave requests
// @Tags Employees
// @Produce json
// @Param page query int false "Page"
// @Param page_size query int false "Page size"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/leave-requests [get]
func (h *EmployeeHandler) ListLeaveRequests(c *gin.Context) {
	schoolID := middleware.GetSchoolID(c)
	filter := getPagination(c)
	items, total, err := h.employeeService.ListLeaveRequests(c.Request.Context(), schoolID, filter)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "success", dto.NewPaginatedResponse(items, total, filter.Page, filter.PageSize)))
}

// SubmitLeaveRequest godoc
// @Summary Submit leave request
// @Tags Employees
// @Accept json
// @Produce json
// @Param request body dto.CreateLeaveRequest true "Leave request data"
// @Success 201 {object} dto.APIResponse
// @Router /api/v1/leave-requests [post]
func (h *EmployeeHandler) SubmitLeaveRequest(c *gin.Context) {
	userID := middleware.GetUserID(c)
	schoolID := middleware.GetSchoolID(c)
	var req dto.CreateLeaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "invalid request", err.Error()))
		return
	}
	resp, err := h.employeeService.SubmitLeaveRequest(c.Request.Context(), userID, schoolID, req)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusCreated, dto.NewAPIResponse(http.StatusCreated, "leave request submitted", resp))
}

// ApproveLeave godoc
// @Summary Approve leave request
// @Tags Employees
// @Produce json
// @Param id path string true "Leave Request ID"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/leave-requests/{id}/approve [put]
func (h *EmployeeHandler) ApproveLeave(c *gin.Context) {
	id := c.Param("id")
	approverID := middleware.GetUserID(c)
	if err := h.employeeService.ApproveLeave(c.Request.Context(), id, approverID); err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "leave approved", nil))
}

// RejectLeave godoc
// @Summary Reject leave request
// @Tags Employees
// @Produce json
// @Param id path string true "Leave Request ID"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/leave-requests/{id}/reject [put]
func (h *EmployeeHandler) RejectLeave(c *gin.Context) {
	id := c.Param("id")
	approverID := middleware.GetUserID(c)
	if err := h.employeeService.RejectLeave(c.Request.Context(), id, approverID); err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "leave rejected", nil))
}
