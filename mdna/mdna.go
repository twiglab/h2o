package mdna

import (
	"context"
	"encoding"
	"encoding/json/v2"
	"time"

	"github.com/twiglab/h2o/dbcli/ent"
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

type StatusChangeNoticeBody struct {
	Now    time.Time
	Code   string
	Op     string
	Device *ent.Device
}

type Speaker interface {
	StatusChangeNotice(ctx context.Context, n StatusChangeNoticeBody) error
}
