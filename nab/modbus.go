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
	"github.com/twiglab/h2o/nab/orm/ent"
	"github.com/twiglab/h2o/proto"
)

type DeviceCollector interface {
	Collect(ctx context.Context, op ModbusOperator, data DeviceData) (proto.MeterValue, error)
}

type DeviceOn interface {
	On(ctx context.Context, op ModbusOperator, data DeviceData) error
}

type DeviceOff interface {
	Off(ctx context.Context, op ModbusOperator, data DeviceData) error
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

func (c *ModbusCli) Do(f func(ModbusOperator) error) (err error) {
	c.glock.Lock()
	defer c.glock.Unlock()

	if err = c.setup(); err != nil {
		return err
	}
	defer c.reset()

	if err = f(mbcli{c.cli}); err != nil {
		if _, ok := isConnErr(err); ok {
			_ = c.close()
		}
	}
	return
}

type ModbusOperator interface {
	FC03(addr uint16) (uint16, error)
	FC04(addr uint16) (uint16, error)

	FC03Uint32(addr uint16) (uint32, error)
	FC04Uint32(addr uint16) (uint32, error)

	FC16(addr uint16, values []uint16) error
	FC06(addr uint16, value uint16) error
}

type mbcli struct {
	cli *modbus.ModbusClient
}

func (c mbcli) FC03(addr uint16) (v uint16, err error) {
	v, err = c.cli.ReadRegister(addr, modbus.HOLDING_REGISTER)
	return
}

func (c mbcli) FC04(addr uint16) (v uint16, err error) {
	v, err = c.cli.ReadRegister(addr, modbus.INPUT_REGISTER)
	return
}

func (c mbcli) FC03Uint32(addr uint16) (v uint32, err error) {
	v, err = c.cli.ReadUint32(addr, modbus.HOLDING_REGISTER)
	return
}

func (c mbcli) FC04Uint32(addr uint16) (v uint32, err error) {
	v, err = c.cli.ReadUint32(addr, modbus.INPUT_REGISTER)
	return
}

func (c mbcli) FC06(addr uint16, value uint16) error {
	return c.cli.WriteRegister(addr, value)
}

func (c mbcli) FC16(addr uint16, values []uint16) error {
	return c.cli.WriteRegisters(addr, values)
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
	Box    *Box
	Record *ent.Dev
	Logger *slog.Logger
}

type PortOperator interface {
	Do(func(ModbusOperator) error) error
	Close() error
}

// nullOptx 终端仅用于无需客户端的设备
// 例如 TimeTick设备
// 如果配置问题，导致错误使用的nullOptx，直接panic，没有也无需任他任何保证
type nullOptx struct {
}

func (c nullOptx) Close() error { return nil }
func (c nullOptx) Do(f func(ModbusOperator) error) error {
	return f(nil)
}

var nullOpt nullOptx
