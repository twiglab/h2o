package fake

import (
	"context"
	"log/slog"
	"time"

	"github.com/simonvetter/modbus"
	"github.com/twiglab/h2o/bak/pkg/common"
	"github.com/twiglab/h2o/nab"
	"github.com/twiglab/h2o/proto"
)

type Fake struct {
	optStatus  int64
	alwaysOpen bool
}

func New(alwaysOpen bool) *Fake {
	return &Fake{
		optStatus:  common.OPT_STATUS_ON,
		alwaysOpen: alwaysOpen,
	}
}

func (e Fake) Collect(ctx context.Context, _ *modbus.ModbusClient, data nab.DeviceData) error {

	now := time.Now()

	var meter nab.Meter
	meter.Data.DataValue = now.Unix()
	meter.Data.OptStatus = e.optStatus

	meter.Code = data.Record.Code
	meter.Type = data.Record.Typ
	meter.DataCode = proto.NewDataCode()
	meter.DataTime = now

	meter.Pos.Project = data.Global.Project

	meter.Gateway.Code = data.Global.Box
	meter.Gateway.Type = proto.GATEWAY_NH
	meter.Gateway.UnitID = data.Record.UnitID

	return data.Sender.SendData(ctx, meter)
}

func (e *Fake) On(ctx context.Context, cli *modbus.ModbusClient, data nab.DeviceData) error {
	data.Logger.DebugContext(ctx, "FAKE DEVICE ON", slog.String("boxCode", data.Global.Box), slog.Any("Record", data.Record))
	e.optStatus = proto.OPT_STATUS_ON
	return nil
}

func (e *Fake) Off(ctx context.Context, cli *modbus.ModbusClient, data nab.DeviceData) error {
	data.Logger.DebugContext(ctx, "FAKE DEVICE OFF", slog.String("boxCode", data.Global.Box), slog.Any("Record", data.Record))
	if !e.alwaysOpen {
		e.optStatus = proto.OPT_STATUS_OFF
	}
	return nil
}
