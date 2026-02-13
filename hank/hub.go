package hank

import (
	"context"
	"fmt"
	"log"
	"log/slog"

	"github.com/twiglab/h2o/pkg/data"
)

type Sender interface {
	SendData(ctx context.Context, data data.Device) error
}

type Hub struct {
	DataLog *slog.Logger
	InfoLog *slog.Logger
	Enh     *Enh
	Sender  Sender
}

func (h *Hub) HandleDeviceStatus(ctx context.Context, data DeviceStatus) error {
	fmt.Println(data)
	return nil
}

func (h *Hub) HandleUploadGatewayInfo(ctx context.Context, data GatewayInfo) error {
	fmt.Println(data)
	return nil
}

func (h *Hub) HandleDeviceData(ctx context.Context, data DeviceData) error {
	fmt.Println(data.Type, data.No, data.DataCode, data.DataTime, data.DataJson.DataValue)
	return nil
}

func doHandleDeviceStatusList(ctx context.Context, dsl DeviceStatusList, h *Hub) {
	go func() {
		for _, ds := range dsl {
			if err := h.HandleDeviceStatus(ctx, ds); err != nil {
				log.Println(err)
			}
		}
	}()
}

func doHandleDeviceDataList(ctx context.Context, ddl DeviceDataList, h *Hub) {
	go func() {
		for _, dd := range ddl {
			if err := h.HandleDeviceData(ctx, dd); err != nil {
				log.Println(err)
			}
		}
	}()
}
