package nab

import (
	"context"
	"log/slog"
	"time"

	"github.com/twiglab/h2o/nab/equlib"
	"github.com/twiglab/h2o/nab/orm"
	"github.com/twiglab/h2o/nab/orm/ent"
)

type CollectTask struct {
	Global    Global
	IDB       *orm.IDB
	ClientMgr *ClientMgr
	Data      []*ent.Dev
	Sender    Sender
	Logger    *slog.Logger

	Delay time.Duration
}

func (t CollectTask) Run() {
	for _, dev := range t.Data {
		coltor, err := equlib.From[DeviceCollector](dev.Clazz)
		if err != nil {
			t.Logger.Error("CollectTask.NotFoundClazz", slog.String("code", dev.Code),
				slog.String("clazz", dev.Clazz),
				slog.String("box", t.Global.Box),
				slog.Any("record", dev))
		}
		client, err := t.ClientMgr.ClientByCode(dev.Cli)
		if err != nil {
			t.Logger.Error("CollectTask.NotFoundClient", slog.String("code", dev.Code),
				slog.String("cli", dev.Clazz),
				slog.String("box", t.Global.Box),
				slog.Any("record", dev))
		}

		err = client.DoCollect(context.Background(), coltor, DeviceData{
			Record: dev,
			Global: t.Global,
			Sender: t.Sender,
			Logger: t.Logger,
		})

		if err != nil {
			t.Logger.Error("CollectTask.CollectError", slog.String("code", dev.Code), slog.Any("error", err), slog.Any("record", dev))
		}
		time.Sleep(t.Delay)
	}
}
