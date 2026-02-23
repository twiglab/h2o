package service

import (
	"time"

	"electric-billing/internal/repository"
	"electric-billing/internal/storage"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// RecoveryService 故障恢复服务
type RecoveryService struct {
	wal         *storage.WAL
	fileStore   *storage.FileStore
	meterRepo   *repository.MeterRepository
	readingRepo *repository.ReadingRepository
	accountRepo *repository.AccountRepository
	db          *gorm.DB
	logger      *zap.Logger
}

// NewRecoveryService 创建恢复服务
func NewRecoveryService(
	wal *storage.WAL,
	fileStore *storage.FileStore,
	meterRepo *repository.MeterRepository,
	readingRepo *repository.ReadingRepository,
	accountRepo *repository.AccountRepository,
	db *gorm.DB,
	logger *zap.Logger,
) *RecoveryService {
	return &RecoveryService{
		wal:         wal,
		fileStore:   fileStore,
		meterRepo:   meterRepo,
		readingRepo: readingRepo,
		accountRepo: accountRepo,
		db:          db,
		logger:      logger,
	}
}

// RecoverFromWAL 从 WAL 恢复未完成的操作
// 服务启动时调用
func (s *RecoveryService) RecoverFromWAL() error {
	entries, err := s.wal.GetPendingEntries()
	if err != nil {
		return err
	}

	if len(entries) == 0 {
		s.logger.Info("no pending WAL entries to recover")
		return nil
	}

	s.logger.Info("recovering from WAL", zap.Int("pending_entries", len(entries)))

	recovered := 0
	failed := 0

	for _, entry := range entries {
		switch entry.Type {
		case storage.WalEntryTypeReading:
			if err := s.recoverReading(entry); err != nil {
				s.logger.Error("failed to recover reading",
					zap.String("id", entry.ID),
					zap.Error(err))
				s.wal.MarkFailed(entry.ID, err.Error())
				failed++
			} else {
				s.wal.MarkCompleted(entry.ID)
				recovered++
			}

		case storage.WalEntryTypeDeduction:
			// 扣费恢复比较复杂，需要检查数据库状态
			if err := s.recoverDeduction(entry); err != nil {
				s.logger.Error("failed to recover deduction",
					zap.String("id", entry.ID),
					zap.Error(err))
				s.wal.MarkFailed(entry.ID, err.Error())
				failed++
			} else {
				s.wal.MarkCompleted(entry.ID)
				recovered++
			}
		}
	}

	s.logger.Info("WAL recovery completed",
		zap.Int("recovered", recovered),
		zap.Int("failed", failed))

	return nil
}

// recoverReading 恢复读数记录
func (s *RecoveryService) recoverReading(entry storage.WalEntry) error {
	// 检查数据库是否已有该读数记录
	// 通过 meter_id + reading_time + reading_value 判断唯一性
	exists, err := s.readingExists(entry.MeterID, entry.Reading, entry.Timestamp)
	if err != nil {
		return err
	}

	if exists {
		s.logger.Debug("reading already exists in database",
			zap.String("id", entry.ID),
			zap.Int64("meter_id", entry.MeterID))
		return nil
	}

	// 如果不存在，记录日志但不重新处理
	// 因为重新处理需要完整的计费服务上下文
	s.logger.Warn("reading not found in database, may need manual recovery",
		zap.String("id", entry.ID),
		zap.Int64("meter_id", entry.MeterID),
		zap.String("reading", entry.Reading.String()),
		zap.Time("timestamp", entry.Timestamp))

	return nil
}

// readingExists 检查电表读数是否已存在
func (s *RecoveryService) readingExists(meterID int64, reading decimal.Decimal, readingTime time.Time) (bool, error) {
	var count int64
	err := s.db.Table("col_electric_reading").
		Where("meter_id = ? AND reading_value = ? AND reading_time = ?",
			meterID, reading, readingTime).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// recoverDeduction 恢复扣费记录
func (s *RecoveryService) recoverDeduction(entry storage.WalEntry) error {
	// 扣费恢复通常需要人工介入，因为涉及账户余额
	s.logger.Warn("deduction recovery requires manual intervention",
		zap.String("id", entry.ID),
		zap.Int64("meter_id", entry.MeterID),
		zap.String("amount", entry.Amount.String()))

	return nil
}

// RecoverFromFiles 从本地文件恢复扣费记录到数据库
// 用于数据库故障后的数据恢复
func (s *RecoveryService) RecoverFromFiles(startDate, endDate time.Time) error {
	records, err := s.fileStore.ReadDeductions(startDate, endDate)
	if err != nil {
		return err
	}

	if len(records) == 0 {
		s.logger.Info("no deduction records to recover from files")
		return nil
	}

	s.logger.Info("recovering deductions from files",
		zap.Int("record_count", len(records)),
		zap.Time("start_date", startDate),
		zap.Time("end_date", endDate))

	recovered := 0
	skipped := 0

	for _, record := range records {
		exists, err := s.deductionExists(record.DeductionNo)
		if err != nil {
			s.logger.Error("failed to check deduction existence",
				zap.String("deduction_no", record.DeductionNo),
				zap.Error(err))
			continue
		}

		if exists {
			skipped++
			continue
		}

		// 插入到数据库
		if err := s.createDeductionFromRecord(record); err != nil {
			s.logger.Error("failed to create deduction from file",
				zap.String("deduction_no", record.DeductionNo),
				zap.Error(err))
			continue
		}

		recovered++
	}

	s.logger.Info("file recovery completed",
		zap.Int("recovered", recovered),
		zap.Int("skipped", skipped))

	return nil
}

// deductionExists 检查扣费记录是否已存在
func (s *RecoveryService) deductionExists(deductionNo string) (bool, error) {
	var count int64
	err := s.db.Table("fin_electric_deduction").
		Where("deduction_no = ?", deductionNo).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// createDeductionFromRecord 从文件记录创建数据库记录
func (s *RecoveryService) createDeductionFromRecord(record storage.DeductionRecord) error {
	return s.db.Exec(`
		INSERT INTO fin_electric_deduction
		(deduction_no, account_id, consumption_id, meter_id, consumption, unit_price,
		 amount, balance_before, balance_after, deduction_time, status, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, record.DeductionNo, record.AccountID, record.ConsumptionID, record.MeterID,
		record.Consumption, record.UnitPrice, record.Amount, record.BalanceBefore,
		record.BalanceAfter, record.CreatedAt, record.Status,
		record.CreatedAt).Error
}

// CheckDatabaseConnection 检查数据库连接
func (s *RecoveryService) CheckDatabaseConnection() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

// GetRecoveryStatus 获取恢复状态
func (s *RecoveryService) GetRecoveryStatus() (map[string]interface{}, error) {
	status := make(map[string]interface{})

	// 获取 WAL 待处理条目数
	pendingEntries, err := s.wal.GetPendingEntries()
	if err != nil {
		return nil, err
	}
	status["wal_pending_count"] = len(pendingEntries)

	// 获取文件存储统计
	fileStats, err := s.fileStore.GetStats()
	if err != nil {
		s.logger.Warn("failed to get file store stats", zap.Error(err))
	} else {
		for k, v := range fileStats {
			status["file_"+k] = v
		}
	}

	// 检查数据库连接
	if err := s.CheckDatabaseConnection(); err != nil {
		status["database_connected"] = false
		status["database_error"] = err.Error()
	} else {
		status["database_connected"] = true
	}

	return status, nil
}
