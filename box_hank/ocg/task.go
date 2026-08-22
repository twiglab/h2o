package ocg

import (
	"context"
	"time"

	"github.com/simonvetter/modbus"
	"github.com/twiglab/h2o/box"
	"github.com/twiglab/h2o/pkg/common"
)

type TaskX struct {
	Client  *modbus.ModbusClient
	Addr    uint16
	RegType modbus.RegType
	Sender  box.Sender

	Code string
	Type string

	Project string
}

func (t *TaskX) Run(ctx context.Context) error {
	if err := t.Client.SetEncoding(modbus.BIG_ENDIAN, modbus.LOW_WORD_FIRST); err != nil {
		return err
	}

	val, err := t.Client.ReadUint32(t.Addr, t.RegType)
	if err != nil {
		return err
	}

	now := time.Now()

	var em box.ElectricityMeter
	em.Data.DataValue = int64(val)
	em.DataCode = common.NewDataCode()
	em.DataTime = now
	em.DataTs = common.Ts(now)
	em.Code = t.Code
	em.Type = t.Type

	em.Pos.Project = t.Project

	return t.Sender.SendData(ctx, em)
}

func TaskChain(ctx context.Context, t ...*TaskX) error {
	var err error
	for _, x := range t {
		if err = x.Run(ctx); err != nil {
			return err
		}
	}
	return nil
}
