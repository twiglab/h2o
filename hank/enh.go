package hank

import (
	"strconv"

	"github.com/twiglab/h2o/pkg/data"
)

type Enh struct {
}

func (e *Enh) Convert(dd DeviceData) data.Device {
	d := data.Device{
		Code: dd.No,
		Type: dd.Type,
		Name: dd.No,

		Time: dd.DataTime,
		UUID: dd.DataCode,
	}

	if dd.Type == ELECTRICITY {
		d.Data.Electricity.DataValue = str2I64(dd.DataJson.DataValue, 100)
	}

	return d
}

func str2I64(s string, i float64) int64 {
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return int64(f * i)
	}
	return -1
}
