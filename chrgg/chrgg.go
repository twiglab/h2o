package chrgg

import (
	"context"
	"encoding"
	"encoding/json/v2"
	"time"

	"github.com/twiglab/h2o/chrgg/orm/ent"
	"github.com/twiglab/h2o/pkg/common"
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

type Charge struct {
	Code       string    `json:"code"`
	Top        int64     `json:"top"`
	Current    int64     `json:"current"`
	Stock      int64     `json:"stock"`
	Incr       int64     `json:"incr"`
	Amount     int64     `json:"amount"`
	UnitPrice  int64     `json:"unit_price"`
	ChargeTime time.Time `json:"charge_time"`
}

type OnOffMessage struct {
	common.Device
	Gateway common.Modbus `json:"gateway,omitzero"`
	Op      string        `json:"op"`
	Charge  Charge        `json:"charge,omitzero"`
}

func NewOnOffMessage(md Meter, vc *ent.ValueCharge, op string) OnOffMessage {
	return OnOffMessage{
		Op:      op,
		Device:  md.Device,
		Gateway: md.Gateway,
		Charge: Charge{
			Code:       vc.Code,
			Top:        vc.Top,
			Current:    md.Data.DataValue,
			Stock:      vc.Stock,
			Incr:       vc.Incr,
			Amount:     vc.Amount,
			UnitPrice:  vc.UnitPrice,
			ChargeTime: vc.ChargeTime,
		},
	}
}

func (o OnOffMessage) Topic() string {
	return common.H2O + "/onoff/" + o.Gateway.Code + "/" + o.Code + "/" + o.Op
}

func (d OnOffMessage) MarshalBinary() ([]byte, error) {
	return json.Marshal(d)
}
