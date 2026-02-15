package chrgg

import (
	"context"

	"github.com/twiglab/h2o/chrgg/orm"
	"github.com/twiglab/h2o/chrgg/orm/ent"
	"github.com/twiglab/h2o/pkg/common"
)

type Data struct {
	DataValue int64
}

type ChargeData struct {
	common.Device
	Pos  common.Pos  `json:"pos,omitzero"`
	Data Data        `json:"data"`
	Flag common.Flag `json:"flag,omitzero"`
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
	}
}

type Ploy struct {
	PloyID  string
	RulerID string
	UnitFee int64
}

func cr2(last *ent.CDR, cd ChargeData, ploy Ploy) CDR {
	return CDR{
		DeviceCode: cd.Code,
		DeviceType: cd.Type,

		LastDataTime: last.DataTime,
		DataTime:     cd.DataTime,

		LastDataValue: last.DataValue,
		DataValue:     cd.DataValue,

		LastDataCode: last.DataCode,
		DataCode:     cd.DataCode,

		PloyID: ploy.PloyID,
		RuleID: ploy.RulerID,

		Value:   cd.DataValue - last.DataValue,
		UnitFee: ploy.UnitFee,
		Fee:     (cd.DataValue - last.DataValue) * ploy.UnitFee,
	}
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
