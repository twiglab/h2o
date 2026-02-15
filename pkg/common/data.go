package common

import (
	"time"
)

type Device struct {
	SN   string `json:"sn,omitempty"`   // 仪表的序列号,仪表上有个条形码,如果没有就是空,或者自定义
	Code string `json:"code"`           // 设备code,业务全局唯一
	Type string `json:"type"`           // 设备类型
	Name string `json:"name,omitempty"` // 设备名称,可以为空

	Time     time.Time `json:"time"`      // 采集时间
	DataCode string    `json:"data_code"` // 采集的唯一标识,全局唯一单调递增

	Status int `json:"status"` // 设备状态, 网关,采集程序或设备自定义
}

func (d Device) Topic() string {
	return "h2o/" + d.Code + "/" + d.Type
}

type Pos struct {
	Project   string `json:"project,omitempty"`  // 所属项目编号
	PosCode   string `json:"pos_code,omitempty"` // 位置编号
	Building  string `json:"building,omitempty"` // 大楼
	FloorCode string `json:"floor_code,omitempty"`
	AreaCode  string `json:"area_code,omitempty"`

	PUID string `json:"puid,omitempty"` // 全局唯一的poscode, 理论上 = project + poscode
}

type Electricity struct {
	VoltageA int64 `json:"voltage_a,omitempty"`
	VoltageB int64 `json:"voltage_b,omitempty"`
	VoltageC int64 `json:"voltage_c,omitempty"`

	CurrentA int64 `json:"current_a,omitempty"`
	CurrentB int64 `json:"current_b,omitempty"`
	CurrentC int64 `json:"current_c,omitempty"`

	DataValue int64 `json:"data_value,omitempty"` // 表显读数

	/*
		Frequency int64 `json:"frequency"` // 频率
		TotalActivePower   int64 `json:"total_active_power"`   // 总有功功率  P
		TotalReactivePower int64 `json:"total_reactive_power"` // 总无功功率  Q
		TotalApperentPower int64 `json:"total_apperent_power"` // 总视在功率  S
		TotalPowerFactor   int64 `json:"total_power_factor"`   // 功率因数 PF = p/s
	*/
}

type Water struct {
	DataValue int64 `json:"data_value,omitempty"` // 表显读数
	OptStatus int64 `json:"opt_status,omitempty"` // 开合状态
}

type Flag struct {
	P1 string `json:"p1,omitempty"`
	P2 string `json:"p2,omitempty"`
	P3 string `json:"p3,omitempty"`
	P4 string `json:"p4,omitempty"`
	P5 string `json:"p5,omitempty"`
}
