package hank

import (
	"context"
	"io"
	"log/slog"
	"os"

	"github.com/twiglab/h2o/pkg/data"
	"gopkg.in/natefinch/lumberjack.v2"
)

type LogAction struct {
	log *slog.Logger
}

func NewLogAction() LogAction {
	return LogAction{
		log: NewLog("console", slog.LevelDebug),
	}
}

func (c LogAction) SendData(ctx context.Context, data data.Device) error {
	c.log.DebugContext(ctx, "sendData", slog.Any("data", data))
	return nil
}

func isConsole(logFile string) bool {
	if logFile == "" || logFile == "console" {
		return true
	}
	return false
}

func NewLog(logFile string, level slog.Level) *slog.Logger {
	var out io.Writer = os.Stdout
	if !isConsole(logFile) {
		out = &lumberjack.Logger{
			Filename:   logFile,
			MaxSize:    10, // megabytes
			MaxBackups: 10,
			MaxAge:     10, //days
		}
	}
	h := slog.NewJSONHandler(out, &slog.HandlerOptions{Level: level})
	return slog.New(h)
}

func NewDataLog(logF string) *slog.Logger {
	logFile := logF

	out := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    10, // megabytes
		MaxBackups: 10,
		MaxAge:     10, //days
	}
	h := slog.NewJSONHandler(out, nil)

	return slog.New(h)
}
