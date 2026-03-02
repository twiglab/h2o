package chrgg

import (
	"fmt"
	"time"
)

type ChargeErr struct {
	Code    string
	Type    string
	Message string
}

func (e ChargeErr) Error() string {
	return fmt.Sprintf("Charge Error: code = %s, type = %s, message = %s", e.Code, e.Type, e.Message)
}

type ChargeData struct {
	MeterData
}

func MinPerDay(t time.Time) int {
	h, m, _ := t.Clock()
	return h*60 + m
}
