package handler

import (
	"admin-server/internal/service"
	"shared/middleware"
	"shared/response"
	"shared/utils"

	"github.com/gin-gonic/gin"
)

type AccountHandler struct {
	accountService *service.AccountService
}

func NewAccountHandler(accountService *service.AccountService) *AccountHandler {
	return &AccountHandler{accountService: accountService}
}

// List 获取账户列表
func (h *AccountHandler) List(c *gin.Context) {
	page, pageSize := utils.GetPage(c)

	accountNo := c.Query("account_no")
	accountName := c.Query("account_name")
	merchantID := utils.ParseInt64Query(c, "merchant_id")
	status := utils.ParseInt8Query(c, "status")

	accounts, total, err := h.accountService.List(accountNo, accountName, merchantID, status, page, pageSize)
	if err != nil {
		response.ServerError(c, "查询失败")
		return
	}

	response.Page(c, accounts, total, page, pageSize)
}

// ListAll 获取所有账户（下拉选择用）
func (h *AccountHandler) ListAll(c *gin.Context) {
	merchantID := utils.ParseInt64Query(c, "merchant_id")

	accounts, err := h.accountService.ListAll(merchantID)
	if err != nil {
		response.ServerError(c, "查询失败")
		return
	}

	response.Success(c, accounts)
}

// Get 获取账户详情
func (h *AccountHandler) Get(c *gin.Context) {
	id, err := utils.ParseInt64Param(c, "id")
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	account, err := h.accountService.GetByID(id)
	if err != nil {
		response.NotFound(c, "账户不存在")
		return
	}

	response.Success(c, account)
}

type CreateAccountRequest struct {
	AccountName *string `json:"account_name"`
	MerchantID  int64   `json:"merchant_id" binding:"required"`
	Status      int8    `json:"status"`
	Remark      *string `json:"remark"`
}

// Create 创建账户
func (h *AccountHandler) Create(c *gin.Context) {
	var req CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误：商户ID不能为空")
		return
	}

	input := &service.CreateAccountInput{
		AccountName: req.AccountName,
		MerchantID:  req.MerchantID,
		Status:      req.Status,
		Remark:      req.Remark,
	}

	account, err := h.accountService.Create(input)
	if err != nil {
		response.ServerError(c, "创建失败")
		return
	}

	response.Success(c, account)
}

type UpdateAccountRequest struct {
	AccountName *string `json:"account_name"`
	Status      *int8   `json:"status"`
	Remark      *string `json:"remark"`
}

// Update 更新账户
func (h *AccountHandler) Update(c *gin.Context) {
	id, err := utils.ParseInt64Param(c, "id")
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	var req UpdateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	input := &service.UpdateAccountInput{
		AccountName: req.AccountName,
		Status:      req.Status,
		Remark:      req.Remark,
	}

	account, err := h.accountService.Update(id, input)
	if err != nil {
		response.ServerError(c, "更新失败")
		return
	}

	response.Success(c, account)
}

type RechargeRequest struct {
	Amount        float64 `json:"amount" binding:"required,gt=0"`
	PaymentMethod int8    `json:"payment_method"`
	Remark        *string `json:"remark"`
}

// Recharge 账户充值
func (h *AccountHandler) Recharge(c *gin.Context) {
	id, err := utils.ParseInt64Param(c, "id")
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	var req RechargeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	userID := middleware.GetUserID(c)
	realName := middleware.GetRealName(c)
	if realName == "" {
		realName = middleware.GetUsername(c)
	}

	input := &service.RechargeInput{
		AccountID:     id,
		Amount:        req.Amount,
		PaymentMethod: req.PaymentMethod,
		OperatorID:    userID,
		OperatorName:  realName,
		Remark:        req.Remark,
	}

	result, err := h.accountService.Recharge(input)
	if err != nil {
		if err == service.ErrAccountNotFound {
			response.NotFound(c, "账户不存在")
			return
		}
		response.ServerError(c, "充值失败")
		return
	}

	response.Success(c, map[string]interface{}{
		"recharge_no":    result.RechargeNo,
		"amount":         result.Amount,
		"balance_before": result.BalanceBefore,
		"balance_after":  result.BalanceAfter,
	})
}

// GetRecharges 获取账户充值记录
func (h *AccountHandler) GetRecharges(c *gin.Context) {
	id, err := utils.ParseInt64Param(c, "id")
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	page, pageSize := utils.GetPage(c)
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	recharges, total, err := h.accountService.GetRecharges(id, startDate, endDate, page, pageSize)
	if err != nil {
		response.ServerError(c, "查询失败")
		return
	}

	response.Page(c, recharges, total, page, pageSize)
}

// ListRecharges 获取所有充值记录
func (h *AccountHandler) ListRecharges(c *gin.Context) {
	page, pageSize := utils.GetPage(c)

	keyword := c.Query("keyword")
	accountID := utils.ParseInt64Query(c, "account_id")
	merchantID := utils.ParseInt64Query(c, "merchant_id")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	recharges, total, err := h.accountService.ListRecharges(accountID, merchantID, keyword, startDate, endDate, page, pageSize)
	if err != nil {
		response.ServerError(c, "查询失败")
		return
	}

	response.Page(c, recharges, total, page, pageSize)
}

// GetElectricDeductions 获取账户电费扣费记录
func (h *AccountHandler) GetElectricDeductions(c *gin.Context) {
	id, err := utils.ParseInt64Param(c, "id")
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	page, pageSize := utils.GetPage(c)
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	deductions, total, err := h.accountService.GetElectricDeductions(id, startDate, endDate, page, pageSize)
	if err != nil {
		response.ServerError(c, "查询失败")
		return
	}

	response.Page(c, deductions, total, page, pageSize)
}

// GetWaterDeductions 获取账户水费扣费记录
func (h *AccountHandler) GetWaterDeductions(c *gin.Context) {
	id, err := utils.ParseInt64Param(c, "id")
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	page, pageSize := utils.GetPage(c)
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	deductions, total, err := h.accountService.GetWaterDeductions(id, startDate, endDate, page, pageSize)
	if err != nil {
		response.ServerError(c, "查询失败")
		return
	}

	response.Page(c, deductions, total, page, pageSize)
}

// ListElectricDeductions 获取所有电费扣费记录
func (h *AccountHandler) ListElectricDeductions(c *gin.Context) {
	page, pageSize := utils.GetPage(c)

	keyword := c.Query("keyword")
	accountID := utils.ParseInt64Query(c, "account_id")
	meterID := utils.ParseInt64Query(c, "meter_id")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	deductions, total, err := h.accountService.ListElectricDeductions(accountID, meterID, keyword, startDate, endDate, page, pageSize)
	if err != nil {
		response.ServerError(c, "查询失败")
		return
	}

	response.Page(c, deductions, total, page, pageSize)
}

// ListWaterDeductions 获取所有水费扣费记录
func (h *AccountHandler) ListWaterDeductions(c *gin.Context) {
	page, pageSize := utils.GetPage(c)

	keyword := c.Query("keyword")
	accountID := utils.ParseInt64Query(c, "account_id")
	meterID := utils.ParseInt64Query(c, "meter_id")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	deductions, total, err := h.accountService.ListWaterDeductions(accountID, meterID, keyword, startDate, endDate, page, pageSize)
	if err != nil {
		response.ServerError(c, "查询失败")
		return
	}

	response.Page(c, deductions, total, page, pageSize)
}

// GetArrears 获取欠费账户列表
func (h *AccountHandler) GetArrears(c *gin.Context) {
	page, pageSize := utils.GetPage(c)

	accounts, total, err := h.accountService.ListArrears(page, pageSize)
	if err != nil {
		response.ServerError(c, "查询失败")
		return
	}

	response.Page(c, accounts, total, page, pageSize)
}
