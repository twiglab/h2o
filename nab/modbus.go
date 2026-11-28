package nab

import (
	"context"
	"errors"
	"io"
	"net"
	"sync"
	"syscall"
	"time"

	"github.com/simonvetter/modbus"
	"github.com/twiglab/h2o/nab/orm/ent"
	"github.com/twiglab/h2o/proto"
)

type DeviceCollector interface {
	Collect(ctx context.Context, op ModbusOperator) (proto.MeterValue, error)
}

type DeviceOn interface {
	On(ctx context.Context, op ModbusOperator) error
}

type DeviceOff interface {
	Off(ctx context.Context, op ModbusOperator) error
}

// 初始化设备地址，字节顺序，大小端
func initDev(cli *modbus.ModbusClient, e uint, o uint, u uint8) error {
	if e != 0 && o != 0 {
		if err := cli.SetEncoding(Endian(e), WordOrder(o)); err != nil {
			return err
		}
	}
	return cli.SetUnitId(u)
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

func Endian(v uint) modbus.Endianness {
	if v == 2 {
		return modbus.LITTLE_ENDIAN
	}
	return modbus.BIG_ENDIAN
}

func WordOrder(v uint) modbus.WordOrder {
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

	e := Endian(record.Endian)
	wo := WordOrder(record.WordOrder)
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

func (c *ModbusCli) Close() error {
	return c.close()
}

func (c *ModbusCli) setup() error {
	// 已经打开，直接返回, 避免产生新连接
	if c.isOpen {
		return nil
	}
	return c.open()
}

func (c *ModbusCli) reset() {
	if c.isOpen {
		// 恢复client设置
		_ = c.cli.SetEncoding(c.endian, c.wordorder)
	}
}

func (c *ModbusCli) Do(data *ent.Dev, f MBFunc) (mv proto.MeterValue, err error) {
	c.glock.Lock()
	defer c.glock.Unlock()

	if err = c.setup(); err != nil {
		return
	}
	defer c.reset()

	if err = initDev(c.cli, data.Endian, data.WordOrder, data.UnitID); err != nil {
		return
	}

	mv, err = f(mbcli{c.cli})
	if err != nil {
		if _, ok := isConnErr(err); ok {
			_ = c.close()
		}
	}
	return
}

type MBFunc func(ModbusOperator) (proto.MeterValue, error)

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

type PortOperator interface {
	Do(*ent.Dev, MBFunc) (proto.MeterValue, error)
	Close() error
}
