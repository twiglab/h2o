package chrgg

import (
	"context"
	"encoding"
	"encoding/json/v2"
	"time"

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
	Gateway common.Modbus     `jaon:"gateway,omitzero"`
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

type OnOffMsg struct {
	common.Device
	Gateway common.Modbus `json:"gateway,omitzero"`
	OP      string        `json:"op"`
	Charge  Charge        `json:"charge,omitzero"`
}

func (o OnOffMsg) Topic() string {
	return common.H2O + "/onoff/" + o.Gateway.Code + "/" + o.Code + "/" + o.OP
	// t := fmt.Sprintf("h2o/onoff/%s/%s/%s", o.Gateway.Code, o.Code, o.OP)
	// return t
}

func (d OnOffMsg) MarshalBinary() ([]byte, error) {
	return json.Marshal(d)
}
