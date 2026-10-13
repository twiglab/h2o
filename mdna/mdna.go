package mdna

import (
	"context"
	"encoding"
	"encoding/json/v2"

	"github.com/twiglab/h2o/proto"
)

const CLIENT_ID = "mdna"

type OnOffMessage struct {
	proto.Device
	Gateway proto.Modbus `json:"gateway,omitzero"`
	Op      string       `json:"op"`
	Top     proto.Top    `json:"top,omitzero"`
}

func (d *OnOffMessage) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, d)
}

type ChrggMessage struct {
	proto.Device
	Pos     proto.Pos        `json:"pos,omitzero"`
	Data    proto.MeterValue `json:"data,omitzero"`
	Gateway proto.Modbus     `json:"gateway,omitzero"`
	Top     proto.Top        `json:"top,omitzero"`
	encoding.BinaryUnmarshaler
}

func (d *ChrggMessage) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, d)
}

type Notice struct {
	Device *ent.Device
	Type   string
	Charge ChrggMessage
	Op     string
}

type Speaker interface {
	Notice(ctx context.Context, n Notice) error
}
