package wecom

import (
	"cmp"
	"fmt"
	"text/template"
	"time"

	"github.com/twiglab/h2o/proto"
)

const help_cmd_tpl = `
支持下列命令
1. 使用 ? 或者 h 命令显示帮助
2. p [<铺位号>]，列出对应铺位中表具，如省略铺位号，则列出所有
  例如：p F1_04
3. u <表号>，列出指定表号的使用情况
  例如：u 909_32
4. r 查看当前用量报表
-------
以下命令只能指定人员在指定群组执行
1. t <表号> <增量> [<起点>]
  指定表号设定限额，限额计算公式为 增量 + 上次限额
  例如：t 909_32 50 在上次限额上增加50个表显
        t 909_32 50 100 设置起点为100，并增加50个表显示, 设置完成后，限额为150
  **特殊用法**：t 909_32 0 100 增量为0，起点为100，限额为100，用于设置绝对限额，如果起点设置很大，相当于**强制合闸**
            t 909_32 0 0   增量为0，起点为0，限额为0，相当于**强制拉闸**
  **t 命令，不受 2 分钟限制，即用即效**

2. m <表号> <金额> [<起点>]
  指定表号设定限额，限额计算公式为 金额 / 单价 / 倍率
  **使用 m 命令前，必须正确设置倍率和单价，否则报错！**
  例如：m 909_32 100 12345
  给909_32 充值100元，起点为12345（表显），起点可以省略，如果省略为上一次限额，如果是第一次则为 0
  此命令会自动计算当次充值的限额，如果没有设置倍率和单价，则提示错误，单价可用 p 指令查询
  **优先考虑使用 m 指令，为了防止重复充值，m 命令使用后 2 分钟内不得对同一个表号再次使用**
`

const charge_cmd_tpl = `
# 表号：{{ .Code }} 铺位号：{{ .Device.PosCode }}（{{ .Device.PosName }}）
>表具：**{{.Device.DeviceSn}}**
>倍率：**{{.Device.Rate}}**
>表显：**{{ .Nh.DataValue | du}}**，采集时间 {{ .Nh.DataTime | datetime}}

>增量：**{{.Incr | du}}** ，存量：**{{.Stock | du}}**
>限额：**{{ .TopVal | du}}**

>可用电量 [(限额 - 当前表显) * 倍率]：**{{ eng .TopVal .Nh.DataValue .Device.Rate | du}}** 度

>本次充值金额：**{{ .Amount | yuan }}** 元
>充值时间：{{.TopAfter.ChargeTime | datetime}}
>充值代码：{{ .TopAfter.Code }}

充值后可用 u {{ .Code }} 命令查看使用详情
`

const pos_cmd_tpl = `
## 表具列表
| 铺位号 | 商铺名称 | 表号 | 表具号 | 类型 | 倍率 | 单价（元） |
| :-----: | :----: | :----: | :-----: | :-----: | :-----: | :-----: |
{{ range .Devices -}}
| {{ .PosCode }} | {{ .PosName }} |**{{.DeviceCode}}** | **{{.DeviceSn}}** | {{.DeviceType | dt}} | {{.Rate}} | {{.UnitPrice|yuan}} |
{{ end }}
`

const usage_cmd_tpl = `
## 用量列表 表号：{{ .Code }}，铺位号：{{.Device.PosCode}}（{{.Device.PosName}}）
>表具：**{{.Device.DeviceSn}}**
>倍率：**{{.Device.Rate}}**
>最后一次采集时间：{{.Record.DataTime | datetime}}
{{ if .Record.DataTime | since12h }}
>**最后一次采集大于12小时，表具可能存在错误**
{{ end }}
{{ if .Top }}
>可用电量 [(限额 - 当前表显) * 倍率]：**{{ eng .Top.Top .Record.DataValue .Device.Rate | du}}** 度，供电状态 **{{ .Top.Status | usage }}**
{{else}}
>无限额记录
{{ end }}
## 最近充限额记录
| 时间 | 限额(表显) | 增量(表显) | 起点(表显) | 金额 |
| :-----: | :----: | :----: | :-----: | :-----: |
{{ range .Tops -}}
| {{.ChargeTime | datetime}} | {{.Top | du}} | **{{.Incr | du}}** | {{.Stock | du}} | {{.Amount | yuan}} |
{{ end }}
## 最近采集记录
| 表号 | 表显 | 类型 | 时间 |
| :----: | :----: | :-----: | :-----: |
{{ range .Records -}}
| **{{.DeviceCode}}** | **{{.DataValue | du}}** | {{.DeviceType | dt}} | {{.DataTime | datetime}} |
{{end}}
`

const report_cmd_tpl = `
## 用量报告 时间：{{ .Now | datetime }}
| 商铺名称 | 预警 |   剩余电量（度） |   预警电量（度）  | 铺位编号 |   表号  |   类型  |  采集时间 |
| :------: |:---: | :------: | :-----: | :-----: | :----: | :----: | :-----: |
{{ range .Items -}}
| {{ .PosName }} | {{ cui .DiffEng .Alarm }} | {{eng .DiffEng 0 .Rate | du}} ({{tian .DiffEng .Alarm}}天) |{{eng .Alarm 0 .Rate | du}} | {{.PosCode}} | {{.DeviceCode}} | {{.DeviceType | dt}} | {{.DataTime | datetime}} |
{{ end }}
`

const (
	template_help   = "help"
	template_pos    = "pos"
	template_charge = "charge"
	template_usage  = "usage"
	template_report = "report"
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
			"eng":      eng,
			"cui":      cui,
			"tian":     tian,

			"since12h": since12h,
		},
	)
	template.Must(t.New(template_help).Parse(help_cmd_tpl))
	template.Must(t.New(template_charge).Parse(charge_cmd_tpl))
	template.Must(t.New(template_pos).Parse(pos_cmd_tpl))
	template.Must(t.New(template_usage).Parse(usage_cmd_tpl))
	template.Must(t.New(template_report).Parse(report_cmd_tpl))
	return t
}

func du(i int64) string {
	x := float64(i) * 0.01
	return fmt.Sprintf("%.2f", x)
}

func eng(t, v, r int64) int64 {
	return (t - v) * r
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
	case 3:
		return "正常"
	case 9:
		return "结束"
	}

	return "非正常状态"
}

func cui(a, b int64) string {
	if cmp.Less(a, b/2) {
		return "催"
	}
	return ""
}

func tian(a, b int64) int64 {
	if b == 0 {
		return 0
	}
	return a / (b / 30)
}

func since12h(t time.Time) bool {
	return time.Since(t) > 12*time.Hour
}
