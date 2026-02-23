package handler

import (
	"time"

	"admin-server/internal/service"
	"shared/response"
	"shared/utils"

	"github.com/gin-gonic/gin"
)

type RateHandler struct {
	rateService *service.RateService
}

func NewRateHandler(rateService *service.RateService) *RateHandler {
	return &RateHandler{rateService: rateService}
}

// List 获取费率列表
func (h *RateHandler) List(c *gin.Context) {
	page, pageSize := utils.GetPage(c)
	keyword := c.Query("keyword")
	scope := utils.ParseInt8Query(c, "scope")
	status := utils.ParseInt8Query(c, "status")
	calcMode := utils.ParseInt8Query(c, "calc_mode")

	rates, total, err := h.rateService.List(keyword, scope, status, calcMode, page, pageSize)
	if err != nil {
		response.ServerError(c, "查询失败")
		return
	}

	response.Page(c, rates, total, page, pageSize)
}

// Get 获取费率详情
func (h *RateHandler) Get(c *gin.Context) {
	id, err := utils.ParseInt64Param(c, "id")
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	rate, err := h.rateService.GetByID(id)
	if err != nil {
		response.NotFound(c, "费率不存在")
		return
	}

	response.Success(c, rate)
}

type CreateRateRequest struct {
	RateCode      string  `json:"rate_code" binding:"required"`
	RateName      string  `json:"rate_name" binding:"required"`
	Scope         int8    `json:"scope"`
	MerchantID    *int64  `json:"merchant_id"`
	CalcMode      int8    `json:"calc_mode"`
	UnitPrice     float64 `json:"unit_price"`
	EffectiveDate string  `json:"effective_date"`
	Status        int8    `json:"status"`
	Remark        *string `json:"remark"`
	TOUDetails    []struct {
		PeriodName string  `json:"period_name"`
		StartTime  string  `json:"start_time"`
		EndTime    string  `json:"end_time"`
		UnitPrice  float64 `json:"unit_price"`
	} `json:"tou_details"`
	ServiceFees []struct {
		FeeName  string  `json:"fee_name"`
		FeeType  int8    `json:"fee_type"`
		FeeValue float64 `json:"fee_value"`
	} `json:"service_fees"`
}

// Create 创建费率
func (h *RateHandler) Create(c *gin.Context) {
	var req CreateRateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	// 处理生效日期
	var effectiveDate *time.Time
	if req.EffectiveDate != "" {
		t, err := time.Parse("2006-01-02", req.EffectiveDate)
		if err != nil {
			response.BadRequest(c, "日期格式错误")
			return
		}
		effectiveDate = &t
	}

	// 转换 TOU 详情
	var touDetails []service.TOUInput
	for _, tou := range req.TOUDetails {
		touDetails = append(touDetails, service.TOUInput{
			PeriodName: tou.PeriodName,
			StartTime:  tou.StartTime,
			EndTime:    tou.EndTime,
			UnitPrice:  tou.UnitPrice,
		})
	}

	// 转换服务费
	var serviceFees []service.ServiceFeeInput
	for _, sf := range req.ServiceFees {
		serviceFees = append(serviceFees, service.ServiceFeeInput{
			FeeName:  sf.FeeName,
			FeeType:  sf.FeeType,
			FeeValue: sf.FeeValue,
		})
	}

	input := &service.CreateRateInput{
		RateCode:      req.RateCode,
		RateName:      req.RateName,
		Scope:         req.Scope,
		MerchantID:    req.MerchantID,
		CalcMode:      req.CalcMode,
		UnitPrice:     req.UnitPrice,
		EffectiveDate: effectiveDate,
		Status:        req.Status,
		Remark:        req.Remark,
		TOUDetails:    touDetails,
		ServiceFees:   serviceFees,
	}

	rate, err := h.rateService.Create(input)
	if err != nil {
		if err == service.ErrRateCodeExists {
			response.Error(c, 400, "费率编码已存在")
			return
		}
		response.ServerError(c, "创建失败")
		return
	}

	response.Success(c, rate)
}

type UpdateRateRequest struct {
	RateName  *string  `json:"rate_name"`
	UnitPrice *float64 `json:"unit_price"`
	Status    *int8    `json:"status"`
	Remark    *string  `json:"remark"`
}

// Update 更新费率
func (h *RateHandler) Update(c *gin.Context) {
	id, err := utils.ParseInt64Param(c, "id")
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	var req UpdateRateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	input := &service.UpdateRateInput{
		RateName:  req.RateName,
		UnitPrice: req.UnitPrice,
		Status:    req.Status,
		Remark:    req.Remark,
	}

	rate, err := h.rateService.Update(id, input)
	if err != nil {
		if err == service.ErrRateNotFound {
			response.NotFound(c, "费率不存在")
			return
		}
		response.ServerError(c, "更新失败")
		return
	}

	response.Success(c, rate)
}

// Delete 删除费率
func (h *RateHandler) Delete(c *gin.Context) {
	id, err := utils.ParseInt64Param(c, "id")
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	if err := h.rateService.Delete(id); err != nil {
		if err == service.ErrRateInUse {
			response.Error(c, 400, "该费率已被电表使用，无法删除")
			return
		}
		response.ServerError(c, "删除失败")
		return
	}

	response.SuccessWithMessage(c, "删除成功", nil)
}
