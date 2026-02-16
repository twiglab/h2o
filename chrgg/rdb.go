package chrgg

import (
	"context"

	"github.com/twiglab/h2o/chrgg/cdr"
	"github.com/twiglab/h2o/chrgg/ploy"
)

type RulerEngine interface {
	GetRuler(context.Context, cdr.ChargeData) ploy.Ruler
}

type defRE struct {
}

func (defRE) GetRuler(context.Context, cdr.ChargeData) ploy.Ruler {
	return ploy.Ruler{
		ID:     "NIL",
		PloyID: "NIL",
		Descr:  "NIL",
	}
}
