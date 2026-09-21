package nab

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/simonvetter/modbus"
	"github.com/twiglab/h2o/nab/orm"
	"github.com/twiglab/h2o/nab/orm/ent"
)

type DeviceCollector interface {
	Collect(ctx context.Context, cli *modbus.ModbusClient, data DeviceData) error
}

type DeviceOn interface {
	On(ctx context.Context, cli *modbus.ModbusClient, data DeviceData) error
}

type DeviceOff interface {
	Off(ctx context.Context, cli *modbus.ModbusClient, data DeviceData) error
}

/*
	mc.unitId     = 1
	mc.endianness = BIG_ENDIAN
	mc.wordOrder  = HIGH_WORD_FIRST
*/

type ModbusCli struct {
	Record *ent.Cli
	cli    *modbus.ModbusClient
	glock  sync.Mutex

	Code string

	endina    modbus.Endianness
	wordorder modbus.WordOrder
}

func endian(v uint) modbus.Endianness {
	if v == 2 {
		return modbus.LITTLE_ENDIAN
	}
	return modbus.BIG_ENDIAN
}

func wordorder(v uint) modbus.WordOrder {
	if v == 2 {
		return modbus.LOW_WORD_FIRST
	}
	return modbus.HIGH_WORD_FIRST
}

func NewModbusCli(record *ent.Cli) (client *ModbusCli, err error) {
	cfg := &modbus.ClientConfiguration{
		URL:      record.URL,
		Speed:    record.Speed,
		DataBits: record.DataBits,
		Parity:   record.Parity,
		StopBits: record.StopBits,
		Timeout:  time.Second,
	}

	cli, err := modbus.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	e := endian(record.Endian)
	wo := wordorder(record.WordOrder)

	if err = cli.SetEncoding(e, wo); err != nil {
		return nil, err
	}

	if err = cli.Open(); err != nil {
		return nil, err
	}

	return &ModbusCli{
		Record: record,
		cli:    cli,
		Code:   record.Code,

		endina:    e,
		wordorder: wo,
	}, nil
}

func (c *ModbusCli) setup(data DeviceData) error {
	c.glock.Lock()

	if data.Record.Endian != 0 && data.Record.WordOrder != 0 {
		if err := c.cli.SetEncoding(modbus.Endianness(data.Record.Endian), modbus.WordOrder(data.Record.WordOrder)); err != nil {
			return err
		}
	}
	if err := c.cli.SetUnitId(data.Record.UnitID); err != nil {
		return err
	}
	return nil
}
func (c *ModbusCli) reset() error {
	defer c.glock.Unlock()
	_ = c.cli.SetEncoding(c.endina, c.wordorder)
	_ = c.cli.SetUnitId(1)
	return nil
}

func (c *ModbusCli) DoCollect(ctx context.Context, coll DeviceCollector, data DeviceData) error {
	if err := c.setup(data); err != nil {
		return err
	}
	defer c.reset()

	return coll.Collect(ctx, c.cli, data)
}

func (c *ModbusCli) DoOn(ctx context.Context, on DeviceOn, data DeviceData) error {
	if err := c.setup(data); err != nil {
		return err
	}
	defer c.reset()

	return on.On(ctx, c.cli, data)
}

func (c *ModbusCli) DoOff(ctx context.Context, off DeviceOff, data DeviceData) error {
	if err := c.setup(data); err != nil {
		return err
	}
	defer c.reset()

	return off.Off(ctx, c.cli, data)
}

func (c *ModbusCli) DoFunc(ctx context.Context, data DeviceData, f func(ctx context.Context, cli *modbus.ModbusClient, data DeviceData) error) error {
	return f(ctx, c.cli, data)
}

type DeviceData struct {
	Global Global
	Record *ent.Dev
	Sender Sender
	IDB    *orm.IDB
	Logger *slog.Logger
}
