package workwx

import (
	"fmt"
	"text/template"

	"github.com/twiglab/h2o/pkg/common"
)

const chrggTempl = `
截至 {{ .Now.Format "2006.01.02 15:04" }}
铺位号：{{ .Notice.Device.PosCode }}
表号：{{ .Notice.Device.Code }}/{{ .Notice.Device.DeviceType | deviceType }}
表具：{{ .Notice.Device.DeviceSn }}
限额：{{ .Notice.Charge.Top.Top | du }}
当前：{{ .Notice.Charge.Data.DataValue | du }}
** 请注意提前充值 **
充值记录 {{ .Notice.Charge.Top.Code }}
`

// # 小于 {{ (.Notice.Charge.Data.DataValue - .Notice.Charge.Top.Top) | du }} 度能耗提醒

const optTempl = `
# 拉闸消息
截至 {{ .Now.Format "2006.01.02 15:04" }}
铺位：{{ .Device.PosCode }}
表号：{{ .Device.DeviceCode }}/{{ .Device.DeviceType | deviceType }}
表具：{{ .Device.DeviceSn }}
**超出限额, 已拉闸断开**（请用 t {{ .Device.DeviceCode }} 指令查看限额信息）
`

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
