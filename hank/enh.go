package hank

import (
	"context"
	"strconv"
	"time"

	"github.com/twiglab/h2o/pkg/common"
)

type Enh struct {
	DDB Cache[string, MetaData]
}

func (e *Enh) ToWater(dd DeviceData) WaterMeter {
	meta, _, _ := e.DDB.Get(context.Background(), dd.No)
	d := WaterMeter{
		Meter: Meter{
			Device: common.Device{
				Code: dd.No,
				Type: dd.Type,
				Name: dd.No,

				DataTime: parseTime(dd.DataTime),
				DataCode: dd.DataCode,

				Status: 0,
			},
			Pos: common.Pos{
				Project:   meta.Project,
				PosCode:   meta.PosCode,
				Building:  meta.Building,
				FloorCode: meta.FloorCode,
				AreaCode:  meta.AreaCode,
			},

			Flag: common.Flag{
				F1: meta.F1,
				F2: meta.F2,
				F3: meta.F3,
				F4: meta.F4,
				F5: meta.F5,
			},
		},
		Data: waterData(dd.DataJson),
	}
	return d
}

func waterData(dm DataMix) common.Water {
	return common.Water{
		MeterValue: common.MeterValue{
			DataValue: str2I64(dm.DataValue, 100),
		},
	}
}

func (e *Enh) ToElectricity(dd DeviceData) ElectricityMeter {
	meta, _, _ := e.DDB.Get(context.Background(), dd.No)

	d := ElectricityMeter{
		Meter: Meter{
			Device: common.Device{
				SN:   meta.SN,
				Code: dd.No,
				Type: dd.Type,
				Name: meta.Name,

				DataTime: parseTime(dd.DataTime),
				DataCode: dd.DataCode,

				Status: 0,
			},

			Pos: common.Pos{
				Project:   meta.Project,
				PosCode:   meta.PosCode,
				Building:  meta.Building,
				FloorCode: meta.FloorCode,
				AreaCode:  meta.AreaCode,
			},

			Flag: common.Flag{
				F1: meta.F1,
				F2: meta.F2,
				F3: meta.F3,
				F4: meta.F4,
				F5: meta.F5,
			},
		},

		Data: electricityData(dd.DataJson),
	}
	return d
}

func electricityData(dm DataMix) common.Electricity {
	var ele common.Electricity

	ele.DataValue = str2I64(dm.DataValue, 100)

	if v, ok := dm.ExtraData["voltage-a"]; ok {
		ele.VoltageA = str2I64(v, 10)
	}
	if v, ok := dm.ExtraData["voltage-b"]; ok {
		ele.VoltageB = str2I64(v, 10)
	}
	if v, ok := dm.ExtraData["voltage-c"]; ok {
		ele.VoltageC = str2I64(v, 10)
	}

	if v, ok := dm.ExtraData["current-a"]; ok {
		ele.CurrentA = str2I64(v, 1)
	}
	if v, ok := dm.ExtraData["current-b"]; ok {
		ele.CurrentB = str2I64(v, 1)
	}
	if v, ok := dm.ExtraData["current-c"]; ok {
		ele.CurrentC = str2I64(v, 1)
	}

	if v, ok := dm.ExtraData["total-active-power"]; ok {
		ele.TotalActivePower = str2I64(v, 100)
	}

	return ele
}

func str2I64(s string, i float64) int64 {
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return int64(f * i)
	}
	return -1
}

var xdate = time.Date(2000, 0, 0, 0, 0, 0, 0, time.Local)

func parseTime(s string) time.Time {
	t, err := time.ParseInLocation(time.DateTime, s, time.Local)
	if err != nil {
		return xdate
	}
	return t
}
