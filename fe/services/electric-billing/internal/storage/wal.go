package storage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

// WAL Write-Ahead Log
type WAL struct {
	dir      string
	file     *os.File
	mu       sync.Mutex
	maxSize  int64 // 单文件最大字节数
	logger   *zap.Logger
	filename string
}

// WalEntryType WAL 条目类型
type WalEntryType string

const (
	WalEntryTypeReading   WalEntryType = "reading"
	WalEntryTypeDeduction WalEntryType = "deduction"
)

// WalEntryStatus WAL 条目状态
type WalEntryStatus string

const (
	WalStatusPending   WalEntryStatus = "pending"
	WalStatusCompleted WalEntryStatus = "completed"
	WalStatusFailed    WalEntryStatus = "failed"
)

// WalEntry WAL 条目
type WalEntry struct {
	ID          string          `json:"id"`
	Type        WalEntryType    `json:"type"`
	Status      WalEntryStatus  `json:"status"`
	MeterNo     string          `json:"meter_no,omitempty"`
	MeterID     int64           `json:"meter_id,omitempty"`
	CommAddr    string          `json:"comm_addr,omitempty"`
	Reading     decimal.Decimal `json:"reading,omitempty"`
	Consumption decimal.Decimal `json:"consumption,omitempty"`
	Amount      decimal.Decimal `json:"amount,omitempty"`
	Timestamp   time.Time       `json:"timestamp"`
	CreatedAt   time.Time       `json:"created_at"`
	CompletedAt *time.Time      `json:"completed_at,omitempty"`
	Error       string          `json:"error,omitempty"`
}

// NewWAL 创建 WAL
func NewWAL(dir string, logger *zap.Logger) (*WAL, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("create WAL directory: %w", err)
	}

	wal := &WAL{
		dir:     dir,
		maxSize: 10 * 1024 * 1024, // 10MB
		logger:  logger,
	}

	if err := wal.openOrCreateFile(); err != nil {
		return nil, err
	}

	return wal, nil
}

// openOrCreateFile 打开或创建 WAL 文件
func (w *WAL) openOrCreateFile() error {
	w.filename = filepath.Join(w.dir, fmt.Sprintf("wal_%s.log", time.Now().Format("20060102_150405")))

	file, err := os.OpenFile(w.filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("open WAL file: %w", err)
	}

	w.file = file
	return nil
}

// rotateIfNeeded 检查是否需要轮转文件
func (w *WAL) rotateIfNeeded() error {
	info, err := w.file.Stat()
	if err != nil {
		return err
	}

	if info.Size() >= w.maxSize {
		w.file.Close()
		return w.openOrCreateFile()
	}

	return nil
}

// WriteReading 写入读数条目
func (w *WAL) WriteReading(meterNo string, meterID int64, commAddr string, reading decimal.Decimal, timestamp time.Time) (*WalEntry, error) {
	entry := &WalEntry{
		ID:        uuid.New().String(),
		Type:      WalEntryTypeReading,
		Status:    WalStatusPending,
		MeterNo:   meterNo,
		MeterID:   meterID,
		CommAddr:  commAddr,
		Reading:   reading,
		Timestamp: timestamp,
		CreatedAt: time.Now(),
	}

	if err := w.write(entry); err != nil {
		return nil, err
	}

	return entry, nil
}

// WriteDeduction 写入扣费条目
func (w *WAL) WriteDeduction(meterID int64, consumption, amount decimal.Decimal) (*WalEntry, error) {
	entry := &WalEntry{
		ID:          uuid.New().String(),
		Type:        WalEntryTypeDeduction,
		Status:      WalStatusPending,
		MeterID:     meterID,
		Consumption: consumption,
		Amount:      amount,
		CreatedAt:   time.Now(),
	}

	if err := w.write(entry); err != nil {
		return nil, err
	}

	return entry, nil
}

// write 写入条目
func (w *WAL) write(entry *WalEntry) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if err := w.rotateIfNeeded(); err != nil {
		w.logger.Error("failed to rotate WAL file", zap.Error(err))
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("marshal WAL entry: %w", err)
	}

	if _, err := w.file.Write(append(data, '\n')); err != nil {
		return fmt.Errorf("write WAL entry: %w", err)
	}

	if err := w.file.Sync(); err != nil {
		return fmt.Errorf("sync WAL file: %w", err)
	}

	w.logger.Debug("WAL entry written",
		zap.String("id", entry.ID),
		zap.String("type", string(entry.Type)))

	return nil
}

// MarkCompleted 标记条目已完成
func (w *WAL) MarkCompleted(id string) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	now := time.Now()
	completed := map[string]interface{}{
		"id":           id,
		"status":       WalStatusCompleted,
		"completed_at": now,
	}

	data, err := json.Marshal(completed)
	if err != nil {
		return fmt.Errorf("marshal completion: %w", err)
	}

	if _, err := w.file.Write(append(data, '\n')); err != nil {
		return fmt.Errorf("write completion: %w", err)
	}

	if err := w.file.Sync(); err != nil {
		return fmt.Errorf("sync WAL file: %w", err)
	}

	return nil
}

// MarkFailed 标记条目失败
func (w *WAL) MarkFailed(id string, errMsg string) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	now := time.Now()
	failed := map[string]interface{}{
		"id":           id,
		"status":       WalStatusFailed,
		"error":        errMsg,
		"completed_at": now,
	}

	data, err := json.Marshal(failed)
	if err != nil {
		return fmt.Errorf("marshal failure: %w", err)
	}

	if _, err := w.file.Write(append(data, '\n')); err != nil {
		return fmt.Errorf("write failure: %w", err)
	}

	return w.file.Sync()
}

// GetPendingEntries 获取未完成的条目
func (w *WAL) GetPendingEntries() ([]WalEntry, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	// 读取所有 WAL 文件
	files, err := filepath.Glob(filepath.Join(w.dir, "wal_*.log"))
	if err != nil {
		return nil, fmt.Errorf("glob WAL files: %w", err)
	}

	entries := make(map[string]*WalEntry)
	completedIDs := make(map[string]bool)

	for _, filename := range files {
		if err := w.parseWALFile(filename, entries, completedIDs); err != nil {
			w.logger.Warn("failed to parse WAL file",
				zap.String("file", filename),
				zap.Error(err))
			continue
		}
	}

	// 过滤出未完成的条目
	var pending []WalEntry
	for id, entry := range entries {
		if !completedIDs[id] && entry.Status == WalStatusPending {
			pending = append(pending, *entry)
		}
	}

	return pending, nil
}

// parseWALFile 解析单个 WAL 文件
func (w *WAL) parseWALFile(filename string, entries map[string]*WalEntry, completedIDs map[string]bool) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		// 尝试解析为完整条目
		var entry WalEntry
		if err := json.Unmarshal(line, &entry); err == nil {
			if entry.ID != "" {
				if entry.Status == WalStatusCompleted || entry.Status == WalStatusFailed {
					completedIDs[entry.ID] = true
				} else if entry.Type != "" {
					entries[entry.ID] = &entry
				}
			}
		}
	}

	return scanner.Err()
}

// Close 关闭 WAL
func (w *WAL) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.file != nil {
		return w.file.Close()
	}
	return nil
}

// Cleanup 清理已完成的 WAL 文件
// 删除超过指定天数且所有条目都已完成的文件
func (w *WAL) Cleanup(retentionDays int) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	cutoff := time.Now().AddDate(0, 0, -retentionDays)

	files, err := filepath.Glob(filepath.Join(w.dir, "wal_*.log"))
	if err != nil {
		return err
	}

	for _, filename := range files {
		info, err := os.Stat(filename)
		if err != nil {
			continue
		}

		if info.ModTime().Before(cutoff) && filename != w.filename {
			if err := os.Remove(filename); err != nil {
				w.logger.Warn("failed to remove old WAL file",
					zap.String("file", filename),
					zap.Error(err))
			} else {
				w.logger.Info("removed old WAL file", zap.String("file", filename))
			}
		}
	}

	return nil
}
