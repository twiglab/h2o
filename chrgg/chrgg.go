package chrgg

import (
	"context"
	"encoding"
	"encoding/json/v2"
	"time"

	"github.com/twiglab/h2o/dbcli/ent"
	"github.com/twiglab/h2o/proto"
)

const (
	STATUS_END     = 1
	STATUS_BEGIN   = 0
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

type Top struct {
	Code      string `json:"code"`
	Top       int64  `json:"top"`
	Current   int64  `json:"current"`
	Stock     int64  `json:"stock"`
	Incr      int64  `json:"incr"`
	Amount    int64  `json:"amount"`
	UnitPrice int64  `json:"unit_price"`

	ChargeTime time.Time `json:"charge_time"`

	Alarm     int        `json:"alarm"`
	AlarmTime *time.Time `json:"alarm_time,omitzero"`

	Status  int        `json:"status"`
	EndTime *time.Time `json:"end_time,omitzero"`
}

type OnOffMessage struct {
	proto.Device
	Data    proto.MeterValue `json:"data,omitzero"`
	Gateway proto.Modbus     `json:"gateway,omitzero"`
	Top     *ent.Top         `json:"top,omitzero"`

	Op string `json:"op"`
}

func (o OnOffMessage) Topic() string {
	return proto.H2O + "/" + proto.ONOFF + "/" + o.Gateway.Code + "/" + o.Code + "/" + o.Op
}

func (d OnOffMessage) MarshalBinary() ([]byte, error) {
	return json.Marshal(d)
}

type ChrggMessage struct {
	proto.Device
	Pos     proto.Pos        `json:"pos,omitzero"`
	Data    proto.MeterValue `json:"data,omitzero"`
	Gateway proto.Modbus     `json:"gateway,omitzero"`
	Top     *ent.Top         `json:"top,omitzero"`
}

func (o ChrggMessage) Topic() string {
	return proto.H2O + "/" + proto.QUOTA + "/" + o.Code + "/" + o.Type
}

func (d ChrggMessage) MarshalBinary() ([]byte, error) {
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
