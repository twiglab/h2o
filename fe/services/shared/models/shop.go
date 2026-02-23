package models

import "time"

// 店铺状态
const (
	ShopStatusDisabled = 0 // 停用
	ShopStatusNormal   = 1 // 正常
)

// Shop 店铺表(服务使用者)
type Shop struct {
	ID         int64  `gorm:"column:id;primaryKey" json:"id"`
	ShopNo     string `gorm:"column:shop_no" json:"shop_no"`
	ShopName   string `gorm:"column:shop_name" json:"shop_name"`
	MerchantID int64  `gorm:"column:merchant_id" json:"merchant_id"`

	// 位置信息
	Building *string `gorm:"column:building" json:"building"`
	Floor    *string `gorm:"column:floor" json:"floor"`
	RoomNo   *string `gorm:"column:room_no" json:"room_no"`

	// 联系人
	ContactName  *string `gorm:"column:contact_name" json:"contact_name"`
	ContactPhone *string `gorm:"column:contact_phone" json:"contact_phone"`

	Status    int8       `gorm:"column:status" json:"status"`
	Remark    *string    `gorm:"column:remark" json:"remark"`
	CreatedAt time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at,omitempty"`

	// Relations
	Merchant       *Merchant       `gorm:"foreignKey:MerchantID" json:"merchant,omitempty"`
	ElectricMeters []ElectricMeter `gorm:"foreignKey:ShopID" json:"electric_meters,omitempty"`
	WaterMeters    []WaterMeter    `gorm:"foreignKey:ShopID" json:"water_meters,omitempty"`
}

func (Shop) TableName() string { return "biz_shop" }

// IsEnabled 是否启用
func (s *Shop) IsEnabled() bool {
	return s.Status == ShopStatusNormal
}

// GetLocation 获取位置信息
func (s *Shop) GetLocation() string {
	var location string
	if s.Building != nil {
		location += *s.Building
	}
	if s.Floor != nil {
		if location != "" {
			location += "-"
		}
		location += *s.Floor + "F"
	}
	if s.RoomNo != nil {
		if location != "" {
			location += "-"
		}
		location += *s.RoomNo
	}
	return location
}
