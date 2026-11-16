package wecom

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"text/template"

	"github.com/twiglab/h2o/dbcli/ent"
)

type balanceCmd struct {
	Code1 string
	Code2 string
	Value int64

	Device1 *ent.Device
	Device2 *ent.Device

	Top1 *ent.Top
	Top2 *ent.Top

	tmpl *template.Template
	dbx  *DBx
}

func (c *balanceCmd) Args(args ...string) error {
	if len(args) < 4 {
		return errors.New("b 命令参数为：表号1 表号2 电量")
	}

	c.Code1 = strings.ToUpper(args[1])
	c.Code2 = strings.ToUpper(args[2])

	x, err := strconv.ParseFloat(args[3], 64)
	if err != nil {
		return err
	}
	c.Value = int64(x * 100)

	return nil
}

func (c *balanceCmd) Do(ctx context.Context) (err error) {
	if c.Device1, err = c.dbx.LoadDevice(ctx, c.Code1); err != nil {
		return
	}
	if c.Device2, err = c.dbx.LoadDevice(ctx, c.Code2); err != nil {
		return
	}

	if c.Device1.PosCode != c.Device2.PosCode {
		err = ErrPosCodeNotSame
		return
	}

	if c.Top1, err = c.dbx.LoadTopLast(ctx, c.Code1); err != nil {
		return err
	}
	if c.Top2, err = c.dbx.LoadTopLast(ctx, c.Code2); err != nil {
		return err
	}

	top1 := c.Top1.Top - (c.Value / c.Device1.Rate)
	top2 := c.Top2.Top + (c.Value / c.Device2.Rate)
	return c.dbx.WithTx(ctx, xx(ctx, c.Device1, c.Top1, top1), xx(ctx, c.Device2, c.Top2, top2))
}

func (x balanceCmd) ToString() string {
	var sb strings.Builder
	_ = x.tmpl.ExecuteTemplate(&sb, template_balance, x)
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
		dbx:  &DBx{cli: cfg.Client},
		tmpl: cfg.Template,
	}
	return c, c.Args(args...)
}

func xx(ctx context.Context, d *ent.Device, _ *ent.Top, top int64) func(*ent.Tx) error {
	return func(tx *ent.Tx) error {
		cr := tx.Top.Create()

		cr.SetDeviceCode(d.DeviceCode)
		cr.SetDeviceType(d.DeviceType)
		cr.SetPosCode(d.PosCode)
		cr.SetProject(d.Project)

		cr.SetTop(top)
		cr.SetStatus(0)

		return cr.Exec(ctx)
	}
}
