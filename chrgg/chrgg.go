package chrgg

import (
	"context"
	"log/slog"

	"github.com/twiglab/h2o/chrgg/cdr"
	"github.com/twiglab/h2o/chrgg/orm"
	"github.com/twiglab/h2o/chrgg/orm/ent"
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
	)

	r, notfound, err = s.dbx.LoadLast(ctx, cd.Code, cd.Type)

	if notfound {
		newCDR = cdr.FirstCDR(cd)
		err = nil
		return
	}

	if err != nil {
		return
	}

	ru := s.re.GetRuler(ctx, cd)

	newCDR = cdr.NewCDR(r, cd, ru)

	// log cdr

	_, err = s.dbx.SaveCurrent(ctx, newCDR)

	return
}
