package orm

import (
	"context"

	"github.com/google/uuid"
	"github.com/twiglab/h2o/chrgg/orm/ent"
	"github.com/twiglab/h2o/chrgg/orm/ent/cdr"
)

//go:generate go tool ent generate ./schema --target ./ent

func LastCDR(ctx context.Context, cli *ent.Client, code, typ string) (r *ent.CDR, notfound bool, err error) {
	q := cli.CDR.Query()

	q.Where(cdr.DeviceCodeEQ(code), cdr.DeviceTypeEQ(typ))
	q.Limit(1)

	r, err = q.Only(ctx)
	notfound = ent.IsNotFound(err)

	return
}

func CDRID() string {
	u, _ := uuid.NewV7()
	return u.String()
}
