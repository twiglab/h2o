package chrgg

import (
	"context"
)

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
