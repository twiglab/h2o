package chrgg

import (
	"context"
	"strconv"

	"github.com/twiglab/h2o/abm"
)

var RulAloneErr = zr{t: "error", c: "alone"}
var RulAloneNoFound = zr{t: "notfound", c: "alone"}

type AloneRuler struct {
	Code    string
	FeeFen  int64  `db:"fee_fen"`
	PosCode string `db:"pos_code"`
}

func (l AloneRuler) UnitFeeFen() int64 {
	return l.FeeFen
}

func (l AloneRuler) ID() string {
	return l.Code
}

func (l AloneRuler) Type() string {
	return l.PosCode
}

func (l AloneRuler) Category() string {
	return "alone"
}

func (l AloneRuler) Memo() string {
	return l.PosCode
}

func (l AloneRuler) ToStrings() []string {
	return []string{l.Code, strconv.FormatInt(l.FeeFen, 10), l.PosCode}
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
		return RulAloneErr, err
	}

	if ok {
		return a, nil
	}
	return RulAloneNoFound, nil
}
