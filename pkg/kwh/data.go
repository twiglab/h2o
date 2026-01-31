package kwh

import (
	"time"
)

type Device struct {
	Code string `json:"code"` // 设备code，业务全局唯一
	Type string `json:"type"` // 设备类型
	Name string `json:"name"` // 设备名称, 空

	Time time.Time `json:"time,format:RFC3339Nano"` // 采集时间
	UUID string    `json:"uuid"`                    // 采集的唯一标识, 全局唯一单调递增

	Status string `json:"status"` // 设备状态, 网关,采集程序或设备自定义

	Data Data `json:"data"`

	Pos Pos `json:"pos,omitzero"`
}

type Pos struct {
	Project   string `json:"project"`  // 所属项目编号
	PosCode   string `json:"pos_code"` // 位置编号
	Building  string `json:"building"` // 大楼
	FloorCode string `json:"floor_code"`
	AreaCode  string `json:"area_code"`

	PUID string `json:"puid"` // 全局唯一的poscode, 理论上 = project + poscode
}

type Voltage struct {
	// 三相电压， 单位 1/1000 V
	VolA uint64 `json:"vol_a,omitzero"`
	VolB uint64 `json:"vol_b,omitzero"`
	VolC uint64 `json:"vol_c,omitzero"`
}

type Amp struct {
	// 三相电流， 单位 1/1000 A
	AmpA uint64 `json:"apm_a,omitzero"`
	AmpB uint64 `json:"apm_b,omitzero"`
	AmpC uint64 `json:"apm_c,omitzero"`
}

type Power struct {
	// 总有功 1wh (kwh 的 1/1000)
	PowerP uint64 `json:"power_p,omitzero"`
	// Q uint64
	// S uint64
}

type Data struct {
	Voltage `json:",inline"`
	Amp     `json:",inline"`
	Power   `json:",inline"`
}
