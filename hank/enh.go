package hank

import (
	"strconv"
	"time"

	"github.com/twiglab/h2o/pkg/data"
)

type Enh struct {
}

func (e *Enh) Convert(dd DeviceData) data.Device {
	d := data.Device{
		Code: dd.No,
		Type: dd.Type,
		Name: dd.No,

		Time: parseTime(dd.DataTime),
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

func parseTime(s string) time.Time {
	t, err := time.ParseInLocation(time.DateTime, s, time.Local)
	if err != nil {
		return time.Date(2000, 0, 0, 0, 0, 0, 0, time.Local)
	}
	return t
}
