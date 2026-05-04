package hank

import "github.com/twiglab/h2o/pkg/common"

type Meter struct {
	common.Device
	Pos  common.Pos  `json:"pos,omitzero"`
	Flag common.Flag `json:"flag,omitzero"`
}

func (m Meter) Topic() string {
	return common.Topic(m.Device)
}

type ElectricityMeter struct {
	Meter
	common.ElectricityCT
	Data common.Electricity `json:"data"`
}

func (m ElectricityMeter) MarshalBinary() ([]byte, error) {
	return marshal(m)
}

type WaterMeter struct {
	Meter
	Data common.Water `json:"data"`
}

func (m WaterMeter) MarshalBinary() (data []byte, err error) {
	return marshal(m)
}

type GasMeter struct {
	Meter
	Data common.Water `json:"data"`
}

func (m GasMeter) MarshalBinary() (data []byte, err error) {
	return marshal(m)
}
