package orm

import (
	"context"

	"github.com/twiglab/h2o/chrgg/orm/ent"
	"github.com/twiglab/h2o/chrgg/orm/ent/top"
)

type DBx struct {
	Cli *ent.Client
}

func (d *DBx) LoadLast(ctx context.Context, code, typ string) (r *ent.Top, notfound bool, err error) {
	q := d.Cli.Top.Query()

	q.Where(
		top.DeviceCodeEQ(code),
		top.DeviceTypeEQ(typ),
		top.IsDelEQ(0),
	)

	q.Order(ent.Desc(top.FieldChargeTime))

	r, err = q.First(ctx)
	notfound = ent.IsNotFound(err)

	return
}

// func (d *DBx) NewTop()
