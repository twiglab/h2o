package kh

import (
	"context"
	"time"

	"github.com/avast/retry-go/v3"
	"github.com/simonvetter/modbus"
	"github.com/twiglab/h2o/nab"
	"github.com/twiglab/h2o/pkg/common"
)

// const KH_DTZY2851 = "KH-DTZY2851"
const KH_DTZY2851_CD7_1 = "KH-DTZY2851-CD7-1"

type DTZY struct {
}

func (e DTZY) Collect(ctx context.Context, cli *modbus.ModbusClient, data nab.DeviceData) error {
	optStatus, err := cli.ReadRegister(0x22, modbus.INPUT_REGISTER)
	if err != nil {
		return err
	}

	dataVal, err := cli.ReadUint32(0x60, modbus.INPUT_REGISTER)
	if err != nil {
		return err
	}

	now := time.Now()

	var meter nab.Meter
	meter.Data.DataValue = int64(dataVal)
	meter.Data.OptStatus = int64(optStatus)

	meter.Code = data.Record.Code
	meter.Type = data.Record.Typ
	meter.DataCode = common.NewDataCode()
	meter.DataTime = now
	meter.DataTs = common.Ts(now)

	meter.Pos.Project = data.Global.Project

	meter.Gateway.Code = data.Global.BoxCode
	meter.Gateway.Type = common.GATEWAY_NH
	meter.Gateway.UnitID = data.Record.UnitID

	return data.Sender.SendData(ctx, meter)
}

func (e DTZY) On(ctx context.Context, cli *modbus.ModbusClient, data nab.DeviceData) error {
	err := retry.Do(
		func() error {
			return cli.WriteRegisters(0x22, []uint16{0x2})
		},
		retry.Attempts(3),
	)
	return err
}

func (e DTZY) Off(ctx context.Context, cli *modbus.ModbusClient, data nab.DeviceData) error {
	err := retry.Do(
		func() error {
			return cli.WriteRegisters(0x22, []uint16{0x1})
		},
		retry.Attempts(3),
	)
	return err
}
