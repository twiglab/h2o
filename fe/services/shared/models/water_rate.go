package models

import (
	"time"

	"github.com/shopspring/decimal"
)

// 水费计算模式
const (
	WaterCalcModeFixed  = 1 // 固定单价
	WaterCalcModeTiered = 2 // 阶梯水价
)

// 水费费率范围
const (
	WaterRateScopeGlobal = 1 // 全局默认
	WaterRateScopeCustom = 2 // 商户个性化
)

// 费率状态
const (
	WaterRateStatusDisabled = 0 // 停用
	WaterRateStatusNormal   = 1 // 正常
)

// WaterRate 水费费率表
type WaterRate struct {
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
	Merchant    *Merchant             `gorm:"foreignKey:MerchantID" json:"merchant,omitempty"`
	TieredRates []WaterRateTiered     `gorm:"foreignKey:RateID" json:"tiered_rates,omitempty"`
	ServiceFees []WaterRateServiceFee `gorm:"foreignKey:RateID" json:"service_fees,omitempty"`
}

func (WaterRate) TableName() string { return "biz_water_rate" }

// IsGlobal 是否为全局费率
func (r *WaterRate) IsGlobal() bool {
	return r.Scope == WaterRateScopeGlobal && r.MerchantID == nil
}

// IsMerchantRate 是否为商户个性化费率
func (r *WaterRate) IsMerchantRate() bool {
	return r.Scope == WaterRateScopeCustom && r.MerchantID != nil
}

// IsFixedPrice 是否为固定单价
func (r *WaterRate) IsFixedPrice() bool {
	return r.CalcMode == WaterCalcModeFixed
}

// IsTiered 是否为阶梯水价
func (r *WaterRate) IsTiered() bool {
	return r.CalcMode == WaterCalcModeTiered
}

// IsEnabled 是否启用
func (r *WaterRate) IsEnabled() bool {
	return r.Status == WaterRateStatusNormal
}

// IsEffective 是否在有效期内
func (r *WaterRate) IsEffective() bool {
	now := time.Now()
	if now.Before(r.EffectiveDate) {
		return false
	}
	if r.ExpireDate != nil && now.After(*r.ExpireDate) {
		return false
	}
	return true
}

// WaterRateTiered 阶梯水价详情表
type WaterRateTiered struct {
	ID         int64            `gorm:"column:id;primaryKey" json:"id"`
	RateID     int64            `gorm:"column:rate_id" json:"rate_id"`
	TierLevel  int              `gorm:"column:tier_level" json:"tier_level"`
	StartValue decimal.Decimal  `gorm:"column:start_value" json:"start_value"`
	EndValue   *decimal.Decimal `gorm:"column:end_value" json:"end_value"`
	UnitPrice  decimal.Decimal  `gorm:"column:unit_price" json:"unit_price"`
	CreatedAt  time.Time        `gorm:"column:created_at" json:"created_at"`
}

func (WaterRateTiered) TableName() string { return "biz_water_rate_tiered" }

// IsUnlimited 是否为无上限阶梯
func (t *WaterRateTiered) IsUnlimited() bool {
	return t.EndValue == nil
}

// Contains 是否包含指定用量
func (t *WaterRateTiered) Contains(value decimal.Decimal) bool {
	if value.LessThan(t.StartValue) {
		return false
	}
	if t.EndValue != nil && value.GreaterThanOrEqual(*t.EndValue) {
		return false
	}
	return true
}

// WaterRateServiceFee 水费服务费配置表
type WaterRateServiceFee struct {
	ID        int64           `gorm:"column:id;primaryKey" json:"id"`
	RateID    int64           `gorm:"column:rate_id" json:"rate_id"`
	FeeName   string          `gorm:"column:fee_name" json:"fee_name"`
	FeeType   int8            `gorm:"column:fee_type" json:"fee_type"`
	FeeValue  decimal.Decimal `gorm:"column:fee_value" json:"fee_value"`
	CreatedAt time.Time       `gorm:"column:created_at" json:"created_at"`
}

func (WaterRateServiceFee) TableName() string { return "biz_water_rate_service_fee" }

// IsFixed 是否为固定金额
func (f *WaterRateServiceFee) IsFixed() bool {
	return f.FeeType == ServiceFeeTypeFixed
}

// IsPercentage 是否为百分比
func (f *WaterRateServiceFee) IsPercentage() bool {
	return f.FeeType == ServiceFeeTypeRate
}

// Calculate 计算服务费
func (f *WaterRateServiceFee) Calculate(baseAmount decimal.Decimal) decimal.Decimal {
	if f.IsFixed() {
		return f.FeeValue
	}
	// 百分比计算
	return baseAmount.Mul(f.FeeValue).Div(decimal.NewFromInt(100))
}
