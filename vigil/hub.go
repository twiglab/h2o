package vigil

import (
	"context"
	"log/slog"
)

type ElectyMeterView interface {
	Merge(data ElectricityMeter)
}

type Hub struct {
	ElectyMeterView ElectyMeterView
	Recorder        Recorder
	Logger          *slog.Logger

	BaseContext func(h *Hub) context.Context
}

func (h *Hub) HandleElectricity(ctx context.Context, data ElectricityMeter) error {
	h.ElectyMeterView.Merge(data)
	return h.Recorder.Tabb(ctx, data)
}
