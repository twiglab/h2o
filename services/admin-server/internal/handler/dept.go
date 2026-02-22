package handler

import (
	"admin-server/internal/service"
	"shared/response"
	"shared/utils"

	"github.com/gin-gonic/gin"
)

type DeptHandler struct {
	deptService *service.DeptService
}

func NewDeptHandler(deptService *service.DeptService) *DeptHandler {
	return &DeptHandler{deptService: deptService}
}

// List 获取部门列表（树形）
func (h *DeptHandler) List(c *gin.Context) {
	keyword := c.Query("keyword")
	status := utils.ParseInt8Query(c, "status")

	depts, err := h.deptService.List(keyword, status)
	if err != nil {
		response.ServerError(c, "查询失败")
		return
	}

	response.Success(c, depts)
}

// Get 获取部门详情
func (h *DeptHandler) Get(c *gin.Context) {
	id, err := utils.ParseInt64Param(c, "id")
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	dept, err := h.deptService.GetByID(id)
	if err != nil {
		response.NotFound(c, "部门不存在")
		return
	}

	response.Success(c, dept)
}

type CreateDeptRequest struct {
	ParentID  int64  `json:"parent_id"`
	DeptCode  string `json:"dept_code" binding:"required"`
	DeptName  string `json:"dept_name" binding:"required"`
	LeaderID  *int64 `json:"leader_id"`
	SortOrder int    `json:"sort_order"`
	Status    int8   `json:"status"`
}

// Create 创建部门
func (h *DeptHandler) Create(c *gin.Context) {
	var req CreateDeptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	input := &service.CreateDeptInput{
		ParentID:  req.ParentID,
		DeptCode:  req.DeptCode,
		DeptName:  req.DeptName,
		LeaderID:  req.LeaderID,
		SortOrder: req.SortOrder,
		Status:    req.Status,
	}

	dept, err := h.deptService.Create(input)
	if err != nil {
		if err == service.ErrDeptCodeExists {
			response.BadRequest(c, "部门编码已存在")
			return
		}
		if err == service.ErrDeptNotFound {
			response.BadRequest(c, "父部门不存在")
			return
		}
		response.ServerError(c, "创建失败")
		return
	}

	response.Success(c, dept)
}

type UpdateDeptRequest struct {
	ParentID  *int64  `json:"parent_id"`
	DeptCode  *string `json:"dept_code"`
	DeptName  *string `json:"dept_name"`
	LeaderID  *int64  `json:"leader_id"`
	SortOrder *int    `json:"sort_order"`
	Status    *int8   `json:"status"`
}

// Update 更新部门
func (h *DeptHandler) Update(c *gin.Context) {
	id, err := utils.ParseInt64Param(c, "id")
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	var req UpdateDeptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	input := &service.UpdateDeptInput{
		ParentID:  req.ParentID,
		DeptCode:  req.DeptCode,
		DeptName:  req.DeptName,
		LeaderID:  req.LeaderID,
		SortOrder: req.SortOrder,
		Status:    req.Status,
	}

	dept, err := h.deptService.Update(id, input)
	if err != nil {
		if err == service.ErrDeptCodeExists {
			response.BadRequest(c, "部门编码已存在")
			return
		}
		if err == service.ErrDeptNotFound {
			response.NotFound(c, "部门不存在")
			return
		}
		if err == service.ErrCannotSetSelfParent {
			response.BadRequest(c, "不能将自己设置为父部门")
			return
		}
		response.ServerError(c, "更新失败")
		return
	}

	response.Success(c, dept)
}

// Delete 删除部门
func (h *DeptHandler) Delete(c *gin.Context) {
	id, err := utils.ParseInt64Param(c, "id")
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	if err := h.deptService.Delete(id); err != nil {
		if err == service.ErrDeptHasChildren {
			response.BadRequest(c, "该部门下存在子部门，无法删除")
			return
		}
		response.ServerError(c, "删除失败")
		return
	}

	response.Success(c, nil)
}
