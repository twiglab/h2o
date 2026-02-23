package service

import (
	"errors"
	"fmt"
	"time"

	"admin-server/internal/repository"
	"shared/models"
	"shared/utils"

	"github.com/shopspring/decimal"
)

var (
	ErrAccountNotFound = errors.New("账户不存在")
)

// AccountInfo 账户信息响应
type AccountInfo struct {
	ID               int64           `json:"id"`
	AccountNo        string          `json:"account_no"`
	AccountName      *string         `json:"account_name"`
	MerchantID       int64           `json:"merchant_id"`
	MerchantName     string          `json:"merchant_name,omitempty"`
	Balance          decimal.Decimal `json:"balance"`
	TotalRecharge    decimal.Decimal `json:"total_recharge"`
	TotalConsumption decimal.Decimal `json:"total_consumption"`
	Status           int8            `json:"status"`
	Remark           *string         `json:"remark"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
}

// AccountService 账户服务
type AccountService struct {
	accountRepo *repository.AccountRepository
	meterRepo   *repository.MeterRepository
}

// NewAccountService 创建账户服务
func NewAccountService(accountRepo *repository.AccountRepository, meterRepo *repository.MeterRepository) *AccountService {
	return &AccountService{
		accountRepo: accountRepo,
		meterRepo:   meterRepo,
	}
}

// accountToInfo 转换账户模型到信息响应
func (s *AccountService) accountToInfo(a *models.Account) AccountInfo {
	info := AccountInfo{
		ID:               a.ID,
		AccountNo:        a.AccountNo,
		AccountName:      a.AccountName,
		MerchantID:       a.MerchantID,
		Balance:          a.Balance,
		TotalRecharge:    a.TotalRecharge,
		TotalConsumption: a.TotalConsumption,
		Status:           a.Status,
		Remark:           a.Remark,
		CreatedAt:        a.CreatedAt,
		UpdatedAt:        a.UpdatedAt,
	}
	if a.Merchant != nil {
		info.MerchantName = a.Merchant.MerchantName
	}
	return info
}

// List 获取账户列表
func (s *AccountService) List(accountNo, accountName string, merchantID *int64, status *int8, page, pageSize int) ([]AccountInfo, int64, error) {
	offset := (page - 1) * pageSize
	accounts, total, err := s.accountRepo.ListAccounts(accountNo, accountName, merchantID, status, offset, pageSize)
	if err != nil {
		return nil, 0, err
	}

	var result []AccountInfo
	for _, a := range accounts {
		result = append(result, s.accountToInfo(&a))
	}

	return result, total, nil
}

// ListAll 获取所有账户（下拉选择用）
func (s *AccountService) ListAll(merchantID *int64) ([]AccountInfo, error) {
	accounts, err := s.accountRepo.ListAllAccounts(merchantID)
	if err != nil {
		return nil, err
	}

	var result []AccountInfo
	for _, a := range accounts {
		result = append(result, s.accountToInfo(&a))
	}

	return result, nil
}

// GetByID 获取账户详情
func (s *AccountService) GetByID(id int64) (*AccountInfo, error) {
	account, err := s.accountRepo.FindByIDWithMerchant(id)
	if err != nil {
		return nil, err
	}

	info := s.accountToInfo(account)
	return &info, nil
}

// CreateAccountInput 创建账户输入
type CreateAccountInput struct {
	AccountName *string
	MerchantID  int64
	Status      int8
	Remark      *string
}

// Create 创建账户
func (s *AccountService) Create(input *CreateAccountInput) (*AccountInfo, error) {
	// 自动生成账户编号
	accountNo := fmt.Sprintf("A%s", time.Now().Format("20060102150405"))

	status := input.Status
	if status == 0 {
		status = models.AccountStatusNormal
	}

	account := &models.Account{
		AccountNo:        accountNo,
		AccountName:      input.AccountName,
		MerchantID:       input.MerchantID,
		Balance:          decimal.Zero,
		TotalRecharge:    decimal.Zero,
		TotalConsumption: decimal.Zero,
		Status:           status,
		Remark:           input.Remark,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if err := s.accountRepo.Create(account); err != nil {
		return nil, err
	}

	info := s.accountToInfo(account)
	return &info, nil
}

// UpdateAccountInput 更新账户输入
type UpdateAccountInput struct {
	AccountName *string
	Status      *int8
	Remark      *string
}

// Update 更新账户
func (s *AccountService) Update(id int64, input *UpdateAccountInput) (*AccountInfo, error) {
	account, err := s.accountRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	updates := map[string]interface{}{
		"updated_at": time.Now(),
	}

	if input.AccountName != nil {
		updates["account_name"] = *input.AccountName
	}
	if input.Status != nil {
		updates["status"] = *input.Status
	}
	if input.Remark != nil {
		updates["remark"] = *input.Remark
	}

	if err := s.accountRepo.Update(nil, account, updates); err != nil {
		return nil, err
	}

	return s.GetByID(id)
}

// ListArrears 获取欠费账户列表
func (s *AccountService) ListArrears(page, pageSize int) ([]AccountInfo, int64, error) {
	offset := (page - 1) * pageSize
	accounts, total, err := s.accountRepo.ListArrears(offset, pageSize)
	if err != nil {
		return nil, 0, err
	}

	var result []AccountInfo
	for _, a := range accounts {
		result = append(result, s.accountToInfo(&a))
	}

	return result, total, nil
}

// RechargeInput 充值输入
type RechargeInput struct {
	AccountID     int64
	Amount        float64
	PaymentMethod int8
	OperatorID    int64
	OperatorName  string
	Remark        *string
}

// RechargeResult 充值结果
type RechargeResult struct {
	RechargeNo    string
	Amount        decimal.Decimal
	BalanceBefore decimal.Decimal
	BalanceAfter  decimal.Decimal
}

// Recharge 账户充值
func (s *AccountService) Recharge(input *RechargeInput) (*RechargeResult, error) {
	amount := decimal.NewFromFloat(input.Amount)

	tx := s.accountRepo.DB().Begin()

	// 锁定账户
	account, err := s.accountRepo.FindByIDForUpdate(tx, input.AccountID)
	if err != nil {
		tx.Rollback()
		return nil, ErrAccountNotFound
	}

	balanceBefore := account.Balance
	balanceAfter := account.Balance.Add(amount)

	// 更新账户状态 - 如果余额从负变正，恢复正常状态
	status := account.Status
	if balanceAfter.GreaterThanOrEqual(decimal.Zero) && account.Status == models.AccountStatusArrears {
		status = models.AccountStatusNormal
	}

	if err := s.accountRepo.Update(tx, account, map[string]interface{}{
		"balance":        balanceAfter,
		"total_recharge": account.TotalRecharge.Add(amount),
		"status":         status,
		"updated_at":     time.Now(),
	}); err != nil {
		tx.Rollback()
		return nil, err
	}

	// 创建充值记录
	paymentMethod := input.PaymentMethod
	if paymentMethod == 0 {
		paymentMethod = 1
	}

	operatorName := input.OperatorName
	recharge := &models.Recharge{
		RechargeNo:    utils.GenerateNo("R"),
		AccountID:     account.ID,
		Amount:        amount,
		BalanceBefore: balanceBefore,
		BalanceAfter:  balanceAfter,
		PaymentMethod: paymentMethod,
		OperatorID:    input.OperatorID,
		OperatorName:  &operatorName,
		Remark:        input.Remark,
		CreatedAt:     time.Now(),
	}

	if err := s.accountRepo.CreateRecharge(tx, recharge); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &RechargeResult{
		RechargeNo:    recharge.RechargeNo,
		Amount:        amount,
		BalanceBefore: balanceBefore,
		BalanceAfter:  balanceAfter,
	}, nil
}

// GetRecharges 获取账户充值记录
func (s *AccountService) GetRecharges(accountID int64, startDate, endDate string, page, pageSize int) ([]models.Recharge, int64, error) {
	offset := (page - 1) * pageSize
	return s.accountRepo.ListRecharges(&accountID, nil, "", startDate, endDate, offset, pageSize)
}

// ListRecharges 获取所有充值记录
func (s *AccountService) ListRecharges(accountID, merchantID *int64, keyword, startDate, endDate string, page, pageSize int) ([]models.Recharge, int64, error) {
	offset := (page - 1) * pageSize
	return s.accountRepo.ListRecharges(accountID, merchantID, keyword, startDate, endDate, offset, pageSize)
}

// GetElectricDeductions 获取账户电费扣费记录
func (s *AccountService) GetElectricDeductions(accountID int64, startDate, endDate string, page, pageSize int) ([]models.ElectricDeduction, int64, error) {
	offset := (page - 1) * pageSize
	return s.accountRepo.ListElectricDeductions(&accountID, nil, "", startDate, endDate, offset, pageSize)
}

// GetWaterDeductions 获取账户水费扣费记录
func (s *AccountService) GetWaterDeductions(accountID int64, startDate, endDate string, page, pageSize int) ([]models.WaterDeduction, int64, error) {
	offset := (page - 1) * pageSize
	return s.accountRepo.ListWaterDeductions(&accountID, nil, "", startDate, endDate, offset, pageSize)
}

// ListElectricDeductions 获取所有电费扣费记录
func (s *AccountService) ListElectricDeductions(accountID, meterID *int64, keyword, startDate, endDate string, page, pageSize int) ([]models.ElectricDeduction, int64, error) {
	offset := (page - 1) * pageSize
	return s.accountRepo.ListElectricDeductions(accountID, meterID, keyword, startDate, endDate, offset, pageSize)
}

// ListWaterDeductions 获取所有水费扣费记录
func (s *AccountService) ListWaterDeductions(accountID, meterID *int64, keyword, startDate, endDate string, page, pageSize int) ([]models.WaterDeduction, int64, error) {
	offset := (page - 1) * pageSize
	return s.accountRepo.ListWaterDeductions(accountID, meterID, keyword, startDate, endDate, offset, pageSize)
}
