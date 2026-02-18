package chrgg

import (
	"context"

	"github.com/google/uuid"
	"github.com/twiglab/h2o/chrgg/orm/ent"
	"github.com/twiglab/h2o/chrgg/orm/ent/cdr"
)

type DBx struct {
	Cli *ent.Client
}

func (d *DBx) LoadLast(ctx context.Context, code, typ string) (r *ent.CDR, notfound bool, err error) {
	q := d.Cli.CDR.Query()

	q.Where(cdr.DeviceCodeEQ(code), cdr.DeviceTypeEQ(typ))
	q.Limit(1)
	q.Order(ent.Desc(cdr.FieldDataTime))

	if r, err = q.Only(ctx); ent.IsNotFound(err) {
		notfound = true
		err = nil
	}

	return
}

func (d *DBx) SaveCurrent(ctx context.Context, cdr CDR) (r *ent.CDR, err error) {
	cr := d.Cli.CDR.Create()

	cr.SetID(cdrid())
	cr.SetDeviceCode(cdr.DeviceCode)
	cr.SetDeviceType(cdr.DeviceType)

	cr.SetLastDataCode(cdr.LastDataCode)
	cr.SetDataCode(cdr.DataCode)

	cr.SetLastDataTime(cdr.LastDataTime)
	cr.SetDataTime(cdr.DataTime)

	cr.SetLastDataValue(cdr.LastDataValue)
	cr.SetDataValue(cdr.DataValue)

	cr.SetValue(cdr.Value)

	cr.SetPloyID(cdr.PloyID)
	cr.SetRuleID(cdr.RuleID)
	cr.SetUnitFee(cdr.UnitFee)

	cr.SetPosCode(cdr.PosCode)
	cr.SetProject(cdr.Project)

	cr.SetMemo(cdr.Momo)

	r, err = cr.Save(ctx)
	return
}

func cdrid() string {
	u, _ := uuid.NewV7()
	return u.String()
}
