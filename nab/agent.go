package nab

import (
	"cmp"
	"context"
	"log/slog"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/twiglab/h2o/nab/equlib"
	"github.com/twiglab/h2o/nab/orm"
	"github.com/twiglab/h2o/nab/orm/ent"
	"github.com/twiglab/h2o/pkg/common"
)

type OnOffer interface {
	DeviceOn
	DeviceOff
}

type Agent struct {
	Sender Sender
	Global Global
	Logger *slog.Logger
	IDB    *orm.IDB
	CliMgr *ClientMgr
	MCli   mqtt.Client
}

func (a Agent) Start() error {
	t := a.MCli.Subscribe(SubscriptTopic(a.Global.Box), 0x01, a.msgHandle())
	t.Wait()
	return t.Error()
}

func (a Agent) On(ctx context.Context, code string) error {
	dev, mcli, onoff, err := a.get(code)
	if err != nil {
		return nil
	}

	err = mcli.DoOn(ctx, onoff,
		DeviceData{
			Record: dev,
			Sender: a.Sender,
			Global: a.Global,
			IDB:    a.IDB,
			Logger: a.Logger},
	)
	if err != nil {
		return err
	}

	oc := OptChange{
		Box:       a.Global.Box,
		Code:      dev.Code,
		Type:      dev.Typ,
		DataCode:  common.NewDataCode(),
		Op:        common.OFF,
		OptStatus: common.OPT_STATUS_OFF,
		DataTime:  time.Now(),
	}
	return a.Sender.SendData(ctx, oc)
}

func (a Agent) Off(ctx context.Context, code string) error {
	dev, mcli, onoff, err := a.get(code)
	if err != nil {
		return nil
	}

	err = mcli.DoOff(ctx, onoff,
		DeviceData{
			Record: dev,
			Sender: a.Sender,
			Global: a.Global,
			IDB:    a.IDB,
			Logger: a.Logger},
	)
	if err != nil {
		return err
	}
	oc := OptChange{
		Box:       a.Global.Box,
		Code:      dev.Code,
		Type:      dev.Typ,
		DataCode:  common.NewDataCode(),
		Op:        common.OFF,
		OptStatus: common.OPT_STATUS_OFF,
		DataTime:  time.Now(),
	}
	return a.Sender.SendData(ctx, oc)
}

func (a Agent) get(devCode string) (*ent.Dev, *ModbusCli, OnOffer, error) {
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

func (a Agent) msgHandle() mqtt.MessageHandler {
	return func(cli mqtt.Client, msg mqtt.Message) {
		if msg.Duplicate() {
			return
		}
		defer msg.Ack()

		boxCode, devCode, op := topicPart(msg.Topic())

		if cmp.Compare(boxCode, a.Global.Box) != 0 {
			a.Logger.Error("onoff box code error",
				slog.String("box", boxCode),
				slog.String("code", devCode),
				slog.String("global.Box", a.Global.Box),
				slog.String("op", op),
			)
		}

		ctx := context.Background()

		switch op {
		case common.ON:
			if err := a.On(ctx, devCode); err != nil {
				a.Logger.Error("onoff error", slog.String("code", devCode),
					slog.String("op", op),
					slog.Any("error", err),
				)
			}
		case common.OFF:
			if err := a.Off(ctx, devCode); err != nil {
				a.Logger.Error("onoff error", slog.String("code", devCode),
					slog.String("op", op),
					slog.Any("error", err),
				)
				return
			}
		}
	}
}
