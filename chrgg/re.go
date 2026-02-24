package chrgg

import (
	"context"

	"github.com/twiglab/h2o/abm"
)

type ChargeRuler interface {
	ID() string
	Type() string
	Category() string
	UnitFeeFen() int64
	Memo() string
}

type ChargeEngine interface {
	GetRuler(context.Context, ChargeData) (ChargeRuler, error)
}

type ZeroRuler string

func (z ZeroRuler) GetRuler(_ context.Context, _ ChargeData) (ChargeRuler, error) {
	return z, nil
}

func (ZeroRuler) UnitFeeFen() int64 {
	return 0
}

func (ZeroRuler) ID() string {
	return "zero"
}

func (ZeroRuler) Type() string {
	return "zero"
}

func (ZeroRuler) Category() string {
	return "zero"
}

func (z ZeroRuler) Memo() string {
	return string(z)
}

type AloneRuler struct {
	Code       string
	UnitFeeFen uint64
}

type LocAloneEng struct {
	knowledge *abm.DuckABM[string, AloneRuler]
}

func (l *LocAloneEng) GetRuler(ctx context.Context, cd ChargeData) (ChargeRuler, error) {
	a, ok, err := l.knowledge.Get(ctx, cd.Code)
	if err != nil {
		return ZeroRuler("err"), nil
	}

	if ok {
		return &AloneRuler{Code: a.Code, UnitFeeFen: a.UnitFeeFen}, nil
	}
	return ZeroRuler("err"), nil
}
