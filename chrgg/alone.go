package chrgg

import (
	"context"

	"github.com/twiglab/h2o/abm"
)

type AloneRuler struct {
	Code    string
	Fee     int64
	PosCode string `db:"pos_code"`
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

type AloneEngine struct {
	knowledge *abm.DuckABM[string, AloneRuler]
}

func NewAloneEngine(knowledge *abm.DuckABM[string, AloneRuler]) *AloneEngine {
	return &AloneEngine{knowledge: knowledge}
}

func (l *AloneEngine) GetRuler(ctx context.Context, cd ChargeData) (ChargeRuler, error) {
	a, ok, err := l.knowledge.Get(ctx, cd.Code)
	if err != nil {
		return ZeroRuler("err"), nil
	}

	if ok {
		return a, nil
	}
	return ZeroRuler("err"), nil
}
