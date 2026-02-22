package repository

import (
	"shared/models"
	"time"

	"gorm.io/gorm"
)

// ReadingRepository 电表读数数据仓库
type ReadingRepository struct {
	db *gorm.DB
}

// NewReadingRepository 创建读数仓库
func NewReadingRepository(db *gorm.DB) *ReadingRepository {
	return &ReadingRepository{db: db}
}

// CreateReading 创建电表读数记录
func (r *ReadingRepository) CreateReading(tx *gorm.DB, reading *models.ElectricReading) error {
	return tx.Create(reading).Error
}

// FindLastReading 查找最后一条有效电表读数
func (r *ReadingRepository) FindLastReading(meterID int64) (*models.ElectricReading, error) {
	var reading models.ElectricReading
	err := r.db.Where("meter_id = ? AND status = ?", meterID, models.ReadingStatusNormal).
		Order("reading_time DESC").
		First(&reading).Error
	if err != nil {
		return nil, err
	}
	return &reading, nil
}

// FindLastReadingBefore 查找指定时间之前的最后一条有效电表读数
func (r *ReadingRepository) FindLastReadingBefore(meterID int64, before time.Time) (*models.ElectricReading, error) {
	var reading models.ElectricReading
	err := r.db.Where("meter_id = ? AND reading_time < ? AND status = ?",
		meterID, before, models.ReadingStatusNormal).
		Order("reading_time DESC").
		First(&reading).Error
	if err != nil {
		return nil, err
	}
	return &reading, nil
}

// FindByID 根据ID查找电表读数
func (r *ReadingRepository) FindByID(id int64) (*models.ElectricReading, error) {
	var reading models.ElectricReading
	err := r.db.First(&reading, id).Error
	if err != nil {
		return nil, err
	}
	return &reading, nil
}

// CreateConsumption 创建电表用量记录
func (r *ReadingRepository) CreateConsumption(tx *gorm.DB, consumption *models.ElectricConsumption) error {
	return tx.Create(consumption).Error
}

// FindConsumptionByID 根据ID查找电表用量记录
func (r *ReadingRepository) FindConsumptionByID(id int64) (*models.ElectricConsumption, error) {
	var consumption models.ElectricConsumption
	err := r.db.First(&consumption, id).Error
	if err != nil {
		return nil, err
	}
	return &consumption, nil
}

// UpdateConsumptionStatus 更新电表用量记录状态
func (r *ReadingRepository) UpdateConsumptionStatus(tx *gorm.DB, id int64, status int16) error {
	return tx.Model(&models.ElectricConsumption{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// FindReadingsByMeterIDAndTimeRange 查找时间范围内的电表读数
func (r *ReadingRepository) FindReadingsByMeterIDAndTimeRange(meterID int64, start, end time.Time) ([]models.ElectricReading, error) {
	var readings []models.ElectricReading
	err := r.db.Where("meter_id = ? AND reading_time BETWEEN ? AND ? AND status = ?",
		meterID, start, end, models.ReadingStatusNormal).
		Order("reading_time ASC").
		Find(&readings).Error
	return readings, err
}
