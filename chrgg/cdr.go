package chrgg

import "time"

type CDR struct {
	DeviceCode string
	Type       string

	DataCode     string
	LastDataCode string

	DataTime     time.Time
	LastDataTime time.Time

	DataValue     int64 // 当前表显
	LastDataValue int64 // 上次表显

	PloyID string
	RuleID string

	Value int64 // 两次表显的差值, 用于计算费用的数值

	Amount int64
	Unit   int64
	Total  int64
}
