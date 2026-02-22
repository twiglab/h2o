package repository

import (
	"shared/models"

	"gorm.io/gorm"
)

// MerchantRepository 商户数据仓库
type MerchantRepository struct {
	db *gorm.DB
}

// NewMerchantRepository 创建商户仓库
func NewMerchantRepository(db *gorm.DB) *MerchantRepository {
	return &MerchantRepository{db: db}
}

// FindByID 根据ID查找商户
func (r *MerchantRepository) FindByID(id int64) (*models.Merchant, error) {
	var merchant models.Merchant
	if err := r.db.First(&merchant, id).Error; err != nil {
		return nil, err
	}
	return &merchant, nil
}

// FindByMerchantNo 根据商户编号查找商户
func (r *MerchantRepository) FindByMerchantNo(merchantNo string) (*models.Merchant, error) {
	var merchant models.Merchant
	if err := r.db.Where("merchant_no = ?", merchantNo).First(&merchant).Error; err != nil {
		return nil, err
	}
	return &merchant, nil
}

// ExistsByMerchantNo 检查商户编号是否存在
func (r *MerchantRepository) ExistsByMerchantNo(merchantNo string) (bool, error) {
	var count int64
	if err := r.db.Model(&models.Merchant{}).Where("merchant_no = ?", merchantNo).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// List 获取商户列表
func (r *MerchantRepository) List(merchantNo, merchantName string, status *int8, offset, limit int) ([]models.Merchant, int64, error) {
	var merchants []models.Merchant
	var total int64

	query := r.db.Model(&models.Merchant{})

	if merchantNo != "" {
		query = query.Where("merchant_no LIKE ?", "%"+merchantNo+"%")
	}
	if merchantName != "" {
		query = query.Where("merchant_name LIKE ?", "%"+merchantName+"%")
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("id DESC").Offset(offset).Limit(limit).Find(&merchants).Error; err != nil {
		return nil, 0, err
	}

	return merchants, total, nil
}

// ListAll 获取所有商户（下拉选择用）
func (r *MerchantRepository) ListAll(status *int8) ([]models.Merchant, error) {
	var merchants []models.Merchant
	query := r.db.Model(&models.Merchant{})

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if err := query.Order("merchant_name ASC").Find(&merchants).Error; err != nil {
		return nil, err
	}

	return merchants, nil
}

// Create 创建商户
func (r *MerchantRepository) Create(merchant *models.Merchant) error {
	return r.db.Create(merchant).Error
}

// Update 更新商户
func (r *MerchantRepository) Update(merchant *models.Merchant, updates map[string]interface{}) error {
	return r.db.Model(merchant).Updates(updates).Error
}

// Delete 删除商户
func (r *MerchantRepository) Delete(id int64) error {
	return r.db.Delete(&models.Merchant{}, id).Error
}

// GetStats 获取商户统计信息
func (r *MerchantRepository) GetStats(merchantID int64) (shopCount, electricMeterCount, waterMeterCount int64, err error) {
	// 店铺数量
	if err = r.db.Model(&models.Shop{}).Where("merchant_id = ?", merchantID).Count(&shopCount).Error; err != nil {
		return
	}

	// 电表数量
	if err = r.db.Model(&models.ElectricMeter{}).Where("merchant_id = ?", merchantID).Count(&electricMeterCount).Error; err != nil {
		return
	}

	// 水表数量
	if err = r.db.Model(&models.WaterMeter{}).Where("merchant_id = ?", merchantID).Count(&waterMeterCount).Error; err != nil {
		return
	}

	return
}
