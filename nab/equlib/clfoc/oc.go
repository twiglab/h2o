package clfoc

import (
	"context"
	"time"

	"github.com/simonvetter/modbus"
	"github.com/twiglab/h2o/nab"
	"github.com/twiglab/h2o/pkg/common"
)

const CLF_OC_PT = "CLF-OC-PT"

type CLFOC struct {
}

func (e CLFOC) Collect(ctx context.Context, cli *modbus.ModbusClient, data nab.DeviceData) error {

	dataVal, err := cli.ReadUint32(data.Record.Addr, modbus.INPUT_REGISTER)
	if err != nil {
		return err
	}

	now := time.Now()

	var meter nab.Meter
	meter.Data.DataValue = int64(dataVal)
	meter.Data.OptStatus = 0x00

	meter.Code = data.Record.Code
	meter.Type = data.Record.Typ
	meter.DataCode = common.NewDataCode()
	meter.DataTime = now
	meter.DataTs = common.Ts(now)

	meter.Pos.Project = data.Config.Project

	meter.Gateway.Code = data.Config.BoxCode
	meter.Gateway.Type = common.GATEWAY_NH
	meter.Gateway.UnitID = data.Record.UnitID

	return data.Sender.SendData(ctx, meter)
}

func (e CLFOC) On(ctx context.Context, cli *modbus.ModbusClient, data nab.DeviceData) error {
	return nil
}

func (e CLFOC) Off(ctx context.Context, cli *modbus.ModbusClient, data nab.DeviceData) error {
	return nil
}
