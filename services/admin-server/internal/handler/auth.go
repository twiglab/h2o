package handler

import (
	"admin-server/internal/service"
	"shared/middleware"
	"shared/response"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login 用户登录
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	result, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		if err == service.ErrInvalidCredentials {
			response.Error(c, 401, "用户名或密码错误")
			return
		}
		response.ServerError(c, "登录失败")
		return
	}

	response.Success(c, map[string]interface{}{
		"access_token":  result.AccessToken,
		"refresh_token": result.RefreshToken,
		"user": map[string]interface{}{
			"id":        result.User.ID,
			"username":  result.User.Username,
			"real_name": result.User.RealName,
		},
	})
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// Refresh 刷新令牌
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	accessToken, refreshToken, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		response.Unauthorized(c, "令牌无效")
		return
	}

	response.Success(c, map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// Profile 获取当前用户信息
func (h *AuthHandler) Profile(c *gin.Context) {
	userID := middleware.GetUserID(c)

	user, err := h.authService.GetProfile(userID)
	if err != nil {
		response.NotFound(c, "用户不存在")
		return
	}

	response.Success(c, map[string]interface{}{
		"id":        user.ID,
		"username":  user.Username,
		"real_name": user.RealName,
		"phone":     user.Phone,
		"email":     user.Email,
		"status":    user.Status,
	})
}

// Logout 用户登出
func (h *AuthHandler) Logout(c *gin.Context) {
	response.SuccessWithMessage(c, "登出成功", nil)
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// ChangePassword 修改密码
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	userID := middleware.GetUserID(c)

	err := h.authService.ChangePassword(userID, req.OldPassword, req.NewPassword)
	if err != nil {
		if err == service.ErrOldPasswordWrong {
			response.Error(c, 400, "原密码错误")
			return
		}
		if err == service.ErrUserNotFound {
			response.NotFound(c, "用户不存在")
			return
		}
		response.ServerError(c, "修改密码失败")
		return
	}

	response.SuccessWithMessage(c, "密码修改成功", nil)
}
