package nab

import (
	"context"
	"log/slog"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/twiglab/h2o/proto"
)

type OnOffer interface {
	DeviceOn
	DeviceOff
}

type Agent struct {
	Box    *Box
	Logger *slog.Logger

	OnoffLogger *slog.Logger
}

func AgentHandle(a *Agent) mqtt.MessageHandler {
	return func(cli mqtt.Client, msg mqtt.Message) {
		// 如果是retained消息，不在处理，直接返回
		if msg.Retained() {
			a.Logger.Warn("onoff retained message ignored",
				slog.String("topic", msg.Topic()))
			return
		}

		box, code, op, err := parseOnoffTopicPart(msg.Topic())
		if err != nil {
			a.Logger.Error("topic error", slog.String("topic", msg.Topic()), slog.Any("error", err))
			return
		}

		if box != a.Box.Code {
			a.Logger.Warn("onoff box mismatch",
				slog.String("expect", a.Box.Code),
				slog.String("got", box),
				slog.String("code", code))
			return
		}

		a.doOnOff(code, op)
	}
}

func (a *Agent) On(ctx context.Context, code string) error {
	dev, cli, on, err := a.Box.GetDevCliEqu[DeviceOn](code)
	if err != nil {
		return err
	}

	_, err = cli.Do(dev,
		func(op ModbusOperator) (mv proto.MeterValue, err error) {
			err = on.On(ctx, op)
			return
		})

	return err
}

func (a *Agent) Off(ctx context.Context, code string) error {
	dev, cli, off, err := a.Box.GetDevCliEqu[DeviceOff](code)
	if err != nil {
		return err
	}

	_, err = cli.Do(dev,
		func(op ModbusOperator) (mv proto.MeterValue, err error) {
			err = off.Off(ctx, op)
			return
		})

	return err
}

func (a *Agent) doOnOff(code, op string) {
	ctx, cancalFn := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancalFn()

	var err error
	switch op {
	case proto.ON:
		err = a.On(ctx, code)
	case proto.OFF:
		err = a.Off(ctx, code)
	default:
		a.Logger.Error("onoff unknow op", slog.String("code", code), slog.String("op", op))
		return
	}

	if err != nil {
		a.Logger.ErrorContext(ctx, "onoff error", slog.String("code", code),
			slog.String("op", op),
			slog.Any("error", err))
		return
	}

	a.OnoffLogger.InfoContext(ctx, "onoff",
		slog.String("device", code),
		slog.String("op", op),
	)
}
