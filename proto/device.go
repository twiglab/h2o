package proto

import (
	"time"
	"uuid"
)

type Device struct {
	Code string `json:"code"` // 设备code,业务全局唯一
	Type string `json:"type"` // 设备类型

	SN string `json:"sn,omitempty"` // 仪表的序列号,仪表上有个条形码,如果没有就是空,或者自定义

	DataTime time.Time `json:"data_time"` // 采集时间
	DataCode string    `json:"data_code"` // 采集的唯一标识,全局唯一单调递增

	Status int `json:"status"` // 设备业务状态, 网关,采集程序或设备自定义, 0表示正常
}

func NewDataCode() string {
	id := uuid.NewV7()
	return id.String()
}
