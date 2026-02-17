package hank

import (
	"context"
	"log/slog"
)

type Hub struct {
	DataLog *slog.Logger
	Sender  Sender
}

func (h *Hub) HandleDeviceStatus(ctx context.Context, data DeviceStatus) error {
	return nil
}

func (h *Hub) HandleElectricity(ctx context.Context, data ElectricityMeter) error {
	h.DataLog.InfoContext(ctx, "deviceData", slog.String("log", "datalog"), slog.String("type", data.Type), slog.Any("data", data))
	return h.Sender.SendData(ctx, data)
}

func (h *Hub) HandleWater(ctx context.Context, data WaterMeter) error {
	h.DataLog.InfoContext(ctx, "deviceData", slog.String("log", "datalog"), slog.String("type", data.Type), slog.Any("data", data))
	return h.Sender.SendData(ctx, data)
}
