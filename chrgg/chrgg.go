package chrgg

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/twiglab/h2o/chrgg/orm/ent"
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

var firstCDRDay = time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local)

type LastCDR struct {
	lastcdr   *ent.CDR
	DataValue int64
	DataCode  string
	DataTime  time.Time
	IsFirst   bool
}

func MakeLast(lcdr *ent.CDR) LastCDR {
	if lcdr == nil {
		return LastCDR{
			DataCode: firstDataCode(),
			DataTime: firstCDRDay,
			IsFirst:  true}
	}

	return LastCDR{
		lastcdr:   lcdr,
		DataValue: lcdr.DataValue,
		DataCode:  lcdr.DataCode,
		DataTime:  lcdr.DataTime,
	}
}

func MinOfDay(h, m int) int {
	return h*60 + m
}

func hourMin(t time.Time) (h, m int) {
	h, m, _ = t.Clock()
	return
}

func firstDataCode() string {
	u, _ := uuid.NewV7()
	return u.String()
}
