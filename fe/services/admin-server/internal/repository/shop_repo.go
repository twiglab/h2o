package repository

import (
	"shared/models"

	"gorm.io/gorm"
)

// ShopRepository 店铺数据仓库
type ShopRepository struct {
	db *gorm.DB
}

// NewShopRepository 创建店铺仓库
func NewShopRepository(db *gorm.DB) *ShopRepository {
	return &ShopRepository{db: db}
}

// FindByID 根据ID查找店铺
func (r *ShopRepository) FindByID(id int64) (*models.Shop, error) {
	var shop models.Shop
	if err := r.db.Preload("Merchant").First(&shop, id).Error; err != nil {
		return nil, err
	}
	return &shop, nil
}

// FindByShopNo 根据店铺编号查找店铺
func (r *ShopRepository) FindByShopNo(shopNo string) (*models.Shop, error) {
	var shop models.Shop
	if err := r.db.Where("shop_no = ?", shopNo).First(&shop).Error; err != nil {
		return nil, err
	}
	return &shop, nil
}

// ExistsByShopNo 检查店铺编号是否存在
func (r *ShopRepository) ExistsByShopNo(shopNo string) (bool, error) {
	var count int64
	if err := r.db.Model(&models.Shop{}).Where("shop_no = ?", shopNo).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// List 获取店铺列表
func (r *ShopRepository) List(shopNo, shopName string, merchantID *int64, status *int8, offset, limit int) ([]models.Shop, int64, error) {
	var shops []models.Shop
	var total int64

	query := r.db.Model(&models.Shop{})

	if shopNo != "" {
		query = query.Where("shop_no LIKE ?", "%"+shopNo+"%")
	}
	if shopName != "" {
		query = query.Where("shop_name LIKE ?", "%"+shopName+"%")
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

	if err := query.Preload("Merchant").Order("id DESC").Offset(offset).Limit(limit).Find(&shops).Error; err != nil {
		return nil, 0, err
	}

	return shops, total, nil
}

// ListAll 获取所有店铺（下拉选择用）
func (r *ShopRepository) ListAll(merchantID *int64, status *int8) ([]models.Shop, error) {
	var shops []models.Shop
	query := r.db.Model(&models.Shop{})

	if merchantID != nil {
		query = query.Where("merchant_id = ?", *merchantID)
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if err := query.Order("shop_name ASC").Find(&shops).Error; err != nil {
		return nil, err
	}

	return shops, nil
}

// Create 创建店铺
func (r *ShopRepository) Create(shop *models.Shop) error {
	return r.db.Create(shop).Error
}

// Update 更新店铺
func (r *ShopRepository) Update(shop *models.Shop, updates map[string]interface{}) error {
	return r.db.Model(shop).Updates(updates).Error
}

// Delete 删除店铺
func (r *ShopRepository) Delete(id int64) error {
	return r.db.Delete(&models.Shop{}, id).Error
}

// GetStats 获取店铺统计信息
func (r *ShopRepository) GetStats(shopID int64) (electricMeterCount, waterMeterCount int64, err error) {
	// 电表数量
	if err = r.db.Model(&models.ElectricMeter{}).Where("shop_id = ?", shopID).Count(&electricMeterCount).Error; err != nil {
		return
	}

	// 水表数量
	if err = r.db.Model(&models.WaterMeter{}).Where("shop_id = ?", shopID).Count(&waterMeterCount).Error; err != nil {
		return
	}

	return
}
