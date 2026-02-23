package handler

import (
	"admin-server/internal/service"
	"shared/response"
	"shared/utils"

	"github.com/gin-gonic/gin"
)

type ShopHandler struct {
	shopService *service.ShopService
}

func NewShopHandler(shopService *service.ShopService) *ShopHandler {
	return &ShopHandler{shopService: shopService}
}

// List 获取店铺列表
func (h *ShopHandler) List(c *gin.Context) {
	page, pageSize := utils.GetPage(c)
	shopNo := c.Query("shop_no")
	shopName := c.Query("shop_name")
	merchantID := utils.ParseInt64Query(c, "merchant_id")
	status := utils.ParseInt8Query(c, "status")

	shops, total, err := h.shopService.List(shopNo, shopName, merchantID, status, page, pageSize)
	if err != nil {
		response.ServerError(c, "查询失败")
		return
	}

	response.Page(c, shops, total, page, pageSize)
}

// ListAll 获取所有店铺（下拉选择用）
func (h *ShopHandler) ListAll(c *gin.Context) {
	merchantID := utils.ParseInt64Query(c, "merchant_id")

	shops, err := h.shopService.ListAll(merchantID)
	if err != nil {
		response.ServerError(c, "查询失败")
		return
	}

	response.Success(c, shops)
}

// Get 获取店铺详情
func (h *ShopHandler) Get(c *gin.Context) {
	id, err := utils.ParseInt64Param(c, "id")
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	shop, err := h.shopService.GetByID(id)
	if err != nil {
		response.NotFound(c, "店铺不存在")
		return
	}

	response.Success(c, shop)
}

type CreateShopRequest struct {
	ShopNo       string  `json:"shop_no"`
	ShopName     string  `json:"shop_name" binding:"required"`
	MerchantID   int64   `json:"merchant_id" binding:"required"`
	Building     *string `json:"building"`
	Floor        *string `json:"floor"`
	RoomNo       *string `json:"room_no"`
	ContactName  *string `json:"contact_name"`
	ContactPhone *string `json:"contact_phone"`
	Status       int8    `json:"status"`
	Remark       *string `json:"remark"`
}

// Create 创建店铺
func (h *ShopHandler) Create(c *gin.Context) {
	var req CreateShopRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误：店铺名称和商户不能为空")
		return
	}

	input := &service.CreateShopInput{
		ShopNo:       req.ShopNo,
		ShopName:     req.ShopName,
		MerchantID:   req.MerchantID,
		Building:     req.Building,
		Floor:        req.Floor,
		RoomNo:       req.RoomNo,
		ContactName:  req.ContactName,
		ContactPhone: req.ContactPhone,
		Status:       req.Status,
		Remark:       req.Remark,
	}

	shop, err := h.shopService.Create(input)
	if err != nil {
		if err == service.ErrShopNoExists {
			response.BadRequest(c, "店铺编号已存在")
			return
		}
		response.ServerError(c, "创建失败")
		return
	}

	response.Success(c, shop)
}

type UpdateShopRequest struct {
	ShopName     *string `json:"shop_name"`
	MerchantID   *int64  `json:"merchant_id"`
	Building     *string `json:"building"`
	Floor        *string `json:"floor"`
	RoomNo       *string `json:"room_no"`
	ContactName  *string `json:"contact_name"`
	ContactPhone *string `json:"contact_phone"`
	Status       *int8   `json:"status"`
	Remark       *string `json:"remark"`
}

// Update 更新店铺
func (h *ShopHandler) Update(c *gin.Context) {
	id, err := utils.ParseInt64Param(c, "id")
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	var req UpdateShopRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	input := &service.UpdateShopInput{
		ShopName:     req.ShopName,
		MerchantID:   req.MerchantID,
		Building:     req.Building,
		Floor:        req.Floor,
		RoomNo:       req.RoomNo,
		ContactName:  req.ContactName,
		ContactPhone: req.ContactPhone,
		Status:       req.Status,
		Remark:       req.Remark,
	}

	shop, err := h.shopService.Update(id, input)
	if err != nil {
		response.ServerError(c, "更新失败")
		return
	}

	response.Success(c, shop)
}

// Delete 删除店铺
func (h *ShopHandler) Delete(c *gin.Context) {
	id, err := utils.ParseInt64Param(c, "id")
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	if err := h.shopService.Delete(id); err != nil {
		response.ServerError(c, "删除失败")
		return
	}

	response.Success(c, nil)
}

// GetStats 获取店铺统计
func (h *ShopHandler) GetStats(c *gin.Context) {
	id, err := utils.ParseInt64Param(c, "id")
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	stats, err := h.shopService.GetStats(id)
	if err != nil {
		response.ServerError(c, "查询失败")
		return
	}

	response.Success(c, stats)
}
