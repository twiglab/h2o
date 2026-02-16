package chrgg

import (
	"time"

	"github.com/twiglab/h2o/chrgg/orm/ent"
)

var Nil CDR

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

	PloyID  string
	RuleID  string
	UnitFee int64

	Fee int64

	PosCode string
	Project string

	Remark string
}

func FirstCDR(cd ChargeData) CDR {
	return CDR{
		DeviceCode: cd.Code,
		DeviceType: cd.Type,

		LastDataTime: cd.Time,
		DataTime:     cd.Time,

		LastDataValue: cd.Data.DataValue,
		DataValue:     cd.Data.DataValue,

		LastDataCode: cd.DataCode,
		DataCode:     cd.DataCode,

		Value:   0,
		UnitFee: 0,
		Fee:     0,

		PosCode: cd.Pos.PosCode,
		Project: cd.Pos.Project,
	}
}

func CalcCDR(last *ent.CDR, cd ChargeData, ru Ruler) CDR {
	value, fee := calc(last.DataValue, cd.Data.DataValue, ru.UnitFee)
	return CDR{
		DeviceCode: cd.Code,
		DeviceType: cd.Type,

		LastDataTime: last.DataTime,
		DataTime:     cd.Time,

		LastDataValue: last.DataValue,
		DataValue:     cd.Data.DataValue,

		LastDataCode: last.DataCode,
		DataCode:     cd.DataCode,

		PloyID: ru.PloyID,
		RuleID: ru.ID,

		Value:   value,
		UnitFee: ru.UnitFee,
		Fee:     fee,

		PosCode: cd.Pos.PosCode,
		Project: cd.Pos.Project,
	}
}

func calc(lastV, currV int64, unit int64) (v int64, f int64) {
	v = currV - lastV
	f = v * unit / 100
	return
}
