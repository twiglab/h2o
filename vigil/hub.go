package vigil

import (
	"context"
)

type ElectyMeterView interface {
	Merge(data ElectricityMeter)
}

type Hub struct {
	ElectyMeterView ElectyMeterView
}

func (h *Hub) HandleElectricity(ctx context.Context, data ElectricityMeter) error {
	h.ElectyMeterView.Merge(data)
	return nil
}
