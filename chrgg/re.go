package chrgg

import (
	"context"
)

type Ruler struct {
	RulerID string
	PloyID  string

	UnitFee int64

	Memo string
}

type Ploy interface {
	GetResult(context.Context, ChargeData) (Ruler, error)
}

type ZeroPloy struct {
}

func (ZeroPloy) GetRuler(context.Context, ChargeData) (Ruler, error) {
	return Ruler{
		RulerID: "zero",
		PloyID:  "zero",
		Memo:    "zero",
	}, nil
}
