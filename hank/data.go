package hank

import (
	"encoding/json/jsontext"
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

type DeviceInfo struct {
	BrandName    string `json:"brandName"`
	BuildingCode string `json:"buildingCode"`

	//CreateTime time.Time `json:"createTime,format:DateTime"`
	CreateTime string `json:"createTime"`

	Address string `json:"deviceAddress"`
	Name    string `json:"deviceName"`
	No      string `json:"deviceNo"`
	Type    string `json:"deviceType"`

	GatewayIP string `json:"gatewayIp"`
	GatewayNo string `json:"gatewayNo"`

	IsControl bool   `json:"isControl"`
	ModelName string `json:"modelName"`
	OptStatus string `json:"optStatus"`

	ProjectNo string `json:"projectNo"`

	Status string `json:"status"`

	UnknownFields map[string]any `json:",unknown"`
}

type DeviceList []DeviceInfo

type GatewayInfo struct {
	// 网关编号，示例值为"g100"
	GatewayNo string
	// 网关名称，示例值为"网关 12"
	GatewayName string
	// 建筑编码，示例值为"123"

	BuildingCode string `json:"buildingCode"`
	// 网关 IP 地址，示例值为"192.168.3.101"
	GatewayAddress string

	UnknownFields map[string]any `json:",unknown"`
}

type DataMix struct {
	DataValue string            `json:"data-value"`
	ExtraData map[string]string `json:",unknown"`
}

type DeviceData struct {
	No   string `json:"deviceNo"`
	Type string `json:"deviceType"`

	Money    float64 `json:"dataMoney"`
	DataTime string  `json:"dataTime"` // 数据记录时间
	Usage    float64 `json:"usage"`

	BuildingCode string `json:"buildingCode"`

	LastDataTime string `json:"lastDataTime"` // 上一次数据记录时间

	DataCode string `json:"dataCode"`

	DataJson DataMix `json:"dataJson,omitzero"`

	UnknownFields map[string]any `json:",unknown"`
}

type DeviceDataList []DeviceData

type SyncData struct {
	Type string         `json:"type"`
	Data jsontext.Value `json:"data"`

	UnknownFields map[string]any `json:",unknown"`
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
