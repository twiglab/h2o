package chrgg

import (
	"context"
	"log/slog"

	"github.com/twiglab/h2o/chrgg/cdr"
	"github.com/twiglab/h2o/chrgg/orm"
	"github.com/twiglab/h2o/chrgg/orm/ent"
	"github.com/twiglab/h2o/chrgg/ploy"
)

type ChangeServer struct {
	dbx    *orm.DBx
	cdrLog *slog.Logger

	re RulerEngine
}

func (s *ChangeServer) DoChange(ctx context.Context, cd cdr.ChargeData) (newCDR cdr.CDR, err error) {
	var (
		r        *ent.CDR
		notfound bool
		ru       ploy.Ruler
	)

	r, notfound, err = s.dbx.LoadLast(ctx, cd.Code, cd.Type)

	if err != nil {
		return
	}

	if notfound {
		newCDR = cdr.FirstCDR(cd)
		return
	}

	if ru, err = s.re.GetResult(ctx, cd); err != nil {
		return
	}

	newCDR = cdr.NewCDR(r, cd, ru)

	// log cdr

	_, err = s.dbx.SaveCurrent(ctx, newCDR)

	return
}
