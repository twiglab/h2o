package chrgg

import (
	"context"
)

type SkipReturn struct {
	Message string
	Skip    bool
}

func SkipOK(msg string) SkipReturn {
	return SkipReturn{Message: msg, Skip: true}
}

var noSkip = SkipReturn{Skip: false}

func NoSkip() SkipReturn {
	return noSkip
}

func (v SkipReturn) String() string {
	return v.Message
}

type SkipFunc func(ctx context.Context, last LastCDR, cd ChargeData) SkipReturn

func NewSkipChain(fs ...SkipFunc) SkipFunc {
	return func(ctx context.Context, last LastCDR, cd ChargeData) SkipReturn {
		for _, f := range fs {
			if r := f(ctx, last, cd); r.Skip {
				return r
			}
		}
		return NoSkip()
	}
}

func DefaultSkipChain() SkipFunc {
	return NewSkipChain(
		Skip22h45m,
		SkipValueLess,
	)
}

const tm_22h45m = 1365 // 22:45分的分钟数

func Skip22h45m(ctx context.Context, last LastCDR, cd ChargeData) SkipReturn {
	if MinPerDay(cd.DataTime) >= tm_22h45m {
		if MinPerDay(last.DataTime) >= tm_22h45m && !IsValueChangeLeeway(last, cd, 100) {
			return SkipOK("上一条已经存在，且无变化")
		}
	}
	return NoSkip()
}

func SkipValueLess(ctx context.Context, last LastCDR, cd ChargeData) SkipReturn {
	if !IsValueChangeLeeway(last, cd, 100) {
		return SkipOK("小于一个读数")
	}
	return NoSkip()
}

func IsValueChangeLeeway(last LastCDR, cd ChargeData, leeway int64) bool {
	return (cd.Data.DataValue - last.DataValue) > leeway
}
