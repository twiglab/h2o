package chrgg

import (
	"fmt"
	"time"
)

const CLIENT_ID = "chrgg"

type ChargeErr struct {
	Code    string
	Type    string
	Message string
}

func (e ChargeErr) Error() string {
	return fmt.Sprintf("Charge Error: code = %s, type = %s, message = %s", e.Code, e.Type, e.Message)
}

var ErrDataCodeDup = &ChargeErr{Code: "check-datacode", Type: "check", Message: "DataCode重复"}
var ErrTimeBefore = &ChargeErr{Code: "check-datatime", Type: "check", Message: "时间小于之前"}

type ChargeData struct {
	MeterData
}

func MinOfDay(h, m int) int {
	return h*60 + m
}

func hourMin(t time.Time) (h, m int) {
	h, m, _ = t.Clock()
	return
}
