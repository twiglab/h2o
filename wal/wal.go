package wal

import (
	"log/slog"

	"gopkg.in/natefinch/lumberjack.v2"
)

func NewLog(logFile string) *slog.Logger {
	out := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    10, // megabytes
		MaxBackups: 10,
		MaxAge:     10, //days
	}
	h := slog.NewJSONHandler(out, nil)
	return slog.New(h)
}
