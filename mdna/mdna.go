package mdna

import (
	"context"
	"encoding"
	"encoding/json/v2"
	"time"

	"github.com/twiglab/h2o/dbx/ent"
	"github.com/twiglab/h2o/pkg/common"
)

const CLIENT_ID = "mdna"

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
	Op      string        `json:"op"`
	Top     Top           `json:"top,omitzero"`
}

func (d *OnOffMessage) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, d)
}

type ChrggMessage struct {
	common.Device
	Pos     common.Pos        `json:"pos,omitzero"`
	Data    common.MeterValue `json:"data,omitzero"`
	Gateway common.Modbus     `json:"gateway,omitzero"`
	Top     Top               `json:"top,omitzero"`
	encoding.BinaryUnmarshaler
}

func (d *ChrggMessage) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, d)
}

type Notice struct {
	Device       *ent.Device
	Type         string
	Charge ChrggMessage
	Op           string
}

type Speaker interface {
	Notice(ctx context.Context, n Notice) error
}
