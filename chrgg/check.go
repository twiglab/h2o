package chrgg

import "context"

var ErrDataCodeDup = &ChargeErr{Code: "check-datacode", Type: "check", Message: "DataCode重复"}
var ErrTimeBefore = &ChargeErr{Code: "check-datatime", Type: "check", Message: "时间小于之前"}

type CheckFunc func(context.Context, LastCDR, ChargeData) error

func DefaultCheck(ctx context.Context, last LastCDR, cd ChargeData) error {
	if cd.DataCode == last.DataCode {
		return ErrDataCodeDup
	}
	if !cd.DataTime.After(last.DataTime) {
		return ErrTimeBefore
	}
	return nil
}

var pass = SkipReturn{}

type SkipReturn struct {
	Message string
	OK      bool
}

func NewSkipReturn(msg string) SkipReturn {
	return SkipReturn{Message: msg, OK: true}
}

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
			if r := f(ctx, last, cd); r.OK {
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
