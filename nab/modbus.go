package nab

import (
	"context"
	"log/slog"

	"github.com/simonvetter/modbus"
	"github.com/twiglab/h2o/nab/orm/ent"
)

type Collector interface {
	Collect(ctx context.Context, cli *modbus.ModbusClient, data DeviceData) error
}

type DeviceOn interface {
	On(ctx context.Context, cli *modbus.ModbusClient, data DeviceData) error
}

type DeviceOff interface {
	Off(ctx context.Context, cli *modbus.ModbusClient, data DeviceData) error
}

type ModbusCli struct {
	Record *ent.Cli
	cli    *modbus.ModbusClient
	// glock  sync.Mutex

	Code string
}

func NewModbusCli(record *ent.Cli) (client *ModbusCli, err error) {
	cfg := &modbus.ClientConfiguration{
		URL:      record.URL,
		Speed:    record.Speed,
		DataBits: record.DataBits,
		Parity:   record.Parity,
		StopBits: record.StopBits,
		// Timeout:  time.Millisecond * 300,
	}

	cli, err := modbus.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	err = cli.Open()

	return &ModbusCli{
		Record: record,
		cli:    cli,
		Code:   record.Code,
	}, err
}

func (c *ModbusCli) setup(data DeviceData) error {
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
	/*
		mc.unitId     = 1
		mc.endianness = BIG_ENDIAN
		mc.wordOrder  = HIGH_WORD_FIRST
	*/
	_ = c.cli.SetEncoding(modbus.BIG_ENDIAN, modbus.HIGH_WORD_FIRST)
	_ = c.cli.SetUnitId(1)
	return nil
}

func (c *ModbusCli) DoCollect(ctx context.Context, coll Collector, data DeviceData) error {
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
	Record *ent.Dev
	Sender Sender
	Config Global
	Logger *slog.Logger
}
