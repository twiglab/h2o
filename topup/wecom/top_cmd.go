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

type topCmd struct {
	Code string

	Incr   int64
	Stock  int64 // 输入的存量
	Amount int64 // 0, 仅被模板使用

	TopVal int64

	Device   *ent.Device
	Nh       *ent.NhRecord
	Top      *ent.Top
	TopAfter *ent.Top

	tmpl *template.Template
	dbx  *DBx
}

func (c *topCmd) Args(args ...string) error {
	if len(args) < 3 {
		return errors.New("t 命令必须提供表号和增量")
	}

	// 表号
	c.Code = args[1]

	// 增量
	incr, err := strconv.ParseFloat(args[2], 64)
	if err != nil {
		return err
	}
	c.Incr = int64(math.Ceil(incr * 100))

	if len(args) > 3 {
		// 起点
		stock, err := strconv.ParseFloat(args[3], 64)
		if err != nil {
			return err
		}
		c.Stock = int64(stock * 100) // 截断，存量支持到0.01
	}

	return nil
}

func (x *topCmd) Do(ctx context.Context) (err error) {
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

	// 没有输入存量，存量就是上一次的top
	if x.Stock == 0 {
		x.Stock = x.Top.Top
	} else if x.Stock < 0 {
		// 输入是负数，取当前表显示
		x.Stock = x.Nh.DataValue
	}

	// 当前限额 永远等于存量+incr
	x.TopVal = x.Stock + x.Incr

	x.TopAfter, err = x.dbx.Save(ctx, TopRecord{
		DeviceCode: x.Device.DeviceCode,
		DeviceType: x.Device.DeviceType,
		PosCode:    x.Device.PosCode,
		Project:    x.Device.Project,

		Top:   x.TopVal,
		Stock: x.Stock,
		Incr:  x.Incr,

		Amount:    0,
		UnitPrice: 0,
	})

	return
}

func (x topCmd) ToString() string {
	var sb strings.Builder
	_ = x.tmpl.ExecuteTemplate(&sb, template_charge, x)
	return sb.String()
}

func topCmdFn(cfg Global, args ...string) (Commander, error) {
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

	c := &topCmd{
		tmpl: cfg.Template,
		dbx:  &DBx{cli: cfg.Client},
	}
	return c, c.Args(args...)
}
