package chrgg

import (
	"context"
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
