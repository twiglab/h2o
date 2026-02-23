package modbus

import (
	"encoding/binary"
	"errors"
	"math"

	"github.com/shopspring/decimal"
)

// 功能码定义
const (
	FuncReadCoils            = 0x01 // 读线圈
	FuncReadDiscreteInputs   = 0x02 // 读离散输入
	FuncReadHoldingRegisters = 0x03 // 读保持寄存器
	FuncReadInputRegisters   = 0x04 // 读输入寄存器
	FuncWriteSingleCoil      = 0x05 // 写单个线圈
	FuncWriteSingleRegister  = 0x06 // 写单个寄存器
	FuncWriteMultipleCoils   = 0x0F // 写多个线圈
	FuncWriteMultipleRegs    = 0x10 // 写多个寄存器
)

// 错误定义
var (
	ErrFrameTooShort   = errors.New("modbus frame too short")
	ErrInvalidCRC      = errors.New("invalid CRC checksum")
	ErrInvalidFunction = errors.New("invalid function code")
	ErrInvalidData     = errors.New("invalid data")
	ErrNoData          = errors.New("no data in response")
)

// ModbusFrame Modbus RTU 帧结构
type ModbusFrame struct {
	Address  byte   // 从站地址
	Function byte   // 功能码
	Data     []byte // 数据部分
}

// ParseFrame 解析 Modbus RTU 帧
// 帧结构: [地址(1B)] [功能码(1B)] [数据(NB)] [CRC(2B)]
func ParseFrame(data []byte) (*ModbusFrame, error) {
	// 最小长度: 地址(1) + 功能码(1) + CRC(2) = 4
	if len(data) < 4 {
		return nil, ErrFrameTooShort
	}

	// 验证 CRC
	if !ValidateCRC(data) {
		return nil, ErrInvalidCRC
	}

	frame := &ModbusFrame{
		Address:  data[0],
		Function: data[1],
	}

	// 提取数据部分 (去掉地址、功能码和CRC)
	if len(data) > 4 {
		frame.Data = data[2 : len(data)-2]
	}

	return frame, nil
}

// ReadResponseData 解析读取响应帧中的数据
// 响应格式: [地址][功能码][字节数][数据...][CRC]
type ReadResponseData struct {
	ByteCount byte
	Registers []uint16
}

// ParseReadResponse 解析读取响应
func ParseReadResponse(frame *ModbusFrame) (*ReadResponseData, error) {
	if frame.Function != FuncReadHoldingRegisters && frame.Function != FuncReadInputRegisters {
		return nil, ErrInvalidFunction
	}

	if len(frame.Data) < 1 {
		return nil, ErrNoData
	}

	byteCount := frame.Data[0]
	if int(byteCount) != len(frame.Data)-1 {
		return nil, ErrInvalidData
	}

	// 每个寄存器2字节
	regCount := byteCount / 2
	registers := make([]uint16, regCount)

	for i := 0; i < int(regCount); i++ {
		// 大端序
		registers[i] = binary.BigEndian.Uint16(frame.Data[1+i*2 : 3+i*2])
	}

	return &ReadResponseData{
		ByteCount: byteCount,
		Registers: registers,
	}, nil
}

// ExtractReading 从帧中提取电能读数
// 电能值通常存储为32位整数或IEEE 754浮点数
func ExtractReading(frame *ModbusFrame) (decimal.Decimal, error) {
	resp, err := ParseReadResponse(frame)
	if err != nil {
		return decimal.Zero, err
	}

	if len(resp.Registers) < 2 {
		return decimal.Zero, ErrInvalidData
	}

	// 尝试解析为32位整数 (常见电表格式)
	// 高位在前，低位在后
	value := uint32(resp.Registers[0])<<16 | uint32(resp.Registers[1])

	// 根据电表类型，可能需要除以精度因子（如100或1000）
	// 这里假设电表返回的是 kWh * 100 的整数值
	reading := decimal.NewFromInt(int64(value)).Div(decimal.NewFromInt(100))

	return reading, nil
}

// ExtractReadingFloat32 从帧中提取IEEE 754浮点格式的电能读数
func ExtractReadingFloat32(frame *ModbusFrame) (decimal.Decimal, error) {
	resp, err := ParseReadResponse(frame)
	if err != nil {
		return decimal.Zero, err
	}

	if len(resp.Registers) < 2 {
		return decimal.Zero, ErrInvalidData
	}

	// IEEE 754 单精度浮点
	bits := uint32(resp.Registers[0])<<16 | uint32(resp.Registers[1])
	value := math.Float32frombits(bits)

	// 检查是否为有效值
	if math.IsNaN(float64(value)) || math.IsInf(float64(value), 0) {
		return decimal.Zero, ErrInvalidData
	}

	return decimal.NewFromFloat(float64(value)), nil
}

// ExtractAddress 从帧中提取从站地址
func ExtractAddress(frame *ModbusFrame) byte {
	return frame.Address
}

// ExtractAddressFromData 从原始数据中提取从站地址（不验证CRC）
func ExtractAddressFromData(data []byte) (byte, error) {
	if len(data) < 1 {
		return 0, ErrFrameTooShort
	}
	return data[0], nil
}

// BuildReadRequest 构建读取寄存器请求帧
func BuildReadRequest(address byte, function byte, startReg uint16, regCount uint16) []byte {
	data := make([]byte, 6)
	data[0] = address
	data[1] = function
	binary.BigEndian.PutUint16(data[2:4], startReg)
	binary.BigEndian.PutUint16(data[4:6], regCount)
	return AppendCRC(data)
}
