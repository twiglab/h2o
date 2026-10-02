package wecom

import (
	"context"
	"fmt"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/twiglab/h2o/dbcli/ent"
	"github.com/twiglab/h2o/dbcli/ent/nhrecord"
)

type UsageCmd struct {
	Cli *ent.Client

	Code string

	all []*ent.NhRecord
}

func (c *UsageCmd) Args(args ...string) error {
	c.Code = args[1]
	return nil
}

func (c *UsageCmd) Do(ctx context.Context) (err error) {
	c.all, err = c.Cli.NhRecord.Query().
		Where(nhrecord.DeviceCodeEQ(c.Code)).
		Order(ent.Desc(nhrecord.FieldDataTime)).
		Limit(10).
		All(ctx)

	return
}

func (c UsageCmd) ToString() string {
	var sb strings.Builder
	table := tablewriter.NewTable(&sb,
		tablewriter.WithRenderer(renderer.NewMarkdown()),
	)

	table.Header([]string{"编号", "度数", "时间"})

	for _, r := range c.all {
		table.Append([]string{r.DeviceCode, fmt.Sprint(r.DataValue), r.DataTs})
	}
	table.Render()
	return sb.String()
}

func usageCmdFn(cfg CmdCfg, args ...string) (Commander, error) {
	u := &UsageCmd{Cli: cfg.Cli}
	err := u.Args(args...)
	return u, err
}
