package wecom

import (
	"context"
	"errors"
	"strings"
	"text/template"
)

type helpCmd struct {
	Template *template.Template
}

func (h helpCmd) Do(ctx context.Context) error {
	return nil
}

func (h helpCmd) ToString() string {
	var sb strings.Builder
	_ = h.Template.ExecuteTemplate(&sb, template_help, nil)
	return sb.String()
}

func helpCmdFn(cfg Global, args ...string) (Commander, error) {
	help := helpCmd{Template: cfg.Template}
	return help, nil
}

func chrggCmdFn(cfg Global, args ...string) (Commander, error) {
	return nil, errors.New("c 命令已被删除，请用 m 命令充值，或用新的 t 命令设置限额，请用 h 查看帮助")
}

func repCmdFn(cfg Global, args ...string) (Commander, error) {
	return nil, errors.New("report 命令已经被删除，请用 r 命令查看用量报表")
}
