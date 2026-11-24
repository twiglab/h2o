package quota

import (
	"context"

	"github.com/twiglab/h2o/dbcli/ent"
	"github.com/twiglab/h2o/dbcli/ent/top"
)

type DBx struct {
	Cli *ent.Client
}

func (d *DBx) LoadLast(ctx context.Context, code, project string) (r *ent.Top, err error) {
	q := d.Cli.Top.Query()

	q.Where(
		top.DeviceCodeEQ(code),
		top.ProjectEQ(project),
		top.IsDelEQ(0),
	)

	q.Order(ent.Desc(top.FieldChargeTime))

	r, err = q.First(ctx)
	return
}
