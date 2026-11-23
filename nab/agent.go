package nab

import (
	"context"
	"log/slog"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/panjf2000/ants/v2"
	"github.com/twiglab/h2o/clog/wal"
	"github.com/twiglab/h2o/nab/equlib"
	"github.com/twiglab/h2o/nab/orm"
	"github.com/twiglab/h2o/nab/orm/ent"
	"github.com/twiglab/h2o/proto"
)

type OnOffer interface {
	DeviceOn
	DeviceOff
}

type Agent struct {
	Global Global
	Logger *slog.Logger
	IDB    *orm.IDB
	CliMgr *ClientMgr
	WAL    *wal.WAL

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

		box, code, op, err := onoffTopicPart(msg.Topic())
		if err != nil {
			return
		}

		if op != proto.ON && op != proto.OFF {
			a.Logger.Error("onoff unknown op")
			return
		}

		if box != a.Global.Box {
			a.Logger.Warn("onoff box mismatch",
				slog.String("expect", a.Global.Box),
				slog.String("got", box),
				slog.String("code", code))
			return
		}

		// 这里用到了线程池
		// 主要作用是防止设备大范围的开合操作
		// topic已经提供了操作所需的所有字段，payload不做解析，直接忽略
		// 所以这里即便ack早于执行，也不会有影响
		if err := a.Pool.Submit(func() {
			a.doOnOff(box, code, op)
		}); err != nil {
			// pool size为32，如果pool满了，记录日志，直接丢弃，等上游程序发送下一轮指令
			a.Logger.Error("submit errorr",
				slog.String("box", box),
				slog.String("code", code),
				slog.Any("error", err),
			)
		}
	}
}

func (a *Agent) On(ctx context.Context, code string) error {
	dev, mcli, onoff, err := a.get(code)
	if err != nil {
		return err
	}

	return mcli.DoOn(ctx, onoff,
		DeviceData{
			Record: dev,
			Global: a.Global,
			IDB:    a.IDB,
			Logger: a.Logger,
		},
	)
}

func (a *Agent) Off(ctx context.Context, code string) error {
	dev, mcli, onoff, err := a.get(code)
	if err != nil {
		return err
	}

	return mcli.DoOff(ctx, onoff,
		DeviceData{
			Record: dev,
			Global: a.Global,
			IDB:    a.IDB,
			Logger: a.Logger,
		},
	)
}

func (a *Agent) get(devCode string) (*ent.Dev, *ModbusCli, OnOffer, error) {
	dev, err := a.IDB.GetDev(context.Background(), devCode)
	if err != nil {
		a.Logger.Error("OnOff.GetDev.Error",
			slog.String("code", devCode),
			slog.String("global.Box", a.Global.Box),
		)
		return nil, nil, nil, err
	}

	mcli, err := a.CliMgr.ClientByCode(dev.Cli)
	if err != nil {
		a.Logger.Error("OnOff.ClientByCode.Error",
			slog.String("code", devCode),
			slog.String("global.Box", a.Global.Box),
			slog.Any("error", err),
		)
		return nil, nil, nil, err
	}

	onoff, err := equlib.From[OnOffer](dev.Clazz)
	if err != nil {
		a.Logger.Error("OnOff.Clazz.Error",
			slog.Any("error", err),
		)
		return nil, nil, nil, err
	}
	return dev, mcli, onoff, nil
}

func (a *Agent) doOnOff(boxCode, devCode, op string) {
	ctx := context.Background()

	// 无论执行是否成功，都先记录
	a.WAL.WriteLogContext(ctx,
		wal.String("box", boxCode),
		wal.String("device", devCode),
		wal.String("op", op),
	)

	switch op {
	case proto.ON:
		if err := a.On(ctx, devCode); err != nil {
			a.Logger.Error("onoff error", slog.String("code", devCode),
				slog.String("op", op),
				slog.Any("error", err),
			)
		}
	case proto.OFF:
		if err := a.Off(ctx, devCode); err != nil {
			a.Logger.Error("onoff error", slog.String("code", devCode),
				slog.String("op", op),
				slog.Any("error", err),
			)
		}
	default:
		a.Logger.Error("onoff unknow op", slog.String("code", devCode), slog.String("op", op))
	}
}
