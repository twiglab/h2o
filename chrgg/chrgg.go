package chrgg

import (
	"context"
	"encoding"
	"encoding/json/v2"

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
	Gateway common.Gateway    `jaon:"gateway,omitzero"`
}

func (d *Meter) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, d)
}

type OnOffMsg struct {
	common.Device
	Gateway common.Gateway `jaon:"gateway,omitzero"`
	OP      string
}

func (o OnOffMsg) Topic() string {
	return common.H2O + "/box/onoff/" + o.Gateway.Code + "/" + o.Code + "/" + o.OP
}

func (d OnOffMsg) MarshalBinary() ([]byte, error) {
	return json.Marshal(d)
}
