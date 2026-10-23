package wecom

import (
	"context"

	"github.com/twiglab/h2o/dbcli/ent"
	"github.com/twiglab/h2o/dbcli/ent/device"
	"github.com/twiglab/h2o/dbcli/ent/nhrecord"
	"github.com/twiglab/h2o/dbcli/ent/top"
)

type DBx struct {
	cli *ent.Client
}

func (d *DBx) LoadDevice(ctx context.Context, code string) (*ent.Device, error) {
	dev, err := d.cli.Device.Query().
		Where(device.IsDelEQ(0)).
		Where(device.DeviceCodeEQ(code)).
		Only(ctx)

	return dev, err
}

func (d *DBx) LoadRecordLimit(ctx context.Context, code string, limit int) ([]*ent.NhRecord, error) {
	all, err := d.cli.NhRecord.Query().
		Where(nhrecord.DeviceCodeEQ(code)).
		Order(ent.Desc(nhrecord.FieldDataTime)).
		Limit(limit).
		All(ctx)
	return all, err
}

func (d *DBx) LoadRecordLast(ctx context.Context, code string) (*ent.NhRecord, error) {
	a, err := d.cli.NhRecord.Query().
		Where(nhrecord.DeviceCodeEQ(code)).
		Order(ent.Desc(nhrecord.FieldDataTime)).
		First(ctx)
	return a, err
}

func (d *DBx) LoadTopLimit(ctx context.Context, code string, limit int) ([]*ent.Top, error) {
	a, err := d.cli.Top.Query().
		Where(top.DeviceCodeEQ(code), top.IsDelEQ(0)).
		Order(ent.Desc(top.FieldChargeTime)).
		Limit(limit).
		All(ctx)
	return a, err
}

func (d *DBx) LoadTopLast(ctx context.Context, code string) (*ent.Top, error) {
	a, err := d.cli.Top.Query().
		Where(top.DeviceCodeEQ(code), top.IsDelEQ(0)).
		Order(ent.Desc(top.FieldChargeTime)).
		First(ctx)
	return a, err
}
