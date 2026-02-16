package chrgg

import (
	"context"

	"github.com/twiglab/h2o/chrgg/cdr"
	"github.com/twiglab/h2o/chrgg/orm"
	"github.com/twiglab/h2o/chrgg/ploy"
)

type ChangeServer struct {
	dbx *orm.DBx
}

func (s *ChangeServer) DoChange(ctx context.Context, cd cdr.ChargeData) (cdr.CDR, error) {
	r, notfound, err := s.dbx.LoadLast(ctx, cd.Code, cd.Type)
	if notfound {
		return cdr.FirstCDR(cd), nil
	}

	if err != nil {
		return cdr.Nil, err
	}

	curr := cdr.CurrentCDR(r, cd, ploy.Ploy{})

	return curr, nil
}
