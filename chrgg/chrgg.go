package chrgg

import (
	"context"
	"log/slog"

	"github.com/twiglab/h2o/chrgg/cdr"
	"github.com/twiglab/h2o/chrgg/orm"
)

type ChangeServer struct {
	dbx    *orm.DBx
	cdrLog *slog.Logger

	re RulerEngine
}

func (s *ChangeServer) calc(ctx context.Context, cd cdr.ChargeData) (cdr.CDR, error) {
	last, notfound, err := s.dbx.LoadLast(ctx, cd.Code, cd.Type)
	if err != nil {
		return cdr.Nil, err
	}

	if notfound {
		return cdr.FirstCDR(cd), nil
	}

	ru, err := s.re.GetResult(ctx, cd)
	if err != nil {
		return cdr.Nil, err
	}

	return cdr.CalcCDR(last, cd, ru), nil
}

func (s *ChangeServer) DoChange(ctx context.Context, cd cdr.ChargeData) (nc cdr.CDR, err error) {
	if nc, err = s.calc(ctx, cd); err == nil { // err == nil
		_, err = s.dbx.SaveCurrent(ctx, nc)
	}
	return
}
