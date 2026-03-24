package vigil

import (
	"encoding/json/v2"
	"time"

	"github.com/twiglab/h2o/pkg/common"
)

type Meter struct {
	common.Device
	Pos common.Pos `json:"pos,omitzero"`
}

type ElectricityMeter struct {
	Meter
	Data common.Electricity `json:"data,omitzero"`
}

func (d *ElectricityMeter) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, d)
}

func (d ElectricityMeter) GetSN() string {
	return d.SN
}

func (d ElectricityMeter) GetCode() string {
	return d.Code
}

func (d ElectricityMeter) GetType() string {
	return d.Type
}

func (d ElectricityMeter) GetName() string {
	return d.Name
}

func (d ElectricityMeter) GetDataCode() string {
	return d.DataCode
}

func (d ElectricityMeter) GetDataTime() time.Time {
	return d.DataTime
}

func (d ElectricityMeter) GetDataValue() int64 {
	return d.Data.DataValue
}

func (d ElectricityMeter) GetPosCode() string {
	return d.Pos.PosCode
}

func (d ElectricityMeter) GetProject() string {
	return d.Pos.Project
}
