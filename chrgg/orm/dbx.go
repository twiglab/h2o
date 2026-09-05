package orm

import (
	"context"

	"github.com/twiglab/h2o/chrgg/orm/ent"
	"github.com/twiglab/h2o/chrgg/orm/ent/vvc"
)

type DBx struct {
	Cli *ent.Client
}

func (d *DBx) LoadLast(ctx context.Context, code, typ string) (r *ent.VVC, notfound bool, err error) {
	q := d.Cli.VVC.Query()

	q.Where(vvc.DeviceCodeEQ(code), vvc.DeviceTypeEQ(typ))
	q.Limit(1)
	//q.Order(ent.Desc(cdr.FieldDataTime))

	r, err = q.First(ctx)
	notfound = ent.IsNotFound(err)

	return
}
