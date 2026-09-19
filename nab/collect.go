package nab

import (
	"context"
	"log/slog"

	"github.com/twiglab/h2o/nab/equlib"
	"github.com/twiglab/h2o/nab/orm/ent"
)

type CollectTask struct {
	Global Global
	Data   []*ent.Dev
	Sender Sender
	Logger *slog.Logger
}

func (t CollectTask) Run() {
	for _, dev := range t.Data {
		collect := equlib.From[Collector](dev.Clazz)
		client := t.Global.MustGetClient(dev.Cli)

		err := client.DoCollect(context.Background(), collect, DeviceData{
			Record: dev,
			Config: t.Global,
			Sender: t.Sender,
			Logger: t.Logger,
		})

		if err != nil {
			// TODO
			// logger
		}
	}
}
