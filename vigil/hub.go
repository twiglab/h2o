package vigil

import (
	"context"
)

type ElectyMeterView interface {
	Merge(data ElectricityMeter)
}

type Hub struct {
	ElectyMeterView ElectyMeterView
	Recorder        Recorder
}

func (h *Hub) HandleElectricity(ctx context.Context, data ElectricityMeter) error {
	h.ElectyMeterView.Merge(data)
	return h.Recorder.Tabb(ctx, data)
}
