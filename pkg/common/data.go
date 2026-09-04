package common

import (
	"time"
	"uuid"
)

type Device struct {
	SN   string `json:"sn,omitempty"`   // 仪表的序列号,仪表上有个条形码,如果没有就是空,或者自定义
	Code string `json:"code"`           // 设备code,业务全局唯一
	Type string `json:"type"`           // 设备类型
	Name string `json:"name,omitempty"` // 设备名称,可以为空

	DataTime time.Time `json:"data_time"` // 采集时间
	DataTs   string    `json:"data_ts"`   // 采集时间，字符串，到秒, 20060102150405
	DataCode string    `json:"data_code"` // 采集的唯一标识,全局唯一单调递增

	Status int `json:"status"` // 设备状态, 网关,采集程序或设备自定义, 0表示正常
}

// 点位信息
type Pos struct {
	Project string `json:"project,omitempty"`  // 所属项目编号
	PosCode string `json:"pos_code,omitempty"` // 位置编号
	Owner   string `json:"owner,omitempty"`    // 归属方
}

type Electricity struct {
	MeterValue

	VoltageA int64 `json:"voltage_a,omitempty"`
	VoltageB int64 `json:"voltage_b,omitempty"`
	VoltageC int64 `json:"voltage_c,omitempty"`

	CurrentA int64 `json:"current_a,omitempty"`
	CurrentB int64 `json:"current_b,omitempty"`
	CurrentC int64 `json:"current_c,omitempty"`

	ActivePowerTotal int64 `json:"active_power_total,omitempty"` // 总有功功率  P
	// ReactivePowerTotal int64 `json:"reactive_power_total,omitempty"` // 总无功功率  Q
	// ApparentPowerTotal int64 `json:"apparent_power_total,omitempty"` // 总视在功率  S

	Frequency int64 `json:"frequency,omitempty"` // 频率

}

type Water struct {
	MeterValue
}

type MeterValue struct {
	DataValue int64 `json:"data_value,omitempty"` // 表显读数
	OptStatus int64 `json:"opt_status,omitempty"` // 开合状态
}

func NewDataCode() string {
	id := uuid.NewV7()
	return id.String()
}

const f = "20060102150405"

func Ts(now time.Time) string {
	return now.Format(f)
}
