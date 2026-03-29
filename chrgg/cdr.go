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
	lastcdr *ent.CDR `json:"-"`

	DataValue int64     `json:"data_value"`
	DataCode  string    `json:"data_code"`
	DataTime  time.Time `json:"data_time"`
	Value     int64     `json:"value"`

	IsFirst bool `json:"is_first"`
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

