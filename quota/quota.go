package quota

import (
	"context"
	"encoding"
	"encoding/json/v2"

	"github.com/twiglab/h2o/dbcli/ent"
	"github.com/twiglab/h2o/proto"
)

const (
	STATUS_END   = 9
	STATUS_BEGIN = 3
	// 1, 2 保留，不允许使用
	STATUS_ZERO    = 0
	STATUS_INVALID = -1
)

type SendObject interface {
	encoding.BinaryMarshaler
	Topic() string
}

type Sender interface {
	SendData(ctx context.Context, obj SendObject) error
}

type Meter struct {
	proto.Device
	Pos     proto.Pos        `json:"pos,omitzero"`
	Data    proto.MeterValue `json:"data,omitzero"`
	Gateway proto.Modbus     `json:"gateway,omitzero"`
}

func (d *Meter) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, d)
}

type OnOffMessage struct {
	proto.Device
	Data    proto.MeterValue `json:"data,omitzero"`
	Gateway proto.Modbus     `json:"gateway,omitzero"`
	Top     *ent.Top         `json:"top,omitzero"`

	Op string `json:"op"`
}

func (o OnOffMessage) Topic() string {
	return proto.MakeTopic(proto.ONOFF, o.Gateway.Code, o.Code, o.Op)
}

func (d OnOffMessage) MarshalBinary() ([]byte, error) {
	return json.Marshal(d)
}

func newOnOffMessage(md Meter, vc *ent.Top, op string) OnOffMessage {
	return OnOffMessage{
		Device:  md.Device,
		Data:    md.Data,
		Gateway: md.Gateway,

		Top: vc,

		Op: op,
	}
}
