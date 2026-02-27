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

type zr struct {
	t, c string
}

func (z zr) GetRuler(_ context.Context, _ ChargeData) (ChargeRuler, error) {
	return z, nil
}

func (zr) UnitFeeFen() int64 {
	return 0
}

func (zr) ID() string {
	return "zero"
}

func (z zr) Type() string {
	return z.t
}

func (z zr) Category() string {
	return z.c
}

func (z zr) Memo() string {
	return "zero"
}

var EngZ = zr{t: "z", c: "z"}
