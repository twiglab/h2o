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
	// 先记录在监控
	if err := h.Recorder.Tabb(ctx, data); err != nil {
		h.Logger.ErrorContext(ctx, "handle Electy error", slog.Any("data", data), slog.Any("error", err))
		return err
	}
	h.ElectyMeterView.Merge(data)
	return nil
}
