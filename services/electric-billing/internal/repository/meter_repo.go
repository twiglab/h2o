package repository

import (
	"shared/models"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// MeterRepository 电表数据仓库
type MeterRepository struct {
	db *gorm.DB
}

// NewMeterRepository 创建电表仓库
func NewMeterRepository(db *gorm.DB) *MeterRepository {
	return &MeterRepository{db: db}
}

// FindByMeterNo 根据表计编号查找
func (r *MeterRepository) FindByMeterNo(meterNo string) (*models.ElectricMeter, error) {
	var meter models.ElectricMeter
	err := r.db.Where("meter_no = ?", meterNo).First(&meter).Error
	if err != nil {
		return nil, err
	}
	return &meter, nil
}

// FindByMeterNoWithRate 根据表计编号查找并预加载费率
func (r *MeterRepository) FindByMeterNoWithRate(meterNo string) (*models.ElectricMeter, error) {
	var meter models.ElectricMeter
	err := r.db.Preload("ElectricRate").
		Preload("ElectricRate.TOUDetails").
		Preload("ElectricRate.ServiceFees").
		Where("meter_no = ?", meterNo).
		First(&meter).Error
	if err != nil {
		return nil, err
	}
	return &meter, nil
}

// FindByID 根据ID查找
func (r *MeterRepository) FindByID(id int64) (*models.ElectricMeter, error) {
	var meter models.ElectricMeter
	err := r.db.First(&meter, id).Error
	if err != nil {
		return nil, err
	}
	return &meter, nil
}

// FindByCommAddr 根据通信地址查找
func (r *MeterRepository) FindByCommAddr(commAddr string) (*models.ElectricMeter, error) {
	var meter models.ElectricMeter
	err := r.db.Where("comm_addr = ?", commAddr).First(&meter).Error
	if err != nil {
		return nil, err
	}
	return &meter, nil
}

// FindByCommAddrWithRate 根据通信地址查找并预加载费率
func (r *MeterRepository) FindByCommAddrWithRate(commAddr string) (*models.ElectricMeter, error) {
	var meter models.ElectricMeter
	err := r.db.Preload("ElectricRate").
		Preload("ElectricRate.TOUDetails").
		Preload("ElectricRate.ServiceFees").
		Where("comm_addr = ?", commAddr).
		First(&meter).Error
	if err != nil {
		return nil, err
	}
	return &meter, nil
}

// FindByIDForUpdate 根据ID查找并锁定(用于更新)
func (r *MeterRepository) FindByIDForUpdate(tx *gorm.DB, id int64) (*models.ElectricMeter, error) {
	var meter models.ElectricMeter
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&meter, id).Error
	if err != nil {
		return nil, err
	}
	return &meter, nil
}

// UpdateReading 更新电表读数和状态
func (r *MeterRepository) UpdateReading(tx *gorm.DB, id int64, reading decimal.Decimal, collectTime time.Time) error {
	return tx.Model(&models.ElectricMeter{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"current_reading": reading,
			"last_collect_at": collectTime,
			"online_status":   models.MeterOnline,
			"updated_at":      time.Now(),
		}).Error
}

// FindByIDWithRelations 根据ID查找并预加载所有关联数据（用于扣费记录快照）
func (r *MeterRepository) FindByIDWithRelations(id int64) (*models.ElectricMeter, error) {
	var meter models.ElectricMeter
	err := r.db.Preload("Merchant").
		Preload("Shop").
		Preload("Account").
		Preload("ElectricRate").
		Preload("ElectricRate.TOUDetails").
		First(&meter, id).Error
	if err != nil {
		return nil, err
	}
	return &meter, nil
}
