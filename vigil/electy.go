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
}

type ElectricityMeter struct {
	Meter
	Data  common.Electricity      `json:"data,omitzero"`
	Param common.ElectricityParam `json:"param,omitzero"`
}

func (d *ElectricityMeter) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, d)
}

func (d ElectricityMeter) PCode() string {
	return d.Pos.Project + "@" + d.Code
}

func (d ElectricityMeter) Ts() string {
	return d.DataTime.Format(f)
}

func (d ElectricityMeter) XDataValue() int64 {
	return d.Data.DataValue * int64(d.Param.Factor)
}

func (d ElectricityMeter) STD() (r float64) {
	fd := stats.LoadRawData([]int64{d.Data.CurrentA, d.Data.CurrentB, d.Data.CurrentC})
	r, _ = fd.StandardDeviationPopulation()
	return
}
