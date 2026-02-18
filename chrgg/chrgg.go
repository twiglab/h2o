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
	h, m, _ := t.Clock()
	if h*60+m > 1365 { // > 22:45
		return true
	}
	return false
}
