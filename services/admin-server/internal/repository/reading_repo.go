package repository

import (
	"shared/models"

	"gorm.io/gorm"
)

// ReadingRepository 读数数据仓库
type ReadingRepository struct {
	db *gorm.DB
}

// NewReadingRepository 创建读数仓库
func NewReadingRepository(db *gorm.DB) *ReadingRepository {
	return &ReadingRepository{db: db}
}

// ListElectricByMeterID 获取电表读数列表
func (r *ReadingRepository) ListElectricByMeterID(meterID int64, startDate, endDate string, offset, limit int) ([]models.ElectricReading, int64, error) {
	var readings []models.ElectricReading
	var total int64

	query := r.db.Model(&models.ElectricReading{}).Where("meter_id = ?", meterID)

	if startDate != "" {
		query = query.Where("reading_time >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("reading_time <= ?", endDate+" 23:59:59")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("reading_time DESC").Offset(offset).Limit(limit).Find(&readings).Error; err != nil {
		return nil, 0, err
	}

	return readings, total, nil
}

// ListWaterByMeterID 获取水表读数列表
func (r *ReadingRepository) ListWaterByMeterID(meterID int64, startDate, endDate string, offset, limit int) ([]models.WaterReading, int64, error) {
	var readings []models.WaterReading
	var total int64

	query := r.db.Model(&models.WaterReading{}).Where("meter_id = ?", meterID)

	if startDate != "" {
		query = query.Where("reading_time >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("reading_time <= ?", endDate+" 23:59:59")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("reading_time DESC").Offset(offset).Limit(limit).Find(&readings).Error; err != nil {
		return nil, 0, err
	}

	return readings, total, nil
}

// ListRecentElectricByMeterID 获取电表最近的读数
func (r *ReadingRepository) ListRecentElectricByMeterID(meterID int64, limit int) ([]models.ElectricReading, error) {
	var readings []models.ElectricReading
	if err := r.db.Where("meter_id = ?", meterID).Order("reading_time DESC").Limit(limit).Find(&readings).Error; err != nil {
		return nil, err
	}
	return readings, nil
}

// ListRecentWaterByMeterID 获取水表最近的读数
func (r *ReadingRepository) ListRecentWaterByMeterID(meterID int64, limit int) ([]models.WaterReading, error) {
	var readings []models.WaterReading
	if err := r.db.Where("meter_id = ?", meterID).Order("reading_time DESC").Limit(limit).Find(&readings).Error; err != nil {
		return nil, err
	}
	return readings, nil
}

// CreateElectricReading 创建电表读数记录
func (r *ReadingRepository) CreateElectricReading(tx *gorm.DB, reading *models.ElectricReading) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	return db.Create(reading).Error
}

// CreateWaterReading 创建水表读数记录
func (r *ReadingRepository) CreateWaterReading(tx *gorm.DB, reading *models.WaterReading) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	return db.Create(reading).Error
}

// CreateElectricManualReading 创建电表手工抄表记录
func (r *ReadingRepository) CreateElectricManualReading(tx *gorm.DB, reading *models.ElectricManualReading) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	return db.Create(reading).Error
}

// CreateWaterManualReading 创建水表手工抄表记录
func (r *ReadingRepository) CreateWaterManualReading(tx *gorm.DB, reading *models.WaterManualReading) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	return db.Create(reading).Error
}

// ListElectricManualReadings 获取电表手工抄表记录列表
func (r *ReadingRepository) ListElectricManualReadings(meterNo, startDate, endDate string, offset, limit int) ([]models.ElectricManualReading, int64, error) {
	var readings []models.ElectricManualReading
	var total int64

	query := r.db.Model(&models.ElectricManualReading{}).Where("status = ?", 1)

	if meterNo != "" {
		query = query.Where("meter_no LIKE ?", "%"+meterNo+"%")
	}
	if startDate != "" {
		query = query.Where("reading_time >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("reading_time <= ?", endDate+" 23:59:59")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("Meter").Order("reading_time DESC").Offset(offset).Limit(limit).Find(&readings).Error; err != nil {
		return nil, 0, err
	}

	return readings, total, nil
}

// ListWaterManualReadings 获取水表手工抄表记录列表
func (r *ReadingRepository) ListWaterManualReadings(meterNo, startDate, endDate string, offset, limit int) ([]models.WaterManualReading, int64, error) {
	var readings []models.WaterManualReading
	var total int64

	query := r.db.Model(&models.WaterManualReading{}).Where("status = ?", 1)

	if meterNo != "" {
		query = query.Where("meter_no LIKE ?", "%"+meterNo+"%")
	}
	if startDate != "" {
		query = query.Where("reading_time >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("reading_time <= ?", endDate+" 23:59:59")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("Meter").Order("reading_time DESC").Offset(offset).Limit(limit).Find(&readings).Error; err != nil {
		return nil, 0, err
	}

	return readings, total, nil
}

// CountElectricByDate 统计指定日期的电表读数数量
func (r *ReadingRepository) CountElectricByDate(date string) (int64, error) {
	var count int64
	if err := r.db.Model(&models.ElectricReading{}).Where("DATE(reading_time) = ?", date).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CountWaterByDate 统计指定日期的水表读数数量
func (r *ReadingRepository) CountWaterByDate(date string) (int64, error) {
	var count int64
	if err := r.db.Model(&models.WaterReading{}).Where("DATE(reading_time) = ?", date).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CountElectricCollectedByDate 统计指定日期采集的电表数量（去重）
func (r *ReadingRepository) CountElectricCollectedByDate(date string) (int64, error) {
	var count int64
	if err := r.db.Model(&models.ElectricReading{}).Where("DATE(reading_time) = ?", date).Distinct("meter_id").Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CountWaterCollectedByDate 统计指定日期采集的水表数量（去重）
func (r *ReadingRepository) CountWaterCollectedByDate(date string) (int64, error) {
	var count int64
	if err := r.db.Model(&models.WaterReading{}).Where("DATE(reading_time) = ?", date).Distinct("meter_id").Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// FindLastElectricReading 查找电表最后一条有效读数
func (r *ReadingRepository) FindLastElectricReading(meterID int64) (*models.ElectricReading, error) {
	var reading models.ElectricReading
	err := r.db.Where("meter_id = ? AND status = ?", meterID, models.ReadingStatusNormal).
		Order("reading_time DESC").
		First(&reading).Error
	if err != nil {
		return nil, err
	}
	return &reading, nil
}

// FindLastWaterReading 查找水表最后一条有效读数
func (r *ReadingRepository) FindLastWaterReading(meterID int64) (*models.WaterReading, error) {
	var reading models.WaterReading
	err := r.db.Where("meter_id = ? AND status = ?", meterID, models.ReadingStatusNormal).
		Order("reading_time DESC").
		First(&reading).Error
	if err != nil {
		return nil, err
	}
	return &reading, nil
}

// DB 返回数据库连接（用于事务）
func (r *ReadingRepository) DB() *gorm.DB {
	return r.db
}
