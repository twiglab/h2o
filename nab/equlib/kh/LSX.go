package kh

import (
	"context"

	"github.com/twiglab/h2o/nab"
	"github.com/twiglab/h2o/proto"
)

const KH_LXS_40_MI = "KH-LXS-40-MI"

// LXS 水表，0x00地址读取用量（即表显）
type LXSNoValve struct {
}

func (e LXSNoValve) Collect(ctx context.Context, cli nab.ModbusOperator, data nab.DeviceData) (proto.MeterValue, error) {
	mv := proto.NewMeterValue()
	mv.OptStatus = proto.OPT_STATUS_UNKNOW

	v, err := cli.FC04Uint32(0x00)
	if err != nil {
		return mv, err
	}
	mv.DataValue = int64(v)

	return mv, nil
}
