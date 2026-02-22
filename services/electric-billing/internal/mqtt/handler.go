package mqtt

import (
	"encoding/hex"
	"regexp"
	"strings"
	"time"

	"electric-billing/internal/modbus"

	"go.uber.org/zap"
)

// ReadingProcessor 读数处理器接口
type ReadingProcessor interface {
	ProcessReading(dtuID string, meterAddr byte, readingValue []byte, readingTime time.Time) error
}

// Handler MQTT消息处理器
type Handler struct {
	logger    *zap.Logger
	processor ReadingProcessor
	topicRe   *regexp.Regexp
}

// NewHandler 创建消息处理器
func NewHandler(logger *zap.Logger, processor ReadingProcessor) *Handler {
	// 匹配主题格式: dtu/{dtu_id}/data
	re := regexp.MustCompile(`^dtu/([^/]+)/data$`)

	return &Handler{
		logger:    logger,
		processor: processor,
		topicRe:   re,
	}
}

// HandleMessage 处理MQTT消息
func (h *Handler) HandleMessage(topic string, payload []byte) {
	h.logger.Debug("received MQTT message",
		zap.String("topic", topic),
		zap.Int("payload_len", len(payload)),
		zap.String("payload_hex", hex.EncodeToString(payload)))

	// 解析主题获取DTU ID
	matches := h.topicRe.FindStringSubmatch(topic)
	if matches == nil {
		h.logger.Warn("invalid topic format", zap.String("topic", topic))
		return
	}
	dtuID := matches[1]

	// 解析Modbus帧
	frame, err := modbus.ParseFrame(payload)
	if err != nil {
		h.logger.Error("failed to parse Modbus frame",
			zap.String("dtu_id", dtuID),
			zap.Error(err),
			zap.String("payload_hex", hex.EncodeToString(payload)))
		return
	}

	h.logger.Debug("parsed Modbus frame",
		zap.String("dtu_id", dtuID),
		zap.Uint8("address", frame.Address),
		zap.Uint8("function", frame.Function),
		zap.Int("data_len", len(frame.Data)))

	// 只处理读取响应 (功能码 03 或 04)
	if frame.Function != modbus.FuncReadHoldingRegisters && frame.Function != modbus.FuncReadInputRegisters {
		h.logger.Debug("ignoring non-read response",
			zap.Uint8("function", frame.Function))
		return
	}

	// 调用处理器处理读数
	readingTime := time.Now()
	if err := h.processor.ProcessReading(dtuID, frame.Address, payload, readingTime); err != nil {
		h.logger.Error("failed to process reading",
			zap.String("dtu_id", dtuID),
			zap.Uint8("meter_addr", frame.Address),
			zap.Error(err))
		return
	}

	h.logger.Info("reading processed successfully",
		zap.String("dtu_id", dtuID),
		zap.Uint8("meter_addr", frame.Address))
}

// ParseDTUIDFromTopic 从主题中解析DTU ID
func ParseDTUIDFromTopic(topic string) string {
	// 主题格式: dtu/{dtu_id}/data
	parts := strings.Split(topic, "/")
	if len(parts) >= 3 && parts[0] == "dtu" && parts[2] == "data" {
		return parts[1]
	}
	return ""
}
