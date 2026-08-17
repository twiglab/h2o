package pick

import (
	"context"
	"time"

	"github.com/simonvetter/modbus"
	"github.com/twiglab/h2o/hank"
)

type TaskX struct {
	Client  *modbus.ModbusClient
	Addr    uint16
	RegType modbus.RegType
	Sender  hank.Sender

	SN   string
	Code string
	Type string

	Project string
	PosCode string
}

func (t *TaskX) Run(ctx context.Context) error {
	val, err := t.Client.ReadUint32(t.Addr, t.RegType)
	if err != nil {
		return err
	}

	var em hank.ElectricityMeter
	em.Data.DataValue = int64(val)
	em.DataCode = NewDataCode()
	em.DataTime = time.Now()
	em.SN = t.SN
	em.Code = t.Code
	em.Type = hank.ELECTRICITY

	return t.Sender.SendData(ctx, em)
}
