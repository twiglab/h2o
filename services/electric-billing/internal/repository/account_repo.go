package repository

import (
	"shared/models"
	"time"

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

// FindByID 根据ID查找账户
func (r *AccountRepository) FindByID(id int64) (*models.Account, error) {
	var account models.Account
	err := r.db.First(&account, id).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

// FindByIDForUpdate 根据ID查找账户并锁定
func (r *AccountRepository) FindByIDForUpdate(tx *gorm.DB, id int64) (*models.Account, error) {
	var account models.Account
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&account, id).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

// UpdateBalance 更新账户余额和消费金额
func (r *AccountRepository) UpdateBalance(tx *gorm.DB, accountID int64, balance decimal.Decimal, consumption decimal.Decimal, status int8) error {
	return tx.Model(&models.Account{}).
		Where("id = ?", accountID).
		Updates(map[string]interface{}{
			"balance":           balance,
			"total_consumption": consumption,
			"status":            status,
			"updated_at":        time.Now(),
		}).Error
}

// CreateElectricDeduction 创建电费扣费记录
func (r *AccountRepository) CreateElectricDeduction(tx *gorm.DB, deduction *models.ElectricDeduction) error {
	return tx.Create(deduction).Error
}

// FindElectricDeductionByNo 根据扣费流水号查找电费扣费记录
func (r *AccountRepository) FindElectricDeductionByNo(deductionNo string) (*models.ElectricDeduction, error) {
	var deduction models.ElectricDeduction
	err := r.db.Where("deduction_no = ?", deductionNo).First(&deduction).Error
	if err != nil {
		return nil, err
	}
	return &deduction, nil
}
