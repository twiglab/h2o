package wecom

import (
	"context"
	"errors"
	"math/big"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/twiglab/h2o/dbcli/ent"
	"github.com/twiglab/h2o/dbcli/ent/device"
	"github.com/twiglab/h2o/dbcli/ent/nhrecord"
	"github.com/twiglab/h2o/dbcli/ent/top"
)

type moneyCmd struct {
	Code string

	Incr   int64
	Stock  int64 // 输入的存量
	Amount int64

	TopVal int64

	Device   *ent.Device
	Nh       *ent.NhRecord
	Top      *ent.Top
	TopAfter *ent.Top

	cli  *ent.Client
	tmpl *template.Template
}

func (c *moneyCmd) Args(args ...string) error {
	c.Code = args[1]

	// 金额
	amount, err := strconv.ParseFloat(args[2], 64)
	if err != nil {
		return err
	}
	c.Amount = int64(amount * 100) // 截断，钱支持到分

	return nil
}

func (x *moneyCmd) Do(ctx context.Context) (err error) {
	x.Device, err = x.cli.Device.Query().
		Where(device.IsDelEQ(0)).
		Where(device.DeviceCodeEQ(x.Code)).
		Only(ctx)

	if err != nil {
		err = errors.New("表号不存在")
		return
	}

	x.Nh, err = x.cli.NhRecord.Query().
		Where(nhrecord.DeviceCodeEQ(x.Code)).
		Order(ent.Desc(nhrecord.FieldDataTime)).
		First(ctx)
	if err != nil {
		err = errors.New("没有采集记录")
		return
	}

	if time.Since(x.Nh.DataTime) > 12*time.Hour {
		err = errors.New("最后一次采集大于12小时，表具可能存在错误")
		return
	}

	x.Top, err = x.cli.Top.Query().
		Where(top.DeviceCodeEQ(x.Code)).
		Where(top.IsDelEQ(0)).
		Order(ent.Desc(top.FieldChargeTime)).
		First(ctx)
	if err != nil {
		if !ent.IsNotFound(err) {
			return
		}
		x.Top = &ent.Top{}
	}

	if time.Since(x.Top.ChargeTime) < 5*time.Minute {
		err = errors.New("充值间隔必须大于5分钟")
		return
	}

	x.Stock = x.Top.Top

	a := big.NewInt(x.Amount)
	r := big.NewInt(x.Device.Rate)
	u := big.NewInt(x.Device.UnitPrice)

	z := a.Div(a, r)
	z = z.Div(z, u)

	x.Incr = z.Int64() + 1 // 多给 0.01
	x.TopVal = x.Stock + x.Incr

	cr := x.cli.Top.Create()
	cr.SetDeviceCode(x.Device.DeviceCode)
	cr.SetDeviceType(x.Device.DeviceType)
	cr.SetPosCode(x.Device.PosCode)
	cr.SetProject(x.Device.Project)

	cr.SetTop(x.TopVal)
	cr.SetStock(x.Stock)
	cr.SetIncr(x.Incr)

	cr.SetAmount(x.Amount)

	x.TopAfter, err = cr.Save(ctx)
	return
}

func (x moneyCmd) ToString() string {
	var sb strings.Builder
	_ = x.tmpl.ExecuteTemplate(&sb, template_charge, x)
	return sb.String()
}

func moneyCmdFn(cfg Global, args ...string) (Commander, error) {
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

	c := &moneyCmd{
		cli:  cfg.Client,
		tmpl: cfg.Template,
	}
	return c, c.Args(args...)
}
