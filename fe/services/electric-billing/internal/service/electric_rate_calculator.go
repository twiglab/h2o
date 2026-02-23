package service

import (
	"encoding/json"
	"shared/models"

	"github.com/shopspring/decimal"
)

// ElectricRateCalculator 电费费率计算器
type ElectricRateCalculator struct{}

// NewElectricRateCalculator 创建电费费率计算器
func NewElectricRateCalculator() *ElectricRateCalculator {
	return &ElectricRateCalculator{}
}

// CalculateFee 计算电费
func (c *ElectricRateCalculator) CalculateFee(rate *models.ElectricRate, consumption decimal.Decimal, includeServiceFee bool) decimal.Decimal {
	if consumption.LessThanOrEqual(decimal.Zero) {
		return decimal.Zero
	}

	// 计算基础费用
	var baseFee decimal.Decimal

	switch rate.CalcMode {
	case models.ElectricCalcModeFixed:
		// 固定单价
		baseFee = c.calculateFixedFee(rate, consumption)
	case models.ElectricCalcModeTOU:
		// 分时电价
		baseFee = c.calculateTOUFee(rate, consumption)
	default:
		baseFee = c.calculateFixedFee(rate, consumption)
	}

	if !includeServiceFee {
		return baseFee
	}

	// 计算服务费
	serviceFee := c.calculateServiceFee(rate, baseFee)

	return baseFee.Add(serviceFee)
}

// calculateFixedFee 计算固定单价费用
func (c *ElectricRateCalculator) calculateFixedFee(rate *models.ElectricRate, consumption decimal.Decimal) decimal.Decimal {
	if rate.UnitPrice == nil {
		return decimal.Zero
	}
	return consumption.Mul(*rate.UnitPrice)
}

// calculateTOUFee 计算分时电价费用
// 简化处理：均匀分摊到各时段
func (c *ElectricRateCalculator) calculateTOUFee(rate *models.ElectricRate, consumption decimal.Decimal) decimal.Decimal {
	if len(rate.TOUDetails) == 0 {
		// 没有分时详情，使用固定单价
		return c.calculateFixedFee(rate, consumption)
	}

	// 简化处理：将用量均匀分摊到各时段
	totalFee := decimal.Zero
	periodCount := decimal.NewFromInt(int64(len(rate.TOUDetails)))
	periodConsumption := consumption.Div(periodCount)

	for _, tou := range rate.TOUDetails {
		periodFee := periodConsumption.Mul(tou.UnitPrice)
		totalFee = totalFee.Add(periodFee)
	}

	return totalFee
}

// calculateServiceFee 计算服务费
func (c *ElectricRateCalculator) calculateServiceFee(rate *models.ElectricRate, baseFee decimal.Decimal) decimal.Decimal {
	serviceFeeTotal := decimal.Zero

	for _, sf := range rate.ServiceFees {
		switch sf.FeeType {
		case 1: // 固定金额
			serviceFeeTotal = serviceFeeTotal.Add(sf.FeeValue)
		case 2: // 费率百分比
			rateFee := baseFee.Mul(sf.FeeValue).Div(decimal.NewFromInt(100))
			serviceFeeTotal = serviceFeeTotal.Add(rateFee)
		}
	}

	return serviceFeeTotal
}

// GetUnitPrice 获取单价（用于记录）
func (c *ElectricRateCalculator) GetUnitPrice(rate *models.ElectricRate) decimal.Decimal {
	if rate.UnitPrice != nil {
		return *rate.UnitPrice
	}

	// 分时电价取平均值
	if rate.CalcMode == models.ElectricCalcModeTOU && len(rate.TOUDetails) > 0 {
		total := decimal.Zero
		for _, tou := range rate.TOUDetails {
			total = total.Add(tou.UnitPrice)
		}
		return total.Div(decimal.NewFromInt(int64(len(rate.TOUDetails))))
	}

	return decimal.Zero
}

// CalculateBaseFee 计算基础电费（不含服务费）
func (c *ElectricRateCalculator) CalculateBaseFee(rate *models.ElectricRate, consumption decimal.Decimal) decimal.Decimal {
	if consumption.LessThanOrEqual(decimal.Zero) {
		return decimal.Zero
	}

	switch rate.CalcMode {
	case models.ElectricCalcModeFixed:
		return c.calculateFixedFee(rate, consumption)
	case models.ElectricCalcModeTOU:
		return c.calculateTOUFee(rate, consumption)
	default:
		return c.calculateFixedFee(rate, consumption)
	}
}

// CalculateOnlyServiceFee 只计算服务费
func (c *ElectricRateCalculator) CalculateOnlyServiceFee(rate *models.ElectricRate, baseFee decimal.Decimal) decimal.Decimal {
	return c.calculateServiceFee(rate, baseFee)
}

// BuildTOUDetailJSON 构建分时电价详情JSON（用于扣费记录快照）
func (c *ElectricRateCalculator) BuildTOUDetailJSON(rate *models.ElectricRate, consumption decimal.Decimal) string {
	if rate.CalcMode != models.ElectricCalcModeTOU || len(rate.TOUDetails) == 0 {
		return ""
	}

	// 简化处理：将用量均匀分摊到各时段
	periodCount := decimal.NewFromInt(int64(len(rate.TOUDetails)))
	periodConsumption := consumption.Div(periodCount)

	items := make([]models.TOUDetailItem, 0, len(rate.TOUDetails))
	for _, tou := range rate.TOUDetails {
		periodFee := periodConsumption.Mul(tou.UnitPrice)

		item := models.TOUDetailItem{
			PeriodName:  tou.PeriodName,
			StartTime:   tou.StartTime,
			EndTime:     tou.EndTime,
			Consumption: periodConsumption.InexactFloat64(),
			UnitPrice:   tou.UnitPrice.InexactFloat64(),
			Amount:      periodFee.InexactFloat64(),
		}
		items = append(items, item)
	}

	jsonBytes, err := json.Marshal(items)
	if err != nil {
		return ""
	}

	return string(jsonBytes)
}
