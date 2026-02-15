package hank

import (
	"context"
	"encoding"
	"log/slog"

	"github.com/twiglab/h2o/pkg/common"
)

type Sender interface {
	SendData(ctx context.Context, data encoding.BinaryMarshaler) error
}

type Hub struct {
	DataLog *slog.Logger
	Sender  Sender
}

func (h *Hub) HandleDeviceStatus(ctx context.Context, data DeviceStatus) error {
	return nil
}

func (h *Hub) HandleElectricity(ctx context.Context, data common.ElectricityMeter) error {
	h.DataLog.InfoContext(ctx, "deviceData", slog.String("log", "datalog"), slog.String("type", data.Type), slog.Any("data", data))
	return h.Sender.SendData(ctx, data)
}

func (h *Hub) HandleWater(ctx context.Context, data common.WaterMeter) error {
	h.DataLog.InfoContext(ctx, "deviceData", slog.String("log", "datalog"), slog.String("type", data.Type), slog.Any("data", data))
	return h.Sender.SendData(ctx, data)
}
