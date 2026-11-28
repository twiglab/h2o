package nab

import (
	"encoding/json/v2"
	"time"

	"github.com/twiglab/h2o/nab/orm/ent"
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

func MakeMeter(box *Box, dev *ent.Dev, mv proto.MeterValue) Meter {
	return Meter{
		Code:     dev.Code,
		Type:     dev.Typ,
		SN:       dev.Sn,
		DataTime: time.Now(),
		DataCode: proto.NewDataCode(),

		// status 为业务状态，这个状态在当前v1版本中无法获取
		// 留到v2版本中解决此问题
		Status: proto.STATUS_UNKNOW,
		Data:   mv,
		Pos: proto.Pos{
			Project: box.Project,
		},
		Gateway: proto.Modbus{
			Code:   box.Code,
			Type:   proto.GATEWAY_NH,
			UnitID: dev.UnitID,
		},
	}
}
