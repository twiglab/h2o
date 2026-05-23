package vigil

import (
	"context"
	"log/slog"
)

type Hub struct {
	TSDB Recorder
	DB   Recorder

	Logger *slog.Logger

	BaseContext func(h *Hub) context.Context
}

func (h *Hub) HandleWater(ctx context.Context, data WaterMeter) error {
	h.TSDB.TabbWater(ctx, data)

	if err := h.DB.TabbWater(ctx, data); err != nil {
		h.Logger.ErrorContext(ctx, "Record Water error", slog.Any("data", data), slog.Any("error", err))
		return err
	}
	return nil
}

func (h *Hub) HandleElecty(ctx context.Context, data ElectricityMeter) error {
	h.TSDB.TabbElecty(ctx, data)

	if err := h.DB.TabbElecty(ctx, data); err != nil {
		h.Logger.ErrorContext(ctx, "Record Electy error", slog.Any("data", data), slog.Any("error", err))
		return err
	}
	return nil
}
