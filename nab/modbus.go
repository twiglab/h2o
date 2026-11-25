package nab

import (
	"context"
	"errors"
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

		endina:    e,
		wordorder: wo,

		cfg: cfg,

		isOpen: false,
	}

}

func (c *ModbusCli) Open() error {
	if c.isOpen {
		return nil
	}

	cli, err := modbus.NewClient(c.cfg)
	if err != nil {
		return err
	}
	if err = cli.SetEncoding(c.endina, c.wordorder); err != nil {
		return err
	}
	if err = cli.Open(); err != nil {
		return err
	}
	c.cli = cli // 这里设置ModbusClient
	c.isOpen = true
	return nil
}

func (c *ModbusCli) Close() error {
	c.isOpen = false
	return c.cli.Close()
}

func (c *ModbusCli) setup(data DeviceData) error {
	// open 方法会判断isOpen条件
	if err := c.Open(); err != nil {
		return err
	}

	// 必须两者都为0才设置，如果只有一个为0，那就不设置，最大容忍配置错误
	// 此问题已经被记录，留到v2版本中去解决
	// TODO
	if data.Record.Endian != 0 && data.Record.WordOrder != 0 {
		if err := c.cli.SetEncoding(modbus.Endianness(data.Record.Endian), modbus.WordOrder(data.Record.WordOrder)); err != nil {
			return err
		}
	}

	return c.cli.SetUnitId(data.Record.UnitID)
}

func (c *ModbusCli) reset() error {
	// 恢复client设置
	_ = c.cli.SetEncoding(c.endina, c.wordorder)
	return nil
}

func (c *ModbusCli) DoCollect(ctx context.Context, coll DeviceCollector, data DeviceData) (err error) {
	c.glock.Lock()
	defer c.glock.Unlock()

	if err = c.setup(data); err != nil {
		return
	}

	if err = coll.Collect(ctx, c.cli, data); err != nil {
		//  底层库只提供的这个错误用于判断modbus连接是否有问题
		if errors.Is(err, modbus.ErrRequestTimedOut) {
			// 这里就直接认为modbus连接除了问题，直接关闭
			// 其余的事情交给setup方法去处理
			_ = c.Close()
		}
	}

	// 即便连接是关闭的，调用reset代码也不会有影响
	_ = c.reset()
	return
}

func (c *ModbusCli) DoOn(ctx context.Context, on DeviceOn, data DeviceData) (err error) {
	c.glock.Lock()
	defer c.glock.Unlock()

	if err = c.setup(data); err != nil {
		return err
	}

	err = on.On(ctx, c.cli, data)
	_ = c.reset()
	return
}

func (c *ModbusCli) DoOff(ctx context.Context, off DeviceOff, data DeviceData) (err error) {
	c.glock.Lock()
	defer c.glock.Unlock()

	if err = c.setup(data); err != nil {
		return err
	}

	err = off.Off(ctx, c.cli, data)
	_ = c.reset()
	return

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
