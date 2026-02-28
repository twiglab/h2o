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
	value, fee := calcFee(last.DataValue, cd.Data.DataValue, cr.UnitFeeFen())
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

func calcFee(lastV, currV int64, unitFen int64) (
	valuePer100 int64, /* 表差 1/100 计量单位*/
	feeFen int64, /* 1 计量单位多少分 */
) {
	// 注意这里的表差是 1/100 计量单位
	// 为了规避误差，这里先做乘法再做除法
	// 最后除以100是把 1/100 计量单位，转成 1 计量单位
	// 不是分转元! 所以最后得出的结果是 1 个计量单位多少分
	valuePer100 = currV - lastV
	feeFen = (valuePer100 * unitFen) / 100
	return
}
