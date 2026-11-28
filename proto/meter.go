package proto

const (
	// 开合状态定义, unknow 必须是保证为0
	OPT_STATUS_OFF    = 1
	OPT_STATUS_ON     = 2
	OPT_STATUS_UNKNOW = 0
)

const (
	KEY_DATA_VALUE = "data_value" // 表显读数
	KEY_OPT_STATUS = "opt_status" // 开合状态
)

// MeterValue 采集结果
// Extra中的key不能为data_value, opt_status，否则报错
// 这是已知问题，也是明确的设计要求
type MeterValue struct {
	DataValue int64            `json:"data_value"`
	OptStatus int64            `json:"opt_status"`
	Extra     map[string]int64 `json:",embed"`
}

func (m *MeterValue) SetValue(key string, value int64) {
	switch key {
	case "":
		return
	case KEY_DATA_VALUE, KEY_OPT_STATUS:
		m.Extra["x_"+key] = value
	default:
		m.Extra[key] = value
	}
}

func (m MeterValue) GetValue(key string) (v int64, ok bool) {
	v, ok = m.Extra[key]
	return
}

func NewMeterValue() MeterValue {
	return MeterValue{
		Extra: make(map[string]int64),
	}
}
