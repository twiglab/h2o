package proto

const (
	// 开合状态定义, unknow 必须是保证为0
	OPT_STATUS_OFF    = 1
	OPT_STATUS_ON     = 2
	OPT_STATUS_UNKNOW = 0
)

type MeterValue struct {
	DataValue int64          `json:"data_value,omitempty"` // 表显读数
	OptStatus int64          `json:"opt_status,omitempty"` // 开合状态
	Extra     map[string]any `json:",embed"`
}
