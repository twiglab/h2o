package hank

import (
	"context"
	"strconv"
	"time"

	"github.com/twiglab/h2o/pkg/common"
)

type Enh struct {
	DDB *DuckDB
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
				PUID:      puid(meta.Project, meta.PosCode),
			},

			Flag: common.Flag{
				F1: meta.F1,
				F2: meta.F2,
				F3: meta.F3,
				F4: meta.F4,
				F5: meta.F5,
			},
		},
		Data: common.Water{
			MeterValue: common.MeterValue{
				DataValue: str2I64(dd.DataJson.DataValue, 100),
			},
		},
	}

	return d
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
				PUID:      puid(meta.Project, meta.PosCode),
			},

			Flag: common.Flag{
				F1: meta.F1,
				F2: meta.F2,
				F3: meta.F3,
				F4: meta.F4,
				F5: meta.F5,
			},
		},

		Data: common.Electricity{
			MeterValue: common.MeterValue{
				DataValue: str2I64(dd.DataJson.DataValue, 100),
			},
		},
	}

	return d
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

func puid(project, posCode string) string {
	return posCode + "@" + project
}

func now() time.Time {
	return time.Now()
}
