package kh

import (
	"context"
	"time"

	"github.com/simonvetter/modbus"
	"github.com/twiglab/h2o/nab"
	"github.com/twiglab/h2o/proto"
)

const KH_LXS_40_MI = "KH-LXS-40-MI"

// LXS 水表，0x00地址读取用量（即表显）
type LXSNoValve struct {
}

func (e LXSNoValve) Collect(ctx context.Context, cli *modbus.ModbusClient, data nab.DeviceData) error {
	dataVal, err := cli.ReadUint32(0x00, modbus.INPUT_REGISTER)
	if err != nil {
		return err
	}

	now := time.Now()

	meter := nab.NewMeter()
	meter.Data.DataValue = int64(dataVal)

	// 无电磁阀水表，状态为常开
	meter.Data.OptStatus = proto.OPT_STATUS_ON

	meter.Code = data.Record.Code
	meter.Type = data.Record.Typ
	meter.SN = data.Record.Sn
	meter.DataCode = proto.NewDataCode()
	meter.DataTime = now

	meter.Pos.Project = data.Global.Project

	meter.Gateway.Code = data.Global.Box
	meter.Gateway.Type = proto.GATEWAY_NH
	meter.Gateway.UnitID = data.Record.UnitID

	return data.Sender.SendData(ctx, meter)
}

func (e LXSNoValve) On(ctx context.Context, cli *modbus.ModbusClient, data nab.DeviceData) error {
	// 无电磁阀水表，不支持开关操作
	return nil
}

func (e LXSNoValve) Off(ctx context.Context, cli *modbus.ModbusClient, data nab.DeviceData) error {
	// 无电磁阀水表，不支持开关操作
	return nil
}
