package clfoc

import (
	"context"

	"github.com/twiglab/h2o/nab"
	"github.com/twiglab/h2o/proto"
)

const CLF_OC_PT = "CLF-OC-PT"

type CLFOC struct {
}

func (e CLFOC) Collect(ctx context.Context, cli nab.ModbusOperator, data nab.DeviceData) (proto.MeterValue, error) {
	mv := proto.NewMeterValue()
	mv.SetValue(proto.KEY_OPT_STATUS, proto.OPT_STATUS_UNKNOW)

	v, err := cli.FC04Uint32(data.Record.Addr)
	if err != nil {
		return mv, err
	}
	mv.DataValue = int64(v)
	return mv, nil
}
