package wal

import (
	"log/slog"

	"gopkg.in/natefinch/lumberjack.v2"
)

func NewHandle(logFile string) slog.Handler {
	out := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    10, // megabytes
		MaxBackups: 10,
		MaxAge:     10, //days
	}
	return slog.NewJSONHandler(out, nil)
}
