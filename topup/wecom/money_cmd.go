package wecom

import (
	"context"
	"errors"
	"math"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/twiglab/h2o/dbcli/ent"
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
	dbx  *DBx
}

func (c *moneyCmd) Args(args ...string) error {
	c.Code = args[1]

	// 金额
	amount, err := strconv.ParseFloat(args[2], 64)
	if err != nil {
		return err
	}
	c.Amount = int64(amount * 100) // 截断，钱支持到分

	if len(args) >= 4 {
		// 起点
		stock, err := strconv.ParseFloat(args[3], 64)
		if err != nil {
			return err
		}
		c.Stock = int64(stock * 100) // 截断，存量支持到 0.01
	}

	return nil
}

func (x *moneyCmd) Do(ctx context.Context) (err error) {
	x.Device, err = x.dbx.LoadDevice(ctx, x.Code)

	if err != nil {
		err = errors.New("表号不存在")
		return
	}

	x.Nh, err = x.dbx.LoadRecordLast(ctx, x.Code)
	if err != nil {
		err = errors.New("没有采集记录")
		return
	}

	if time.Since(x.Nh.DataTime) > 12*time.Hour {
		err = errors.New("最后一次采集大于12小时，表具可能存在错误")
		return
	}

	x.Top, err = x.dbx.LoadTopLast(ctx, x.Code)
	if ent.IsNotFound(err) {
		x.Top = &ent.Top{}
	}
	if ent.MaskNotFound(err) != nil {
		return err
	}

	if time.Since(x.Top.ChargeTime) < 5*time.Minute {
		err = errors.New("充值间隔必须大于5分钟")
		return
	}

	if x.Device.Rate == 0 || x.Device.UnitPrice == 0 {
		err = errors.New("表具倍率或单价未配置")
		return
	}

	// 没有输入存量，存量就是上一次的top
	if x.Stock == 0 {
		x.Stock = x.Top.Top
	}

	x.Incr = calcIncr(x.Amount, x.Device.Rate, x.Device.UnitPrice)
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
		dbx:  &DBx{cli: cfg.Client},
	}
	return c, c.Args(args...)
}

func calcIncr(a, r, u int64) int64 {
	//var a, r, u int64 = 10000, 1, 121
	if r == 0 || u == 0 {
		return 0
	}
	return int64(math.Ceil(float64(a) / float64(r) / float64(u) * 100))
}
