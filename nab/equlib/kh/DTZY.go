package kh

import (
	"context"
	"time"

	"github.com/simonvetter/modbus"
	"github.com/twiglab/h2o/nab"
	"github.com/twiglab/h2o/nab/equlib"
	"github.com/twiglab/h2o/proto"
)

const KH_DTZY2851_CD7_1 = "KH-DTZY2851-CD7-1"

// DTZY 电表，0x60地址读取有功总电能（即表显）
// 0x22 地址读取开合状态，写入执行对应的操作（0x01拉闸, 0x02合闸）
// 注意，对接文档要求，写入必须用modbus 0x10 （WriteRegisters）写入
type DTZY struct {
}

func (e DTZY) Collect(ctx context.Context, cli *modbus.ModbusClient, data nab.DeviceData) error {
	dataVal, err := cli.ReadUint32(0x60, modbus.INPUT_REGISTER)
	if err != nil {
		return err
	}

	// 这里阻塞了循环采集，但在实际项目中验证，无影响
	// 后续根据程序实际运行情况，考虑减少休眠时间
	// 堵塞问题留存到v2 版本中解决，目前此逻辑维持不变
	// TODO
	time.Sleep(300 * time.Millisecond)

	opt, err := cli.ReadRegister(0x22, modbus.INPUT_REGISTER)
	if err != nil {
		return err
	}

	now := time.Now()

	meter := nab.NewMeter()
	meter.Data.DataValue = int64(dataVal)
	meter.Data.OptStatus = optStatus(opt)

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

func (e DTZY) On(ctx context.Context, cli *modbus.ModbusClient, data nab.DeviceData) error {
	return equlib.TryDo(ctx, 3, 300*time.Millisecond,
		func() error {
			return cli.WriteRegisters(0x22, []uint16{0x02})
		},
	)
}

func (e DTZY) Off(ctx context.Context, cli *modbus.ModbusClient, data nab.DeviceData) error {
	return equlib.TryDo(ctx, 3, 300*time.Millisecond,
		func() error {
			return cli.WriteRegisters(0x22, []uint16{0x01})
		},
	)
}
