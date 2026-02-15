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

func (e *Enh) ToWater(dd DeviceData) common.WaterMeter {
	meta, _, _ := e.DDB.Get(context.Background(), dd.No)
	d := common.WaterMeter{
		Device: common.Device{
			Code: dd.No,
			Type: dd.Type,
			Name: dd.No,

			Time:     parseTime(dd.DataTime),
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

		Data: common.Water{
			DataValue: str2I64(dd.DataJson.DataValue, 100),
		},
	}

	return d
}

func (e *Enh) ToElectricity(dd DeviceData) common.ElectricityMeter {
	meta, _, _ := e.DDB.Get(context.Background(), dd.No)

	d := common.ElectricityMeter{
		Device: common.Device{
			SN:   meta.SN,
			Code: dd.No,
			Type: dd.Type,
			Name: meta.Name,

			Time:     parseTime(dd.DataTime),
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
			P1: meta.P1,
			P2: meta.P2,
			P3: meta.P3,
			P4: meta.P4,
			P5: meta.P5,
		},

		Data: common.Electricity{
			DataValue: str2I64(dd.DataJson.DataValue, 100),
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
