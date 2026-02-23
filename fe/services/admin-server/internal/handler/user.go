package handler

import (
	"admin-server/internal/service"
	"shared/response"
	"shared/utils"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// List 获取用户列表
func (h *UserHandler) List(c *gin.Context) {
	page, pageSize := utils.GetPage(c)
	keyword := c.Query("keyword")
	status := utils.ParseInt8Query(c, "status")

	users, total, err := h.userService.List(keyword, status, page, pageSize)
	if err != nil {
		response.ServerError(c, "查询失败")
		return
	}

	response.Page(c, users, total, page, pageSize)
}

// Get 获取用户详情
func (h *UserHandler) Get(c *gin.Context) {
	id, err := utils.ParseInt64Param(c, "id")
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	user, err := h.userService.GetByID(id)
	if err != nil {
		response.NotFound(c, "用户不存在")
		return
	}

	response.Success(c, user)
}

type CreateUserRequest struct {
	Username string  `json:"username" binding:"required"`
	Password string  `json:"password" binding:"required,min=6"`
	RealName *string `json:"real_name"`
	Phone    *string `json:"phone"`
	Email    *string `json:"email"`
	DeptID   *int64  `json:"dept_id"`
	UserType int8    `json:"user_type"`
	Status   int8    `json:"status"`
}

// Create 创建用户
func (h *UserHandler) Create(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	input := &service.CreateUserInput{
		Username: req.Username,
		Password: req.Password,
		RealName: req.RealName,
		Phone:    req.Phone,
		Email:    req.Email,
		DeptID:   req.DeptID,
		UserType: req.UserType,
		Status:   req.Status,
	}

	user, err := h.userService.Create(input)
	if err != nil {
		if err == service.ErrUsernameExists {
			response.BadRequest(c, "用户名已存在")
			return
		}
		response.ServerError(c, "创建失败")
		return
	}

	response.Success(c, user)
}

type UpdateUserRequest struct {
	RealName *string `json:"real_name"`
	Phone    *string `json:"phone"`
	Email    *string `json:"email"`
	DeptID   *int64  `json:"dept_id"`
	Status   *int8   `json:"status"`
}

// Update 更新用户
func (h *UserHandler) Update(c *gin.Context) {
	id, err := utils.ParseInt64Param(c, "id")
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	input := &service.UpdateUserInput{
		RealName: req.RealName,
		Phone:    req.Phone,
		Email:    req.Email,
		DeptID:   req.DeptID,
		Status:   req.Status,
	}

	user, err := h.userService.Update(id, input)
	if err != nil {
		response.ServerError(c, "更新失败")
		return
	}

	response.Success(c, user)
}

// Delete 删除用户
func (h *UserHandler) Delete(c *gin.Context) {
	id, err := utils.ParseInt64Param(c, "id")
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	if err := h.userService.Delete(id); err != nil {
		if err == service.ErrCannotDeleteAdmin {
			response.BadRequest(c, "不能删除管理员账号")
			return
		}
		response.ServerError(c, "删除失败")
		return
	}

	response.Success(c, nil)
}

type ResetPasswordRequest struct {
	Password string `json:"password" binding:"required,min=6"`
}

// ResetPassword 重置用户密码
func (h *UserHandler) ResetPassword(c *gin.Context) {
	id, err := utils.ParseInt64Param(c, "id")
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "密码至少6位")
		return
	}

	if err := h.userService.ResetPassword(id, req.Password); err != nil {
		response.ServerError(c, "操作失败")
		return
	}

	response.Success(c, nil)
}
