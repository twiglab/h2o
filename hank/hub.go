package hank

import (
	"context"

	"github.com/twiglab/h2o/wal"
)

type Hub struct {
	WAL    *wal.WAL
	Sender Sender
}

func (h *Hub) HandleDeviceStatus(ctx context.Context, data DeviceStatus) error {
	return nil
}

func (h *Hub) HandleElectricity(ctx context.Context, data ElectricityMeter) error {
	h.WAL.WriteLogContext(ctx,
		wal.String("type", data.Type),
		wal.Any("data", data),
		wal.String("topic", data.Topic()))

	return h.Sender.SendData(ctx, data)
}

func (h *Hub) HandleWater(ctx context.Context, data WaterMeter) error {
	h.WAL.WriteLogContext(ctx,
		wal.String("type", data.Type),
		wal.Any("data", data),
		wal.String("topic", data.Topic()))

	return h.Sender.SendData(ctx, data)
}
