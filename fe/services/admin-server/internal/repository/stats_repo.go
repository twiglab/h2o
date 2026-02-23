package repository

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// StatsRepository 统计数据仓库
type StatsRepository struct {
	db *gorm.DB
}

// NewStatsRepository 创建统计仓库
func NewStatsRepository(db *gorm.DB) *StatsRepository {
	return &StatsRepository{db: db}
}

// ConsumptionStat 用量统计
type ConsumptionStat struct {
	Date                string          `json:"date"`
	ElectricConsumption decimal.Decimal `json:"electric_consumption"`
	WaterConsumption    decimal.Decimal `json:"water_consumption"`
	ElectricCount       int64           `json:"electric_count"`
	WaterCount          int64           `json:"water_count"`
}

// RevenueStat 收入统计
type RevenueStat struct {
	Date            string          `json:"date"`
	RechargeAmount  decimal.Decimal `json:"recharge_amount"`
	DeductionAmount decimal.Decimal `json:"deduction_amount"`
}

// electricConsumptionResult 电量统计结果
type electricConsumptionResult struct {
	Date        string          `gorm:"column:date"`
	Consumption decimal.Decimal `gorm:"column:consumption"`
	Count       int64           `gorm:"column:count"`
}

// waterConsumptionResult 水量统计结果
type waterConsumptionResult struct {
	Date        string          `gorm:"column:date"`
	Consumption decimal.Decimal `gorm:"column:consumption"`
	Count       int64           `gorm:"column:count"`
}

// GetConsumptionStats 获取用量统计（电表和水表）
func (r *StatsRepository) GetConsumptionStats(startDate, endDate, groupBy string) ([]ConsumptionStat, error) {
	var dateFormat string
	if groupBy == "month" {
		dateFormat = "%Y-%m"
	} else {
		dateFormat = "%Y-%m-%d"
	}

	// 分别查询电量和水量
	var electricResults []electricConsumptionResult
	if err := r.db.Raw(`
		SELECT
			DATE_FORMAT(period_end, ?) as date,
			COALESCE(SUM(consumption), 0) as consumption,
			COUNT(*) as count
		FROM fin_electric_consumption
		WHERE period_end BETWEEN ? AND ?
		GROUP BY date
		ORDER BY date
	`, dateFormat, startDate, endDate+" 23:59:59").Scan(&electricResults).Error; err != nil {
		return nil, err
	}

	var waterResults []waterConsumptionResult
	if err := r.db.Raw(`
		SELECT
			DATE_FORMAT(period_end, ?) as date,
			COALESCE(SUM(consumption), 0) as consumption,
			COUNT(*) as count
		FROM fin_water_consumption
		WHERE period_end BETWEEN ? AND ?
		GROUP BY date
		ORDER BY date
	`, dateFormat, startDate, endDate+" 23:59:59").Scan(&waterResults).Error; err != nil {
		return nil, err
	}

	// 合并结果
	dateMap := make(map[string]*ConsumptionStat)

	for _, e := range electricResults {
		if _, ok := dateMap[e.Date]; !ok {
			dateMap[e.Date] = &ConsumptionStat{Date: e.Date}
		}
		dateMap[e.Date].ElectricConsumption = e.Consumption
		dateMap[e.Date].ElectricCount = e.Count
	}

	for _, w := range waterResults {
		if _, ok := dateMap[w.Date]; !ok {
			dateMap[w.Date] = &ConsumptionStat{Date: w.Date}
		}
		dateMap[w.Date].WaterConsumption = w.Consumption
		dateMap[w.Date].WaterCount = w.Count
	}

	// 转换为切片并排序
	stats := make([]ConsumptionStat, 0, len(dateMap))
	for _, v := range dateMap {
		stats = append(stats, *v)
	}

	// 按日期排序
	for i := 0; i < len(stats)-1; i++ {
		for j := i + 1; j < len(stats); j++ {
			if stats[i].Date > stats[j].Date {
				stats[i], stats[j] = stats[j], stats[i]
			}
		}
	}

	return stats, nil
}

// rechargeResult 充值统计结果
type rechargeResult struct {
	Date   string          `gorm:"column:date"`
	Amount decimal.Decimal `gorm:"column:amount"`
}

// deductionResult 扣费统计结果
type deductionResult struct {
	Date   string          `gorm:"column:date"`
	Amount decimal.Decimal `gorm:"column:amount"`
}

// GetRevenueStats 获取收入统计
func (r *StatsRepository) GetRevenueStats(startDate, endDate, groupBy string) ([]RevenueStat, error) {
	var dateFormat string
	if groupBy == "month" {
		dateFormat = "%Y-%m"
	} else {
		dateFormat = "%Y-%m-%d"
	}

	// 查询充值
	var rechargeResults []rechargeResult
	if err := r.db.Raw(`
		SELECT
			DATE_FORMAT(created_at, ?) as date,
			COALESCE(SUM(amount), 0) as amount
		FROM biz_recharge
		WHERE created_at BETWEEN ? AND ?
		GROUP BY date
		ORDER BY date
	`, dateFormat, startDate, endDate+" 23:59:59").Scan(&rechargeResults).Error; err != nil {
		return nil, err
	}

	// 查询电费扣费
	var electricDeductionResults []deductionResult
	if err := r.db.Raw(`
		SELECT
			DATE_FORMAT(deduction_time, ?) as date,
			COALESCE(SUM(amount), 0) as amount
		FROM biz_electric_deduction
		WHERE deduction_time BETWEEN ? AND ?
		GROUP BY date
		ORDER BY date
	`, dateFormat, startDate, endDate+" 23:59:59").Scan(&electricDeductionResults).Error; err != nil {
		return nil, err
	}

	// 查询水费扣费
	var waterDeductionResults []deductionResult
	if err := r.db.Raw(`
		SELECT
			DATE_FORMAT(deduction_time, ?) as date,
			COALESCE(SUM(amount), 0) as amount
		FROM biz_water_deduction
		WHERE deduction_time BETWEEN ? AND ?
		GROUP BY date
		ORDER BY date
	`, dateFormat, startDate, endDate+" 23:59:59").Scan(&waterDeductionResults).Error; err != nil {
		return nil, err
	}

	// 合并结果
	dateMap := make(map[string]*RevenueStat)

	for _, r := range rechargeResults {
		if _, ok := dateMap[r.Date]; !ok {
			dateMap[r.Date] = &RevenueStat{Date: r.Date}
		}
		dateMap[r.Date].RechargeAmount = r.Amount
	}

	for _, d := range electricDeductionResults {
		if _, ok := dateMap[d.Date]; !ok {
			dateMap[d.Date] = &RevenueStat{Date: d.Date}
		}
		dateMap[d.Date].DeductionAmount = dateMap[d.Date].DeductionAmount.Add(d.Amount)
	}

	for _, d := range waterDeductionResults {
		if _, ok := dateMap[d.Date]; !ok {
			dateMap[d.Date] = &RevenueStat{Date: d.Date}
		}
		dateMap[d.Date].DeductionAmount = dateMap[d.Date].DeductionAmount.Add(d.Amount)
	}

	// 转换为切片并排序
	stats := make([]RevenueStat, 0, len(dateMap))
	for _, v := range dateMap {
		stats = append(stats, *v)
	}

	// 按日期排序
	for i := 0; i < len(stats)-1; i++ {
		for j := i + 1; j < len(stats); j++ {
			if stats[i].Date > stats[j].Date {
				stats[i], stats[j] = stats[j], stats[i]
			}
		}
	}

	return stats, nil
}
