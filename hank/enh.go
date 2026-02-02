package hank

import (
	"strconv"

	"github.com/twiglab/h2o/pkg/kwh"
)

type Enh struct {
}

func (e *Enh) Convert(dd DeviceData) kwh.Device {
	return kwh.Device{
		Code: dd.No,
		Type: dd.Type,
		Name: dd.No,

		Time: dd.DataTime,
		UUID: dd.DataCode,

		Data: kwh.Data{
			Voltage: kwh.Voltage{
				VolA: str2I64(dd.DataJson.VoltageA, 100),
			},
		},
	}
}

func str2I64(s string, i float64) int64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return -1
	}

	return int64(f * i)
}
