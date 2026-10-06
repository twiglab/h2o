package wecom

import (
	"context"
	"strings"
	"text/template"

	"github.com/twiglab/h2o/dbcli/ent"
	"github.com/twiglab/h2o/dbcli/ent/nhrecord"
)

type UsageCmd struct {
	cli   *ent.Client
	templ *template.Template

	Code string

	All []*ent.NhRecord
}

func (c *UsageCmd) Args(args ...string) error {
	c.Code = args[1]
	return nil
}

func (c *UsageCmd) Do(ctx context.Context) (err error) {
	c.All, err = c.cli.NhRecord.Query().
		Where(nhrecord.DeviceCodeEQ(c.Code)).
		Order(ent.Desc(nhrecord.FieldDataTime)).
		Limit(10).
		All(ctx)

	return
}

func (c UsageCmd) ToString() string {
	var sb strings.Builder
	_ = c.templ.ExecuteTemplate(&sb, template_usage, c)
	return sb.String()
}

func usageCmdFn(cfg Global, args ...string) (Commander, error) {
	u := &UsageCmd{cli: cfg.Client, templ: cfg.Template}
	err := u.Args(args...)
	return u, err
}
