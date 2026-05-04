package hank

import (
	"cmp"
	"context"
	"strconv"
	"time"

	"github.com/twiglab/h2o/cache"
	"github.com/twiglab/h2o/pkg/common"
)

type Enh struct {
	Cache cache.Cache[string, MetaData]
}

func (e *Enh) ToWater(dd DeviceData) (WaterMeter, error) {
	data, err := WaterData(dd.DataJson)
	if err != nil {
		return WaterMeter{}, err
	}

	meta, _, _ := e.Cache.Get(context.Background(), dd.No)

	return WaterMeter{
		Meter: Meter{
			Device: common.Device{
				Code: dd.No,
				Type: common.WATER,
				Name: cmp.Or(meta.Name, dd.No),

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
		Data: data,
	}, nil
}

func (e *Enh) ToElecty(dd DeviceData) (ElectricityMeter, error) {
	data, err := ElectyData(dd.DataJson)
	if err != nil {
		return ElectricityMeter{}, err
	}

	meta, _, _ := e.Cache.Get(context.Background(), dd.No)

	return ElectricityMeter{
		Meter: Meter{
			Device: common.Device{
				SN:   meta.SN,
				Code: dd.No,
				Type: common.ELECTRICITY,
				Name: cmp.Or(meta.Name, dd.No),

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

		Data: data,
	}, nil
}

func (e *Enh) ToGas(dd DeviceData) (gm GasMeter, err error) {
	return
}

func WaterData(dm DataMix) (cd common.Water, err error) {
	cd.DataValue, err = str2I64E(dm.DataValue, 1)
	return
}

func ElectyData(dm DataMix) (cd common.Electricity, err error) {
	if cd.DataValue, err = str2I64E(dm.DataValue, 1); err != nil {
		return
	}

	if cd.VoltageA, err = str2I64E(dm.VoltageA, 10); err != nil {
		return
	}
	if cd.VoltageB, err = str2I64E(dm.VoltageB, 10); err != nil {
		return
	}
	if cd.VoltageC, err = str2I64E(dm.VoltageC, 10); err != nil {
		return
	}

	if cd.CurrentA, err = str2I64E(dm.CurrentA, 1); err != nil {
		return
	}
	if cd.CurrentB, err = str2I64E(dm.CurrentB, 1); err != nil {
		return
	}
	if cd.CurrentC, err = str2I64E(dm.CurrentC, 1); err != nil {
		return
	}

	if cd.ActivePowerTotal, err = str2I64E(dm.ActivePowerTotal, 1); err != nil {
		return
	}

	if cd.Frequency, err = str2I64E(dm.Frequency, 1); err != nil {
		return
	}

	return
}

func str2I64E(s string, i int64) (v int64, err error) {
	if v, err = strconv.ParseInt(s, 10, 64); err != nil {
		return
	}
	v = v * i
	return
}

var xdate = time.Date(2000, 0, 0, 0, 0, 0, 0, time.Local)

func parseTime(s string) time.Time {
	t, err := time.ParseInLocation(time.DateTime, s, time.Local)
	if err != nil {
		return xdate
	}
	return t
}
