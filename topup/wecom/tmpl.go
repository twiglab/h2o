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

/*
const top_cmd_tpl = `
# 表号：{{ .Code }}
>当前表显 **{{.Curr.DataValue | du}}** 度，采集时间**{{.Curr.DataTime | datetime}}**
## 充值记录
| 时间 | 度数 | 限额 | 金额 | 铺位号 |
| :-----: | :----: | :----: | :-----: | :-----: |
{{ range .Tops -}}
| {{.ChargeTime | datetime}} | **{{.Incr | du}}** | **{{.Top | du}}** | {{.Amount | du}} | {{.PosCode}} |
{{ end }}
`
*/

const top_cmd_tpl = `
# 表号：{{ .Code }}
>当前表显 **{{.Curr.DataValue | du}}** 度，采集时间**{{.Curr.DataTime | datetime}}**
{{ if .Tops }}
## 充值记录
| 时间 | 度数 | 限额 | 金额 | 铺位号 |
| :-----: | :----: | :----: | :-----: | :-----: |
{{ range .Tops -}}
| {{.ChargeTime | datetime}} | **{{.Incr | du}}** | **{{.Top | du}}** | {{.Amount | du}} | {{.PosCode}} |
{{ end }}
>剩余约**{{ sub .Last.Top .Curr.DataValue | du}}** 度
{{ end }}
`

// >剩余{{ . }} 度

func build() *template.Template {
	t := template.New("cmd").Funcs(
		template.FuncMap{
			"du":       du,
			"datetime": fmtTime,
			"sub":      sub,
		},
	)
	template.Must(t.New("help").Parse(help_cmd_tpl))
	template.Must(t.New("top").Parse(top_cmd_tpl))
	return t
}
