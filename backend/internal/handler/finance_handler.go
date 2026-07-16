package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opencode/erp-school-backend/internal/dto"
	"github.com/opencode/erp-school-backend/internal/middleware"
	"github.com/opencode/erp-school-backend/internal/service"
)

type FinanceHandler struct {
	financeService service.FinanceService
}

func NewFinanceHandler(financeService service.FinanceService) *FinanceHandler {
	return &FinanceHandler{financeService: financeService}
}

// ListFeeTypes godoc
// @Summary List fee types
// @Tags Finance
// @Produce json
// @Param page query int false "Page"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/fees/types [get]
func (h *FinanceHandler) ListFeeTypes(c *gin.Context) {
	schoolID := middleware.GetSchoolID(c)
	filter := getPagination(c)
	items, total, err := h.financeService.ListFeeTypes(c.Request.Context(), schoolID, filter)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "success", dto.NewPaginatedResponse(items, total, filter.Page, filter.PageSize)))
}

// CreateFeeType godoc
// @Summary Create fee type
// @Tags Finance
// @Accept json
// @Produce json
// @Param request body dto.CreateFeeTypeRequest true "Fee type data"
// @Success 201 {object} dto.APIResponse
// @Router /api/v1/fees/types [post]
func (h *FinanceHandler) CreateFeeType(c *gin.Context) {
	schoolID := middleware.GetSchoolID(c)
	var req dto.CreateFeeTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "invalid request", err.Error()))
		return
	}
	resp, err := h.financeService.CreateFeeType(c.Request.Context(), schoolID, req)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusCreated, dto.NewAPIResponse(http.StatusCreated, "fee type created", resp))
}

// ListInvoices godoc
// @Summary List invoices
// @Tags Finance
// @Produce json
// @Param page query int false "Page"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/invoices [get]
func (h *FinanceHandler) ListInvoices(c *gin.Context) {
	schoolID := middleware.GetSchoolID(c)
	filter := getPagination(c)
	items, total, err := h.financeService.ListInvoices(c.Request.Context(), schoolID, filter)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "success", dto.NewPaginatedResponse(items, total, filter.Page, filter.PageSize)))
}

// CreateInvoice godoc
// @Summary Create invoice
// @Tags Finance
// @Accept json
// @Produce json
// @Param request body dto.CreateInvoiceRequest true "Invoice data"
// @Success 201 {object} dto.APIResponse
// @Router /api/v1/invoices [post]
func (h *FinanceHandler) CreateInvoice(c *gin.Context) {
	schoolID := middleware.GetSchoolID(c)
	createdBy := middleware.GetUserID(c)
	var req dto.CreateInvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "invalid request", err.Error()))
		return
	}
	resp, err := h.financeService.CreateInvoice(c.Request.Context(), schoolID, createdBy, req)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusCreated, dto.NewAPIResponse(http.StatusCreated, "invoice created", resp))
}

// GetInvoice godoc
// @Summary Get invoice
// @Tags Finance
// @Produce json
// @Param id path string true "Invoice ID"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/invoices/{id} [get]
func (h *FinanceHandler) GetInvoice(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.financeService.GetInvoice(c.Request.Context(), id)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "success", resp))
}

// ListPayments godoc
// @Summary List payments
// @Tags Finance
// @Produce json
// @Param page query int false "Page"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/payments [get]
func (h *FinanceHandler) ListPayments(c *gin.Context) {
	schoolID := middleware.GetSchoolID(c)
	filter := getPagination(c)
	items, total, err := h.financeService.ListPayments(c.Request.Context(), schoolID, filter)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "success", dto.NewPaginatedResponse(items, total, filter.Page, filter.PageSize)))
}

// CreatePayment godoc
// @Summary Create payment
// @Tags Finance
// @Accept json
// @Produce json
// @Param request body dto.CreatePaymentRequest true "Payment data"
// @Success 201 {object} dto.APIResponse
// @Router /api/v1/payments [post]
func (h *FinanceHandler) CreatePayment(c *gin.Context) {
	var req dto.CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "invalid request", err.Error()))
		return
	}
	resp, err := h.financeService.CreatePayment(c.Request.Context(), req)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusCreated, dto.NewAPIResponse(http.StatusCreated, "payment created", resp))
}

// VerifyPayment godoc
// @Summary Verify payment
// @Tags Finance
// @Accept json
// @Produce json
// @Param id path string true "Payment ID"
// @Param request body dto.VerifyPaymentRequest true "Verification data"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/payments/{id}/verify [put]
func (h *FinanceHandler) VerifyPayment(c *gin.Context) {
	id := c.Param("id")
	verifiedBy := middleware.GetUserID(c)
	var req dto.VerifyPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "invalid request", err.Error()))
		return
	}
	if err := h.financeService.VerifyPayment(c.Request.Context(), id, verifiedBy, req); err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "payment verified", nil))
}

// ListJournals godoc
// @Summary List journals
// @Tags Finance
// @Produce json
// @Param page query int false "Page"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/journals [get]
func (h *FinanceHandler) ListJournals(c *gin.Context) {
	schoolID := middleware.GetSchoolID(c)
	filter := getPagination(c)
	items, total, err := h.financeService.ListJournals(c.Request.Context(), schoolID, filter)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "success", dto.NewPaginatedResponse(items, total, filter.Page, filter.PageSize)))
}

// CreateJournal godoc
// @Summary Create journal
// @Tags Finance
// @Accept json
// @Produce json
// @Param request body dto.CreateJournalRequest true "Journal data"
// @Success 201 {object} dto.APIResponse
// @Router /api/v1/journals [post]
func (h *FinanceHandler) CreateJournal(c *gin.Context) {
	schoolID := middleware.GetSchoolID(c)
	createdBy := middleware.GetUserID(c)
	var req dto.CreateJournalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "invalid request", err.Error()))
		return
	}
	resp, err := h.financeService.CreateJournal(c.Request.Context(), schoolID, createdBy, req)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusCreated, dto.NewAPIResponse(http.StatusCreated, "journal created", resp))
}

// ListLedger godoc
// @Summary List ledger
// @Tags Finance
// @Produce json
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/ledger [get]
func (h *FinanceHandler) ListLedger(c *gin.Context) {
	schoolID := middleware.GetSchoolID(c)
	accountCode := c.Query("account_code")
	items, err := h.financeService.ListLedger(c.Request.Context(), schoolID, accountCode)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "success", items))
}

// ListPayrollPeriods godoc
// @Summary List payroll periods
// @Tags Finance
// @Produce json
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/payroll/periods [get]
func (h *FinanceHandler) ListPayrollPeriods(c *gin.Context) {
	schoolID := middleware.GetSchoolID(c)
	items, err := h.financeService.ListPayrollPeriods(c.Request.Context(), schoolID)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "success", items))
}

// ListPayrollDetails godoc
// @Summary List payroll details
// @Tags Finance
// @Produce json
// @Param period query string true "Period"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/payroll/details [get]
func (h *FinanceHandler) ListPayrollDetails(c *gin.Context) {
	schoolID := middleware.GetSchoolID(c)
	period := c.Query("period")
	items, err := h.financeService.ListPayrollDetails(c.Request.Context(), schoolID, period)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "success", items))
}

// ProcessPayroll godoc
// @Summary Process payroll
// @Tags Finance
// @Accept json
// @Produce json
// @Param request body dto.CreatePayrollRequest true "Payroll data"
// @Success 201 {object} dto.APIResponse
// @Router /api/v1/payroll/process [post]
func (h *FinanceHandler) ProcessPayroll(c *gin.Context) {
	schoolID := middleware.GetSchoolID(c)
	var req dto.CreatePayrollRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "invalid request", err.Error()))
		return
	}
	resp, err := h.financeService.ProcessPayroll(c.Request.Context(), schoolID, req)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusCreated, dto.NewAPIResponse(http.StatusCreated, "payroll processed", resp))
}
