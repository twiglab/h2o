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

const (
	data_value                     = "data-value"                     //	总电量
	voltage_a                      = "voltage-a"                      //	A相电压
	voltage_b                      = "voltage-b"                      //	B相电压
	voltage_c                      = "voltage-c"                      //	C相电压
	current_a                      = "current-a"                      //	A相电流
	current_b                      = "current-b"                      //	B相电流
	current_c                      = "current-c"                      //	C相电流
	combined_active_total          = "combined-active-total"          //	(当前) 组合有功总电能
	forward_active_total           = "forward-active-total"           //	(当前) 正向有功总电能
	forward_active_rate1           = "forward-active-rate1"           //	(当前) 正向有功费率1电能
	forward_active_rate2           = "forward-active-rate2"           //	(当前) 正向有功费率2电能
	forward_active_rate3           = "forward-active-rate3"           //	(当前) 正向有功费率3电能
	forward_active_rate4           = "forward-active-rate4"           //	(当前) 正向有功费率4电能
	reverse_active_total           = "reverse-active-total"           //	(当前) 反向有功总电能
	combined_reactive1_total       = "combined-reactive1-total"       //	(当前) 组合无功1总电能
	forward_apparent_total         = "forward-apparent-total"         //	(当前) 正向视在总电能
	total_active_power             = "total-active-power"             //	瞬时总有功功率
	active_power_a                 = "active-power-a"                 //	瞬时A相有功功率
	active_power_b                 = "active-power-b"                 //	瞬时B相有功功率
	active_power_c                 = "active-power-c"                 //	瞬时C相有功功率
	total_reactive_power           = "total-reactive-power"           //	瞬时总无功功率
	reactive_power_a               = "reactive-power-a"               //	瞬时A相无功功率
	reactive_power_b               = "reactive-power-b"               //	瞬时B相无功功率
	reactive_power_c               = "reactive-power-c"               //	瞬时C相无功功率
	total_apparent_power           = "total-apparent-power"           //	瞬时总视在功率
	apparent_power_a               = "apparent-power-a"               //	瞬时A相视在功率
	apparent_power_b               = "apparent-power-b"               //	瞬时B相视在功率
	apparent_power_c               = "apparent-power-c"               //	瞬时C相视在功率
	power_factor                   = "power-factor"                   //	总功率因数
	frequency                      = "frequency"                      //	电网频率
	display_time                   = "display-time"                   //	每屏显示时间
	communication_address          = "communication-address"          //	通信地址
	meter_number                   = "meter-number"                   //	表号
	run_status_1                   = "run-status-1"                   //	电表运行状态字 1
	run_status_2                   = "run-status-2"                   //	电表运行状态字 2
	run_status_3                   = "run-status-3"                   //	电表运行状态字 3（操作类
	run_status_4                   = "run-status-4"                   //	电表运行状态字 4（A 相故障状态）
	run_status_5                   = "run-status-5"                   //	电表运行状态字 5（B 相故障状态）
	run_status_6                   = "run-status-6"                   //	电表运行状态字 6（C 相故障状态
	run_status_7                   = "run-status-7"                   //	电表运行状态字 7（合相故障状态）
	last_freeze_time               = "last-freeze-time"               //	（上 1 次）定时冻结时间
	last_active_energy_block       = "last-active-energy-block"       //	（上 1 次）正向有功电能数据块
	last_31_freeze_time            = "last-31-freeze-time"            //	（上 31 次）定时冻结时间
	opt_status                     = "opt-status"                     //	开合状态
	pulse_constant                 = "pulse-constant"                 //	脉冲常数
	voltage_ab                     = "voltage-ab"                     //	A-B 线电压
	voltage_cb                     = "voltage-cb"                     //	C-B 线电压
	voltage_ca                     = "voltage-ca"                     //	C-A 线电压
	total_active_energy            = "total-active-energy"            //	当前组合有功总电能
	peak_active_energy             = "peak-active-energy"             //	当前组合有功尖电能
	peak_active_energy_h           = "peak-active-energy-h"           //	当前组合有功峰电能
	flat_active_energy             = "flat-active-energy"             //	当前组合有功平电能
	valley_active_energy           = "valley-active-energy"           //	当前组合有功谷电能
	total_active_energy_forward    = "total-active-energy-forward"    //	当前正向总有功电能
	peak_active_energy_forward     = "peak-active-energy-forward"     //	当前正向有功尖电能
	peak_active_energy_h_forward   = "peak-active-energy-h-forward"   //	当前正向有功峰电能
	flat_active_energy_forward     = "flat-active-energy-forward"     //	当前正向有功平电能
	valley_active_energy_forward   = "valley-active-energy-forward"   //	当前正向有功谷电能
	total_active_energy_reverse    = "total-active-energy-reverse"    //	当前反向总有功电能
	peak_active_energy_reverse     = "eak-active-energy-reverse"      //	当前反向有功尖电能
	peak_active_energy_h_reverse   = "peak-active-energy-h-reverse"   //	当前反向有功峰电能
	flat_active_energy_reverse     = "flat-active-energy-reverse"     //	当前反向有功平电能
	valley_active_energy_reverse   = "valley-active-energy-reverse"   //	当前反向有功谷电能
	total_reactive_energy          = "total-reactive-energy"          //	当前组合无功总电能
	peak_reactive_energy           = "peak-reactive-energy"           //	当前组合无功尖电能
	peak_reactive_energy_h         = "peak-reactive-energy-h"         //	当前组合无功峰电能
	flat_reactive_energy           = "flat-reactive-energy"           //	当前组合无功平电能
	valley_reactive_energy         = "valley-reactive-energy"         //	当前组合无功谷电能
	total_reactive_energy_forward  = "total-reactive-energy-forward"  //	当前正向总无功电能
	peak_reactive_energy_forward   = "peak-reactive-energy-forward"   //	当前正向无功尖电能
	peak_reactive_energy_h_forward = "peak-reactive-energy-h-forward" //	当前正向无功峰电能
	flat_reactive_energy_forward   = "flat-reactive-energy-forward"   //	当前正向无功平电能
	valley_reactive_energy_forward = "valley-reactive-energy-forward" //	当前正向无功谷电能
	total_reactive_energy_reverse  = "total-reactive-energy-reverse"  //	当前反向总无功电能
	peak_reactive_energy_reverse   = "peak-reactive-energy-reverse"   //	当前反向无功尖电能
	peak_reactive_energy_h_reverse = "peak-reactive-energy-h-reverse" //	当前反向无功峰电能
	flat_reactive_energy_reverse   = "flat-reactive-energy-reverse"   //	当前反向无功平电能
	valley_reactive_energy_reverse = "valley-reactive-energy-reverse" //	当前反向无功谷电能
	forward_active_max_demand      = "forward-active-max-demand"      //	正向有功最大需量
	reverse_active_max_demand      = "reverse-active-max-demand"      //	反向有功最大需量
	forward_reactive_max_demand    = "forward-reactive-max-demand"    //	正向无功最大需量
	reverse_reactive_max_demand    = "reverse-reactive-max-demand"    //	反向无功最大需量
	active_energy_a_forward        = "active-energy-a-forward"        //	A 相正向有功电能
	active_energy_b_forward        = "active-energy-b-forward"        //	B 相正向有功电能
	active_energy_c_forward        = "active-energy-c-forward"        //	C 相正向有功电能
	pt_ratio                       = "pt-ratio"                       //	电压变比 PT
	ct_ratio                       = "ct-ratio"                       //	电流变比 CT
	zero_sequence_current          = "zero-sequence-current"          //	零序电流
	voltage_unbalance              = "voltage-unbalance"              //	电压不平衡度
	current_unbalance              = "current-unbalance"              //	电流不平衡度
	power_factor_a                 = "power-factor-a"                 //	A 相功率因数
	power_factor_b                 = "power-factor-b"                 //	B 相功率因数
	power_factor_c                 = "power-factor-c"                 //	C 相功率因数
	active_power_total             = "active-power-total"             //	总有功功率
	reactive_power_total           = "reactive-power-total"           //	总无功功率
	apparent_power_total           = "apparent-power-total"           //	总视在功率
	active_energy_total            = "active-energy-total"            //	总有功电能
	active_energy_forward          = "active-energy-forward"          //	正向有功电能
	active_energy_reverse          = "active-energy-reverse"          //	反向有功电能
	reactive_energy_total          = "reactive-energy-total"          //	总无功电能
	reactive_energy_forward        = "reactive-energy-forward"        //	正向无功电能
	reactive_energy_reverse        = "reactive-energy-reverse"        //	反向无功电能
	remaining_amount               = "remaining-amount"               //	电费剩余金额
	total_amount                   = "total-amount"                   //	总购电费金额
	current_rate                   = "current-rate"                   //	当前费率
	current_price                  = "current-price"                  //	当前电价
	system_status                  = "system-status"                  //	系统状态
	current_n                      = "current-n"                      //	零线电流
)
