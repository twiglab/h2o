package kh

import (
	"context"
	"time"

	"github.com/twiglab/h2o/nab"
	"github.com/twiglab/h2o/nab/equlib"
	"github.com/twiglab/h2o/proto"
)

const KH_DTZY2851_CD7_1 = "KH-DTZY2851-CD7-1"

// DTZY 电表，0x60地址读取有功总电能（即表显）
// 0x22 地址读取开合状态，写入执行对应的操作（0x01拉闸, 0x02合闸）
// 文档要求 开合状态读取和写入均为0x22
// 注意，对接文档要求，写入必须用modbus 0x10 （WriteRegisters）写入
type DTZY struct {
}

func (e DTZY) Collect(ctx context.Context, cli nab.ModbusOperator, data nab.DeviceData) (proto.MeterValue, error) {
	mv := proto.NewMeterValue()

	v, err := cli.FC04Uint32(0x60)
	if err != nil {
		return mv, err
	}

	opt, err := cli.FC04(0x22)
	if err != nil {
		return mv, err
	}

	mv.SetValue(proto.KEY_DATA_VALUE, int64(v))
	mv.SetValue(proto.KEY_OPT_STATUS, int64(opt))

	return mv, nil
}

func (e DTZY) On(ctx context.Context, cli nab.ModbusOperator, data nab.DeviceData) error {
	return equlib.TryDo(ctx, 3, 300*time.Millisecond,
		func() error {
			return cli.FC16(0x22, []uint16{0x02})
		},
	)
}

func (e DTZY) Off(ctx context.Context, cli nab.ModbusOperator, data nab.DeviceData) error {
	return equlib.TryDo(ctx, 3, 300*time.Millisecond,
		func() error {
			return cli.FC16(0x22, []uint16{0x01})
		},
	)
}
