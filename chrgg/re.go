package chrgg

import (
	"context"

	"github.com/twiglab/h2o/chrgg/cdr"
	"github.com/twiglab/h2o/chrgg/ploy"
)

type RulerEngine interface {
	GetResult(context.Context, cdr.ChargeData) (ploy.Ruler, error)
}

type defRE struct {
}

func (defRE) GetRuler(context.Context, cdr.ChargeData) (ploy.Ruler, error) {
	return ploy.Ruler{
		ID:     "NIL",
		PloyID: "NIL",
		Descr:  "NIL",
	}, nil
}
