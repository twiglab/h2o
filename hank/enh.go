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
	DDB cache.Cache[string, MetaData]
}

func (e *Enh) ToWater(dd DeviceData) (WaterMeter, error) {
	data, err := WaterData(dd.DataJson)
	if err != nil {
		return WaterMeter{}, err
	}

	meta, _, _ := e.DDB.Get(context.Background(), dd.No)

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

	meta, _, _ := e.DDB.Get(context.Background(), dd.No)

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

func WaterData(dm DataMix) (w common.Water, err error) {
	w.DataValue, err = str2I64E(dm.DataValue, 100)
	return
}

func ElectyData(dm DataMix) (ele common.Electricity, err error) {

	if ele.DataValue, err = str2I64E(dm.ExtraData[data_value_old], 1); err != nil {
		return
	}

	if v, ok := dm.ExtraData[voltage_a]; ok {
		if ele.VoltageA, err = str2I64E(v, 100); err != nil {
			return
		}
	}
	if v, ok := dm.ExtraData[voltage_b]; ok {
		if ele.VoltageB, err = str2I64E(v, 100); err != nil {
			return
		}
	}
	if v, ok := dm.ExtraData[voltage_c]; ok {
		if ele.VoltageC, err = str2I64E(v, 100); err != nil {
			return
		}
	}

	if v, ok := dm.ExtraData[current_a]; ok {
		if ele.CurrentA, err = str2I64E(v, 100); err != nil {
			return
		}
	}
	if v, ok := dm.ExtraData[current_b]; ok {
		if ele.CurrentB, err = str2I64E(v, 100); err != nil {
			return
		}
	}
	if v, ok := dm.ExtraData[current_c]; ok {
		if ele.CurrentC, err = str2I64E(v, 100); err != nil {
			return
		}
	}

	if v, ok := dm.ExtraData[active_power_total]; ok {
		if ele.ActivePowerTotal, err = str2I64E(v, 100); err != nil {
			return
		}
	}
	return
}

func str2I64E(s string, i float64) (v int64, err error) {
	var f float64
	if f, err = strconv.ParseFloat(s, 64); err != nil {
		return
	}
	v = int64(f * i)
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
