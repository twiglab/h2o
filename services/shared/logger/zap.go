package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Options 日志选项
type Options struct {
	Level       string // debug, info, warn, error
	Encoding    string // json, console
	Development bool
}

// DefaultOptions 返回默认日志选项
func DefaultOptions() Options {
	return Options{
		Level:       "info",
		Encoding:    "json",
		Development: false,
	}
}

// New 创建新的 zap.Logger 实例
func New(opts ...Options) *zap.Logger {
	opt := DefaultOptions()
	if len(opts) > 0 {
		opt = opts[0]
	}

	zapLevel := parseLevel(opt.Level)

	config := zap.Config{
		Level:       zap.NewAtomicLevelAt(zapLevel),
		Development: opt.Development,
		Encoding:    opt.Encoding,
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, err := config.Build()
	if err != nil {
		panic(err)
	}

	return logger
}

// NewWithLevel 使用指定日志级别创建 logger
func NewWithLevel(level string) *zap.Logger {
	return New(Options{Level: level, Encoding: "json"})
}

// NewDevelopment 创建开发模式的 logger (console格式)
func NewDevelopment() *zap.Logger {
	return New(Options{
		Level:       "debug",
		Encoding:    "console",
		Development: true,
	})
}

// parseLevel 解析日志级别
func parseLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}
