package orm

import (
	"context"

	"github.com/twiglab/h2o/vigil"
	"github.com/twiglab/h2o/vigil/orm/ent"
)

type DBx struct {
	Client *ent.Client
}

func (d *DBx) Tabb(ctx context.Context, data vigil.MeterData) error {
	cr := d.Client.Record.Create()
	cr.SetDeviceSn(data.GetSN())
	cr.SetDeviceCode(data.GetCode())
	cr.SetDeviceType(data.GetType())
	cr.SetDeviceName(data.GetName())
	cr.SetDataCode(data.GetDataCode())
	cr.SetDataTime(data.GetDataTime())
	cr.SetDataValue(data.GetDataValue())
	cr.SetPosCode(data.GetPosCode())
	cr.SetProject(data.GetProject())
	return cr.Exec(ctx)
}
