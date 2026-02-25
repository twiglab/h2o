package hank

import (
	"context"
	"log/slog"
)

type Hub struct {
	WAL    *slog.Logger
	Sender Sender

	EP *ElectricityPacket
}

func (h *Hub) HandleDeviceStatus(ctx context.Context, data DeviceStatus) error {
	h.EP.SetStatus(data.No, Online(data.Status))
	return nil
}

func (h *Hub) HandleElectricity(ctx context.Context, data ElectricityMeter) error {
	h.WAL.InfoContext(ctx, "deviceData", slog.String("log", "datalog"), slog.String("type", data.Type), slog.Any("data", data))

	h.EP.Add(data)

	return h.Sender.SendData(ctx, data)
}

func (h *Hub) HandleWater(ctx context.Context, data WaterMeter) error {
	h.WAL.InfoContext(ctx, "deviceData", slog.String("log", "datalog"), slog.String("type", data.Type), slog.Any("data", data))

	return h.Sender.SendData(ctx, data)
}
