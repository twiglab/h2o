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
	"github.com/twiglab/h2o/dbcli/ent/device"
	"github.com/twiglab/h2o/dbcli/ent/nhrecord"
	"github.com/twiglab/h2o/dbcli/ent/top"
)

type chargeCmd struct {
	Code string

	Incr  int64
	Stock int64

	TopVal int64

	Device *ent.Device
	Nh     *ent.NhRecord
	Top    *ent.Top

	cli   *ent.Client
	templ *template.Template
}

func (c *chargeCmd) Args(args ...string) error {
	c.Code = args[1]
	i, err := strconv.ParseFloat(args[2], 64)
	if err != nil {
		return err
	}
	c.Incr = int64(math.Ceil(i * 100))
	return nil
}

func (x *chargeCmd) Do(ctx context.Context) (err error) {
	x.Device, err = x.cli.Device.Query().
		Where(device.IsDelEQ(0)).
		Where(device.DeviceCodeEQ(x.Code)).
		Only(ctx)

	if err != nil {
		return
	}

	x.Nh, err = x.cli.NhRecord.Query().
		Where(nhrecord.DeviceCodeEQ(x.Code)).
		Order(ent.Desc(nhrecord.FieldDataTime)).
		First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			x.Nh = &ent.NhRecord{
				// DataValue:  0,
				DataTime: time.Now(),
			}
		}
		return
	}

	if time.Since(x.Nh.DataTime) > 12*time.Hour {
		return errors.New("最后一次采集大于12小时")
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

	if x.Stock == 0 {
		x.Stock = max(x.Top.Top, x.Top.EndStock, x.Nh.DataValue)
	}

	x.TopVal = x.Stock + x.Incr

	cr := x.cli.Top.Create()
	cr.SetDeviceCode(x.Device.DeviceCode)
	cr.SetDeviceType(x.Device.DeviceType)
	cr.SetPosCode(x.Device.PosCode)
	cr.SetProject(x.Device.Project)

	cr.SetTop(x.TopVal)
	cr.SetStock(x.Stock)
	cr.SetIncr(x.Incr)

	return cr.Exec(ctx)
}

func (x chargeCmd) ToString() string {
	var sb strings.Builder
	_ = x.templ.ExecuteTemplate(&sb, "charge", x)
	return sb.String()
}

func chargeCmdFn(cfg CmdCfg, args ...string) (Commander, error) {
	c := &chargeCmd{cli: cfg.Cli, templ: build()}
	c.Args(args...)
	return c, nil
}
