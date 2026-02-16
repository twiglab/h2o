package orm

import (
	"context"

	"github.com/google/uuid"
	ccdr "github.com/twiglab/h2o/chrgg/cdr"
	"github.com/twiglab/h2o/chrgg/orm/ent"
	"github.com/twiglab/h2o/chrgg/orm/ent/cdr"
)

//go:generate go tool ent generate ./schema --target ./ent

type DBx struct {
	cli *ent.Client
}

func (d *DBx) LoadLast(ctx context.Context, code, typ string) (r *ent.CDR, notfound bool, err error) {
	q := d.cli.CDR.Query()

	q.Where(cdr.DeviceCodeEQ(code), cdr.DeviceTypeEQ(typ))
	q.Limit(1)
	q.Order(ent.Desc(cdr.FieldDataTime))

	//q.Order(ent.Desc(cdr.FieldID))

	if r, err = q.Only(ctx); ent.IsNotFound(err) {
		notfound = true
		err = nil
	}

	return
}

func (d *DBx) SaveCurrent(ctx context.Context, ccdr ccdr.CDR) (r *ent.CDR, err error) {
	cr := d.cli.CDR.Create()
	cr.SetID(cdrid())
	r, err = cr.Save(ctx)
	return
}

func cdrid() string {
	u, _ := uuid.NewV7()
	return u.String()
}
