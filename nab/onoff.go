package nab

import (
	"context"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/twiglab/h2o/nab/equlib"
	"github.com/twiglab/h2o/pkg/common"
)

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

func OnOffHandle(global Global) mqtt.MessageHandler {
	return func(cli mqtt.Client, msg mqtt.Message) {
		if msg.Duplicate() {
			return
		}
		defer msg.Ack()

		_, code, op := topic(msg.Topic())

		dev := global.MustGetDev(context.Background(), code)

		mcli := global.MustGetClient(dev.Cli)

		onoff := equlib.From[xx](dev.Clazz)

		switch op {
		case common.ON:
			if err := mcli.DoOn(context.Background(), onoff, DeviceData{
				Record: dev,
			}); err != nil {
				// log err
				return
			}
		case common.OFF:
			if err := mcli.DoOff(context.Background(), onoff, DeviceData{}); err != nil {
				// log err
				return
			}
		}
	}
}

type xx interface {
	DeviceOn
	DeviceOff
}
