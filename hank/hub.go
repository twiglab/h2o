package hank

import (
	"context"
	"fmt"
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
	fmt.Println(data)
	return nil
}

func (h *Hub) HandleDeviceData(ctx context.Context, data data.Device) error {
	fmt.Println(data.Code, data.Type, data.Time, data.UUID)
	return nil
}
