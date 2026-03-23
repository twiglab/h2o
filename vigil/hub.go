package vigil

import (
	"context"
)

type Recorder interface {
	Tabb(ctx context.Context, data ElectricityMeter) error
}

type ElectyMeterView interface {
	Merge(data ElectricityMeter)
}

type Hub struct {
	ElectyMeterView ElectyMeterView
	Recorder        Recorder
}

func (h *Hub) HandleElectricity(ctx context.Context, data ElectricityMeter) error {
	h.ElectyMeterView.Merge(data)
	h.Recorder.Tabb(ctx, data)
	return nil
}
