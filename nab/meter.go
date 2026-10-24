package nab

import (
	"encoding/json/v2"

	"github.com/twiglab/h2o/proto"
)

type Meter struct {
	proto.Device
	Pos     proto.Pos        `json:"pos,omitzero"`
	Data    proto.MeterValue `json:"data,omitzero"`
	Gateway proto.Modbus     `json:"gateway,omitzero"`
}

func (m Meter) MarshalBinary() (data []byte, err error) {
	return json.Marshal(m)
}

func (m Meter) Topic() string {
	return proto.MakeTopic(proto.DATA, m.Gateway.Code, m.Code, m.Type)

}
