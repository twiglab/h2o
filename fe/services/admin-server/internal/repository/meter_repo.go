package repository

import (
	"time"

	"shared/models"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// MeterRepository 表计数据仓库
type MeterRepository struct {
	db *gorm.DB
}

// NewMeterRepository 创建表计仓库
func NewMeterRepository(db *gorm.DB) *MeterRepository {
	return &MeterRepository{db: db}
}

// ========== 电表操作 ==========

// FindElectricByID 根据ID查找电表
func (r *MeterRepository) FindElectricByID(id int64) (*models.ElectricMeter, error) {
	var meter models.ElectricMeter
	if err := r.db.Preload("ElectricRate").Preload("Merchant").Preload("Shop").First(&meter, id).Error; err != nil {
		return nil, err
	}
	return &meter, nil
}

// FindElectricByMeterNo 根据表计编号查找电表
func (r *MeterRepository) FindElectricByMeterNo(meterNo string) (*models.ElectricMeter, error) {
	var meter models.ElectricMeter
	if err := r.db.Where("meter_no = ?", meterNo).First(&meter).Error; err != nil {
		return nil, err
	}
	return &meter, nil
}

// ExistsElectricByMeterNo 检查电表编号是否存在
func (r *MeterRepository) ExistsElectricByMeterNo(meterNo string) (bool, error) {
	var count int64
	if err := r.db.Model(&models.ElectricMeter{}).Where("meter_no = ?", meterNo).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// ListElectric 获取电表列表
func (r *MeterRepository) ListElectric(meterNo string, status, onlineStatus *int16, offset, limit int) ([]models.ElectricMeter, int64, error) {
	var meters []models.ElectricMeter
	var total int64

	query := r.db.Model(&models.ElectricMeter{})

	if meterNo != "" {
		query = query.Where("meter_no LIKE ?", "%"+meterNo+"%")
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	if onlineStatus != nil {
		query = query.Where("online_status = ?", *onlineStatus)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("ElectricRate").Preload("Merchant").Preload("Shop").Order("id DESC").Offset(offset).Limit(limit).Find(&meters).Error; err != nil {
		return nil, 0, err
	}

	return meters, total, nil
}

// CreateElectric 创建电表
func (r *MeterRepository) CreateElectric(meter *models.ElectricMeter) error {
	return r.db.Create(meter).Error
}

// UpdateElectric 更新电表
func (r *MeterRepository) UpdateElectric(meter *models.ElectricMeter, updates map[string]interface{}) error {
	return r.db.Model(meter).Updates(updates).Error
}

// DeleteElectric 删除电表
func (r *MeterRepository) DeleteElectric(id int64) error {
	return r.db.Delete(&models.ElectricMeter{}, id).Error
}

// UpdateElectricReading 更新电表读数
func (r *MeterRepository) UpdateElectricReading(tx *gorm.DB, id int64, reading decimal.Decimal, collectTime time.Time) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	return db.Model(&models.ElectricMeter{}).Where("id = ?", id).Updates(map[string]interface{}{
		"current_reading": reading,
		"last_collect_at": collectTime,
		"updated_at":      time.Now(),
	}).Error
}

// CountElectricByStatus 统计电表数量
func (r *MeterRepository) CountElectricByStatus(status int16) (int64, error) {
	var count int64
	if err := r.db.Model(&models.ElectricMeter{}).Where("status = ?", status).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CountElectricOnline 统计在线电表数量
func (r *MeterRepository) CountElectricOnline() (int64, error) {
	var count int64
	if err := r.db.Model(&models.ElectricMeter{}).Where("status = ? AND online_status = ?", models.MeterStatusNormal, models.MeterOnline).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CountElectricByRateID 统计使用指定费率的电表数量
func (r *MeterRepository) CountElectricByRateID(rateID int64) (int64, error) {
	var count int64
	if err := r.db.Model(&models.ElectricMeter{}).Where("rate_id = ?", rateID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// ========== 水表操作 ==========

// FindWaterByID 根据ID查找水表
func (r *MeterRepository) FindWaterByID(id int64) (*models.WaterMeter, error) {
	var meter models.WaterMeter
	if err := r.db.Preload("WaterRate").Preload("Merchant").Preload("Shop").First(&meter, id).Error; err != nil {
		return nil, err
	}
	return &meter, nil
}

// FindWaterByMeterNo 根据表计编号查找水表
func (r *MeterRepository) FindWaterByMeterNo(meterNo string) (*models.WaterMeter, error) {
	var meter models.WaterMeter
	if err := r.db.Where("meter_no = ?", meterNo).First(&meter).Error; err != nil {
		return nil, err
	}
	return &meter, nil
}

// ExistsWaterByMeterNo 检查水表编号是否存在
func (r *MeterRepository) ExistsWaterByMeterNo(meterNo string) (bool, error) {
	var count int64
	if err := r.db.Model(&models.WaterMeter{}).Where("meter_no = ?", meterNo).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// ListWater 获取水表列表
func (r *MeterRepository) ListWater(meterNo string, status, onlineStatus *int16, offset, limit int) ([]models.WaterMeter, int64, error) {
	var meters []models.WaterMeter
	var total int64

	query := r.db.Model(&models.WaterMeter{})

	if meterNo != "" {
		query = query.Where("meter_no LIKE ?", "%"+meterNo+"%")
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	if onlineStatus != nil {
		query = query.Where("online_status = ?", *onlineStatus)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("WaterRate").Preload("Merchant").Preload("Shop").Order("id DESC").Offset(offset).Limit(limit).Find(&meters).Error; err != nil {
		return nil, 0, err
	}

	return meters, total, nil
}

// CreateWater 创建水表
func (r *MeterRepository) CreateWater(meter *models.WaterMeter) error {
	return r.db.Create(meter).Error
}

// UpdateWater 更新水表
func (r *MeterRepository) UpdateWater(meter *models.WaterMeter, updates map[string]interface{}) error {
	return r.db.Model(meter).Updates(updates).Error
}

// DeleteWater 删除水表
func (r *MeterRepository) DeleteWater(id int64) error {
	return r.db.Delete(&models.WaterMeter{}, id).Error
}

// UpdateWaterReading 更新水表读数
func (r *MeterRepository) UpdateWaterReading(tx *gorm.DB, id int64, reading decimal.Decimal, collectTime time.Time) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	return db.Model(&models.WaterMeter{}).Where("id = ?", id).Updates(map[string]interface{}{
		"current_reading": reading,
		"last_collect_at": collectTime,
		"updated_at":      time.Now(),
	}).Error
}

// CountWaterByStatus 统计水表数量
func (r *MeterRepository) CountWaterByStatus(status int16) (int64, error) {
	var count int64
	if err := r.db.Model(&models.WaterMeter{}).Where("status = ?", status).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CountWaterOnline 统计在线水表数量
func (r *MeterRepository) CountWaterOnline() (int64, error) {
	var count int64
	if err := r.db.Model(&models.WaterMeter{}).Where("status = ? AND online_status = ?", models.MeterStatusNormal, models.MeterOnline).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CountWaterByRateID 统计使用指定费率的水表数量
func (r *MeterRepository) CountWaterByRateID(rateID int64) (int64, error) {
	var count int64
	if err := r.db.Model(&models.WaterMeter{}).Where("rate_id = ?", rateID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// ========== 统计方法 (兼容 Dashboard) ==========

// CountByStatus 统计所有表计数量(电表+水表)
func (r *MeterRepository) CountByStatus(status int16) (int64, error) {
	electricCount, err := r.CountElectricByStatus(status)
	if err != nil {
		return 0, err
	}
	waterCount, err := r.CountWaterByStatus(status)
	if err != nil {
		return 0, err
	}
	return electricCount + waterCount, nil
}

// CountOnline 统计所有在线表计数量(电表+水表)
func (r *MeterRepository) CountOnline() (int64, error) {
	electricOnline, err := r.CountElectricOnline()
	if err != nil {
		return 0, err
	}
	waterOnline, err := r.CountWaterOnline()
	if err != nil {
		return 0, err
	}
	return electricOnline + waterOnline, nil
}
