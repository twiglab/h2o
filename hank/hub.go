package hank

import (
	"context"
	"log/slog"

	"github.com/twiglab/h2o/pkg/data"
)

type Sender interface {
	SendData(ctx context.Context, data data.Device) error
}

type Hub struct {
	DataLog *slog.Logger
	Sender  Sender
}

func (h *Hub) HandleDeviceStatus(ctx context.Context, data DeviceStatus) error {
	return nil
}

func (h *Hub) HandleDeviceData(ctx context.Context, data data.Device) error {
	h.DataLog.InfoContext(ctx, "deviceData", slog.Any("data", data))
	return h.Sender.SendData(ctx, data)
}
