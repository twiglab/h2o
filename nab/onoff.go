package nab

import (
	"cmp"
	"context"
	"encoding/json/v2"
	"log/slog"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/twiglab/h2o/nab/equlib"
	"github.com/twiglab/h2o/pkg/common"
)

type OnOffLite struct {
	BoxCode string
	Code    string
	Op      string
}

func (m OnOffLite) MarshalBinary() (data []byte, err error) {
	return json.Marshal(m)
}

func (o OnOffLite) Topic() string {
	return common.H2O + "/onoff/" + o.BoxCode + "/" + o.Code + "/" + o.Op
}

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

		if cmp.Compare(boxCode, data.Global.Box) != 0 {
			data.Logger.Error("onoff box code error",
				slog.String("box", boxCode),
				slog.String("code", devCode),
				slog.String("global.Box", data.Global.Box),
				slog.String("op", op),
			)
		}

		dev, err := data.Global.IDB.GetDev(context.Background(), devCode)
		if err != nil {
			data.Logger.Error("OnOff.GetDev.Error",
				slog.String("box", boxCode),
				slog.String("code", devCode),
				slog.String("global.Box", data.Global.Box),
				slog.String("op", op),
			)
			return
		}

		mcli, ok := data.Global.ClientByCode(dev.Cli)
		if !ok {
			data.Logger.Error("OnOff.ClientByCode.Error",
				slog.String("box", boxCode),
				slog.String("code", devCode),
				slog.String("global.Box", data.Global.Box),
				slog.String("op", op),
			)
			return
		}

		onoff, ok := equlib.From[OnOffer](dev.Clazz)

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
