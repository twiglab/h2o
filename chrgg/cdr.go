package chrgg

import "time"

type CDR struct {
	DataCode     string
	LastDataCode string

	DataTime     time.Time
	LastDataTime time.Time

	DataValue     int64
	LastDataValue int64

	PloyID string
	RuleID string

	Amount int64
	Unit   int64
	Total  int64
}
