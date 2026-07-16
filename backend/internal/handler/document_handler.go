package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opencode/erp-school-backend/internal/dto"
	"github.com/opencode/erp-school-backend/internal/middleware"
	"github.com/opencode/erp-school-backend/internal/service"
)

type DocumentHandler struct {
	documentService service.DocumentService
}

func NewDocumentHandler(documentService service.DocumentService) *DocumentHandler {
	return &DocumentHandler{documentService: documentService}
}

// ListDocuments godoc
// @Summary List documents
// @Tags Documents
// @Produce json
// @Param page query int false "Page"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/documents [get]
func (h *DocumentHandler) ListDocuments(c *gin.Context) {
	schoolID := middleware.GetSchoolID(c)
	filter := getPagination(c)
	items, total, err := h.documentService.ListDocuments(c.Request.Context(), schoolID, filter)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "success", dto.NewPaginatedResponse(items, total, filter.Page, filter.PageSize)))
}

// UploadDocument godoc
// @Summary Upload document
// @Tags Documents
// @Accept json
// @Produce json
// @Param request body dto.CreateDocumentRequest true "Document data"
// @Success 201 {object} dto.APIResponse
// @Router /api/v1/documents [post]
func (h *DocumentHandler) UploadDocument(c *gin.Context) {
	schoolID := middleware.GetSchoolID(c)
	createdBy := middleware.GetUserID(c)
	var req dto.CreateDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "invalid request", err.Error()))
		return
	}
	resp, err := h.documentService.UploadDocument(c.Request.Context(), schoolID, createdBy, req.Title, req.DocType, "", 0, "")
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusCreated, dto.NewAPIResponse(http.StatusCreated, "document uploaded", resp))
}

// GetDocument godoc
// @Summary Get document
// @Tags Documents
// @Produce json
// @Param id path string true "Document ID"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/documents/{id} [get]
func (h *DocumentHandler) GetDocument(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.documentService.GetDocument(c.Request.Context(), id)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "success", resp))
}

// DeleteDocument godoc
// @Summary Delete document
// @Tags Documents
// @Produce json
// @Param id path string true "Document ID"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/documents/{id} [delete]
func (h *DocumentHandler) DeleteDocument(c *gin.Context) {
	id := c.Param("id")
	if err := h.documentService.DeleteDocument(c.Request.Context(), id); err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "document deleted", nil))
}

// ListLetters godoc
// @Summary List letters
// @Tags Documents
// @Produce json
// @Param page query int false "Page"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/letters [get]
func (h *DocumentHandler) ListLetters(c *gin.Context) {
	schoolID := middleware.GetSchoolID(c)
	filter := getPagination(c)
	items, total, err := h.documentService.ListLetters(c.Request.Context(), schoolID, filter)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "success", dto.NewPaginatedResponse(items, total, filter.Page, filter.PageSize)))
}

// CreateLetter godoc
// @Summary Create letter
// @Tags Documents
// @Accept json
// @Produce json
// @Param request body dto.CreateLetterRequest true "Letter data"
// @Success 201 {object} dto.APIResponse
// @Router /api/v1/letters [post]
func (h *DocumentHandler) CreateLetter(c *gin.Context) {
	schoolID := middleware.GetSchoolID(c)
	createdBy := middleware.GetUserID(c)
	var req dto.CreateLetterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "invalid request", err.Error()))
		return
	}
	resp, err := h.documentService.CreateLetter(c.Request.Context(), schoolID, createdBy, req)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusCreated, dto.NewAPIResponse(http.StatusCreated, "letter created", resp))
}

// GetLetter godoc
// @Summary Get letter
// @Tags Documents
// @Produce json
// @Param id path string true "Letter ID"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/letters/{id} [get]
func (h *DocumentHandler) GetLetter(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.documentService.GetLetter(c.Request.Context(), id)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "success", resp))
}
