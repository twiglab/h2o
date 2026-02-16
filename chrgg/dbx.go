package chrgg

import (
	"context"

	"github.com/google/uuid"
	"github.com/twiglab/h2o/chrgg/orm/ent"
	"github.com/twiglab/h2o/chrgg/orm/ent/cdr"
)

type DBx struct {
	cli *ent.Client
}

func (d *DBx) LoadLast(ctx context.Context, code, typ string) (r *ent.CDR, notfound bool, err error) {
	q := d.cli.CDR.Query()

	q.Where(cdr.DeviceCodeEQ(code), cdr.DeviceTypeEQ(typ))
	q.Limit(1)
	q.Order(ent.Desc(cdr.FieldDataTime))

	if r, err = q.Only(ctx); ent.IsNotFound(err) {
		notfound = true
		err = nil
	}

	return
}

func (d *DBx) SaveCurrent(ctx context.Context, ccdr CDR) (r *ent.CDR, err error) {
	cr := d.cli.CDR.Create()
	cr.SetID(cdrid())
	r, err = cr.Save(ctx)
	return
}

func cdrid() string {
	u, _ := uuid.NewV7()
	return u.String()
}
