package chrgg

import (
	"encoding/json"

	"github.com/twiglab/h2o/pkg/common"
)

type ElectyMeterData struct {
	common.Device
	Pos   common.Pos              `json:"pos,omitzero"`
	Data  common.MeterValue       `json:"data"`
	Param common.ElectricityParam `json:"param"`
	Flag  common.Flag             `json:"flag,omitzero"`

	Topic string `json:"topic"`
}

func (d *ElectyMeterData) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, d)
}

type WaterMeterData struct {
	common.Device
	Pos  common.Pos        `json:"pos,omitzero"`
	Data common.MeterValue `json:"data"`
	Flag common.Flag       `json:"flag,omitzero"`

	Topic string `json:"topic"`
}

func (d *WaterMeterData) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, d)
}
