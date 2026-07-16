package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opencode/erp-school-backend/internal/dto"
	"github.com/opencode/erp-school-backend/internal/middleware"
	"github.com/opencode/erp-school-backend/internal/service"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// ListUsers godoc
// @Summary List users
// @Tags Users
// @Produce json
// @Param page query int false "Page"
// @Param page_size query int false "Page size"
// @Param search query string false "Search"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	schoolID := middleware.GetSchoolID(c)
	filter := getPagination(c)
	items, total, err := h.userService.ListUsers(c.Request.Context(), schoolID, filter)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "success", dto.NewPaginatedResponse(items, total, filter.Page, filter.PageSize)))
}

// CreateUser godoc
// @Summary Create user
// @Tags Users
// @Accept json
// @Produce json
// @Param request body dto.CreateUserRequest true "User data"
// @Success 201 {object} dto.APIResponse
// @Router /api/v1/users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	schoolID := middleware.GetSchoolID(c)
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "invalid request", err.Error()))
		return
	}
	resp, err := h.userService.CreateUser(c.Request.Context(), schoolID, req)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusCreated, dto.NewAPIResponse(http.StatusCreated, "user created", resp))
}

// GetUser godoc
// @Summary Get user
// @Tags Users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.userService.GetUser(c.Request.Context(), id)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "success", resp))
}

// UpdateUser godoc
// @Summary Update user
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param request body dto.UpdateUserRequest true "User data"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "invalid request", err.Error()))
		return
	}
	resp, err := h.userService.UpdateUser(c.Request.Context(), id, req)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "user updated", resp))
}

// DeleteUser godoc
// @Summary Delete user
// @Tags Users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if err := h.userService.DeleteUser(c.Request.Context(), id); err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "user deleted", nil))
}

// ListRoles godoc
// @Summary List roles
// @Tags Users
// @Produce json
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/roles [get]
func (h *UserHandler) ListRoles(c *gin.Context) {
	schoolID := middleware.GetSchoolID(c)
	items, err := h.userService.ListRoles(c.Request.Context(), schoolID)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "success", items))
}

// CreateRole godoc
// @Summary Create role
// @Tags Users
// @Accept json
// @Produce json
// @Param request body dto.CreateRoleRequest true "Role data"
// @Success 201 {object} dto.APIResponse
// @Router /api/v1/roles [post]
func (h *UserHandler) CreateRole(c *gin.Context) {
	schoolID := middleware.GetSchoolID(c)
	var req dto.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "invalid request", err.Error()))
		return
	}
	resp, err := h.userService.CreateRole(c.Request.Context(), schoolID, req)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusCreated, dto.NewAPIResponse(http.StatusCreated, "role created", resp))
}

// UpdateRole godoc
// @Summary Update role
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "Role ID"
// @Param request body dto.UpdateRoleRequest true "Role data"
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/roles/{id} [put]
func (h *UserHandler) UpdateRole(c *gin.Context) {
	id := c.Param("id")
	var req dto.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(http.StatusBadRequest, "invalid request", err.Error()))
		return
	}
	resp, err := h.userService.UpdateRole(c.Request.Context(), id, req)
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "role updated", resp))
}

// ListPermissions godoc
// @Summary List permissions
// @Tags Users
// @Produce json
// @Success 200 {object} dto.APIResponse
// @Router /api/v1/permissions [get]
func (h *UserHandler) ListPermissions(c *gin.Context) {
	items, err := h.userService.ListPermissions(c.Request.Context())
	if err != nil {
		ei := parseError(err)
		c.JSON(ei.Code, dto.NewErrorResponse(ei.Code, ei.Message, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.NewAPIResponse(http.StatusOK, "success", items))
}
