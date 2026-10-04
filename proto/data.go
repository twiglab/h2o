package proto

import (
	"time"
)

const (
	ON  = "ON"
	OFF = "OFF"
)

const (
	// 开合状态定义, unknow 必须是保证为0
	OPT_STATUS_OFF    = 1
	OPT_STATUS_ON     = 2
	OPT_STATUS_UNKNOW = 0
)

type Device struct {
	Code string `json:"code"` // 设备code,业务全局唯一
	Type string `json:"type"` // 设备类型

	SN string `json:"sn,omitempty"` // 仪表的序列号,仪表上有个条形码,如果没有就是空,或者自定义

	DataTime time.Time `json:"data_time"` // 采集时间
	DataCode string    `json:"data_code"` // 采集的唯一标识,全局唯一单调递增

	Status int `json:"status"` // 设备业务状态, 网关,采集程序或设备自定义, 0表示正常
}

type MeterValue struct {
	DataValue int64          `json:"data_value,omitempty"` // 表显读数
	OptStatus int64          `json:"opt_status,omitempty"` // 开合状态
	Extra     map[string]any `json:",embed"`
}

// 点位信息
type Pos struct {
	Project string         `json:"project,omitempty"`  // 所属项目编号
	PosCode string         `json:"pos_code,omitempty"` // 位置编号
	Extra   map[string]any `json:",embed"`
}

const (
	GATEWAY_NH = "NH"
)

type Gateway struct {
	Code string `json:"code"`           // 网关code,业务全局唯一
	Type string `json:"type,omitempty"` // 网关类型
	Pos  string `json:"pos,omitempty"`  // 网关所在位置
}

type Modbus struct {
	Gateway
	UnitID uint8 `json:"unit_id"`
}
