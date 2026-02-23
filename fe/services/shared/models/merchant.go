package models

import "time"

// 商户类型
const (
	MerchantTypeEnterprise = 1 // 企业
	MerchantTypePersonal   = 2 // 个人
)

// 商户状态
const (
	MerchantStatusDisabled = 0 // 停用
	MerchantStatusNormal   = 1 // 正常
)

// Merchant 商户表(签约主体)
type Merchant struct {
	ID           int64  `gorm:"column:id;primaryKey" json:"id"`
	MerchantNo   string `gorm:"column:merchant_no" json:"merchant_no"`
	MerchantName string `gorm:"column:merchant_name" json:"merchant_name"`
	MerchantType int8   `gorm:"column:merchant_type" json:"merchant_type"`

	// 联系人
	ContactName  *string `gorm:"column:contact_name" json:"contact_name"`
	ContactPhone *string `gorm:"column:contact_phone" json:"contact_phone"`

	Status    int8       `gorm:"column:status" json:"status"`
	Remark    *string    `gorm:"column:remark" json:"remark"`
	CreatedAt time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at,omitempty"`

	// Relations
	Shops          []Shop          `gorm:"foreignKey:MerchantID" json:"shops,omitempty"`
	ElectricMeters []ElectricMeter `gorm:"foreignKey:MerchantID" json:"electric_meters,omitempty"`
	WaterMeters    []WaterMeter    `gorm:"foreignKey:MerchantID" json:"water_meters,omitempty"`
}

func (Merchant) TableName() string { return "biz_merchant" }

// IsEnterprise 是否为企业商户
func (m *Merchant) IsEnterprise() bool {
	return m.MerchantType == MerchantTypeEnterprise
}

// IsPersonal 是否为个人商户
func (m *Merchant) IsPersonal() bool {
	return m.MerchantType == MerchantTypePersonal
}

// IsEnabled 是否启用
func (m *Merchant) IsEnabled() bool {
	return m.Status == MerchantStatusNormal
}
