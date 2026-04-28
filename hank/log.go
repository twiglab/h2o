package hank

import (
	"context"
	"fmt"
	"io"
	"log/slog"

	"github.com/twiglab/h2o/clog"
)

type LogAction struct {
}

func (c LogAction) SendData(ctx context.Context, obj SendObject) error {
	slog.DebugContext(ctx, "logAction", slog.Any("data", obj), slog.String("topic", obj.Topic()))
	return nil
}

type PlayBack struct {
	out io.Writer
}

func NewPlayBack(logf string) *PlayBack {
	o := clog.NewLogWriter(logf)
	return &PlayBack{
		out: o,
	}
}

func (p *PlayBack) Record(ctx context.Context, data string) {
	fmt.Fprintln(p.out, data)
}
