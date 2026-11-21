package kh

import (
	"context"
	"time"

	"github.com/avast/retry-go/v5"
	"github.com/simonvetter/modbus"
	"github.com/twiglab/h2o/nab"
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

	var meter nab.Meter
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
	err := retry.New(
		retry.Attempts(3),
		retry.Delay(300*time.Millisecond), // 给串口留出帧间隔 retry.DelayType(retry.CombineDelay(retry.BackOffDelay, retry.RandomDelay)),
		retry.Context(ctx),                // 支持取消
		retry.LastErrorOnly(true),         // 只返回最后一次错误，日志更干净
	).Do(
		func() error {
			return cli.WriteRegisters(0x22, []uint16{0x02})
		},
	)
	return err
}

func (e DTZY) Off(ctx context.Context, cli *modbus.ModbusClient, data nab.DeviceData) error {
	err := retry.New(
		retry.Attempts(3),
		retry.Delay(300*time.Millisecond), // 给串口留出帧间隔 retry.DelayType(retry.CombineDelay(retry.BackOffDelay, retry.RandomDelay)),
		retry.Context(ctx),                // 支持取消
		retry.LastErrorOnly(true),         // 只返回最后一次错误，日志更干净
	).Do(
		func() error {
			// 宽厚公司电表文档明确要求写操作必须用FC16
			return cli.WriteRegisters(0x22, []uint16{0x01})
		},
	)
	return err
}
