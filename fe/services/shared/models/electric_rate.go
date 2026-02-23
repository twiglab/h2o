package models

import (
	"time"

	"github.com/shopspring/decimal"
)

// 电费计算模式
const (
	ElectricCalcModeFixed = 1 // 固定单价
	ElectricCalcModeTOU   = 2 // 分时电价
)

// 电费费率范围
const (
	ElectricRateScopeGlobal = 1 // 全局默认
	ElectricRateScopeCustom = 2 // 商户个性化
)

// 费率状态
const (
	ElectricRateStatusDisabled = 0 // 停用
	ElectricRateStatusNormal   = 1 // 正常
)

// 服务费类型
const (
	ServiceFeeTypeFixed = 1 // 固定金额
	ServiceFeeTypeRate  = 2 // 费率百分比
)

// ElectricRate 电费费率表
type ElectricRate struct {
	ID            int64            `gorm:"column:id;primaryKey" json:"id"`
	RateCode      string           `gorm:"column:rate_code" json:"rate_code"`
	RateName      string           `gorm:"column:rate_name" json:"rate_name"`
	MerchantID    *int64           `gorm:"column:merchant_id" json:"merchant_id"`
	Scope         int8             `gorm:"column:scope" json:"scope"`
	CalcMode      int8             `gorm:"column:calc_mode" json:"calc_mode"`
	UnitPrice     *decimal.Decimal `gorm:"column:unit_price" json:"unit_price"`
	EffectiveDate time.Time        `gorm:"column:effective_date" json:"effective_date"`
	ExpireDate    *time.Time       `gorm:"column:expire_date" json:"expire_date"`
	Status        int8             `gorm:"column:status" json:"status"`
	Remark        *string          `gorm:"column:remark" json:"remark"`
	CreatedAt     time.Time        `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     time.Time        `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt     *time.Time       `gorm:"column:deleted_at" json:"deleted_at,omitempty"`

	// Relations
	Merchant    *Merchant                `gorm:"foreignKey:MerchantID" json:"merchant,omitempty"`
	TOUDetails  []ElectricRateTOU        `gorm:"foreignKey:RateID" json:"tou_details,omitempty"`
	ServiceFees []ElectricRateServiceFee `gorm:"foreignKey:RateID" json:"service_fees,omitempty"`
}

func (ElectricRate) TableName() string { return "biz_electric_rate" }

// IsGlobal 是否为全局费率
func (r *ElectricRate) IsGlobal() bool {
	return r.Scope == ElectricRateScopeGlobal && r.MerchantID == nil
}

// IsMerchantRate 是否为商户个性化费率
func (r *ElectricRate) IsMerchantRate() bool {
	return r.Scope == ElectricRateScopeCustom && r.MerchantID != nil
}

// IsFixedPrice 是否为固定单价
func (r *ElectricRate) IsFixedPrice() bool {
	return r.CalcMode == ElectricCalcModeFixed
}

// IsTOU 是否为分时电价
func (r *ElectricRate) IsTOU() bool {
	return r.CalcMode == ElectricCalcModeTOU
}

// IsEnabled 是否启用
func (r *ElectricRate) IsEnabled() bool {
	return r.Status == ElectricRateStatusNormal
}

// IsEffective 是否在有效期内
func (r *ElectricRate) IsEffective() bool {
	now := time.Now()
	if now.Before(r.EffectiveDate) {
		return false
	}
	if r.ExpireDate != nil && now.After(*r.ExpireDate) {
		return false
	}
	return true
}

// ElectricRateTOU 分时电价详情表
type ElectricRateTOU struct {
	ID         int64           `gorm:"column:id;primaryKey" json:"id"`
	RateID     int64           `gorm:"column:rate_id" json:"rate_id"`
	PeriodName string          `gorm:"column:period_name" json:"period_name"`
	StartTime  string          `gorm:"column:start_time" json:"start_time"`
	EndTime    string          `gorm:"column:end_time" json:"end_time"`
	UnitPrice  decimal.Decimal `gorm:"column:unit_price" json:"unit_price"`
	CreatedAt  time.Time       `gorm:"column:created_at" json:"created_at"`
}

func (ElectricRateTOU) TableName() string { return "biz_electric_rate_tou" }

// ElectricRateServiceFee 电费服务费配置表
type ElectricRateServiceFee struct {
	ID        int64           `gorm:"column:id;primaryKey" json:"id"`
	RateID    int64           `gorm:"column:rate_id" json:"rate_id"`
	FeeName   string          `gorm:"column:fee_name" json:"fee_name"`
	FeeType   int8            `gorm:"column:fee_type" json:"fee_type"`
	FeeValue  decimal.Decimal `gorm:"column:fee_value" json:"fee_value"`
	CreatedAt time.Time       `gorm:"column:created_at" json:"created_at"`
}

func (ElectricRateServiceFee) TableName() string { return "biz_electric_rate_service_fee" }

// IsFixed 是否为固定金额
func (f *ElectricRateServiceFee) IsFixed() bool {
	return f.FeeType == ServiceFeeTypeFixed
}

// IsPercentage 是否为百分比
func (f *ElectricRateServiceFee) IsPercentage() bool {
	return f.FeeType == ServiceFeeTypeRate
}

// Calculate 计算服务费
func (f *ElectricRateServiceFee) Calculate(baseAmount decimal.Decimal) decimal.Decimal {
	if f.IsFixed() {
		return f.FeeValue
	}
	// 百分比计算
	return baseAmount.Mul(f.FeeValue).Div(decimal.NewFromInt(100))
}
