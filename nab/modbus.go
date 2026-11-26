package nab

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net"
	"sync"
	"syscall"
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

	endian    modbus.Endianness
	wordorder modbus.WordOrder

	cfg *modbus.ClientConfiguration

	isOpen bool
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

func NewModbusCli(record *ent.Cli) (client *ModbusCli) {
	timeout := time.Second
	if record.Timeout > 0 && record.Timeout <= 10000 { // 最大10秒
		timeout = time.Duration(record.Timeout) * time.Millisecond
	}

	cfg := &modbus.ClientConfiguration{
		URL:      record.URL,
		Speed:    record.Speed,
		DataBits: record.DataBits,
		Parity:   record.Parity,
		StopBits: record.StopBits,
		Timeout:  timeout,
	}

	e := endian(record.Endian)
	wo := wordorder(record.WordOrder)
	return &ModbusCli{
		Record: record,
		Code:   record.Code,

		endian:    e,
		wordorder: wo,

		cfg: cfg,

		isOpen: false,
	}

}

func (c *ModbusCli) open() error {
	// 已经打开，直接返回, 避免产生新连接
	if c.isOpen {
		return nil
	}

	cli, err := modbus.NewClient(c.cfg)
	if err != nil {
		return err
	}
	if err = cli.SetEncoding(c.endian, c.wordorder); err != nil {
		return err
	}
	if err = cli.Open(); err != nil {
		return err
	}

	// 这里设置ModbusClient
	// Close 中 cli已经被设置成nil
	c.cli = cli
	c.isOpen = true
	return nil
}

func (c *ModbusCli) close() error {
	defer func() {
		c.isOpen = false
		c.cli = nil // 释放 cli
	}()

	if c.cli == nil {
		return nil
	}
	return c.cli.Close()
}

func (c *ModbusCli) setup(data DeviceData) error {
	// Open方法会判断isOpen条件
	if err := c.open(); err != nil {
		return err
	}

	// 必须两者都不为0才设置，如果只有一个为0，那就不设置，最大容忍配置错误
	// 此问题已经被记录，留到v2版本中去解决
	// TODO
	if data.Record.Endian != 0 && data.Record.WordOrder != 0 {
		if err := c.cli.SetEncoding(endian(data.Record.Endian), wordorder(data.Record.WordOrder)); err != nil {
			return err
		}
	}

	return c.cli.SetUnitId(data.Record.UnitID)
}

func (c *ModbusCli) reset() error {
	if c.isOpen {
		// 恢复client设置
		_ = c.cli.SetEncoding(c.endian, c.wordorder)
	}
	return nil
}

func (c *ModbusCli) DoCollect(ctx context.Context, coll DeviceCollector, data DeviceData) (err error) {
	c.glock.Lock()
	defer c.glock.Unlock()

	if err = c.setup(data); err != nil {
		return
	}

	if err = coll.Collect(ctx, c.cli, data); err != nil {
		if _, ok := isConnErr(err); ok {
			_ = c.close()
		}
	}
	_ = c.reset()
	return
}

func (c *ModbusCli) DoOn(ctx context.Context, on DeviceOn, data DeviceData) (err error) {
	c.glock.Lock()
	defer c.glock.Unlock()

	if err = c.setup(data); err != nil {
		return err
	}

	if err = on.On(ctx, c.cli, data); err != nil {
		if _, ok := isConnErr(err); ok {
			_ = c.close()
		}
	}
	_ = c.reset()
	return
}

func (c *ModbusCli) DoOff(ctx context.Context, off DeviceOff, data DeviceData) (err error) {
	c.glock.Lock()
	defer c.glock.Unlock()

	if err = c.setup(data); err != nil {
		return err
	}
	if err = off.Off(ctx, c.cli, data); err != nil {
		if _, ok := isConnErr(err); ok {
			_ = c.close()
		}
	}
	_ = c.reset()
	return

}

func isConnErr(err error) (error, bool) {
	if err == nil {
		return err, false
	}
	if neterr, ok := err.(net.Error); ok {
		return neterr, ok
	}
	return err,
		errors.Is(err, modbus.ErrRequestTimedOut) ||
			errors.Is(err, io.EOF) ||
			errors.Is(err, syscall.ECONNRESET)
}

type DeviceData struct {
	Global Global
	Record *ent.Dev
	Sender Sender
	IDB    *orm.IDB
	Logger *slog.Logger
}

type PortOperator interface {
	DoCollect(ctx context.Context, coll DeviceCollector, data DeviceData) (err error)
	DoOn(ctx context.Context, on DeviceOn, data DeviceData) (err error)
	DoOff(ctx context.Context, off DeviceOff, data DeviceData) (err error)
}

// nullOptx 终端仅用于无需客户端的设备
// 例如 TimeTick设备
// 如果配置问题，导致错误使用的nullOptx，直接panic，没有也无需任他任何保证
type nullOptx struct {
}

func (c nullOptx) DoCollect(ctx context.Context, coll DeviceCollector, data DeviceData) error {
	return coll.Collect(ctx, nil, data)
}

func (c nullOptx) DoOn(ctx context.Context, on DeviceOn, data DeviceData) error {
	return on.On(ctx, nil, data)
}

func (c nullOptx) DoOff(ctx context.Context, off DeviceOff, data DeviceData) error {
	return off.Off(ctx, nil, data)
}

var nullOpt nullOptx
