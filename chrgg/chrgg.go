package chrgg

import (
	"context"
	"encoding"
	"encoding/json/v2"
	"time"

	"github.com/twiglab/h2o/chrgg/orm/ent"
	"github.com/twiglab/h2o/pkg/common"
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
	common.Device
	Pos     common.Pos        `json:"pos,omitzero"`
	Data    common.MeterValue `json:"data,omitzero"`
	Gateway common.Modbus     `json:"gateway,omitzero"`
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
	common.Device
	Gateway common.Modbus `json:"gateway,omitzero"`
	Top     Top           `json:"top,omitzero"`

	Op string `json:"op"`
}

func newOnOffMessage(md Meter, vc *ent.Top, op string) OnOffMessage {
	return OnOffMessage{
		Op:      op,
		Device:  md.Device,
		Gateway: md.Gateway,

		Top: Top{
			Code:       vc.Code,
			Top:        vc.Top,
			Current:    md.Data.DataValue,
			Stock:      vc.Stock,
			Incr:       vc.Incr,
			Amount:     vc.Amount,
			UnitPrice:  vc.UnitPrice,
			ChargeTime: vc.ChargeTime,

			Alarm:     vc.Alarm,
			AlarmTime: vc.AlarmTime,

			Status:  vc.Status,
			EndTime: vc.EndTime,
		},
	}
}

func (o OnOffMessage) Topic() string {
	return common.H2O + "/onoff/" + o.Gateway.Code + "/" + o.Code + "/" + o.Op
}

func (d OnOffMessage) MarshalBinary() ([]byte, error) {
	return json.Marshal(d)
}

type ChrggMessage struct {
	common.Device
	Pos     common.Pos        `json:"pos,omitzero"`
	Data    common.MeterValue `json:"data,omitzero"`
	Gateway common.Modbus     `json:"gateway,omitzero"`
	Top     Top               `json:"top,omitzero"`
}

func (o ChrggMessage) Topic() string {
	return common.H2O + "/chrgg/" + o.Code + "/" + o.Type
}

func (d ChrggMessage) MarshalBinary() ([]byte, error) {
	return json.Marshal(d)
}
