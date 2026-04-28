package clog

import (
	"io"
	"log/slog"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
)

func isConsole(logFile string) bool {
	if logFile == "" || logFile == "console" {
		return true
	}
	return false
}

func NewLog(logFile string, level slog.Level) *slog.Logger {
	var out io.Writer = os.Stdout
	if !isConsole(logFile) {
		out = NewLogWriter(logFile)
	}
	h := slog.NewJSONHandler(out, &slog.HandlerOptions{Level: level})
	return slog.New(h)
}

func NewLogWriter(logf string) io.Writer {
	return &lumberjack.Logger{
		Filename:  logf,
		MaxSize:   10,  // megabytes
		MaxAge:    180, //days
		LocalTime: true,
	}
}
