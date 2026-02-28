package chrgg

import "context"

var ErrDataCodeDup = &ChargeErr{Code: "check-datacode", Type: "check", Message: "DataCode重复"}
var ErrTimeBefore = &ChargeErr{Code: "check-datatime", Type: "check", Message: "时间小于之前"}

type CheckFunc func(context.Context, LastCDR, ChargeData) error
type VerifyFunc func(ctx context.Context, last LastCDR, cd ChargeData) bool

func DefaultCheck(ctx context.Context, last LastCDR, cd ChargeData) error {
	if cd.DataCode == last.DataCode {
		return ErrDataCodeDup
	}
	if !cd.DataTime.After(last.DataTime) {
		return ErrTimeBefore
	}
	return nil
}

func DefaultVerify(ctx context.Context, last LastCDR, cd ChargeData) bool {
	if MinOfDay(hourMin(cd.DataTime)) >= 1365 {
		return true
	}

	if cd.Data.DataValue-last.DataValue < 100 { // 小于一个读数
		return false
	}
	return true
}
