package fake

import (
	"context"
	"time"

	"github.com/simonvetter/modbus"
	"github.com/twiglab/h2o/nab"
	"github.com/twiglab/h2o/nab/equlib"
	"github.com/twiglab/h2o/proto"
)

func init() {
	equlib.Register("time_tick", newTimeTick())
}

// 滴答设备返回当前时间，用于测试
// 滴答设备不需要client
// 滴答设备类型固定为 proto.TYPE_TIME_TICK
type TimeTick struct {
	optStatus int64
}

func newTimeTick() *TimeTick {
	return &TimeTick{
		optStatus: proto.OPT_STATUS_ON,
	}
}

func (e TimeTick) Collect(ctx context.Context, _ *modbus.ModbusClient, data nab.DeviceData) error {

	now := time.Now()

	var meter nab.Meter
	meter.Data.DataValue = now.Unix()
	meter.Data.OptStatus = e.optStatus

	meter.Code = data.Record.Code
	meter.Type = proto.TYPE_TIME_TICK
	meter.DataCode = proto.NewDataCode()
	meter.DataTime = now

	meter.Pos.Project = data.Global.Project

	meter.Gateway.Code = data.Global.Box
	meter.Gateway.Type = proto.GATEWAY_NH
	meter.Gateway.UnitID = 0

	return data.Sender.SendData(ctx, meter)
}

func (e *TimeTick) On(_ context.Context, _ *modbus.ModbusClient, _ nab.DeviceData) error {
	e.optStatus = proto.OPT_STATUS_ON
	return nil
}

func (e *TimeTick) Off(_ context.Context, _ *modbus.ModbusClient, _ nab.DeviceData) error {
	e.optStatus = proto.OPT_STATUS_OFF
	return nil
}
