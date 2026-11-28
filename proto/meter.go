package proto

const (
	// 开合状态定义, unknow 必须是保证为0
	OPT_STATUS_OFF    = 1
	OPT_STATUS_ON     = 2
	OPT_STATUS_UNKNOW = 0
)

// MeterValue 采集结果
// Extra中的key不能为data_value, opt_status，否则报错
// 这是明确的设计要求，用于兼容其他上有应用
// 新协议将在v2版本中提供
type MeterValue struct {
	DataValue int64            `json:"data_value"`
	OptStatus int64            `json:"opt_status"`
	Extra     map[string]int64 `json:",embed"`
}

func NewMeterValue() MeterValue {
	return MeterValue{
		Extra: make(map[string]int64),
	}
}
