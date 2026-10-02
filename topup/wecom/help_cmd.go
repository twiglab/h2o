package wecom

import (
	"context"
	"strings"
	"text/template"
)

var hCmd = h()

func h() helpCmd {
	t, _ := template.New("x").Parse(help_cmd_tpl)
	return helpCmd{
		help: t,
	}
}

type helpCmd struct {
	help *template.Template
}

func (x helpCmd) Do(ctx context.Context) error {
	return nil
}

func (x helpCmd) ToString() string {
	var sb strings.Builder
	_ = x.help.Execute(&sb, nil)
	return sb.String()
}

func helpCmdFn(cfg CmdCfg, args ...string) (Commander, error) {
	return hCmd, nil
}
