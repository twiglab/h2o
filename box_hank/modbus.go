package box

import (
	"context"
	"log/slog"
	"time"

	"github.com/simonvetter/modbus"
	"github.com/twiglab/h2o/pkg/common"
)

func NewModbusClient(config *modbus.ClientConfiguration) (client *modbus.ModbusClient, err error) {
	if client, err = modbus.NewClient(config); err != nil {
		return
	}

	err = client.Open()
	return
}

type ModbusTask struct {
	Client  *modbus.ModbusClient
	Addr    uint16
	RegType modbus.RegType
	UnitID  uint8
	Sender  Sender

	Code string
	Type string

	Project string

	Ctx context.Context

	Logger *slog.Logger
}

func (t *ModbusTask) Run() {
	if err := t.Client.SetUnitId(t.UnitID); err != nil {
		return
	}

	val, err := t.Client.ReadUint32(t.Addr, t.RegType)
	if err != nil {
		return
	}

	now := time.Now()

	var meter Meter
	meter.Data.DataValue = int64(val)
	meter.DataCode = common.NewDataCode()
	meter.DataTime = now
	meter.DataTs = common.Ts(now)
	meter.Code = t.Code
	meter.Type = t.Type
	meter.Pos.Project = t.Project

	if err := t.Sender.SendData(t.Ctx, meter); err != nil {
		return
	}
}
