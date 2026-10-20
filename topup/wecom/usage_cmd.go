package wecom

import (
	"context"
	"errors"
	"slices"
	"strings"
	"text/template"

	"github.com/twiglab/h2o/dbcli/ent"
)

type UsageCmd struct {
	cli   *ent.Client
	templ *template.Template
	dbx   *DBx

	Code string

	Records []*ent.NhRecord
	Record  *ent.NhRecord

	Tops []*ent.Top
	Top  *ent.Top

	Device *ent.Device
}

func (c *UsageCmd) Args(args ...string) error {
	c.Code = args[1]
	return nil
}

func (c *UsageCmd) Do(ctx context.Context) (err error) {
	c.Device, err = c.dbx.LoadDevice(ctx, c.Code)
	if err != nil {
		err = errors.New("表号不存在")
		return
	}

	c.Records, err = c.dbx.LoadRecordLimit(ctx, c.Code, 3)
	if err != nil {
		return
	}

	if len(c.Records) < 1 {
		err = errors.New("没有采集数据，表具可能存在问题")
		return
	}

	c.Record = slices.MaxFunc(c.Records, func(a, b *ent.NhRecord) int {
		return a.DataTime.Compare(b.DataTime)
	})

	c.Tops, err = c.dbx.LoadTopLimit(ctx, c.Code, 3)
	if err == nil && len(c.Tops) > 0 {
		c.Top = slices.MaxFunc(c.Tops, func(a, b *ent.Top) int {
			return a.ChargeTime.Compare(b.ChargeTime)
		})
	}
	return
}

func (c UsageCmd) ToString() string {
	var sb strings.Builder
	_ = c.templ.ExecuteTemplate(&sb, template_usage, c)
	return sb.String()
}

func usageCmdFn(cfg Global, args ...string) (Commander, error) {
	u := &UsageCmd{
		cli:   cfg.Client,
		templ: cfg.Template,
		dbx:   &DBx{cli: cfg.Client},
	}
	err := u.Args(args...)
	return u, err
}
