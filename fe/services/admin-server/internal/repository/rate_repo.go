package repository

import (
	"shared/models"

	"gorm.io/gorm"
)

// RateRepository 费率数据仓库
type RateRepository struct {
	db *gorm.DB
}

// NewRateRepository 创建费率仓库
func NewRateRepository(db *gorm.DB) *RateRepository {
	return &RateRepository{db: db}
}

// FindByID 根据ID查找费率
func (r *RateRepository) FindByID(id int64) (*models.ElectricRate, error) {
	var rate models.ElectricRate
	if err := r.db.Preload("Merchant").Preload("TOUDetails").Preload("ServiceFees").First(&rate, id).Error; err != nil {
		return nil, err
	}
	return &rate, nil
}

// FindByCode 根据编码查找费率
func (r *RateRepository) FindByCode(code string) (*models.ElectricRate, error) {
	var rate models.ElectricRate
	if err := r.db.Where("rate_code = ?", code).First(&rate).Error; err != nil {
		return nil, err
	}
	return &rate, nil
}

// ExistsByCode 检查费率编码是否存在
func (r *RateRepository) ExistsByCode(code string) (bool, error) {
	var count int64
	if err := r.db.Model(&models.ElectricRate{}).Where("rate_code = ?", code).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// List 获取费率列表
func (r *RateRepository) List(keyword string, scope, status, calcMode *int8, offset, limit int) ([]models.ElectricRate, int64, error) {
	var rates []models.ElectricRate
	var total int64

	query := r.db.Model(&models.ElectricRate{})

	if keyword != "" {
		query = query.Where("rate_name LIKE ? OR rate_code LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if scope != nil {
		query = query.Where("scope = ?", *scope)
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	if calcMode != nil {
		query = query.Where("calc_mode = ?", *calcMode)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("Merchant").Order("id DESC").Offset(offset).Limit(limit).Find(&rates).Error; err != nil {
		return nil, 0, err
	}

	// 加载关联数据
	for i := range rates {
		r.db.Where("rate_id = ?", rates[i].ID).Find(&rates[i].TOUDetails)
		r.db.Where("rate_id = ?", rates[i].ID).Find(&rates[i].ServiceFees)
	}

	return rates, total, nil
}

// Create 创建费率
func (r *RateRepository) Create(rate *models.ElectricRate) error {
	return r.db.Create(rate).Error
}

// CreateTOU 创建分时电价
func (r *RateRepository) CreateTOU(tou *models.ElectricRateTOU) error {
	return r.db.Create(tou).Error
}

// CreateServiceFee 创建服务费
func (r *RateRepository) CreateServiceFee(fee *models.ElectricRateServiceFee) error {
	return r.db.Create(fee).Error
}

// Update 更新费率
func (r *RateRepository) Update(rate *models.ElectricRate, updates map[string]interface{}) error {
	return r.db.Model(rate).Updates(updates).Error
}

// Delete 删除费率
func (r *RateRepository) Delete(id int64) error {
	return r.db.Delete(&models.ElectricRate{}, id).Error
}

// DeleteTOUByRateID 删除费率的分时电价
func (r *RateRepository) DeleteTOUByRateID(rateID int64) error {
	return r.db.Where("rate_id = ?", rateID).Delete(&models.ElectricRateTOU{}).Error
}

// DeleteServiceFeeByRateID 删除费率的服务费
func (r *RateRepository) DeleteServiceFeeByRateID(rateID int64) error {
	return r.db.Where("rate_id = ?", rateID).Delete(&models.ElectricRateServiceFee{}).Error
}

// DB 返回数据库连接（用于事务）
func (r *RateRepository) DB() *gorm.DB {
	return r.db
}
