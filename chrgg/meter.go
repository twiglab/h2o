package chrgg

import (
	"encoding/json"

	"github.com/twiglab/h2o/pkg/common"
)

type ElectyMeterData struct {
	common.Device
	Data  common.MeterValue `json:"data"`
	Topic string            `json:"topic"`
}

func (d *ElectyMeterData) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, d)
}
