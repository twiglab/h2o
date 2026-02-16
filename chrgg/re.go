package chrgg

import (
	"context"
)

type Ruler struct {
	ID     string
	PloyID string

	UnitFee int64

	Descr string
}

type RulerEngine interface {
	GetResult(context.Context, ChargeData) (Ruler, error)
}

type defRE struct {
}

func (defRE) GetRuler(context.Context, ChargeData) (Ruler, error) {
	return Ruler{
		ID:     "NIL",
		PloyID: "NIL",
		Descr:  "NIL",
	}, nil
}
