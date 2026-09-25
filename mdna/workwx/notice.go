package workwx

import (
	"fmt"
	"text/template"
	"time"

	"github.com/twiglab/h2o/mdna"
	"github.com/twiglab/h2o/pkg/common"
)

const chrggTempl = `
# 小于 {{ (.Notice.Charge.Data.DataValue - .Notice.Charge.Top.Top) | du }} 度能耗提醒
截至 {{ .Now.Format "2006.01.02 15:04" }}
铺位号：{{ .Notice.Device.PosCode }}
表号：{{ .Notice.Device.Code }}/{{ .Notice.Device.DeviceType | deviceType }}
表具：{{ .Notice.Device.DeviceSn }}
限额：{{ .Notice.Charge.Top.Top | du }}
当前：{{ .Notice.Charge.Data.DataValue | du }}
** 请注意提前充值 **
充值记录 {{ .Notice.Charge.Top.Code }}
`

const optTempl = `
# 拉闸提醒
截至 {{ .Now.Format "2006.01.02 15:04" }}
铺位号：{{ .Notice.Device.PosCode }}
表号：{{ .Notice.Device.Code }}/{{ .Notice.Device.DeviceType | deviceType }}
表具：{{ .Notice.Device.DeviceSn }}
**超出限额, 已拉闸断开**
`

type NotictBox struct {
	Now    time.Time
	Notice mdna.Notice
}

func du(total int64) string {
	f := float64(total) * 0.01
	return fmt.Sprintf("%.2f", f)
}

func fm() template.FuncMap {
	return template.FuncMap{
		"du":         du,
		"deviceType": deviceType,
	}
}

func deviceType(t string) string {
	switch t {
	case common.TYPE_ELECTRICITY:
		return "电表"
	case common.TYPE_WATER:
		return "水表"
	case common.TYPE_GAS:
		return "气表"
	}
	return "unknow"
}
