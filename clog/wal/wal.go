package wal

import (
	"context"
	"log/slog"

	"github.com/twiglab/h2o/clog"
)

type Conf struct {
	Filename string
}

type WAL struct {
	inner *slog.Logger

	c Conf
}

func New(c Conf) *WAL {
	l := clog.NewLog(c.Filename, slog.LevelInfo)
	return &WAL{
		inner: l,
		c:     c,
	}
}

func (w *WAL) WriteLog(a ...any) {
	w.WriteLogContext(context.Background(), a...)
}

func (w *WAL) WriteLogContext(ctx context.Context, a ...any) {
	w.inner.InfoContext(ctx, "wal", a...)
}

func Any(s string, d any) slog.Attr {
	return slog.Any(s, d)
}

func Int64(s string, i int64) slog.Attr {
	return slog.Int64(s, i)
}

func Int(s string, i int) slog.Attr {
	return slog.Int(s, i)
}
func String(k, v string) slog.Attr {
	return slog.String(k, v)
}
