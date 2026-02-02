package data

import (
	"time"
)

type Device struct {
	SN   string `json:"sn"`   // 电表的序列号,电表上有个条形码,如果没有就是空,或者自定义
	Code string `json:"code"` // 设备code,业务全局唯一
	Type string `json:"type"` // 设备类型
	Name string `json:"name"` // 设备名称,可以为空

	Time time.Time `json:"time,format:RFC3339Nano"` // 采集时间
	UUID string    `json:"uuid"`                    // 采集的唯一标识,全局唯一单调递增

	Status string `json:"status"` // 设备状态, 网关,采集程序或设备自定义

	Data DataMix `json:"data"` // 设备数据

	Pos Pos `json:"pos,omitzero"` // 设备所在的位置信息
}

type Pos struct {
	Project   string `json:"project"`  // 所属项目编号
	PosCode   string `json:"pos_code"` // 位置编号
	Building  string `json:"building"` // 大楼
	FloorCode string `json:"floor_code"`
	AreaCode  string `json:"area_code"`

	PUID string `json:"puid"` // 全局唯一的poscode, 理论上 = project + poscode
}

type Electricity struct {
	VoltageA int64 `json:"voltage_a"`
	VoltageB int64 `json:"voltage_b"`
	VoltageC int64 `json:"voltage_c"`

	CurrentA int64 `json:"current_a"`
	CurrentB int64 `json:"current_b"`
	CurrentC int64 `json:"current_c"`

	Frequency int64 `json:"frequency"` // 频率

	TotalActivePower   int64 `json:"total_active_power"`   //总有功功率  P
	TotalReactivePower int64 `json:"total_reactive_power"` //总无功功率  Q
	TotalApperentPower int64 `json:"total_apperent_power"` //总视在功率  S
	TotalPowerFactor   int64 `json:"total_power_factor"`   // 功率因数 PF = p/s
}

type Water struct {
}

type DataMix struct {
	Electricity `json:",inline"`
	Water       `json:",inline"`
}
