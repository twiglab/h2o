package chrgg

import "time"

type CDR struct {
	DeviceCode string
	DeviceType string

	LastDataCode string
	DataCode     string

	LastDataTime time.Time
	DataTime     time.Time

	LastDataValue int64 // 上次表显
	DataValue     int64 // 当前表显

	Value int64 // 两次表显的差值, 用于计算费用的数值

	PloyID  string
	RuleID  string
	UnitFee int64

	Fee    int64
	Remark string
}
