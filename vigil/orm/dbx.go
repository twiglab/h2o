package orm

import (
	"context"

	"github.com/twiglab/h2o/dbcli/ent"
	"github.com/twiglab/h2o/vigil"
)

type DBx struct {
	Client *ent.Client
}

func (d *DBx) TabbElecty(ctx context.Context, data vigil.ElectricityMeter) error {
	cr := d.Client.NhRecord.Create()
	cr.SetDeviceSn(data.SN)
	cr.SetDeviceCode(data.Code)
	cr.SetDeviceType(data.Type)
	cr.SetDataCode(data.DataCode)
	cr.SetDataTime(data.DataTime)
	cr.SetDataValue(data.Data.DataValue)
	cr.SetPosCode(data.Pos.PosCode)
	cr.SetProject(data.Pos.Project)
	cr.SetDataTs(data.DataTs)
	cr.SetOwner(data.Pos.Owner)
	return cr.Exec(ctx)
}

func (d *DBx) TabbWater(ctx context.Context, data vigil.WaterMeter) error {
	cr := d.Client.NhRecord.Create()
	cr.SetDeviceSn(data.SN)
	cr.SetDeviceCode(data.Code)
	cr.SetDeviceType(data.Type)
	cr.SetDataCode(data.DataCode)
	cr.SetDataTime(data.DataTime)
	cr.SetDataValue(data.Data.DataValue)
	cr.SetPosCode(data.Pos.PosCode)
	cr.SetProject(data.Pos.Project)
	cr.SetDataTs(data.DataTs)
	cr.SetOwner(data.Pos.Owner)
	return cr.Exec(ctx)
}
