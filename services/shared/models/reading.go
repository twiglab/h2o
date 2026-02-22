package models

import (
	"time"

	"github.com/shopspring/decimal"
)

// 采集类型
const (
	CollectTypeAuto   = 1 // 自动采集
	CollectTypeManual = 2 // 手工抄表
)

// 读数状态
const (
	ReadingStatusNormal   = 1 // 正常
	ReadingStatusAbnormal = 2 // 异常
	ReadingStatusVoid     = 3 // 作废
)

// 用量状态
const (
	ConsumptionStatusPending  = 1 // 待扣费
	ConsumptionStatusDeducted = 2 // 已扣费
	ConsumptionStatusAbnormal = 3 // 异常
)

// ElectricReading 电表读数记录
type ElectricReading struct {
	ID           int64           `gorm:"column:id;primaryKey" json:"id"`
	MerchantID   int64           `gorm:"column:merchant_id" json:"merchant_id"`
	MeterID      int64           `gorm:"column:meter_id" json:"meter_id"`
	MeterNo      string          `gorm:"column:meter_no" json:"meter_no"`
	ReadingValue decimal.Decimal `gorm:"column:reading_value" json:"reading_value"`
	ReadingTime  time.Time       `gorm:"column:reading_time" json:"reading_time"`
	CollectType  int8            `gorm:"column:collect_type" json:"collect_type"`
	OperatorID   *int64          `gorm:"column:operator_id" json:"operator_id"`
	Status       int8            `gorm:"column:status" json:"status"`
	Remark       *string         `gorm:"column:remark" json:"remark"`
	CreatedAt    time.Time       `gorm:"column:created_at" json:"created_at"`

	// Relations
	Merchant *Merchant      `gorm:"foreignKey:MerchantID" json:"merchant,omitempty"`
	Meter    *ElectricMeter `gorm:"foreignKey:MeterID" json:"meter,omitempty"`
}

func (ElectricReading) TableName() string { return "col_electric_reading" }

// IsNormal 是否正常
func (r *ElectricReading) IsNormal() bool {
	return r.Status == ReadingStatusNormal
}

// IsAuto 是否自动采集
func (r *ElectricReading) IsAuto() bool {
	return r.CollectType == CollectTypeAuto
}

// WaterReading 水表读数记录
type WaterReading struct {
	ID           int64           `gorm:"column:id;primaryKey" json:"id"`
	MerchantID   int64           `gorm:"column:merchant_id" json:"merchant_id"`
	MeterID      int64           `gorm:"column:meter_id" json:"meter_id"`
	MeterNo      string          `gorm:"column:meter_no" json:"meter_no"`
	ReadingValue decimal.Decimal `gorm:"column:reading_value" json:"reading_value"`
	ReadingTime  time.Time       `gorm:"column:reading_time" json:"reading_time"`
	CollectType  int8            `gorm:"column:collect_type" json:"collect_type"`
	OperatorID   *int64          `gorm:"column:operator_id" json:"operator_id"`
	Status       int8            `gorm:"column:status" json:"status"`
	Remark       *string         `gorm:"column:remark" json:"remark"`
	CreatedAt    time.Time       `gorm:"column:created_at" json:"created_at"`

	// Relations
	Merchant *Merchant   `gorm:"foreignKey:MerchantID" json:"merchant,omitempty"`
	Meter    *WaterMeter `gorm:"foreignKey:MeterID" json:"meter,omitempty"`
}

func (WaterReading) TableName() string { return "col_water_reading" }

// IsNormal 是否正常
func (r *WaterReading) IsNormal() bool {
	return r.Status == ReadingStatusNormal
}

// IsAuto 是否自动采集
func (r *WaterReading) IsAuto() bool {
	return r.CollectType == CollectTypeAuto
}

// ElectricManualReading 电表手工抄表记录
type ElectricManualReading struct {
	ID              int64           `gorm:"column:id;primaryKey" json:"id"`
	ReadingNo       string          `gorm:"column:reading_no" json:"reading_no"`
	MerchantID      int64           `gorm:"column:merchant_id" json:"merchant_id"`
	MeterID         int64           `gorm:"column:meter_id" json:"meter_id"`
	MeterNo         string          `gorm:"column:meter_no" json:"meter_no"`
	PreviousReading decimal.Decimal `gorm:"column:previous_reading" json:"previous_reading"`
	CurrentReading  decimal.Decimal `gorm:"column:current_reading" json:"current_reading"`
	ReadingTime     time.Time       `gorm:"column:reading_time" json:"reading_time"`
	OperatorID      *int64          `gorm:"column:operator_id" json:"operator_id"`
	OperatorName    *string         `gorm:"column:operator_name" json:"operator_name"`
	Status          int8            `gorm:"column:status" json:"status"`
	Remark          *string         `gorm:"column:remark" json:"remark"`
	CreatedAt       time.Time       `gorm:"column:created_at" json:"created_at"`
	UpdatedAt       time.Time       `gorm:"column:updated_at" json:"updated_at"`

	// Relations
	Merchant *Merchant      `gorm:"foreignKey:MerchantID" json:"merchant,omitempty"`
	Meter    *ElectricMeter `gorm:"foreignKey:MeterID" json:"meter,omitempty"`
}

func (ElectricManualReading) TableName() string { return "col_electric_manual_reading" }

// GetConsumption 获取用量
func (r *ElectricManualReading) GetConsumption() decimal.Decimal {
	return r.CurrentReading.Sub(r.PreviousReading)
}

// WaterManualReading 水表手工抄表记录
type WaterManualReading struct {
	ID              int64           `gorm:"column:id;primaryKey" json:"id"`
	ReadingNo       string          `gorm:"column:reading_no" json:"reading_no"`
	MerchantID      int64           `gorm:"column:merchant_id" json:"merchant_id"`
	MeterID         int64           `gorm:"column:meter_id" json:"meter_id"`
	MeterNo         string          `gorm:"column:meter_no" json:"meter_no"`
	PreviousReading decimal.Decimal `gorm:"column:previous_reading" json:"previous_reading"`
	CurrentReading  decimal.Decimal `gorm:"column:current_reading" json:"current_reading"`
	ReadingTime     time.Time       `gorm:"column:reading_time" json:"reading_time"`
	OperatorID      *int64          `gorm:"column:operator_id" json:"operator_id"`
	OperatorName    *string         `gorm:"column:operator_name" json:"operator_name"`
	Status          int8            `gorm:"column:status" json:"status"`
	Remark          *string         `gorm:"column:remark" json:"remark"`
	CreatedAt       time.Time       `gorm:"column:created_at" json:"created_at"`
	UpdatedAt       time.Time       `gorm:"column:updated_at" json:"updated_at"`

	// Relations
	Merchant *Merchant   `gorm:"foreignKey:MerchantID" json:"merchant,omitempty"`
	Meter    *WaterMeter `gorm:"foreignKey:MeterID" json:"meter,omitempty"`
}

func (WaterManualReading) TableName() string { return "col_water_manual_reading" }

// GetConsumption 获取用量
func (r *WaterManualReading) GetConsumption() decimal.Decimal {
	return r.CurrentReading.Sub(r.PreviousReading)
}

// ElectricConsumption 电费用量记录
type ElectricConsumption struct {
	ID             int64           `gorm:"column:id;primaryKey" json:"id"`
	MerchantID     int64           `gorm:"column:merchant_id" json:"merchant_id"`
	ShopID         *int64          `gorm:"column:shop_id" json:"shop_id"`
	AccountID      int64           `gorm:"column:account_id" json:"account_id"`
	MeterID        int64           `gorm:"column:meter_id" json:"meter_id"`
	StartReadingID int64           `gorm:"column:start_reading_id" json:"start_reading_id"`
	EndReadingID   int64           `gorm:"column:end_reading_id" json:"end_reading_id"`
	StartReading   decimal.Decimal `gorm:"column:start_reading" json:"start_reading"`
	EndReading     decimal.Decimal `gorm:"column:end_reading" json:"end_reading"`
	Consumption    decimal.Decimal `gorm:"column:consumption" json:"consumption"`
	PeriodStart    time.Time       `gorm:"column:period_start" json:"period_start"`
	PeriodEnd      time.Time       `gorm:"column:period_end" json:"period_end"`
	Status         int8            `gorm:"column:status" json:"status"`
	CreatedAt      time.Time       `gorm:"column:created_at" json:"created_at"`

	// Relations
	Merchant *Merchant      `gorm:"foreignKey:MerchantID" json:"merchant,omitempty"`
	Shop     *Shop          `gorm:"foreignKey:ShopID" json:"shop,omitempty"`
	Account  *Account       `gorm:"foreignKey:AccountID" json:"account,omitempty"`
	Meter    *ElectricMeter `gorm:"foreignKey:MeterID" json:"meter,omitempty"`
}

func (ElectricConsumption) TableName() string { return "fin_electric_consumption" }

// IsPending 是否待扣费
func (c *ElectricConsumption) IsPending() bool {
	return c.Status == ConsumptionStatusPending
}

// IsDeducted 是否已扣费
func (c *ElectricConsumption) IsDeducted() bool {
	return c.Status == ConsumptionStatusDeducted
}

// WaterConsumption 水费用量记录
type WaterConsumption struct {
	ID             int64           `gorm:"column:id;primaryKey" json:"id"`
	MerchantID     int64           `gorm:"column:merchant_id" json:"merchant_id"`
	ShopID         *int64          `gorm:"column:shop_id" json:"shop_id"`
	AccountID      int64           `gorm:"column:account_id" json:"account_id"`
	MeterID        int64           `gorm:"column:meter_id" json:"meter_id"`
	StartReadingID int64           `gorm:"column:start_reading_id" json:"start_reading_id"`
	EndReadingID   int64           `gorm:"column:end_reading_id" json:"end_reading_id"`
	StartReading   decimal.Decimal `gorm:"column:start_reading" json:"start_reading"`
	EndReading     decimal.Decimal `gorm:"column:end_reading" json:"end_reading"`
	Consumption    decimal.Decimal `gorm:"column:consumption" json:"consumption"`
	PeriodStart    time.Time       `gorm:"column:period_start" json:"period_start"`
	PeriodEnd      time.Time       `gorm:"column:period_end" json:"period_end"`
	Status         int8            `gorm:"column:status" json:"status"`
	CreatedAt      time.Time       `gorm:"column:created_at" json:"created_at"`

	// Relations
	Merchant *Merchant   `gorm:"foreignKey:MerchantID" json:"merchant,omitempty"`
	Shop     *Shop       `gorm:"foreignKey:ShopID" json:"shop,omitempty"`
	Account  *Account    `gorm:"foreignKey:AccountID" json:"account,omitempty"`
	Meter    *WaterMeter `gorm:"foreignKey:MeterID" json:"meter,omitempty"`
}

func (WaterConsumption) TableName() string { return "fin_water_consumption" }

// IsPending 是否待扣费
func (c *WaterConsumption) IsPending() bool {
	return c.Status == ConsumptionStatusPending
}

// IsDeducted 是否已扣费
func (c *WaterConsumption) IsDeducted() bool {
	return c.Status == ConsumptionStatusDeducted
}
