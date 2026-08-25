package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/simonvetter/modbus"
	"github.com/twiglab/h2o/box"
	"github.com/twiglab/h2o/pkg/common"
)

type TaskX struct {
	Client  *modbus.ModbusClient
	Addr    uint16
	RegType modbus.RegType
	UnitID  uint8
	Sender  box.Sender

	Code string
	Type string

	Project string
}

func (t *TaskX) Run(ctx context.Context) error {
	//if err := t.Client.SetUnitId(0x21); err != nil {
	if err := t.Client.SetUnitId(t.UnitID); err != nil {
		return err
	}

	//val, err := t.Client.ReadUint32(0, modbus.INPUT_REGISTER)
	val, err := t.Client.ReadUint32(t.Addr, t.RegType)
	if err != nil {
		return err
	}

	now := time.Now()

	var wm box.WaterMeter
	wm.Data.DataValue = int64(val)
	wm.DataCode = common.NewDataCode()
	wm.DataTime = now
	wm.DataTs = common.Ts(now)
	wm.Code = t.Code
	wm.Type = t.Type

	return t.Sender.SendData(ctx, wm)
}

func main() {
	client, err := modbus.NewClient(&modbus.ClientConfiguration{
		URL:      "rtu:///dev/ttyS0",
		Speed:    9600,               // default
		DataBits: 8,                  // default, optional
		Parity:   modbus.PARITY_NONE, // default, optional
		StopBits: 1,                  // default if no parity, optional
		Timeout:  5 * time.Second,
	})

	if err != nil {
		log.Fatal(err)
	}

	if err := client.Open(); err != nil {
		log.Fatal(err)
	}

	defer client.Close()

	client.SetUnitId(0x21)

	val, err := client.ReadUint32(0, modbus.INPUT_REGISTER)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(val)
}
