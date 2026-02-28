package chrgg

import (
	"fmt"
	"time"

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

type LastCDR struct {
	lastcdr   *ent.CDR
	DataValue int64
	DataCode  string
	DataTime  time.Time
	IsFirst   bool
}

func MakeLast(lcdr *ent.CDR) LastCDR {
	if lcdr == nil {
		return LastCDR{DataTime: time.Now(), IsFirst: true}
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
