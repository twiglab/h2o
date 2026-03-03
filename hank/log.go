package hank

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
)

type LogAction struct {
}

func (c LogAction) SendData(ctx context.Context, obj SendObject) error {
	slog.DebugContext(ctx, "logAction", slog.Any("data", obj), slog.String("topic", obj.Topic()))
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

type PlayBack struct {
	out io.Writer
}

func NewPlayBack(logf string) *PlayBack {
	o := &lumberjack.Logger{
		Filename:   logf,
		MaxSize:    10, // megabytes
		MaxBackups: 10,
		MaxAge:     10, //days
	}
	return &PlayBack{
		out: o,
	}
}

func (p *PlayBack) Record(ctx context.Context, data string) {
	fmt.Fprintln(p.out, data)
}
