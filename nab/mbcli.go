package nab

import (
	"github.com/simonvetter/modbus"
	"github.com/twiglab/h2o/nab/orm/ent"
	"github.com/twiglab/h2o/proto"
)

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

// nullOptx 终端仅用于无需客户端的设备
// 例如 TimeTick设备
// 如果配置问题，导致错误使用的nullOptx，直接panic，无需处理
type nullOptx struct {
}

func (c nullOptx) Close() error { return nil }
func (c nullOptx) Do(_ *ent.Dev, f MBFunc) (proto.MeterValue, error) {
	return f(nil)
}

var nullOpt nullOptx
