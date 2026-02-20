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

type ChargeEngine interface {
	GetResult(context.Context, ChargeData) (Ruler, error)
}

type ZeroCe struct {
}

func (ZeroCe) GetResult(context.Context, ChargeData) (Ruler, error) {
	return Ruler{
		RulerID: "zero",
		PloyID:  "zero",
		Memo:    "zero",
	}, nil
}
