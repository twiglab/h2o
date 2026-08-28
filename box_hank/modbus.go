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
	Client *modbus.ModbusClient

	Addr      uint16
	UnitID    uint8
	RegType   modbus.RegType
	Endian    modbus.Endianness
	WordOrder modbus.WordOrder

	Sender Sender

	Code string
	Type string

	Project string

	Ctx context.Context

	Logger *slog.Logger
}

func (t *ModbusTask) Run() {
	if t.Endian != 0 && t.WordOrder != 0 {
		if err := t.Client.SetEncoding(t.Endian, t.WordOrder); err != nil {
			t.Logger.Error("SetEncoding error", slog.String("code", t.Code), slog.Any("error", err))
			return
		}
	}

	if err := t.Client.SetUnitId(t.UnitID); err != nil {
		t.Logger.Error("SetUnitID error", slog.String("code", t.Code), slog.Any("error", err))
		return
	}

	val, err := t.Client.ReadUint32(t.Addr, t.RegType)
	if err != nil {
		t.Logger.Error("Read error", slog.String("code", t.Code), slog.Any("error", err))
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
		t.Logger.Error("SendData error", slog.String("code", t.Code), slog.Any("error", err))
	}
}
