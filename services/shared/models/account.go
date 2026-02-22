package models

import (
	"time"

	"github.com/shopspring/decimal"
)

// 账户状态
const (
	AccountStatusFrozen  = 0 // 冻结
	AccountStatusNormal  = 1 // 正常
	AccountStatusArrears = 2 // 欠费
)

// 支付方式
const (
	PaymentMethodCash   = 1 // 现金
	PaymentMethodBank   = 2 // 银行转账
	PaymentMethodWechat = 3 // 微信
	PaymentMethodAlipay = 4 // 支付宝
)

// 扣费状态
const (
	DeductionStatusSuccess = 1 // 成功
	DeductionStatusPartial = 2 // 部分扣费
	DeductionStatusFailed  = 3 // 失败
)

// Account 预付费账户
// 一个商户可以有多个账户，一个账户可以为多个店铺的表计付费
type Account struct {
	ID          int64   `gorm:"column:id;primaryKey" json:"id"`
	AccountNo   string  `gorm:"column:account_no" json:"account_no"`
	AccountName *string `gorm:"column:account_name" json:"account_name"`
	MerchantID  int64   `gorm:"column:merchant_id" json:"merchant_id"`

	// 账户余额
	Balance decimal.Decimal `gorm:"column:balance" json:"balance"`

	// 累计统计
	TotalRecharge    decimal.Decimal `gorm:"column:total_recharge" json:"total_recharge"`
	TotalConsumption decimal.Decimal `gorm:"column:total_consumption" json:"total_consumption"`

	Status    int8      `gorm:"column:status" json:"status"`
	Remark    *string   `gorm:"column:remark" json:"remark"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`

	// Relations
	Merchant       *Merchant       `gorm:"foreignKey:MerchantID" json:"merchant,omitempty"`
	ElectricMeters []ElectricMeter `gorm:"foreignKey:AccountID" json:"electric_meters,omitempty"`
	WaterMeters    []WaterMeter    `gorm:"foreignKey:AccountID" json:"water_meters,omitempty"`
}

func (Account) TableName() string { return "biz_account" }

// IsNormal 是否正常
func (a *Account) IsNormal() bool {
	return a.Status == AccountStatusNormal
}

// IsFrozen 是否冻结
func (a *Account) IsFrozen() bool {
	return a.Status == AccountStatusFrozen
}

// IsArrears 是否欠费
func (a *Account) IsArrears() bool {
	return a.Status == AccountStatusArrears
}

// Recharge 充值记录
type Recharge struct {
	ID         int64  `gorm:"column:id;primaryKey" json:"id"`
	RechargeNo string `gorm:"column:recharge_no" json:"recharge_no"`
	AccountID  int64  `gorm:"column:account_id" json:"account_id"`

	// 充值信息
	Amount        decimal.Decimal `gorm:"column:amount" json:"amount"`
	BalanceBefore decimal.Decimal `gorm:"column:balance_before" json:"balance_before"`
	BalanceAfter  decimal.Decimal `gorm:"column:balance_after" json:"balance_after"`

	// 支付信息
	PaymentMethod int8    `gorm:"column:payment_method" json:"payment_method"`
	PaymentNo     *string `gorm:"column:payment_no" json:"payment_no"`

	// 操作信息
	OperatorID   int64   `gorm:"column:operator_id" json:"operator_id"`
	OperatorName *string `gorm:"column:operator_name" json:"operator_name"`

	Remark    *string   `gorm:"column:remark" json:"remark"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`

	// Relations
	Account *Account `gorm:"foreignKey:AccountID" json:"account,omitempty"`
}

func (Recharge) TableName() string { return "biz_recharge" }

// TOUDetailItem 分时计费明细项
type TOUDetailItem struct {
	PeriodName  string  `json:"period_name"` // 时段名称
	StartTime   string  `json:"start_time"`  // 开始时间
	EndTime     string  `json:"end_time"`    // 结束时间
	Consumption float64 `json:"consumption"` // 该时段用量
	UnitPrice   float64 `json:"unit_price"`  // 该时段单价
	Amount      float64 `json:"amount"`      // 该时段费用
}

// ElectricDeduction 电费扣费记录（保存完整快照，避免关联数据修改影响历史记录）
type ElectricDeduction struct {
	ID          int64  `gorm:"column:id;primaryKey" json:"id"`
	DeductionNo string `gorm:"column:deduction_no" json:"deduction_no"`

	// 商户快照
	MerchantID   int64  `gorm:"column:merchant_id" json:"merchant_id"`
	MerchantNo   string `gorm:"column:merchant_no" json:"merchant_no"`
	MerchantName string `gorm:"column:merchant_name" json:"merchant_name"`

	// 店铺快照
	ShopID   *int64  `gorm:"column:shop_id" json:"shop_id"`
	ShopNo   *string `gorm:"column:shop_no" json:"shop_no"`
	ShopName *string `gorm:"column:shop_name" json:"shop_name"`

	// 账户快照
	AccountID   int64   `gorm:"column:account_id" json:"account_id"`
	AccountNo   string  `gorm:"column:account_no" json:"account_no"`
	AccountName *string `gorm:"column:account_name" json:"account_name"`

	// 电表快照
	MeterID    int64           `gorm:"column:meter_id" json:"meter_id"`
	MeterNo    string          `gorm:"column:meter_no" json:"meter_no"`
	Multiplier decimal.Decimal `gorm:"column:multiplier" json:"multiplier"`

	// 用量记录关联
	ConsumptionID *int64 `gorm:"column:consumption_id" json:"consumption_id"`

	// 读数信息
	StartReading decimal.Decimal `gorm:"column:start_reading" json:"start_reading"`
	EndReading   decimal.Decimal `gorm:"column:end_reading" json:"end_reading"`
	Consumption  decimal.Decimal `gorm:"column:consumption" json:"consumption"`

	// 计费周期
	PeriodStart time.Time `gorm:"column:period_start" json:"period_start"`
	PeriodEnd   time.Time `gorm:"column:period_end" json:"period_end"`

	// 费率快照
	RateID    *int64           `gorm:"column:rate_id" json:"rate_id"`
	RateCode  *string          `gorm:"column:rate_code" json:"rate_code"`
	RateName  *string          `gorm:"column:rate_name" json:"rate_name"`
	CalcMode  int8             `gorm:"column:calc_mode" json:"calc_mode"`
	UnitPrice *decimal.Decimal `gorm:"column:unit_price" json:"unit_price"`
	TOUDetail *string          `gorm:"column:tou_detail;type:json" json:"tou_detail"`

	// 费用明细
	BaseAmount    decimal.Decimal `gorm:"column:base_amount" json:"base_amount"`
	ServiceAmount decimal.Decimal `gorm:"column:service_amount" json:"service_amount"`
	Amount        decimal.Decimal `gorm:"column:amount" json:"amount"`

	// 余额变化
	BalanceBefore decimal.Decimal `gorm:"column:balance_before" json:"balance_before"`
	BalanceAfter  decimal.Decimal `gorm:"column:balance_after" json:"balance_after"`

	DeductionTime time.Time `gorm:"column:deduction_time" json:"deduction_time"`
	Status        int8      `gorm:"column:status" json:"status"`
	Remark        *string   `gorm:"column:remark" json:"remark"`
	CreatedAt     time.Time `gorm:"column:created_at" json:"created_at"`
}

func (ElectricDeduction) TableName() string { return "biz_electric_deduction" }

// IsSuccess 是否扣费成功
func (d *ElectricDeduction) IsSuccess() bool {
	return d.Status == DeductionStatusSuccess
}

// IsFixed 是否固定单价模式
func (d *ElectricDeduction) IsFixed() bool {
	return d.CalcMode == ElectricCalcModeFixed
}

// IsTOU 是否分时电价模式
func (d *ElectricDeduction) IsTOU() bool {
	return d.CalcMode == ElectricCalcModeTOU
}

// WaterDeduction 水费扣费记录（保存完整快照，避免关联数据修改影响历史记录）
type WaterDeduction struct {
	ID          int64  `gorm:"column:id;primaryKey" json:"id"`
	DeductionNo string `gorm:"column:deduction_no" json:"deduction_no"`

	// 商户快照
	MerchantID   int64  `gorm:"column:merchant_id" json:"merchant_id"`
	MerchantNo   string `gorm:"column:merchant_no" json:"merchant_no"`
	MerchantName string `gorm:"column:merchant_name" json:"merchant_name"`

	// 店铺快照
	ShopID   *int64  `gorm:"column:shop_id" json:"shop_id"`
	ShopNo   *string `gorm:"column:shop_no" json:"shop_no"`
	ShopName *string `gorm:"column:shop_name" json:"shop_name"`

	// 账户快照
	AccountID   int64   `gorm:"column:account_id" json:"account_id"`
	AccountNo   string  `gorm:"column:account_no" json:"account_no"`
	AccountName *string `gorm:"column:account_name" json:"account_name"`

	// 水表快照
	MeterID    int64           `gorm:"column:meter_id" json:"meter_id"`
	MeterNo    string          `gorm:"column:meter_no" json:"meter_no"`
	Multiplier decimal.Decimal `gorm:"column:multiplier" json:"multiplier"`

	// 用量记录关联
	ConsumptionID *int64 `gorm:"column:consumption_id" json:"consumption_id"`

	// 读数信息
	StartReading decimal.Decimal `gorm:"column:start_reading" json:"start_reading"`
	EndReading   decimal.Decimal `gorm:"column:end_reading" json:"end_reading"`
	Consumption  decimal.Decimal `gorm:"column:consumption" json:"consumption"`

	// 计费周期
	PeriodStart time.Time `gorm:"column:period_start" json:"period_start"`
	PeriodEnd   time.Time `gorm:"column:period_end" json:"period_end"`

	// 费率快照
	RateID    *int64           `gorm:"column:rate_id" json:"rate_id"`
	RateCode  *string          `gorm:"column:rate_code" json:"rate_code"`
	RateName  *string          `gorm:"column:rate_name" json:"rate_name"`
	UnitPrice *decimal.Decimal `gorm:"column:unit_price" json:"unit_price"`

	// 费用明细
	BaseAmount    decimal.Decimal `gorm:"column:base_amount" json:"base_amount"`
	ServiceAmount decimal.Decimal `gorm:"column:service_amount" json:"service_amount"`
	Amount        decimal.Decimal `gorm:"column:amount" json:"amount"`

	// 余额变化
	BalanceBefore decimal.Decimal `gorm:"column:balance_before" json:"balance_before"`
	BalanceAfter  decimal.Decimal `gorm:"column:balance_after" json:"balance_after"`

	DeductionTime time.Time `gorm:"column:deduction_time" json:"deduction_time"`
	Status        int8      `gorm:"column:status" json:"status"`
	Remark        *string   `gorm:"column:remark" json:"remark"`
	CreatedAt     time.Time `gorm:"column:created_at" json:"created_at"`
}

func (WaterDeduction) TableName() string { return "biz_water_deduction" }

// IsSuccess 是否扣费成功
func (d *WaterDeduction) IsSuccess() bool {
	return d.Status == DeductionStatusSuccess
}
