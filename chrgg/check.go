package chrgg

import "context"

var ErrDataCodeDup = &ChargeErr{Code: "check-datacode", Type: "check", Message: "DataCode重复"}
var ErrTimeBefore = &ChargeErr{Code: "check-datatime", Type: "check", Message: "时间小于之前"}

type VerifyReturn struct {
	Message string
}

func (v VerifyReturn) String() string {
	return v.Message
}

type CheckFunc func(context.Context, LastCDR, ChargeData) error
type VerifyFunc func(ctx context.Context, last LastCDR, cd ChargeData) (VerifyReturn, bool)

func DefaultCheck(ctx context.Context, last LastCDR, cd ChargeData) error {
	if cd.DataCode == last.DataCode {
		return ErrDataCodeDup
	}
	if !cd.DataTime.After(last.DataTime) {
		return ErrTimeBefore
	}
	return nil
}

func DefaultVerify(ctx context.Context, last LastCDR, cd ChargeData) (VerifyReturn, bool) {
	if MinOfDay(hourMin(cd.DataTime)) >= 1365 {
		return VerifyReturn{}, true
	}

	if cd.Data.DataValue-last.DataValue < 100 { // 小于一个读数
		return VerifyReturn{Message: "小于一个读数"}, false
	}
	return VerifyReturn{}, true
}
