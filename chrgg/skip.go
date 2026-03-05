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

var noSkip = SkipReturn{Message: "正常执行", Skip: false}

func NoSkip() SkipReturn {
	return noSkip
}

func (v SkipReturn) String() string {
	return v.Message
}

type SkipFunc func(ctx context.Context, last LastCDR, cd ChargeData) SkipReturn

const tm_22h45m = 1365 // 22:45分的分钟数

func DefaultSkip(_ context.Context, last LastCDR, cd ChargeData) SkipReturn {
	if MinPerDay(cd.DataTime) < tm_22h45m && !IsValueChangeLeeway(last, cd, 100) {
		return SkipOK("小于一个读数")
	}
	if MinPerDay(last.DataTime) >= tm_22h45m && !IsValueChangeLeeway(last, cd, 100) {
		return SkipOK("当天记录已存在，且小于一个读数")
	}
	return NoSkip()
}

func IsValueChangeLeeway(last LastCDR, cd ChargeData, leeway int64) bool {
	return (cd.Data.DataValue - last.DataValue) > leeway
}
