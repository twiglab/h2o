package chrgg

import (
	"context"

	"github.com/twiglab/h2o/chrgg/orm"
	"github.com/twiglab/h2o/chrgg/orm/ent"
	"github.com/twiglab/h2o/pkg/common"
)

type ChargeData struct {
	common.Device
	Pos  common.Pos   `json:"pos,omitzero"`
	Data common.DataV `json:"data"`
	Flag common.Flag  `json:"flag,omitzero"`
}

func cr1(cd ChargeData) CDR {
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

type Ploy struct {
	PloyID  string
	RulerID string
	UnitFee int64
}

func cr2(last *ent.CDR, cd ChargeData, ploy Ploy) CDR {
	v, f := calc(last.DataValue, cd.Data.DataValue, ploy.UnitFee)
	return CDR{
		DeviceCode: cd.Code,
		DeviceType: cd.Type,

		LastDataTime: last.DataTime,
		DataTime:     cd.Time,

		LastDataValue: last.DataValue,
		DataValue:     cd.Data.DataValue,

		LastDataCode: last.DataCode,
		DataCode:     cd.DataCode,

		PloyID: ploy.PloyID,
		RuleID: ploy.RulerID,

		Value:   v,
		UnitFee: ploy.UnitFee,
		Fee:     f,
	}
}

func calc(lastV, currV int64, unit int64) (v int64, f int64) {
	v = currV - lastV
	f = v * unit / 100
	return
}

type Server struct {
	Client *ent.Client
}

func (s *Server) X(ctx context.Context, cd ChargeData) (CDR, error) {
	c, notfound, err := orm.LastCDR(ctx, s.Client, "", "")

	if err != nil {
		return CDR{}, err
	}

	if notfound {
		return cr1(cd), nil
	}

	return cr2(c, cd, Ploy{}), nil
}
