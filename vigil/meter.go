package vigil

import (
	"encoding/json/v2"

	"github.com/montanaflynn/stats"
	"github.com/twiglab/h2o/pkg/common"
)

const f = "20060102150405"

type Meter struct {
	common.Device
	Pos common.Pos `json:"pos,omitzero"`

	Ts string
}

/*
func (d Meter) Ts() string {
	return d.DataTime.Format(f)
}
*/

type ElectricityMeter struct {
	Meter
	Data  common.Electricity      `json:"data,omitzero"`
	Param common.ElectricityParam `json:"param,omitzero"`

	XDateValue int64
	STD        float64
}

func (d *ElectricityMeter) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, d)
}

func (d *ElectricityMeter) setup() {
	d.Ts = d.DataTime.Format(f)
	d.XDateValue = d.Data.DataValue * int64(d.Param.Factor)

	fd := stats.LoadRawData([]int64{d.Data.CurrentA, d.Data.CurrentB, d.Data.CurrentC})
	d.STD, _ = fd.StandardDeviationPopulation()
}

type WaterMeter struct {
	Meter
	Data common.Water `json:"data,omitzero"`
}

func (d *WaterMeter) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, d)
}

func (d *WaterMeter) setup() {
	d.Ts = d.DataTime.Format(f)
}
