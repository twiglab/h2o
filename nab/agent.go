package nab

import (
	"context"
	"log/slog"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/panjf2000/ants/v2"
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

	Pool *ants.Pool
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

		// 这里用到了线程池
		// 主要作用是防止设备大范围的开合操作
		// topic已经提供了操作所需的所有字段，payload不做解析，直接忽略
		// 所以这里即便ack早于执行，也不会有影响
		if err := a.Pool.Submit(func() {
			a.doOnOff(code, op)
		}); err != nil {
			// pool size为32，如果pool满了，记录日志，直接丢弃
			// optStatus是水平出发的，即便丢弃，下一轮上报上游程序仍然能够获取正确值
			// 上游程序会保证再次下发正确的指令
			a.Logger.Error("submit errorr",
				slog.String("box", box),
				slog.String("code", code),
				slog.Any("error", err),
			)
		}
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
