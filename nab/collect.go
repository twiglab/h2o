package nab

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/twiglab/h2o/nab/orm/ent"
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
				fmt.Println("qqq")
				return
			default:
				t.doCollect(ctx, dev)
			}
		}
	}
}

func (t CollectTask) doCollect(ctx context.Context, dev *ent.Dev) {
	time.Sleep(t.Delay)

	_, cli, coll, err := t.Box.GetDevCliEqu[DeviceCollector](dev.Code)
	if err != nil {
		t.Logger.Error("GetDevCliEqu error", slog.String("code", dev.Code), slog.Any("error", err))
		return
	}
	err = cli.Do(
		func(op ModbusOperator) error {
			mv, err := coll.Collect(ctx, op, DeviceData{Box: t.Box, Record: dev, Logger: t.Logger})
			if err != nil {
				return err
			}
			return t.Pipeline.SendData(ctx, MakeMeter(t.Box, dev, mv))
		},
	)

	if err != nil {
		t.Logger.Error("CollectTask error", slog.String("code", dev.Code), slog.Any("error", err))
	}
}
