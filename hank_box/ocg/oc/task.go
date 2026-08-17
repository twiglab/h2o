package oc

import (
	"context"
	"time"

	"github.com/simonvetter/modbus"
	"github.com/twiglab/h2o/hank"
	"github.com/twiglab/h2o/hank/pick"
)

type TaskX struct {
	Client  *modbus.ModbusClient
	Addr    uint16
	RegType modbus.RegType
	Sender  hank.Sender

	Code string
	Type string

	Project string
}

func (t *TaskX) Run(ctx context.Context) error {
	val, err := t.Client.ReadUint32(t.Addr, t.RegType)
	if err != nil {
		return err
	}

	now := time.Now()

	var em hank.ElectricityMeter
	em.Data.DataValue = int64(val)
	em.DataCode = pick.NewDataCode()
	em.DataTime = now
	em.DataTs = now.Format(f)
	em.Code = t.Code
	em.Type = t.Type

	em.Pos.Project = t.Project

	return t.Sender.SendData(ctx, em)
}

func NewModbusClient(url string) (*modbus.ModbusClient, error) {
	var client *modbus.ModbusClient
	var err error

	client, err = modbus.NewClient(&modbus.ClientConfiguration{
		// URL:     "tcp://192.168.0.100:502",
		URL:     url,
		Timeout: 3 * time.Second,
	})

	if err != nil {
		return client, err
	}

	if err = client.SetEncoding(modbus.BIG_ENDIAN, modbus.LOW_WORD_FIRST); err != nil {
		return client, err
	}

	err = client.Open()
	return client, err
}

const f = "20060102150405"

func ClientID(code string) string {
	now := time.Now()
	ts := now.Format(f)
	return code + "@" + ts
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

