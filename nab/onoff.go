package nab

import (
	"context"
	"log/slog"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/twiglab/h2o/nab/equlib"
	"github.com/twiglab/h2o/pkg/common"
)

type HandleData struct {
	Sender Sender
	Global Global
	Logger *slog.Logger
}

func topic(t string) (string, string, string) {
	ss := strings.Split(t, "/")
	_ = ss[4]

	if ss[0] != common.H2O {
		panic("not h2o topic")
	}
	if ss[1] != "onoff" {
		panic("not h2o onoff")
	}

	return ss[2], ss[3], ss[4]
}

func OnOffHandle(data HandleData) mqtt.MessageHandler {
	return func(cli mqtt.Client, msg mqtt.Message) {
		if msg.Duplicate() {
			return
		}
		defer msg.Ack()

		_, code, op := topic(msg.Topic())

		dev := data.Global.MustGetDev(context.Background(), code)

		mcli := data.Global.MustGetClient(dev.Cli)

		onoff := equlib.From[OnOffer](dev.Clazz)

		switch op {
		case common.ON:
			if err := mcli.DoOn(context.Background(), onoff, DeviceData{
				Record: dev,
				Sender: data.Sender,
				Global: data.Global,
				Logger: data.Logger,
			}); err != nil {
				data.Logger.Error("onoff error", slog.String("code", dev.Code),
					slog.String("op", op),
					slog.Any("error", err), slog.Any("record", dev))
				return
			}
		case common.OFF:
			if err := mcli.DoOff(context.Background(), onoff, DeviceData{
				Record: dev,
				Sender: data.Sender,
				Global: data.Global,
				Logger: data.Logger,
			}); err != nil {
				data.Logger.Error("onoff error", slog.String("code", dev.Code),
					slog.String("op", op),
					slog.Any("error", err), slog.Any("record", dev))
				return
			}
		}
	}
}

type OnOffer interface {
	DeviceOn
	DeviceOff
}
