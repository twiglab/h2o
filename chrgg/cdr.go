package chrgg

import (
	"time"

	"github.com/google/uuid"
	"github.com/twiglab/h2o/chrgg/orm/ent"
)

var nilCDR CDR

var RulNew = zr{t: "new", c: "new"}

var firstCDRDay = time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local)

type LastCDR struct {
	lastcdr *ent.CDR

	DataValue int64
	DataCode  string
	DataTime  time.Time
	Value     int64

	IsFirst bool
}

func MakeLast(lcdr *ent.CDR) LastCDR {
	if lcdr == nil {
		return LastCDR{
			DataCode: firstDataCode(),
			DataTime: firstCDRDay,
			IsFirst:  true,
		}
	}

	return LastCDR{
		lastcdr: lcdr,

		DataValue: lcdr.DataValue,
		DataCode:  lcdr.DataCode,
		DataTime:  lcdr.DataTime,
		Value:     lcdr.Value,
	}
}

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

func firstDataCode() string {
	u, _ := uuid.NewV7()
	return u.String()
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
