package wecom

import (
	"text/template"
)

const help_cmd_tpl = `
支持下列命令
1. 使用?或者h命令显示帮助
2. p <铺位号>，列出对应铺位中表具
  例如：p F1_01
3. t <表号>，查询对应仪表的当前用量和充值记录
  例如：t 909_21
`

const top_cmd_tpl = `
# 表号：{{ .Code }}
>当前表显 **{{.Curr.DataValue | du}}** 度，采集时间**{{.Curr.DataTime | datetime}}**
{{ if .Tops }}
## 充值记录
| 时间 | 度数 | 限额 | 金额 | 铺位号 |
| :-----: | :----: | :----: | :-----: | :-----: |
{{ range .Tops -}}
| {{.ChargeTime | datetime}} | **{{.Incr | du}}** | **{{.Top | du}}** | {{.Amount | yuan}} | {{.PosCode}} |
{{ end }}
>剩余约**{{ sub .Last.Top .Curr.DataValue | du}}** 度
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

func build() *template.Template {
	t := template.New("cmd").Funcs(
		template.FuncMap{
			"du":       du,
			"yuan":     du,
			"datetime": fmtTime,
			"sub":      sub,
			"add":      add,
		},
	)
	template.Must(t.New("help").Parse(help_cmd_tpl))
	template.Must(t.New("top").Parse(top_cmd_tpl))
	template.Must(t.New("charge").Parse(charge_cmd_tpl))
	return t
}
