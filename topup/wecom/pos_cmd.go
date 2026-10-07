package wecom

import (
	"context"
	"strings"
	"text/template"

	"github.com/twiglab/h2o/dbcli/ent"
	"github.com/twiglab/h2o/dbcli/ent/device"
)

type PosQueryCmd struct {
	PosCode string

	Devices []*ent.Device

	cli  *ent.Client
	tmpl *template.Template
}

func (c *PosQueryCmd) Args(args ...string) error {
	if len(args) > 1 {
		c.PosCode = strings.ToUpper(args[1])
	}
	return nil
}

func (c *PosQueryCmd) Do(ctx context.Context) (err error) {
	q := c.cli.Device.Query().
		Where(device.IsDelEQ(0), device.PosCodeNotNil()).
		Order(ent.Asc(device.FieldPosCode))
		//Limit(c.Limit)

	if c.PosCode != "" {
		q.Where(device.PosCodeEQ(c.PosCode))
	}

	c.Devices, err = q.All(ctx)
	return
}

func (c PosQueryCmd) ToString() string {
	var sb strings.Builder
	_ = c.tmpl.ExecuteTemplate(&sb, template_pos, c)
	return sb.String()

}

func posQueryCmdFn(cfg Global, args ...string) (Commander, error) {
	u := &PosQueryCmd{cli: cfg.Client, tmpl: cfg.Template}
	return u, u.Args(args...)
}
