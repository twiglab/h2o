package chrgg

import (
	"context"
)

type SkipReturn struct {
	Message string
	OK      bool
}

func NewSkipReturn(msg string) SkipReturn {
	return SkipReturn{Message: msg, OK: true}
}

var pass = SkipReturn{OK: false}

func SkipPass() SkipReturn {
	return pass
}

func (v SkipReturn) String() string {
	return v.Message
}

type SkipFunc func(ctx context.Context, last LastCDR, cd ChargeData) SkipReturn

func NewSkipChain(fs ...SkipFunc) SkipFunc {
	return func(ctx context.Context, last LastCDR, cd ChargeData) SkipReturn {
		for _, f := range fs {
			if r := f(ctx, last, cd); !r.OK {
				return r
			}
		}
		return SkipPass()
	}
}

func DefaultSkipChain() SkipFunc {
	return NewSkipChain(Skip22h45m, SkipValueLess)
}

const tm_22h45m = 1365 // 22:45分的分钟数

func Skip22h45m(ctx context.Context, last LastCDR, cd ChargeData) SkipReturn {
	if MinPerDay(cd.DataTime) >= tm_22h45m {
		if MinPerDay(last.DataTime) >= tm_22h45m && !IsValueChange(last, cd) {
			return NewSkipReturn("上一条已经存在，且无变化")
		}
	}
	return SkipPass()
}

func SkipValueLess(ctx context.Context, last LastCDR, cd ChargeData) SkipReturn {
	if !IsValueChange(last, cd) {
		return NewSkipReturn("小于一个读数")
	}
	return SkipPass()
}

func IsValueChange(last LastCDR, cd ChargeData) bool {
	return (cd.Data.DataValue - last.DataValue) > 100
}
