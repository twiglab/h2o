package nab

import (
	"context"
	"log/slog"
	"time"

	"github.com/twiglab/h2o/nab/orm/ent"
	"github.com/twiglab/h2o/proto"
)

type CollectTask struct {
	Box    *Box
	Logger *slog.Logger

	Data []*ent.Dev

	Pipeline Sender

	Delay time.Duration
}

func (t CollectTask) Run(ctx context.Context) {
	for {
		for _, dev := range t.Data {
			select {
			case <-ctx.Done():
				t.Logger.DebugContext(ctx, "collect quit")
				return
			default:
				t.doCollect(ctx, dev)
			}
		}
	}
}

func (t CollectTask) doCollect(ctx context.Context, dev *ent.Dev) {
	time.Sleep(t.Delay)

	dev, cli, coll, err := t.Box.GetDevCliEqu[DeviceCollector](dev.Code)
	if err != nil {
		t.Logger.ErrorContext(ctx, "GetDevCliEqu error", slog.String("code", dev.Code), slog.Any("error", err))
		return
	}

	mv, err := cli.Do(dev,
		func(op ModbusOperator) (mv proto.MeterValue, err error) {
			mv, err = coll.Collect(ctx, op)
			return
		},
	)

	if err != nil {
		t.Logger.ErrorContext(ctx, "CollectTask error", slog.String("code", dev.Code), slog.Any("error", err))
		return
	}

	_ = t.Pipeline.SendData(ctx, MakeMeter(t.Box, dev, mv))
}
