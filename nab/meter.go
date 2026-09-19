package nab

import (
	"encoding/json/v2"

	"github.com/twiglab/h2o/pkg/common"
)

type Meter struct {
	common.Device
	Pos     common.Pos        `json:"pos,omitzero"`
	Data    common.MeterValue `json:"data,omitzero"`
	Gateway common.Modbus     `json:"gateway,omitzero"`
}

func (m Meter) MarshalBinary() (data []byte, err error) {
	return json.Marshal(m)
}

func (m Meter) Topic() string {
	return common.DataTopic(m.Device)
}
