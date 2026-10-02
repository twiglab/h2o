package wecom

import (
	"context"
	"fmt"
	"strings"
	"text/template"

	"github.com/twiglab/h2o/dbcli/ent"
	"github.com/twiglab/h2o/dbcli/ent/nhrecord"
	"github.com/twiglab/h2o/dbcli/ent/top"
)

type topCmd struct {
	Cli *ent.Client

	Code string

	Tops []*ent.Top
	Curr *ent.NhRecord
	Last *ent.Top

	tmpl *template.Template
}

func (c *topCmd) Args(args ...string) error {
	c.Code = args[1]
	return nil
}

func (c *topCmd) Do(ctx context.Context) (err error) {
	c.Curr, err = c.Cli.NhRecord.Query().
		Where(nhrecord.DeviceCodeEQ(c.Code)).
		Order(ent.Desc(nhrecord.FieldDataTime)).
		First(ctx)

	if ent.IsNotFound(err) {
		return ErrorCmd{str: fmt.Sprintf("%s 无采集信息", c.Code)}
	}

	if err != nil {
		return err
	}

	c.Tops, err = c.Cli.Top.Query().
		Where(top.DeviceCodeEQ(c.Code)).
		Order(ent.Desc(top.FieldChargeTime)).
		Limit(10).
		All(ctx)

	if err != nil {
		return
	}

	if len(c.Tops) > 0 {
		c.Last = c.Tops[0]
	}

	return
}

func (c topCmd) ToString() string {
	var sb strings.Builder
	if err := c.tmpl.ExecuteTemplate(&sb, "top", c); err != nil {
		return err.Error()
	}
	return sb.String()
}

func topCmdFn(cfg CmdCfg, args ...string) (Commander, error) {
	u := &topCmd{Cli: cfg.Cli, tmpl: build()}
	err := u.Args(args...)
	return u, err
}
