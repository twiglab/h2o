package wecom

import (
	"context"
	"errors"
	"strings"
	"text/template"

	"github.com/twiglab/h2o/dbcli/ent"
	"github.com/twiglab/h2o/dbcli/ent/nhrecord"
	"github.com/twiglab/h2o/dbcli/ent/top"
)

type topCmd struct {
	cli  *ent.Client
	tmpl *template.Template
	dbx  *DBx

	Code string

	Tops []*ent.Top
	Curr *ent.NhRecord
	Last *ent.Top

	Device *ent.Device
}

func (c *topCmd) Args(args ...string) error {
	c.Code = args[1]
	return nil
}

func (c *topCmd) Do(ctx context.Context) (err error) {
	c.Device, err = c.dbx.LoadDeviceByCode(ctx, c.Code)
	if err != nil {
		err = errors.New("表号不存在")
		return
	}

	c.Curr, err = c.cli.NhRecord.Query().
		Where(nhrecord.DeviceCodeEQ(c.Code)).
		Order(ent.Desc(nhrecord.FieldDataTime)).
		First(ctx)

	if err != nil {
		err = errors.New("无采集记录")
		return
	}

	c.Tops, err = c.cli.Top.Query().
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
	if err := c.tmpl.ExecuteTemplate(&sb, template_top, c); err != nil {
		return err.Error()
	}
	return sb.String()
}

func topCmdFn(cfg Global, args ...string) (Commander, error) {
	u := &topCmd{
		cli:  cfg.Client,
		tmpl: cfg.Template,
		dbx:  &DBx{cli: cfg.Client},
	}

	err := u.Args(args...)
	return u, err
}
