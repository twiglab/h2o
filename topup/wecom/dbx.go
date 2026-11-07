package wecom

import (
	"context"

	"github.com/twiglab/h2o/dbcli/ent"
	"github.com/twiglab/h2o/dbcli/ent/device"
	"github.com/twiglab/h2o/dbcli/ent/nhrecord"
	"github.com/twiglab/h2o/dbcli/ent/top"
)

type TopRecord struct {
	DeviceCode string
	DeviceType string
	PosCode    string
	Project    string

	Top   int64
	Stock int64
	Incr  int64

	Amount    int64
	UnitPrice int64
}

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

func (d *DBx) Save(ctx context.Context, record TopRecord) (*ent.Top, error) {
	cr := d.cli.Top.Create()

	cr.SetDeviceCode(record.DeviceCode)
	cr.SetDeviceType(record.DeviceType)
	cr.SetPosCode(record.PosCode)
	cr.SetProject(record.Project)

	cr.SetTop(record.Top)
	cr.SetStock(record.Stock)
	cr.SetIncr(record.Incr)

	cr.SetAmount(record.Amount)
	cr.SetUnitPrice(record.UnitPrice)

	cr.SetStatus(0)

	return cr.Save(ctx)
}
