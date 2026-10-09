package wecom

import (
	"fmt"
	"text/template"
	"time"

	"github.com/twiglab/h2o/proto"
)

const help_cmd_tpl = `
支持下列命令
1. 使用 ? 或者 h 命令显示帮助
2. p <铺位号>，列出对应铺位中表具，如省略铺位号，则列出所有
  例如：p F1_04
3. t <表号>，查询对应仪表的当前用量和充值记录
3. u <表号>，列出指定表号的用量采集列表
  例如：u 909_32
-------
以下命令只能指定人员在指定群组执行
1. c <表号> <增量> [<金额> <起点>]
   指定表号设定限额，限额计算公式为 增量 + max(当前表显示，上次限额)
  例如：c 909_32 50
`

const top_cmd_tpl = `
# 表号：{{ .Code }}
>当前表显 **{{.Curr.DataValue | du}}** 度，采集时间 **{{.Curr.DataTime | datetime}}**
>当前余量（限额 - 表显） **{{ sub .Last.Top .Curr.DataValue | du}}** 度，供电状态 **{{ .Last.Status | usage }}**
{{ if .Tops }}
## 充值记录
| 时间 | 度数 | 限额 | 金额 | 铺位号 |
| :-----: | :----: | :----: | :-----: | :-----: |
{{ range .Tops -}}
| {{.ChargeTime | datetime}} | **{{.Incr | du}}** | **{{.Top | du}}** | {{.Amount | yuan}} | {{.PosCode}} |
{{ end }}
{{ end }}
`

const charge_cmd_tpl = `
# 表号：{{ .Code }} 铺位号：{{ .Device.PosCode }}
>充值度数 **{{.Incr | du}}** 度，起点**{{.Stock | du}}** 度
>充值金额 **{{ .Amount | yuan }}** 元
>充值后限额：**{{ .TopVal | du}}** 度
>当前表显 **{{ .Nh.DataValue | du}}** 度

>充值时间 {{.TopAfter.ChargeTime | datetime}}
>充值代码 {{ .TopAfter.Code }}

充值后可用 t {{ .Code }} 命令查看限额详情
`

const pos_cmd_tpl = `
## 表具列表
| 铺位号 | 表号 | 表具号 | 类型 | 状态 |
| :-----: | :----: | :----: | :-----: | :-----: |
{{ range .Devices -}}
| {{ .PosCode }} | **{{.DeviceCode}}** | **{{.DeviceSn}}** | {{.DeviceType | dt}} | {{.Status}} |
{{ end }}
`

const usage_cmd_tpl = `
## 用量列表 表号：{{ .Code }}
| 表号 | 表显 | 类型 | 时间 |
| :----: | :----: | :-----: | :-----: |
{{ range .All -}}
| **{{.DeviceCode}}** | **{{.DataValue | du}}** | {{.DeviceType | dt}} | {{.DataTime | datetime}} |
{{ end }}
`

const (
	template_help   = "help"
	template_pos    = "pos"
	template_charge = "charge"
	template_top    = "top"
	template_usage  = "usage"
)

func buildTemplate() *template.Template {
	t := template.New("cmd").Funcs(
		template.FuncMap{
			"du":       du,
			"yuan":     du,
			"datetime": fmtTime,
			"dt":       deviceType,
			"sub":      sub,
			"add":      add,
			"usage":    usageStatus,
		},
	)
	template.Must(t.New(template_help).Parse(help_cmd_tpl))
	template.Must(t.New(template_top).Parse(top_cmd_tpl))
	template.Must(t.New(template_charge).Parse(charge_cmd_tpl))
	template.Must(t.New(template_pos).Parse(pos_cmd_tpl))
	template.Must(t.New(template_usage).Parse(usage_cmd_tpl))
	return t
}

func du(i int64) string {
	x := float64(i) * 0.01
	return fmt.Sprintf("%.2f", x)
}

func fmtTimePtr(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(time.DateTime)
}

func fmtTime(t time.Time) string {
	return t.Format(time.DateTime)
}

func sub(a, b int64) int64 {
	return a - b
}

func add(a, b int64) int64 {
	return a + b
}
func deviceType(t string) string {
	switch t {
	case proto.TYPE_ELECTRICITY:
		return "电表"
	case proto.TYPE_WATER:
		return "水表"
	case proto.TYPE_GAS:
		return "气表"
	}
	return "unknow"
}

func usageStatus(status int) string {
	switch status {
	case 0:
		return "正常"
	case 1:
		return "结束"
	}

	return "非正常状态"
}
