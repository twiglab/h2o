package handler

import (
	"admin-server/internal/service"
	"shared/response"
	"shared/utils"

	"github.com/gin-gonic/gin"
)

type PermissionHandler struct {
	permService *service.PermissionService
}

func NewPermissionHandler(permService *service.PermissionService) *PermissionHandler {
	return &PermissionHandler{permService: permService}
}

// GetTree 获取权限树
func (h *PermissionHandler) GetTree(c *gin.Context) {
	status := utils.ParseInt8Query(c, "status")

	tree, err := h.permService.GetTree(status)
	if err != nil {
		response.ServerError(c, "查询失败")
		return
	}

	response.Success(c, tree)
}

// Get 获取权限详情
func (h *PermissionHandler) Get(c *gin.Context) {
	id, err := utils.ParseInt64Param(c, "id")
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	permission, err := h.permService.GetByID(id)
	if err != nil {
		response.NotFound(c, "权限不存在")
		return
	}

	response.Success(c, permission)
}

type CreatePermissionRequest struct {
	ParentID  int64   `json:"parent_id"`
	PermCode  string  `json:"perm_code" binding:"required"`
	PermName  string  `json:"perm_name" binding:"required"`
	PermType  int8    `json:"perm_type"`
	Path      *string `json:"path"`
	Component *string `json:"component"`
	Icon      *string `json:"icon"`
	APIPath   *string `json:"api_path"`
	APIMethod *string `json:"api_method"`
	Visible   int8    `json:"visible"`
	SortOrder int     `json:"sort_order"`
	Status    int8    `json:"status"`
}

// Create 创建权限
func (h *PermissionHandler) Create(c *gin.Context) {
	var req CreatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	input := &service.CreatePermissionInput{
		ParentID:  req.ParentID,
		PermCode:  req.PermCode,
		PermName:  req.PermName,
		PermType:  req.PermType,
		Path:      req.Path,
		Component: req.Component,
		Icon:      req.Icon,
		APIPath:   req.APIPath,
		APIMethod: req.APIMethod,
		Visible:   req.Visible,
		SortOrder: req.SortOrder,
		Status:    req.Status,
	}

	permission, err := h.permService.Create(input)
	if err != nil {
		if err == service.ErrPermCodeExists {
			response.BadRequest(c, "权限编码已存在")
			return
		}
		response.ServerError(c, "创建失败")
		return
	}

	response.Success(c, permission)
}

type UpdatePermissionRequest struct {
	ParentID  *int64  `json:"parent_id"`
	PermName  *string `json:"perm_name"`
	PermType  *int8   `json:"perm_type"`
	Path      *string `json:"path"`
	Component *string `json:"component"`
	Icon      *string `json:"icon"`
	APIPath   *string `json:"api_path"`
	APIMethod *string `json:"api_method"`
	Visible   *int8   `json:"visible"`
	SortOrder *int    `json:"sort_order"`
	Status    *int8   `json:"status"`
}

// Update 更新权限
func (h *PermissionHandler) Update(c *gin.Context) {
	id, err := utils.ParseInt64Param(c, "id")
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	var req UpdatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	input := &service.UpdatePermissionInput{
		ParentID:  req.ParentID,
		PermName:  req.PermName,
		PermType:  req.PermType,
		Path:      req.Path,
		Component: req.Component,
		Icon:      req.Icon,
		APIPath:   req.APIPath,
		APIMethod: req.APIMethod,
		Visible:   req.Visible,
		SortOrder: req.SortOrder,
		Status:    req.Status,
	}

	permission, err := h.permService.Update(id, input)
	if err != nil {
		if err == service.ErrPermNotFound {
			response.NotFound(c, "权限不存在")
			return
		}
		response.ServerError(c, "更新失败")
		return
	}

	response.Success(c, permission)
}

// Delete 删除权限
func (h *PermissionHandler) Delete(c *gin.Context) {
	id, err := utils.ParseInt64Param(c, "id")
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	if err := h.permService.Delete(id); err != nil {
		if err == service.ErrPermHasChildren {
			response.BadRequest(c, "该权限下有子权限，无法删除")
			return
		}
		response.ServerError(c, "删除失败")
		return
	}

	response.SuccessWithMessage(c, "删除成功", nil)
}
