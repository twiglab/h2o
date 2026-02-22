package service

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"electric-billing/internal/modbus"
	"electric-billing/internal/repository"
	"electric-billing/internal/storage"
	"shared/models"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// 错误定义
var (
	ErrMeterNotFound       = errors.New("meter not found")
	ErrMeterNoRate         = errors.New("meter has no rate")
	ErrAccountNotFound     = errors.New("account not found")
	ErrReadingTooLow       = errors.New("reading value lower than current")
	ErrInvalidReading      = errors.New("invalid reading data")
	ErrConsumptionNegative = errors.New("consumption is negative")
)

// BillingService 计费服务
type BillingService struct {
	db               *gorm.DB
	meterRepo        *repository.MeterRepository
	electricRateRepo *repository.ElectricRateRepository
	accountRepo      *repository.AccountRepository
	readingRepo      *repository.ReadingRepository
	calculator       *ElectricRateCalculator
	wal              *storage.WAL
	fileStore        *storage.FileStore
	logger           *zap.Logger
}

// NewBillingService 创建计费服务
func NewBillingService(
	db *gorm.DB,
	meterRepo *repository.MeterRepository,
	electricRateRepo *repository.ElectricRateRepository,
	accountRepo *repository.AccountRepository,
	readingRepo *repository.ReadingRepository,
	wal *storage.WAL,
	fileStore *storage.FileStore,
	logger *zap.Logger,
) *BillingService {
	return &BillingService{
		db:               db,
		meterRepo:        meterRepo,
		electricRateRepo: electricRateRepo,
		accountRepo:      accountRepo,
		readingRepo:      readingRepo,
		calculator:       NewElectricRateCalculator(),
		wal:              wal,
		fileStore:        fileStore,
		logger:           logger,
	}
}

// ProcessReading 处理电表读数（实现 mqtt.ReadingProcessor 接口）
func (s *BillingService) ProcessReading(dtuID string, meterAddr byte, rawData []byte, readingTime time.Time) error {
	// 解析Modbus帧
	frame, err := modbus.ParseFrame(rawData)
	if err != nil {
		return fmt.Errorf("parse modbus frame: %w", err)
	}

	// 提取电能读数
	readingValue, err := modbus.ExtractReading(frame)
	if err != nil {
		// 尝试浮点格式
		readingValue, err = modbus.ExtractReadingFloat32(frame)
		if err != nil {
			return fmt.Errorf("extract reading: %w", err)
		}
	}

	// 构建通信地址 (DTU ID + 从站地址)
	commAddr := fmt.Sprintf("%s_%02X", dtuID, meterAddr)

	s.logger.Info("processing reading",
		zap.String("dtu_id", dtuID),
		zap.Uint8("meter_addr", meterAddr),
		zap.String("comm_addr", commAddr),
		zap.String("reading_value", readingValue.String()))

	// 根据通信地址查找电表
	meter, err := s.meterRepo.FindByCommAddr(commAddr)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 也尝试用纯从站地址查找
			addrStr := fmt.Sprintf("%02X", meterAddr)
			meter, err = s.meterRepo.FindByCommAddr(addrStr)
			if err != nil {
				return fmt.Errorf("%w: comm_addr=%s", ErrMeterNotFound, commAddr)
			}
		} else {
			return fmt.Errorf("find meter: %w", err)
		}
	}

	// 调用核心处理逻辑
	return s.processReadingForMeter(meter, readingValue, readingTime)
}

// processReadingForMeter 处理指定电表的读数
func (s *BillingService) processReadingForMeter(meter *models.ElectricMeter, readingValue decimal.Decimal, readingTime time.Time) error {
	// 验证读数不能小于当前读数
	if readingValue.LessThan(meter.CurrentReading) {
		s.logger.Warn("reading value lower than current",
			zap.String("meter_no", meter.MeterNo),
			zap.String("current", meter.CurrentReading.String()),
			zap.String("new", readingValue.String()))
		return ErrReadingTooLow
	}

	// 1. 先写入 WAL (Write-Ahead Log)
	var walEntry *storage.WalEntry
	var commAddr string
	if meter.CommAddr != nil {
		commAddr = *meter.CommAddr
	}
	if s.wal != nil {
		var err error
		walEntry, err = s.wal.WriteReading(meter.MeterNo, meter.ID, commAddr, readingValue, readingTime)
		if err != nil {
			s.logger.Error("failed to write WAL entry", zap.Error(err))
		}
	}

	// 使用事务处理
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// 2. 锁定电表记录
		lockedMeter, err := s.meterRepo.FindByIDForUpdate(tx, meter.ID)
		if err != nil {
			return fmt.Errorf("lock meter: %w", err)
		}

		// 再次验证读数
		if readingValue.LessThan(lockedMeter.CurrentReading) {
			return ErrReadingTooLow
		}

		// 3. 创建电表读数记录
		electricReading := &models.ElectricReading{
			MeterID:      meter.ID,
			MeterNo:      meter.MeterNo,
			ReadingValue: readingValue,
			ReadingTime:  readingTime,
			CollectType:  models.CollectTypeAuto,
			Status:       models.ReadingStatusNormal,
			CreatedAt:    time.Now(),
		}
		if err := s.readingRepo.CreateReading(tx, electricReading); err != nil {
			return fmt.Errorf("create reading: %w", err)
		}

		// 4. 保存读数到本地文件
		if s.fileStore != nil {
			readingRecord := &storage.ReadingRecord{
				MeterNo:      meter.MeterNo,
				MeterID:      meter.ID,
				ReadingValue: readingValue,
				ReadingTime:  readingTime,
				CollectType:  models.CollectTypeAuto,
				CreatedAt:    time.Now(),
			}
			if err := s.fileStore.SaveReading(readingRecord); err != nil {
				s.logger.Warn("failed to save reading to file", zap.Error(err))
			}
		}

		// 5. 更新电表当前读数
		if err := s.meterRepo.UpdateReading(tx, meter.ID, readingValue, readingTime); err != nil {
			return fmt.Errorf("update meter reading: %w", err)
		}

		// 6. 查找上次读数计算用量
		lastReading, err := s.readingRepo.FindLastReadingBefore(meter.ID, readingTime)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// 首次读数，无法计算用量
				s.logger.Info("first reading, skip consumption calculation",
					zap.String("meter_no", meter.MeterNo))
				return nil
			}
			return fmt.Errorf("find last reading: %w", err)
		}

		// 7. 计算用量
		rawConsumption := readingValue.Sub(lastReading.ReadingValue)
		consumption := rawConsumption.Mul(meter.Multiplier)

		if consumption.LessThan(decimal.Zero) {
			return ErrConsumptionNegative
		}

		// 零用量不处理
		if consumption.Equal(decimal.Zero) {
			s.logger.Debug("zero consumption, skip",
				zap.String("meter_no", meter.MeterNo))
			return nil
		}

		// 8. 创建电表用量记录
		consumptionRecord := &models.ElectricConsumption{
			MeterID:        meter.ID,
			StartReadingID: lastReading.ID,
			EndReadingID:   electricReading.ID,
			StartReading:   lastReading.ReadingValue,
			EndReading:     readingValue,
			Consumption:    consumption,
			PeriodStart:    lastReading.ReadingTime,
			PeriodEnd:      readingTime,
			Status:         models.ConsumptionStatusPending,
			CreatedAt:      time.Now(),
		}
		if err := s.readingRepo.CreateConsumption(tx, consumptionRecord); err != nil {
			return fmt.Errorf("create consumption: %w", err)
		}

		// 9. 获取适用费率
		rate, err := s.electricRateRepo.GetApplicableRate(meter.ID, meter.RateID)
		if err != nil {
			return fmt.Errorf("get applicable rate: %w", err)
		}

		// 10. 计算费用并扣费
		return s.processDeduction(tx, meter, consumptionRecord, rate)
	})

	// 11. 标记 WAL 条目已完成
	if walEntry != nil && s.wal != nil {
		if err == nil {
			s.wal.MarkCompleted(walEntry.ID)
		} else {
			s.wal.MarkFailed(walEntry.ID, err.Error())
		}
	}

	return err
}

// processDeduction 处理扣费
func (s *BillingService) processDeduction(tx *gorm.DB, meter *models.ElectricMeter, consumption *models.ElectricConsumption, rate *models.ElectricRate) error {
	if rate == nil {
		return ErrMeterNoRate
	}

	// 计算费用 (暂不计算服务费)
	fee := s.calculator.CalculateFee(rate, consumption.Consumption, false)
	if fee.LessThanOrEqual(decimal.Zero) {
		return nil
	}

	// 检查电表是否关联了账户
	if meter.AccountID == nil {
		s.logger.Warn("meter has no account", zap.String("meter_no", meter.MeterNo))
		return ErrAccountNotFound
	}

	// 获取电表的完整关联数据（用于快照）
	meterWithRelations, err := s.meterRepo.FindByIDWithRelations(meter.ID)
	if err != nil {
		s.logger.Warn("failed to load meter relations, using basic info", zap.Error(err))
		meterWithRelations = meter
	}

	// 锁定账户（通过电表的 AccountID 查找）
	account, err := s.accountRepo.FindByIDForUpdate(tx, *meter.AccountID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrAccountNotFound
		}
		return fmt.Errorf("find account: %w", err)
	}

	// 执行扣费
	balanceBefore := account.Balance
	var balanceAfter decimal.Decimal
	var deductionStatus int8
	var accountStatus int8
	newConsumption := account.TotalConsumption.Add(fee)

	if account.Balance.GreaterThanOrEqual(fee) {
		// 余额充足，全额扣费
		balanceAfter = account.Balance.Sub(fee)
		deductionStatus = models.DeductionStatusSuccess
		accountStatus = account.Status
	} else if account.Balance.GreaterThan(decimal.Zero) {
		// 余额不足，部分扣费
		balanceAfter = decimal.Zero
		deductionStatus = models.DeductionStatusPartial
		accountStatus = models.AccountStatusArrears
	} else {
		// 无余额
		balanceAfter = account.Balance.Sub(fee) // 余额变负
		deductionStatus = models.DeductionStatusFailed
		accountStatus = models.AccountStatusArrears
	}

	// 更新账户
	if err := s.accountRepo.UpdateBalance(tx, account.ID, balanceAfter, newConsumption, accountStatus); err != nil {
		return fmt.Errorf("update account: %w", err)
	}

	unitPrice := s.calculator.GetUnitPrice(rate)

	// 构建扣费记录（包含完整快照信息）
	deduction := &models.ElectricDeduction{
		DeductionNo: generateNo("D"),

		// 商户快照
		MerchantID:   meter.MerchantID,
		MerchantNo:   "",
		MerchantName: "",

		// 账户快照
		AccountID:   account.ID,
		AccountNo:   account.AccountNo,
		AccountName: account.AccountName,

		// 电表快照
		MeterID:    meter.ID,
		MeterNo:    meter.MeterNo,
		Multiplier: meter.Multiplier,

		// 用量记录关联
		ConsumptionID: &consumption.ID,

		// 读数信息
		StartReading: consumption.StartReading,
		EndReading:   consumption.EndReading,
		Consumption:  consumption.Consumption,

		// 计费周期
		PeriodStart: consumption.PeriodStart,
		PeriodEnd:   consumption.PeriodEnd,

		// 费率快照
		RateID:   &rate.ID,
		RateCode: &rate.RateCode,
		RateName: &rate.RateName,
		CalcMode: rate.CalcMode,

		// 费用明细
		BaseAmount:    fee,
		ServiceAmount: decimal.Zero, // 服务费暂不计算
		Amount:        fee,

		// 余额变化
		BalanceBefore: balanceBefore,
		BalanceAfter:  balanceAfter,

		DeductionTime: time.Now(),
		Status:        deductionStatus,
		CreatedAt:     time.Now(),
	}

	// 填充商户快照
	if meterWithRelations.Merchant != nil {
		deduction.MerchantNo = meterWithRelations.Merchant.MerchantNo
		deduction.MerchantName = meterWithRelations.Merchant.MerchantName
	}

	// 填充店铺快照
	if meterWithRelations.Shop != nil {
		deduction.ShopID = &meterWithRelations.Shop.ID
		deduction.ShopNo = &meterWithRelations.Shop.ShopNo
		deduction.ShopName = &meterWithRelations.Shop.ShopName
	}

	// 填充费率单价
	if rate.CalcMode == models.ElectricCalcModeFixed && rate.UnitPrice != nil {
		deduction.UnitPrice = rate.UnitPrice
	} else {
		deduction.UnitPrice = &unitPrice
	}

	// 分时电价记录详情（JSON格式）
	if rate.CalcMode == models.ElectricCalcModeTOU && len(rate.TOUDetails) > 0 {
		touDetailJSON := s.calculator.BuildTOUDetailJSON(rate, consumption.Consumption)
		if touDetailJSON != "" {
			deduction.TOUDetail = &touDetailJSON
		}
	}

	if err := s.accountRepo.CreateElectricDeduction(tx, deduction); err != nil {
		return fmt.Errorf("create deduction: %w", err)
	}

	// 保存扣费记录到本地文件
	if s.fileStore != nil {
		deductionRecord := &storage.DeductionRecord{
			DeductionNo:   deduction.DeductionNo,
			MeterNo:       meter.MeterNo,
			MeterID:       meter.ID,
			AccountID:     account.ID,
			ConsumptionID: consumption.ID,
			Consumption:   consumption.Consumption,
			UnitPrice:     unitPrice,
			Amount:        fee,
			BalanceBefore: balanceBefore,
			BalanceAfter:  balanceAfter,
			DeductionType: models.MeterTypeElectric,
			Status:        int16(deductionStatus),
			CreatedAt:     time.Now(),
		}
		if err := s.fileStore.SaveDeduction(deductionRecord); err != nil {
			s.logger.Warn("failed to save deduction to file", zap.Error(err))
		}
	}

	// 更新用量记录状态
	if err := s.readingRepo.UpdateConsumptionStatus(tx, consumption.ID, models.ConsumptionStatusDeducted); err != nil {
		return fmt.Errorf("update consumption status: %w", err)
	}

	s.logger.Info("deduction processed",
		zap.String("meter_no", meter.MeterNo),
		zap.String("deduction_no", deduction.DeductionNo),
		zap.String("consumption", consumption.Consumption.String()),
		zap.String("amount", fee.String()),
		zap.String("balance_before", balanceBefore.String()),
		zap.String("balance_after", balanceAfter.String()),
		zap.Int8("status", deductionStatus))

	// 检查余额预警
	s.checkBalanceWarning(account.ID, balanceAfter)

	return nil
}

// checkBalanceWarning 检查余额预警
func (s *BillingService) checkBalanceWarning(accountID int64, balance decimal.Decimal) {
	var alertLevel string
	var message string

	// 低余额预警阈值设为100元
	threshold := decimal.NewFromInt(100)

	if balance.LessThan(decimal.Zero) {
		alertLevel = "critical"
		message = fmt.Sprintf("账户余额为负: %s 元", balance.String())
	} else if balance.Equal(decimal.Zero) {
		alertLevel = "warning"
		message = "账户余额已为0"
	} else if balance.LessThan(threshold) {
		alertLevel = "info"
		message = fmt.Sprintf("账户余额 %s 元，低于预警阈值", balance.String())
	} else {
		return
	}

	s.logger.Warn("balance warning",
		zap.Int64("account_id", accountID),
		zap.String("level", alertLevel),
		zap.String("message", message))

	// TODO: 发送预警通知
}

// generateNo 生成业务编号
// 格式: 前缀 + 日期 + 6位随机数
func generateNo(prefix string) string {
	dateStr := time.Now().Format("20060102")
	randomStr := fmt.Sprintf("%06d", rand.Intn(1000000))
	return prefix + dateStr + randomStr
}

// ProcessReadingByMeterNo 根据电表编号处理读数（用于测试）
func (s *BillingService) ProcessReadingByMeterNo(meterNo string, readingValue decimal.Decimal, readingTime time.Time) error {
	meter, err := s.meterRepo.FindByMeterNo(meterNo)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("%w: meter_no=%s", ErrMeterNotFound, meterNo)
		}
		return fmt.Errorf("find meter: %w", err)
	}

	return s.processReadingForMeter(meter, readingValue, readingTime)
}

// GetMeterByNo 获取电表信息
func (s *BillingService) GetMeterByNo(meterNo string) (*models.ElectricMeter, error) {
	return s.meterRepo.FindByMeterNo(meterNo)
}

// GetMeterByCommAddr 根据通信地址获取电表
func (s *BillingService) GetMeterByCommAddr(commAddr string) (*models.ElectricMeter, error) {
	return s.meterRepo.FindByCommAddr(commAddr)
}

// CalculateFee 计算费用（用于测试）
func (s *BillingService) CalculateFee(rateID int64, consumption decimal.Decimal) (decimal.Decimal, error) {
	rate, err := s.electricRateRepo.FindByIDWithDetails(rateID)
	if err != nil {
		return decimal.Zero, err
	}
	return s.calculator.CalculateFee(rate, consumption, false), nil
}

// init 初始化随机数种子
func init() {
	rand.Seed(time.Now().UnixNano())
}
