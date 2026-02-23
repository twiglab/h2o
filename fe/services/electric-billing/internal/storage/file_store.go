package storage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

// FileStore 本地文件存储
type FileStore struct {
	baseDir string
	mu      sync.Mutex
	logger  *zap.Logger
}

// DeductionRecord 扣费记录（文件存储格式）
type DeductionRecord struct {
	DeductionNo   string          `json:"deduction_no"`
	MeterNo       string          `json:"meter_no"`
	MeterID       int64           `json:"meter_id"`
	AccountID     int64           `json:"account_id"`
	ConsumptionID int64           `json:"consumption_id"`
	Consumption   decimal.Decimal `json:"consumption"`
	UnitPrice     decimal.Decimal `json:"unit_price"`
	Amount        decimal.Decimal `json:"amount"`
	BalanceBefore decimal.Decimal `json:"balance_before"`
	BalanceAfter  decimal.Decimal `json:"balance_after"`
	DeductionType int16           `json:"deduction_type"`
	Status        int16           `json:"status"`
	CreatedAt     time.Time       `json:"created_at"`
}

// ReadingRecord 读数记录（文件存储格式）
type ReadingRecord struct {
	MeterNo      string          `json:"meter_no"`
	MeterID      int64           `json:"meter_id"`
	ReadingValue decimal.Decimal `json:"reading_value"`
	ReadingTime  time.Time       `json:"reading_time"`
	CollectType  int16           `json:"collect_type"`
	CreatedAt    time.Time       `json:"created_at"`
}

// NewFileStore 创建文件存储
func NewFileStore(baseDir string, logger *zap.Logger) (*FileStore, error) {
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, fmt.Errorf("create base directory: %w", err)
	}

	return &FileStore{
		baseDir: baseDir,
		logger:  logger,
	}, nil
}

// SaveDeduction 保存扣费记录
func (s *FileStore) SaveDeduction(record *DeductionRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 按日期分目录: data/deductions/2024/01/23/
	dir := filepath.Join(s.baseDir, "deductions",
		record.CreatedAt.Format("2006"),
		record.CreatedAt.Format("01"),
		record.CreatedAt.Format("02"))

	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create deduction directory: %w", err)
	}

	filename := filepath.Join(dir, fmt.Sprintf("deductions_%s.jsonl", record.CreatedAt.Format("20060102")))

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("open deduction file: %w", err)
	}
	defer file.Close()

	data, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("marshal deduction record: %w", err)
	}

	if _, err := file.Write(append(data, '\n')); err != nil {
		return fmt.Errorf("write deduction record: %w", err)
	}

	s.logger.Debug("deduction record saved to file",
		zap.String("deduction_no", record.DeductionNo),
		zap.String("file", filename))

	return nil
}

// SaveReading 保存读数记录
func (s *FileStore) SaveReading(record *ReadingRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 按日期分目录: data/readings/2024/01/23/
	dir := filepath.Join(s.baseDir, "readings",
		record.CreatedAt.Format("2006"),
		record.CreatedAt.Format("01"),
		record.CreatedAt.Format("02"))

	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create reading directory: %w", err)
	}

	filename := filepath.Join(dir, fmt.Sprintf("readings_%s.jsonl", record.CreatedAt.Format("20060102")))

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("open reading file: %w", err)
	}
	defer file.Close()

	data, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("marshal reading record: %w", err)
	}

	if _, err := file.Write(append(data, '\n')); err != nil {
		return fmt.Errorf("write reading record: %w", err)
	}

	s.logger.Debug("reading record saved to file",
		zap.String("meter_no", record.MeterNo),
		zap.String("file", filename))

	return nil
}

// ReadDeductions 读取指定日期范围的扣费记录
func (s *FileStore) ReadDeductions(startDate, endDate time.Time) ([]DeductionRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var records []DeductionRecord

	// 遍历日期范围
	for date := startDate; !date.After(endDate); date = date.AddDate(0, 0, 1) {
		dir := filepath.Join(s.baseDir, "deductions",
			date.Format("2006"),
			date.Format("01"),
			date.Format("02"))

		filename := filepath.Join(dir, fmt.Sprintf("deductions_%s.jsonl", date.Format("20060102")))

		dayRecords, err := s.readDeductionFile(filename)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			s.logger.Warn("failed to read deduction file",
				zap.String("file", filename),
				zap.Error(err))
			continue
		}

		records = append(records, dayRecords...)
	}

	return records, nil
}

// readDeductionFile 读取单个扣费文件
func (s *FileStore) readDeductionFile(filename string) ([]DeductionRecord, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var records []DeductionRecord
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		var record DeductionRecord
		if err := json.Unmarshal(line, &record); err != nil {
			s.logger.Warn("failed to unmarshal deduction record",
				zap.Error(err),
				zap.String("line", string(line)))
			continue
		}

		records = append(records, record)
	}

	return records, scanner.Err()
}

// ReadReadings 读取指定日期范围的读数记录
func (s *FileStore) ReadReadings(startDate, endDate time.Time) ([]ReadingRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var records []ReadingRecord

	// 遍历日期范围
	for date := startDate; !date.After(endDate); date = date.AddDate(0, 0, 1) {
		dir := filepath.Join(s.baseDir, "readings",
			date.Format("2006"),
			date.Format("01"),
			date.Format("02"))

		filename := filepath.Join(dir, fmt.Sprintf("readings_%s.jsonl", date.Format("20060102")))

		dayRecords, err := s.readReadingFile(filename)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			s.logger.Warn("failed to read reading file",
				zap.String("file", filename),
				zap.Error(err))
			continue
		}

		records = append(records, dayRecords...)
	}

	return records, nil
}

// readReadingFile 读取单个读数文件
func (s *FileStore) readReadingFile(filename string) ([]ReadingRecord, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var records []ReadingRecord
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		var record ReadingRecord
		if err := json.Unmarshal(line, &record); err != nil {
			s.logger.Warn("failed to unmarshal reading record",
				zap.Error(err),
				zap.String("line", string(line)))
			continue
		}

		records = append(records, record)
	}

	return records, scanner.Err()
}

// Cleanup 清理旧文件
func (s *FileStore) Cleanup(retentionDays int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	cutoff := time.Now().AddDate(0, 0, -retentionDays)

	// 清理扣费记录
	if err := s.cleanupDirectory(filepath.Join(s.baseDir, "deductions"), cutoff); err != nil {
		s.logger.Warn("failed to cleanup deductions", zap.Error(err))
	}

	// 清理读数记录
	if err := s.cleanupDirectory(filepath.Join(s.baseDir, "readings"), cutoff); err != nil {
		s.logger.Warn("failed to cleanup readings", zap.Error(err))
	}

	return nil
}

// cleanupDirectory 清理指定目录下的旧文件
func (s *FileStore) cleanupDirectory(dir string, cutoff time.Time) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // 忽略错误继续
		}

		if info.IsDir() {
			return nil
		}

		if info.ModTime().Before(cutoff) {
			if err := os.Remove(path); err != nil {
				s.logger.Warn("failed to remove old file",
					zap.String("file", path),
					zap.Error(err))
			} else {
				s.logger.Info("removed old file", zap.String("file", path))
			}
		}

		return nil
	})
}

// GetStats 获取存储统计信息
func (s *FileStore) GetStats() (map[string]interface{}, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	stats := make(map[string]interface{})

	// 计算扣费记录大小
	deductionsSize, deductionsCount, err := s.getDirectoryStats(filepath.Join(s.baseDir, "deductions"))
	if err == nil {
		stats["deductions_size_bytes"] = deductionsSize
		stats["deductions_file_count"] = deductionsCount
	}

	// 计算读数记录大小
	readingsSize, readingsCount, err := s.getDirectoryStats(filepath.Join(s.baseDir, "readings"))
	if err == nil {
		stats["readings_size_bytes"] = readingsSize
		stats["readings_file_count"] = readingsCount
	}

	return stats, nil
}

// getDirectoryStats 获取目录统计
func (s *FileStore) getDirectoryStats(dir string) (int64, int, error) {
	var totalSize int64
	var fileCount int

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() {
			totalSize += info.Size()
			fileCount++
		}
		return nil
	})

	return totalSize, fileCount, err
}
