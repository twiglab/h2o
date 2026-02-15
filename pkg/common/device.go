package common

import "encoding/json"

type ElectricityMeter struct {
	Device

	Pos  Pos         `json:"pos,omitzero"`
	Data Electricity `json:"data"`

	Flag Flag `json:"flag,omitzero"`
}

func (m ElectricityMeter) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

type WaterMeter struct {
	Device

	Pos  Pos   `json:"pos,omitzero"`
	Data Water `json:"data"`

	Flag Flag `json:"flag,omitzero"`
}

func (m WaterMeter) MarshalBinary() (data []byte, err error) {
	return json.Marshal(m)
}
