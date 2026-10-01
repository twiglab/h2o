package vigil

import (
	"encoding/json/v2"

	"github.com/twiglab/h2o/pkg/common"
	"github.com/twiglab/h2o/vigil/stats"
)

type Meter struct {
	common.Device
	Pos     common.Pos     `json:"pos,omitzero"`
	Gateway common.Gateway `json:"gateway,omitzero"`
}

type ElectricityMeter struct {
	Meter
	Data common.Electricity `json:"data,omitzero"`

	STD float64
}

func (d *ElectricityMeter) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, d)
}

func (d *ElectricityMeter) setup() {
	_, d.STD = stats.MeanAndStdDev([]float64{
		float64(d.Data.CurrentA),
		float64(d.Data.CurrentB),
		float64(d.Data.CurrentC),
	})
}

type WaterMeter struct {
	Meter
	Data common.Water `json:"data,omitzero"`
}

func (d *WaterMeter) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, d)
}

func (d *WaterMeter) setup() {
}
