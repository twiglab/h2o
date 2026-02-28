package chrgg

import (
	"time"
)

var nilCDR CDR

var RulNew = zr{t: "new", c: "new"}

type CDR struct {
	DeviceCode string
	DeviceType string

	LastDataCode string
	DataCode     string

	LastDataTime time.Time
	DataTime     time.Time

	LastDataValue int64 // 上次表显
	DataValue     int64 // 当前表显

	Value int64 // 计量值,两次表显的差值,用于计算费用的数值

	RuleID     string
	RuleCtg    string
	RuleType   string
	UnitFeeFen int64

	FeeFen int64

	PosCode string
	Project string

	Memo string
}

func CalcCDR(last LastCDR, cd ChargeData, cr ChargeRuler) CDR {
	value, fee := calc(last.DataValue, cd.Data.DataValue, cr.UnitFeeFen())
	return CDR{
		DeviceCode: cd.Code,
		DeviceType: cd.Type,

		LastDataTime: last.DataTime,
		DataTime:     cd.DataTime,

		LastDataValue: last.DataValue,
		DataValue:     cd.Data.DataValue,

		LastDataCode: last.DataCode,
		DataCode:     cd.DataCode,

		RuleID:   cr.ID(),
		RuleType: cr.Type(),
		RuleCtg:  cr.Category(),

		Value:      value,
		UnitFeeFen: cr.UnitFeeFen(),
		FeeFen:     fee,

		PosCode: cd.Pos.PosCode,
		Project: cd.Pos.Project,

		Memo: cr.Memo(),
	}
}

func calc(lastV, currV int64, unit int64) (v int64, f int64) {
	v = currV - lastV
	f = v * unit / 100
	return
}
