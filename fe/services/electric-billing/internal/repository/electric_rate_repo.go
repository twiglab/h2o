package repository

import (
	"shared/models"

	"gorm.io/gorm"
)

// ElectricRateRepository 电费费率数据仓库
type ElectricRateRepository struct {
	db *gorm.DB
}

// NewElectricRateRepository 创建电费费率仓库
func NewElectricRateRepository(db *gorm.DB) *ElectricRateRepository {
	return &ElectricRateRepository{db: db}
}

// FindByID 根据ID查找费率
func (r *ElectricRateRepository) FindByID(id int64) (*models.ElectricRate, error) {
	var rate models.ElectricRate
	err := r.db.First(&rate, id).Error
	if err != nil {
		return nil, err
	}
	return &rate, nil
}

// FindByIDWithDetails 根据ID查找费率(包含详情)
func (r *ElectricRateRepository) FindByIDWithDetails(id int64) (*models.ElectricRate, error) {
	var rate models.ElectricRate
	err := r.db.Preload("TOUDetails").
		Preload("ServiceFees").
		First(&rate, id).Error
	if err != nil {
		return nil, err
	}
	return &rate, nil
}

// GetApplicableRate 获取适用费率
// 优先级: meter.rate_id > scope=1 全局默认
func (r *ElectricRateRepository) GetApplicableRate(meterID int64, defaultRateID *int64) (*models.ElectricRate, error) {
	// 1. 使用电表关联的默认费率
	if defaultRateID != nil {
		var defaultRate models.ElectricRate
		err := r.db.Preload("TOUDetails").
			Preload("ServiceFees").
			First(&defaultRate, *defaultRateID).Error
		if err == nil {
			return &defaultRate, nil
		}
	}

	// 2. 使用全局默认费率 (scope=1, status=1)
	var globalRate models.ElectricRate
	err := r.db.Preload("TOUDetails").
		Preload("ServiceFees").
		Where("scope = ? AND status = ?",
			models.ElectricRateScopeGlobal, 1).
		Order("effective_date DESC, id DESC").
		First(&globalRate).Error
	if err != nil {
		return nil, err
	}

	return &globalRate, nil
}

// FindGlobalRate 查找全局默认费率
func (r *ElectricRateRepository) FindGlobalRate() (*models.ElectricRate, error) {
	var rate models.ElectricRate
	err := r.db.Preload("TOUDetails").
		Preload("ServiceFees").
		Where("scope = ? AND status = ?", models.ElectricRateScopeGlobal, 1).
		Order("effective_date DESC, id DESC").
		First(&rate).Error
	if err != nil {
		return nil, err
	}
	return &rate, nil
}

// FindTOUByRateID 查找费率的分时电价详情
func (r *ElectricRateRepository) FindTOUByRateID(rateID int64) ([]models.ElectricRateTOU, error) {
	var touDetails []models.ElectricRateTOU
	err := r.db.Where("rate_id = ?", rateID).
		Order("start_time ASC").
		Find(&touDetails).Error
	return touDetails, err
}

// FindServiceFeesByRateID 查找费率的服务费详情
func (r *ElectricRateRepository) FindServiceFeesByRateID(rateID int64) ([]models.ElectricRateServiceFee, error) {
	var serviceFees []models.ElectricRateServiceFee
	err := r.db.Where("rate_id = ?", rateID).Find(&serviceFees).Error
	return serviceFees, err
}
