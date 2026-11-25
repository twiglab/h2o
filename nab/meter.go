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

func NewMeter() Meter {
	return Meter{
		// status 为业务状态，这个状态在当前v1版本中无法获取
		// 在v1 版本中，所有的设备都是上线状态
		// 留到v2版本中解决此问题
		Status: proto.STATUS_ENABLE,
		Data: proto.MeterValue{
			Extra: make(map[string]int64),
		},
	}
}
