package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/opencode/erp-school-backend/internal/domain"
	"github.com/opencode/erp-school-backend/internal/dto"
)

type errorInfo struct {
	Code    int
	Message string
}

func parseError(err error) errorInfo {
	var de *domain.DomainError
	if errors.As(err, &de) {
		switch de.Code {
		case "NOT_FOUND":
			return errorInfo{Code: http.StatusNotFound, Message: de.Message}
		case "DUPLICATE":
			return errorInfo{Code: http.StatusConflict, Message: de.Message}
		case "FORBIDDEN":
			return errorInfo{Code: http.StatusForbidden, Message: de.Message}
		case "UNAUTHORIZED":
			return errorInfo{Code: http.StatusUnauthorized, Message: de.Message}
		case "INVALID_INPUT", "VALIDATION_ERROR":
			return errorInfo{Code: http.StatusBadRequest, Message: de.Message}
		case "RATE_LIMIT":
			return errorInfo{Code: http.StatusTooManyRequests, Message: de.Message}
		default:
			return errorInfo{Code: http.StatusInternalServerError, Message: de.Message}
		}
	}

	if errors.Is(err, domain.ErrNotFound) {
		return errorInfo{Code: http.StatusNotFound, Message: "resource not found"}
	}
	if errors.Is(err, domain.ErrUnauthorized) {
		return errorInfo{Code: http.StatusUnauthorized, Message: "unauthorized"}
	}
	if errors.Is(err, domain.ErrForbidden) {
		return errorInfo{Code: http.StatusForbidden, Message: "forbidden"}
	}
	if errors.Is(err, domain.ErrInvalidInput) || errors.Is(err, domain.ErrValidation) {
		return errorInfo{Code: http.StatusBadRequest, Message: err.Error()}
	}

	return errorInfo{Code: http.StatusInternalServerError, Message: "internal server error"}
}

func getPagination(c *gin.Context) dto.PaginationRequest {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	search := c.Query("search")
	sortBy := c.Query("sort_by")
	sortDir := c.Query("sort_dir")

	p := dto.PaginationRequest{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
		SortBy:   sortBy,
		SortDir:  sortDir,
	}
	p.Defaults()
	return p
}
