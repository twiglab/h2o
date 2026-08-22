package box

import (
	"encoding/json/v2"

	"github.com/twiglab/h2o/pkg/common"
)

type Meter struct {
	common.Device
	Pos common.Pos `json:"pos,omitzero"`
}

func (m Meter) Topic() string {
	return common.Topic(m.Device)
}

type ElectricityMeter struct {
	Meter
	Data common.Electricity `json:"data"`
}

func (m ElectricityMeter) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

type WaterMeter struct {
	Meter
	Data common.Water `json:"data"`
}

func (m WaterMeter) MarshalBinary() (data []byte, err error) {
	return json.Marshal(m)
}
