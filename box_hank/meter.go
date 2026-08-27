package box

import (
	"encoding/json/v2"

	"github.com/twiglab/h2o/pkg/common"
)

type Meter struct {
	common.Device
	Pos  common.Pos `json:"pos,omitzero"`
	Data common.MeterValue  `json:"data,omitzero"`
}

func (m Meter) MarshalBinary() (data []byte, err error) {
	return json.Marshal(m)
}

func (m Meter) Topic() string {
	return common.Topic(m.Device)
}
