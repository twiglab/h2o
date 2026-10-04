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
	Type    string

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
		Order(ent.Asc(device.FieldDeviceType))
		// Limit(10).

	if c.PosCode != "" {
		q.Where(device.PosCodeEQ(c.PosCode))
	}

	if c.Type != "" {
		q.Where(device.DeviceTypeEQ(c.Type))
	}

	c.Devices, err = q.All(ctx)
	return
}

func (c PosQueryCmd) ToString() string {
	var sb strings.Builder
	_ = c.tmpl.ExecuteTemplate(&sb, "pos", c)
	return sb.String()

}

func posQueryCmdFn(cfg CmdCfg, args ...string) (Commander, error) {
	u := &PosQueryCmd{cli: cfg.Cli}
	err := u.Args(args...)
	return u, err
}
