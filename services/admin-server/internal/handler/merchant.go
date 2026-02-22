package handler

import (
	"admin-server/internal/service"
	"shared/response"
	"shared/utils"

	"github.com/gin-gonic/gin"
)

type MerchantHandler struct {
	merchantService *service.MerchantService
}

func NewMerchantHandler(merchantService *service.MerchantService) *MerchantHandler {
	return &MerchantHandler{merchantService: merchantService}
}

// List 获取商户列表
func (h *MerchantHandler) List(c *gin.Context) {
	page, pageSize := utils.GetPage(c)
	merchantNo := c.Query("merchant_no")
	merchantName := c.Query("merchant_name")
	status := utils.ParseInt8Query(c, "status")

	merchants, total, err := h.merchantService.List(merchantNo, merchantName, status, page, pageSize)
	if err != nil {
		response.ServerError(c, "查询失败")
		return
	}

	response.Page(c, merchants, total, page, pageSize)
}

// ListAll 获取所有商户（下拉选择用）
func (h *MerchantHandler) ListAll(c *gin.Context) {
	merchants, err := h.merchantService.ListAll()
	if err != nil {
		response.ServerError(c, "查询失败")
		return
	}

	response.Success(c, merchants)
}

// Get 获取商户详情
func (h *MerchantHandler) Get(c *gin.Context) {
	id, err := utils.ParseInt64Param(c, "id")
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	merchant, err := h.merchantService.GetByID(id)
	if err != nil {
		response.NotFound(c, "商户不存在")
		return
	}

	response.Success(c, merchant)
}

type CreateMerchantRequest struct {
	MerchantName string  `json:"merchant_name" binding:"required"`
	MerchantType int8    `json:"merchant_type"`
	ContactName  *string `json:"contact_name"`
	ContactPhone *string `json:"contact_phone"`
	Status       int8    `json:"status"`
	Remark       *string `json:"remark"`
}

// Create 创建商户
func (h *MerchantHandler) Create(c *gin.Context) {
	var req CreateMerchantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误：商户名称不能为空")
		return
	}

	input := &service.CreateMerchantInput{
		MerchantName: req.MerchantName,
		MerchantType: req.MerchantType,
		ContactName:  req.ContactName,
		ContactPhone: req.ContactPhone,
		Status:       req.Status,
		Remark:       req.Remark,
	}

	merchant, err := h.merchantService.Create(input)
	if err != nil {
		if err == service.ErrMerchantNoExists {
			response.BadRequest(c, "商户编号已存在")
			return
		}
		response.ServerError(c, "创建失败")
		return
	}

	response.Success(c, merchant)
}

type UpdateMerchantRequest struct {
	MerchantName *string `json:"merchant_name"`
	MerchantType *int8   `json:"merchant_type"`
	ContactName  *string `json:"contact_name"`
	ContactPhone *string `json:"contact_phone"`
	Status       *int8   `json:"status"`
	Remark       *string `json:"remark"`
}

// Update 更新商户
func (h *MerchantHandler) Update(c *gin.Context) {
	id, err := utils.ParseInt64Param(c, "id")
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	var req UpdateMerchantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	input := &service.UpdateMerchantInput{
		MerchantName: req.MerchantName,
		MerchantType: req.MerchantType,
		ContactName:  req.ContactName,
		ContactPhone: req.ContactPhone,
		Status:       req.Status,
		Remark:       req.Remark,
	}

	merchant, err := h.merchantService.Update(id, input)
	if err != nil {
		response.ServerError(c, "更新失败")
		return
	}

	response.Success(c, merchant)
}

// Delete 删除商户
func (h *MerchantHandler) Delete(c *gin.Context) {
	id, err := utils.ParseInt64Param(c, "id")
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	if err := h.merchantService.Delete(id); err != nil {
		response.ServerError(c, "删除失败")
		return
	}

	response.Success(c, nil)
}

// GetStats 获取商户统计
func (h *MerchantHandler) GetStats(c *gin.Context) {
	id, err := utils.ParseInt64Param(c, "id")
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	stats, err := h.merchantService.GetStats(id)
	if err != nil {
		response.ServerError(c, "查询失败")
		return
	}

	response.Success(c, stats)
}
