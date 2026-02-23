package repository

import (
	"shared/models"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// AccountRepository 账户数据仓库
type AccountRepository struct {
	db *gorm.DB
}

// NewAccountRepository 创建账户仓库
func NewAccountRepository(db *gorm.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

// ========== 统一账户操作 ==========

// FindByID 根据ID查找账户
func (r *AccountRepository) FindByID(id int64) (*models.Account, error) {
	var account models.Account
	if err := r.db.First(&account, id).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

// FindByIDWithMerchant 根据ID查找账户（包含商户信息）
func (r *AccountRepository) FindByIDWithMerchant(id int64) (*models.Account, error) {
	var account models.Account
	if err := r.db.Preload("Merchant").First(&account, id).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

// ListAccounts 获取账户列表（支持按账户编号、名称、商户ID、状态筛选）
func (r *AccountRepository) ListAccounts(accountNo, accountName string, merchantID *int64, status *int8, offset, limit int) ([]models.Account, int64, error) {
	var accounts []models.Account
	var total int64

	query := r.db.Model(&models.Account{})

	if accountNo != "" {
		query = query.Where("account_no LIKE ?", "%"+accountNo+"%")
	}
	if accountName != "" {
		query = query.Where("account_name LIKE ?", "%"+accountName+"%")
	}
	if merchantID != nil {
		query = query.Where("merchant_id = ?", *merchantID)
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("Merchant").Order("id DESC").Offset(offset).Limit(limit).Find(&accounts).Error; err != nil {
		return nil, 0, err
	}

	return accounts, total, nil
}

// ListAllAccounts 获取所有账户（下拉选择用）
func (r *AccountRepository) ListAllAccounts(merchantID *int64) ([]models.Account, error) {
	var accounts []models.Account

	query := r.db.Model(&models.Account{}).Where("status = ?", models.AccountStatusNormal)

	if merchantID != nil {
		query = query.Where("merchant_id = ?", *merchantID)
	}

	if err := query.Order("id DESC").Find(&accounts).Error; err != nil {
		return nil, err
	}

	return accounts, nil
}

// FindByIDForUpdate 根据ID查找并锁定账户
func (r *AccountRepository) FindByIDForUpdate(tx *gorm.DB, id int64) (*models.Account, error) {
	var account models.Account
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&account, id).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

// FindByAccountNo 根据账户编号查找账户
func (r *AccountRepository) FindByAccountNo(accountNo string) (*models.Account, error) {
	var account models.Account
	if err := r.db.Where("account_no = ?", accountNo).First(&account).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

// List 获取账户列表
func (r *AccountRepository) List(merchantID *int64, status *int8, hasArrears, lowBalance bool, keyword string, offset, limit int) ([]models.Account, int64, error) {
	var accounts []models.Account
	var total int64

	query := r.db.Model(&models.Account{})

	if merchantID != nil {
		query = query.Where("merchant_id = ?", *merchantID)
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	if hasArrears {
		// 欠费状态
		query = query.Where("status = ?", models.AccountStatusArrears)
	}
	if lowBalance {
		// 低余额
		query = query.Where("balance < 100")
	}
	if keyword != "" {
		query = query.Where("account_no LIKE ? OR account_name LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("Merchant").Order("id DESC").Offset(offset).Limit(limit).Find(&accounts).Error; err != nil {
		return nil, 0, err
	}

	return accounts, total, nil
}

// ListArrears 获取欠费账户列表
func (r *AccountRepository) ListArrears(offset, limit int) ([]models.Account, int64, error) {
	var accounts []models.Account
	var total int64

	query := r.db.Model(&models.Account{}).Where("status = ?", models.AccountStatusArrears)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("Merchant").Order("balance ASC").Offset(offset).Limit(limit).Find(&accounts).Error; err != nil {
		return nil, 0, err
	}

	return accounts, total, nil
}

// Update 更新账户
func (r *AccountRepository) Update(tx *gorm.DB, account *models.Account, updates map[string]interface{}) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	return db.Model(account).Updates(updates).Error
}

// Create 创建账户
func (r *AccountRepository) Create(account *models.Account) error {
	return r.db.Create(account).Error
}

// ========== 充值记录操作 ==========

// CreateRecharge 创建充值记录
func (r *AccountRepository) CreateRecharge(tx *gorm.DB, recharge *models.Recharge) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	return db.Create(recharge).Error
}

// ListRecharges 获取充值记录
func (r *AccountRepository) ListRecharges(accountID *int64, merchantID *int64, keyword, startDate, endDate string, offset, limit int) ([]models.Recharge, int64, error) {
	var recharges []models.Recharge
	var total int64

	query := r.db.Model(&models.Recharge{})

	if accountID != nil {
		query = query.Where("account_id = ?", *accountID)
	}
	if merchantID != nil {
		// 通过子查询过滤商户
		query = query.Where("account_id IN (SELECT id FROM biz_account WHERE merchant_id = ?)", *merchantID)
	}
	if keyword != "" {
		query = query.Where("recharge_no LIKE ?", "%"+keyword+"%")
	}
	if startDate != "" {
		query = query.Where("created_at >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("created_at <= ?", endDate+" 23:59:59")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("Account").Order("created_at DESC").Offset(offset).Limit(limit).Find(&recharges).Error; err != nil {
		return nil, 0, err
	}

	return recharges, total, nil
}

// ========== 电费扣费记录操作 ==========

// ListElectricDeductions 获取电费扣费记录
func (r *AccountRepository) ListElectricDeductions(accountID *int64, meterID *int64, keyword, startDate, endDate string, offset, limit int) ([]models.ElectricDeduction, int64, error) {
	var deductions []models.ElectricDeduction
	var total int64

	query := r.db.Model(&models.ElectricDeduction{})

	if accountID != nil {
		query = query.Where("account_id = ?", *accountID)
	}
	if meterID != nil {
		query = query.Where("meter_id = ?", *meterID)
	}
	if keyword != "" {
		query = query.Where("deduction_no LIKE ?", "%"+keyword+"%")
	}
	if startDate != "" {
		query = query.Where("deduction_time >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("deduction_time <= ?", endDate+" 23:59:59")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("deduction_time DESC").Offset(offset).Limit(limit).Find(&deductions).Error; err != nil {
		return nil, 0, err
	}

	return deductions, total, nil
}

// ========== 水费扣费记录操作 ==========

// ListWaterDeductions 获取水费扣费记录
func (r *AccountRepository) ListWaterDeductions(accountID *int64, meterID *int64, keyword, startDate, endDate string, offset, limit int) ([]models.WaterDeduction, int64, error) {
	var deductions []models.WaterDeduction
	var total int64

	query := r.db.Model(&models.WaterDeduction{})

	if accountID != nil {
		query = query.Where("account_id = ?", *accountID)
	}
	if meterID != nil {
		query = query.Where("meter_id = ?", *meterID)
	}
	if keyword != "" {
		query = query.Where("deduction_no LIKE ?", "%"+keyword+"%")
	}
	if startDate != "" {
		query = query.Where("deduction_time >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("deduction_time <= ?", endDate+" 23:59:59")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("deduction_time DESC").Offset(offset).Limit(limit).Find(&deductions).Error; err != nil {
		return nil, 0, err
	}

	return deductions, total, nil
}

// ========== 统计方法 ==========

// Count 统计账户数量
func (r *AccountRepository) Count() (int64, error) {
	var count int64
	if err := r.db.Model(&models.Account{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CountArrears 统计欠费账户数量
func (r *AccountRepository) CountArrears() (int64, error) {
	var count int64
	if err := r.db.Model(&models.Account{}).Where("status = ?", models.AccountStatusArrears).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// SumBalance 统计总余额
func (r *AccountRepository) SumBalance() (decimal.Decimal, error) {
	var sum decimal.Decimal
	if err := r.db.Model(&models.Account{}).Select("COALESCE(SUM(balance), 0)").Scan(&sum).Error; err != nil {
		return decimal.Zero, err
	}
	return sum, nil
}

// SumNegativeBalance 统计欠费账户的负余额总和
func (r *AccountRepository) SumNegativeBalance() (decimal.Decimal, error) {
	var sum decimal.Decimal
	if err := r.db.Model(&models.Account{}).Where("balance < 0").Select("COALESCE(SUM(ABS(balance)), 0)").Scan(&sum).Error; err != nil {
		return decimal.Zero, err
	}
	return sum, nil
}

// SumRechargeByDate 统计指定日期充值金额
func (r *AccountRepository) SumRechargeByDate(date string) (decimal.Decimal, error) {
	var sum decimal.Decimal
	if err := r.db.Model(&models.Recharge{}).Where("DATE(created_at) = ?", date).Select("COALESCE(SUM(amount), 0)").Scan(&sum).Error; err != nil {
		return decimal.Zero, err
	}
	return sum, nil
}

// SumDeductionByDate 统计指定日期扣费金额（电费+水费）
func (r *AccountRepository) SumDeductionByDate(date string) (decimal.Decimal, error) {
	var electricSum, waterSum decimal.Decimal

	if err := r.db.Model(&models.ElectricDeduction{}).Where("DATE(deduction_time) = ?", date).Select("COALESCE(SUM(amount), 0)").Scan(&electricSum).Error; err != nil {
		return decimal.Zero, err
	}

	if err := r.db.Model(&models.WaterDeduction{}).Where("DATE(deduction_time) = ?", date).Select("COALESCE(SUM(amount), 0)").Scan(&waterSum).Error; err != nil {
		return decimal.Zero, err
	}

	return electricSum.Add(waterSum), nil
}

// DB 返回数据库连接（用于事务）
func (r *AccountRepository) DB() *gorm.DB {
	return r.db
}
