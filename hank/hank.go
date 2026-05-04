package hank

import (
	"encoding/json/jsontext"
)

const (
	CLIENT_ID = "hank-plugin"
)

const (
	ELECTRICITY = "electricity"
	WATER       = "water"
)

const (
	TypeDeviceList   = "deviceList"
	TypeGatewayInfo  = "gatewayInfo"
	TypeRate         = "rate"
	TypeDeviceData   = "deviceData"
	TypeDeviceStatus = "deviceStatus"
	TypeTime         = "time"
)

const (
	offline = "offline"
	online  = "online"
)

type DeviceStatus struct {
	No     string `json:"deviceNo"`
	Type   string `json:"deviceType"`
	Status string `json:"status"`
}

func Online(s string) int {
	if s == online {
		return 0
	}
	return -1
}

type DeviceStatusList []DeviceStatus

type DataMix struct {
	DataValue string `json:"data-valueold,omitempty"`

	VoltageA string `json:"voltage-aold,omitempty"`
	VoltageB string `json:"voltage-bold,omitempty"`
	VoltageC string `json:"voltage-cold,omitempty"`

	CurrentA string `json:"current-aold,omitempty"`
	CurrentB string `json:"current-bold,omitempty"`
	CurrentC string `json:"current-cold,omitempty"`

	ActivePowerTotal string `json:"active-power-totalold,omitempty"` // 总有功功率  P
	Frequency        string `json:"frequency,omitempty"`
}

type DeviceData struct {
	No   string `json:"deviceNo"`
	Type string `json:"deviceType"`

	DataTime     string `json:"dataTime"`     // 数据记录时间
	LastDataTime string `json:"lastDataTime"` // 上一次数据记录时间
	DataCode     string `json:"dataCode"`     // 记录的唯一标识

	DataJson DataMix `json:"dataJson,omitzero"`
}

type DeviceDataList []DeviceData

type SyncData struct {
	Type string         `json:"type"`
	Data jsontext.Value `json:"data"`
}

type ReturnMessage struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func Error(message string) *ReturnMessage {
	return &ReturnMessage{Message: message, Type: "error"}
}

func Success() *ReturnMessage {
	return &ReturnMessage{Type: "success"}
}

var OK = Success()
var ErrNoRate = Error("没有待同步的费率数据")
