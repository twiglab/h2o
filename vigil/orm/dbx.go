package orm

import (
	"context"

	"github.com/twiglab/h2o/vigil"
	"github.com/twiglab/h2o/vigil/orm/ent"
)

type DBx struct {
	Client *ent.Client
}

func (d *DBx) Tabb(ctx context.Context, data vigil.ElectricityMeter) error {
	cr := d.Client.Record.Create()
	cr.SetDataCode(data.Code)
	cr.SetDeviceType(data.Type)
	cr.SetDataCode(data.DataCode)
	cr.SetDataTime(data.DataTime)
	cr.SetDataValue(data.Data.DataValue)
	cr.SetPosCode(data.Pos.PosCode)
	cr.SetProject(data.Pos.Project)
	return cr.Exec(ctx)
}
