package wal

import (
	"context"
	"log/slog"

	"gopkg.in/natefinch/lumberjack.v2"
)

func newHandle(logFile string) slog.Handler {
	out := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    10, // megabytes
		MaxBackups: 10,
		MaxAge:     10, //days
		LocalTime:  true,
	}
	return slog.NewJSONHandler(out, nil)
}

type Conf struct {
	Filename string
}

type WAL struct {
	inner *slog.Logger

	c Conf
}

func New(c Conf) *WAL {
	h := newHandle(c.Filename)
	return &WAL{
		inner: slog.New(h),
		c:     c,
	}
}

func (w *WAL) WriteLog(a ...any) {
	w.WriteLogContext(context.Background(), a...)
}

func (w *WAL) WriteLogContext(ctx context.Context, a ...any) {
	w.inner.InfoContext(ctx, "wal", a...)
}

func Data(d any) slog.Attr {
	return slog.Any("data", d)
}

func Type(t string) slog.Attr {
	return slog.String("type", t)
}
