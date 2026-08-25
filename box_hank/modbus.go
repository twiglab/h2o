package box

import (
	"context"
	"log/slog"
	"time"

	"github.com/simonvetter/modbus"
	"github.com/twiglab/h2o/pkg/common"
)

func NewModbusClient(config *modbus.ClientConfiguration) (*modbus.ModbusClient, error) {
	var client *modbus.ModbusClient
	var err error

	if client, err = modbus.NewClient(config); err != nil {
		return client, err
	}

	return client, client.Open()
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

	var wm WaterMeter
	wm.Data.DataValue = int64(val)
	wm.DataCode = common.NewDataCode()
	wm.DataTime = now
	wm.DataTs = common.Ts(now)
	wm.Code = t.Code
	wm.Type = t.Type

	if err := t.Sender.SendData(t.Ctx, wm); err != nil {
		return
	}

	return
}
