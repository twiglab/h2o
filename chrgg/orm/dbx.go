package orm

import (
	"context"

	"github.com/twiglab/h2o/chrgg/orm/ent"
	"github.com/twiglab/h2o/chrgg/orm/ent/valuecharge"
)

type DBx struct {
	Cli *ent.Client
}

func (d *DBx) LoadLast(ctx context.Context, code, typ string) (r *ent.ValueCharge, notfound bool, err error) {
	q := d.Cli.ValueCharge.Query()

	q.Where(
		valuecharge.DeviceCodeEQ(code),
		valuecharge.DeviceTypeEQ(typ),
		valuecharge.IsDelEQ(0),
	)

	q.Order(ent.Desc(valuecharge.FieldChargeTime))

	r, err = q.First(ctx)
	notfound = ent.IsNotFound(err)

	return
}
