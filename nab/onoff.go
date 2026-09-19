package nab

import (
	"context"
	"log/slog"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/twiglab/h2o/nab/equlib"
	"github.com/twiglab/h2o/pkg/common"
)

type HandleData struct {
	Sender Sender
	Global Global
	Logger *slog.Logger
}

func SubscriptTopic(boxCode string) string {
	return common.H2O + "/onoff/" + boxCode + "/#"
}

func topicPart(t string) (string, string, string) {
	ss := common.TopicPart(t)
	_ = ss[4]

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

		boxCode, devCode, op := topicPart(msg.Topic())
		if boxCode != data.Global.BoxCode {
			data.Logger.Error("onoff box code error",
				slog.String("boxCode", boxCode),
				slog.String("devCode", devCode),
				slog.String("globalBoxCode", data.Global.BoxCode),
				slog.String("op", op),
			)
		}

		dev := data.Global.MustGetDev(context.Background(), devCode)

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
					slog.Any("record", dev),
					slog.Any("error", err),
				)
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
					slog.Any("record", dev),
					slog.Any("error", err),
				)
				return
			}
		}
	}
}

type OnOffer interface {
	DeviceOn
	DeviceOff
}
