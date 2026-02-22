package service

import (
	"time"

	"admin-server/internal/repository"

	"github.com/shopspring/decimal"
)

// DashboardService 仪表盘服务
type DashboardService struct {
	meterRepo   *repository.MeterRepository
	accountRepo *repository.AccountRepository
	readingRepo *repository.ReadingRepository
	statsRepo   *repository.StatsRepository
}

// NewDashboardService 创建仪表盘服务
func NewDashboardService(
	meterRepo *repository.MeterRepository,
	accountRepo *repository.AccountRepository,
	readingRepo *repository.ReadingRepository,
	statsRepo *repository.StatsRepository,
) *DashboardService {
	return &DashboardService{
		meterRepo:   meterRepo,
		accountRepo: accountRepo,
		readingRepo: readingRepo,
		statsRepo:   statsRepo,
	}
}

// DashboardStats 仪表盘统计
type DashboardStats struct {
	Meter   MeterStats   `json:"meter"`
	Account AccountStats `json:"account"`
	Today   TodayStats   `json:"today"`
}

// MeterStats 电表统计
type MeterStats struct {
	Total  int64 `json:"total"`
	Online int64 `json:"online"`
}

// AccountStats 账户统计
type AccountStats struct {
	Total        int64           `json:"total"`
	ArrearsCount int64           `json:"arrears_count"`
	TotalBalance decimal.Decimal `json:"total_balance"`
	TotalArrears decimal.Decimal `json:"total_arrears"`
}

// TodayStats 今日统计
type TodayStats struct {
	Recharge     decimal.Decimal `json:"recharge"`
	Deduction    decimal.Decimal `json:"deduction"`
	ReadingCount int64           `json:"reading_count"`
}

// GetDashboard 获取仪表盘统计
func (s *DashboardService) GetDashboard() (*DashboardStats, error) {
	today := time.Now().Format("2006-01-02")

	// 电表统计
	meterCount, _ := s.meterRepo.CountByStatus(1)
	meterOnline, _ := s.meterRepo.CountOnline()

	// 账户统计
	accountCount, _ := s.accountRepo.Count()
	arrearsCount, _ := s.accountRepo.CountArrears()
	totalBalance, _ := s.accountRepo.SumBalance()
	totalArrears, _ := s.accountRepo.SumNegativeBalance()

	// 今日统计
	todayRecharge, _ := s.accountRepo.SumRechargeByDate(today)
	todayDeduction, _ := s.accountRepo.SumDeductionByDate(today)
	electricCount, _ := s.readingRepo.CountElectricByDate(today)
	waterCount, _ := s.readingRepo.CountWaterByDate(today)
	todayReadingCount := electricCount + waterCount

	return &DashboardStats{
		Meter: MeterStats{
			Total:  meterCount,
			Online: meterOnline,
		},
		Account: AccountStats{
			Total:        accountCount,
			ArrearsCount: arrearsCount,
			TotalBalance: totalBalance,
			TotalArrears: totalArrears,
		},
		Today: TodayStats{
			Recharge:     todayRecharge,
			Deduction:    todayDeduction,
			ReadingCount: todayReadingCount,
		},
	}, nil
}

// GetConsumptionReport 获取用量报表
func (s *DashboardService) GetConsumptionReport(startDate, endDate, groupBy string) ([]repository.ConsumptionStat, error) {
	return s.statsRepo.GetConsumptionStats(startDate, endDate, groupBy)
}

// GetRevenueReport 获取收入报表
func (s *DashboardService) GetRevenueReport(startDate, endDate, groupBy string) ([]repository.RevenueStat, error) {
	return s.statsRepo.GetRevenueStats(startDate, endDate, groupBy)
}

// CollectionStats 采集统计
type CollectionStats struct {
	Date           string  `json:"date"`
	TotalMeters    int64   `json:"total_meters"`
	CollectedCount int64   `json:"collected_count"`
	Uncollected    int64   `json:"uncollected"`
	AutoCount      int64   `json:"auto_count"`
	ManualCount    int64   `json:"manual_count"`
	SuccessRate    float64 `json:"success_rate"`
}

// GetCollectionStats 获取采集统计
func (s *DashboardService) GetCollectionStats(date string) (*CollectionStats, error) {
	totalMeters, _ := s.meterRepo.CountByStatus(1)
	electricCollected, _ := s.readingRepo.CountElectricCollectedByDate(date)
	waterCollected, _ := s.readingRepo.CountWaterCollectedByDate(date)
	collectedCount := electricCollected + waterCollected

	successRate := float64(0)
	if totalMeters > 0 {
		successRate = float64(collectedCount) / float64(totalMeters) * 100
	}

	return &CollectionStats{
		Date:           date,
		TotalMeters:    totalMeters,
		CollectedCount: collectedCount,
		Uncollected:    totalMeters - collectedCount,
		AutoCount:      0, // 需要查询手工抄表表
		ManualCount:    0, // 需要查询手工抄表表
		SuccessRate:    successRate,
	}, nil
}
