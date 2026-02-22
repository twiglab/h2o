package handler

import (
	"shared/response"
	"shared/utils"

	"github.com/gin-gonic/gin"
)

type SystemHandler struct{}

func NewSystemHandler() *SystemHandler {
	return &SystemHandler{}
}

// GetRoles 获取角色列表
func (h *SystemHandler) GetRoles(c *gin.Context) {
	// 暂时返回空列表，后续可以扩展
	response.Success(c, []interface{}{})
}

// GetOperationLogs 获取操作日志
func (h *SystemHandler) GetOperationLogs(c *gin.Context) {
	page, pageSize := utils.GetPage(c)

	// 暂时返回空数据，后续可以扩展
	response.Page(c, []interface{}{}, 0, page, pageSize)
}
