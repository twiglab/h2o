package orm

import (
	"context"

	"github.com/twiglab/h2o/vigil"
	"github.com/twiglab/h2o/vigil/orm/ent"
)

type DBx struct {
	Client *ent.Client
}

func (d *DBx) TabbElecty(ctx context.Context, data vigil.ElectricityMeter) error {
	cr := d.Client.Electy.Create()
	cr.SetPCode(data.PCode())
	cr.SetDeviceSn(data.SN)
	cr.SetDeviceCode(data.Code)
	cr.SetDeviceType(data.Type)
	cr.SetDeviceName(data.Name)
	cr.SetDataCode(data.DataCode)
	cr.SetDataTime(data.DataTime)
	cr.SetDataValue(data.Data.DataValue)
	cr.SetPosCode(data.Pos.PosCode)
	cr.SetProject(data.Pos.Project)
	cr.SetDataTs(data.Ts())
	cr.SetXDataValue(data.XDataValue())
	cr.SetFactor(data.Param.Factor)
	return cr.Exec(ctx)
}

func (d *DBx) TabbWater(ctx context.Context, data vigil.WaterMeter) error {
	cr := d.Client.Water.Create()
	cr.SetPCode(data.PCode())
	cr.SetDeviceSn(data.SN)
	cr.SetDeviceCode(data.Code)
	cr.SetDeviceType(data.Type)
	cr.SetDeviceName(data.Name)
	cr.SetDataCode(data.DataCode)
	cr.SetDataTime(data.DataTime)
	cr.SetDataValue(data.Data.DataValue)
	cr.SetPosCode(data.Pos.PosCode)
	cr.SetProject(data.Pos.Project)
	cr.SetDataTs(data.Ts())
	return cr.Exec(ctx)
}
