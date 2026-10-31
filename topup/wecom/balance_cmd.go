package wecom

import (
	"context"
	"errors"
	"strings"
	"text/template"

	"github.com/twiglab/h2o/dbcli/ent"
)

type balanceCmd struct {
	PosCode1 string
	PosCode2 string

	cli  *ent.Client
	tmpl *template.Template
}

func (c *balanceCmd) Args(args ...string) error {
	if len(args) > 1 {
		c.PosCode1 = strings.ToUpper(args[1])
	}
	return nil
}

func (c *balanceCmd) Do(ctx context.Context) (err error) {
	return
}
func (x balanceCmd) ToString() string {
	var sb strings.Builder
	_ = x.tmpl.ExecuteTemplate(&sb, template_charge, x)
	return sb.String()
}

func balanceCmdFn(cfg Global, args ...string) (Commander, error) {
	ok, err := cfg.Auth.Auth(MessageToken{
		ChatID:   cfg.TextMessage.ChatID,
		UserID:   cfg.TextMessage.From.UserID,
		ChatType: cfg.TextMessage.ChatType,
	})
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("当前命令只允许特定人员在特定的群组执行")
	}

	c := &balanceCmd{
		cli:  cfg.Client,
		tmpl: cfg.Template,
	}
	return c, c.Args(args...)
}
