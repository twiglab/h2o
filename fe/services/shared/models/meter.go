package models

import (
	"time"

	"github.com/shopspring/decimal"
)

// 在线状态
const (
	MeterOffline = 0 // 离线
	MeterOnline  = 1 // 在线
)

// 电表状态
const (
	MeterStatusDisabled = 0 // 停用
	MeterStatusNormal   = 1 // 正常
)

// 表计类型常量 (用于手工抄表等统一场景)
const (
	MeterTypeElectric = 1 // 电表
	MeterTypeWater    = 2 // 水表
)

// 通信协议
const (
	MeterProtocolModbus = "modbus"
	MeterProtocolDLT645 = "dlt645"
	MeterProtocolMQTT   = "mqtt"
)

// ElectricMeter 电表
type ElectricMeter struct {
	ID         int64  `gorm:"column:id;primaryKey" json:"id"`
	MeterNo    string `gorm:"column:meter_no" json:"meter_no"`
	MerchantID int64  `gorm:"column:merchant_id" json:"merchant_id"`
	ShopID     *int64 `gorm:"column:shop_id" json:"shop_id"`
	AccountID  *int64 `gorm:"column:account_id" json:"account_id"`
	RateID     *int64 `gorm:"column:rate_id" json:"rate_id"`

	// 通信配置
	MqttTopic *string `gorm:"column:mqtt_topic" json:"mqtt_topic"`
	CommAddr  *string `gorm:"column:comm_addr" json:"comm_addr"`
	Protocol  string  `gorm:"column:protocol" json:"protocol"`

	// 表计参数
	Multiplier     decimal.Decimal `gorm:"column:multiplier" json:"multiplier"`
	InitReading    decimal.Decimal `gorm:"column:init_reading" json:"init_reading"`
	CurrentReading decimal.Decimal `gorm:"column:current_reading" json:"current_reading"`
	LastCollectAt  *time.Time      `gorm:"column:last_collect_at" json:"last_collect_at"`

	// 状态
	OnlineStatus int8 `gorm:"column:online_status" json:"online_status"`
	Status       int8 `gorm:"column:status" json:"status"`

	// 设备信息
	Brand       *string    `gorm:"column:brand" json:"brand"`
	Model       *string    `gorm:"column:model" json:"model"`
	InstallDate *time.Time `gorm:"column:install_date" json:"install_date"`
	Remark      *string    `gorm:"column:remark" json:"remark"`

	CreatedAt time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at,omitempty"`

	// Relations
	Merchant     *Merchant     `gorm:"foreignKey:MerchantID" json:"merchant,omitempty"`
	Shop         *Shop         `gorm:"foreignKey:ShopID" json:"shop,omitempty"`
	Account      *Account      `gorm:"foreignKey:AccountID" json:"account,omitempty"`
	ElectricRate *ElectricRate `gorm:"foreignKey:RateID" json:"rate,omitempty"`
}

func (ElectricMeter) TableName() string { return "biz_electric_meter" }

// IsOnline 是否在线
func (m *ElectricMeter) IsOnline() bool {
	return m.OnlineStatus == MeterOnline
}

// IsEnabled 是否启用
func (m *ElectricMeter) IsEnabled() bool {
	return m.Status == MeterStatusNormal
}

// GetUsage 获取当前用量(当前读数-初始读数)*倍率
func (m *ElectricMeter) GetUsage() decimal.Decimal {
	return m.CurrentReading.Sub(m.InitReading).Mul(m.Multiplier)
}

// WaterMeter 水表
type WaterMeter struct {
	ID         int64  `gorm:"column:id;primaryKey" json:"id"`
	MeterNo    string `gorm:"column:meter_no" json:"meter_no"`
	MerchantID int64  `gorm:"column:merchant_id" json:"merchant_id"`
	ShopID     *int64 `gorm:"column:shop_id" json:"shop_id"`
	AccountID  *int64 `gorm:"column:account_id" json:"account_id"`
	RateID     *int64 `gorm:"column:rate_id" json:"rate_id"`

	// 通信配置
	MqttTopic *string `gorm:"column:mqtt_topic" json:"mqtt_topic"`
	CommAddr  *string `gorm:"column:comm_addr" json:"comm_addr"`
	Protocol  string  `gorm:"column:protocol" json:"protocol"`

	// 表计参数
	Multiplier     decimal.Decimal `gorm:"column:multiplier" json:"multiplier"`
	InitReading    decimal.Decimal `gorm:"column:init_reading" json:"init_reading"`
	CurrentReading decimal.Decimal `gorm:"column:current_reading" json:"current_reading"`
	LastCollectAt  *time.Time      `gorm:"column:last_collect_at" json:"last_collect_at"`

	// 状态
	OnlineStatus int8 `gorm:"column:online_status" json:"online_status"`
	Status       int8 `gorm:"column:status" json:"status"`

	// 设备信息
	Brand       *string    `gorm:"column:brand" json:"brand"`
	Model       *string    `gorm:"column:model" json:"model"`
	InstallDate *time.Time `gorm:"column:install_date" json:"install_date"`
	Remark      *string    `gorm:"column:remark" json:"remark"`

	CreatedAt time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at,omitempty"`

	// Relations
	Merchant  *Merchant  `gorm:"foreignKey:MerchantID" json:"merchant,omitempty"`
	Shop      *Shop      `gorm:"foreignKey:ShopID" json:"shop,omitempty"`
	Account   *Account   `gorm:"foreignKey:AccountID" json:"account,omitempty"`
	WaterRate *WaterRate `gorm:"foreignKey:RateID" json:"rate,omitempty"`
}

func (WaterMeter) TableName() string { return "biz_water_meter" }

// IsOnline 是否在线
func (m *WaterMeter) IsOnline() bool {
	return m.OnlineStatus == MeterOnline
}

// IsEnabled 是否启用
func (m *WaterMeter) IsEnabled() bool {
	return m.Status == MeterStatusNormal
}

// GetUsage 获取当前用量(当前读数-初始读数)*倍率
func (m *WaterMeter) GetUsage() decimal.Decimal {
	return m.CurrentReading.Sub(m.InitReading).Mul(m.Multiplier)
}
