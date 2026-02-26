package hank

import (
	"context"

	"github.com/twiglab/h2o/wal"
)

type Hub struct {
	WAL    *wal.WAL
	Sender Sender

	EP *ElectricityPacket
}

func (h *Hub) HandleDeviceStatus(ctx context.Context, data DeviceStatus) error {
	h.EP.SetStatus(data.No, Online(data.Status))
	return nil
}

func (h *Hub) HandleElectricity(ctx context.Context, data ElectricityMeter) error {
	h.WAL.WriteLogContext(ctx, wal.Type(data.Type), wal.Data(data))

	h.EP.Merge(data)

	return h.Sender.SendData(ctx, data)
}

func (h *Hub) HandleWater(ctx context.Context, data WaterMeter) error {
	h.WAL.WriteLogContext(ctx, wal.Type(data.Type), wal.Data(data))

	return h.Sender.SendData(ctx, data)
}
