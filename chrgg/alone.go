package chrgg

import (
	"context"

	"github.com/twiglab/h2o/abm"
)

type AloneRuler struct {
	Code string
	Fee  int64
}

func (l AloneRuler) UnitFeeFen() int64 {
	return l.Fee
}

func (l AloneRuler) ID() string {
	return l.Code
}

func (l AloneRuler) Type() string {
	return "alone"
}

func (l AloneRuler) Category() string {
	return l.Code
}

func (l AloneRuler) Memo() string {
	return l.Code
}

type AloneEng struct {
	knowledge *abm.DuckABM[string, AloneRuler]
}

func (l *AloneEng) GetRuler(ctx context.Context, cd ChargeData) (ChargeRuler, error) {
	a, ok, err := l.knowledge.Get(ctx, cd.Code)
	if err != nil {
		return ZeroRuler("err"), nil
	}

	if ok {
		return a, nil
	}
	return ZeroRuler("err"), nil
}
