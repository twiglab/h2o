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

type ChargeData struct {
	MeterData
}

func isInTimeRange(t time.Time) bool {
	if DayMinute(t) > 1365 { // > 22:45
		return true
	}
	return false
}

func DayMinute(t time.Time) int {
	h, m, _ := t.Clock()
	return h*60 + m
}
