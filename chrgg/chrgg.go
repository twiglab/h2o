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

func MinPerDay(t time.Time) int {
	h, m, _ := t.Clock()
	return h*60 + m
}

var m_22_45 = 1365 // 22:45分的分钟数
