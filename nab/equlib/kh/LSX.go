package kh

import (
	"context"
	"time"

	"github.com/simonvetter/modbus"
	"github.com/twiglab/h2o/nab"
	"github.com/twiglab/h2o/pkg/common"
)

const KH_LXS_40_MI = "KH-LXS-40-MI"

type LXS struct {
}

func (e LXS) Collect(ctx context.Context, cli *modbus.ModbusClient, data nab.DeviceData) error {

	dataVal, err := cli.ReadUint32(0x00, modbus.INPUT_REGISTER)
	if err != nil {
		return err
	}

	now := time.Now()

	var meter nab.Meter
	meter.Data.DataValue = int64(dataVal)
	meter.Data.OptStatus = common.OPT_STATUS_ON

	meter.Code = data.Record.Code
	meter.Type = data.Record.Typ
	meter.DataCode = common.NewDataCode()
	meter.DataTime = now
	meter.DataTs = common.Ts(now)

	meter.Pos.Project = data.Global.Project

	meter.Gateway.Code = data.Global.Box
	meter.Gateway.Type = common.GATEWAY_NH
	meter.Gateway.UnitID = data.Record.UnitID

	return data.Sender.SendData(ctx, meter)
}

func (e LXS) On(ctx context.Context, cli *modbus.ModbusClient, data nab.DeviceData) error {
	return nil
}

func (e LXS) Off(ctx context.Context, cli *modbus.ModbusClient, data nab.DeviceData) error {
	return nil
}
