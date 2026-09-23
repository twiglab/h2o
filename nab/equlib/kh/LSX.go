package kh

import (
	"context"
	"time"

	"github.com/avast/retry-go/v5"
	"github.com/simonvetter/modbus"
	"github.com/twiglab/h2o/nab"
	"github.com/twiglab/h2o/pkg/common"
)

const KH_LXS_40_MI = "KH-LXS-40-MI"

func optStatus(op uint16) int64 {
	if op == 0x00 {
		return common.OPT_STATUS_OFF
	}
	return common.OPT_STATUS_ON
}

type LXS struct {
}

func (e LXS) Collect(ctx context.Context, cli *modbus.ModbusClient, data nab.DeviceData) error {
	dataVal, err := cli.ReadUint32(0x00, modbus.INPUT_REGISTER)
	if err != nil {
		return err
	}

	time.Sleep(300 * time.Millisecond)

	op, err := cli.ReadRegister(0x02, modbus.INPUT_REGISTER)
	if err != nil {
		return err
	}

	now := time.Now()

	var meter nab.Meter
	meter.Data.DataValue = int64(dataVal)
	meter.Data.OptStatus = optStatus(op)

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
	err := retry.New(retry.Attempts(3)).Do(
		func() error {
			return cli.WriteRegister(0x02, 0xFF)
		},
	)
	return err
}

func (e LXS) Off(ctx context.Context, cli *modbus.ModbusClient, data nab.DeviceData) error {
	err := retry.New(retry.Attempts(3)).Do(
		func() error {
			return cli.WriteRegister(0x02, 0x00)
		},
	)
	return err
}
