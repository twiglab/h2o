package vigil

import (
	"encoding/json/v2"

	"github.com/twiglab/h2o/proto"
)

type Meter struct {
	proto.Device
	Pos     proto.Pos        `json:"pos,omitzero"`
	Gateway proto.Gateway    `json:"gateway,omitzero"`
	Data    proto.MeterValue `json:"data,omitzero"`
}

func (d *Meter) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, d)
}
