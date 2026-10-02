package wecom

import (
	"context"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/twiglab/h2o/dbcli/ent"
	"github.com/twiglab/h2o/dbcli/ent/device"
)

type PosQueryCmd struct {
	Cli *ent.Client

	PosCode string
	Type    string

	all []*ent.Device
}

func (c *PosQueryCmd) Args(args ...string) error {
	if len(args) > 1 {
		c.PosCode = strings.ToUpper(args[1])
	}
	return nil
}

func (c *PosQueryCmd) Do(ctx context.Context) (err error) {
	q := c.Cli.Device.Query().
		Where(device.IsDelEQ(0), device.PosCodeNotNil()).
		Order(ent.Asc(device.FieldDeviceType))
		// Limit(10).

	if c.PosCode != "" {
		q.Where(device.PosCodeEQ(c.PosCode))
	}

	if c.Type != "" {
		q.Where(device.DeviceTypeEQ(c.Type))
	}

	c.all, err = q.All(ctx)
	return
}

func (c PosQueryCmd) ToString() string {
	var sb strings.Builder
	table := tablewriter.NewTable(&sb,
		tablewriter.WithRenderer(renderer.NewMarkdown()),
	)

	table.Header([]string{"铺位", "编号", "类型"})

	for _, r := range c.all {
		table.Append([]string{r.PosCode, r.DeviceCode, r.DeviceType})
	}
	table.Render()
	return sb.String()
}

func posQueryCmdFn(cfg CmdCfg, args ...string) (Commander, error) {
	u := &PosQueryCmd{Cli: cfg.Cli}
	err := u.Args(args...)
	return u, err
}
