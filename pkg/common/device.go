package common

import "encoding/json"

type ElectricityMeter struct {
	Device

	Pos  Pos         `json:"pos"`
	Data Electricity `json:"data"`

	Flag Flag `json:"flag"`
}

func (m ElectricityMeter) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

type WaterMeter struct {
	Device

	Pos  Pos   `json:"pos"`
	Data Water `json:"data"`

	Flag Flag `json:"flag"`
}

func (m WaterMeter) MarshalBinary() (data []byte, err error) {
	return json.Marshal(m)
}
